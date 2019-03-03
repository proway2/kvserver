package kvstorage

import (
	"container/list"
	"kvserver/element"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestKVStorage_Init(t *testing.T) {
	type fields struct {
		kvstorage   map[string]*element.Element
		mux         sync.Mutex
		queue       list.List
		initialized bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// Test cases.
		{
			name:   "Already initialized",
			fields: fields{initialized: true},
			want:   false,
		},
		{
			name: "Normal run",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStorage{
				kvstorage:   tt.fields.kvstorage,
				queue:       tt.fields.queue,
				initialized: tt.fields.initialized,
			}
			if got := kv.Init(); got != tt.want {
				t.Errorf("KVStorage.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStorage_Set(t *testing.T) {
	var kv KVStorage
	kv.Init()
	kv.Set("key1", "new key1 value")

	type fields struct {
		kvstorage   map[string]*element.Element
		mux         sync.Mutex
		queue       list.List
		initialized bool
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// Test cases.
		{
			name: "Storage is not initialized",
			want: false,
		},
		{
			name: "Empty key",
			args: args{"", "empty key"},
			fields: fields{
				kvstorage:   make(map[string]*element.Element),
				initialized: true,
			},
			want: false,
		},
		{
			name: "Appropriate key, empty value",
			args: args{"key1", ""},
			fields: fields{
				kvstorage:   make(map[string]*element.Element),
				initialized: true,
			},
			want: true,
		},
		{
			name: "Key already in storage",
			args: args{"key1", "new key1 value"},
			fields: fields{
				kvstorage:   kv.kvstorage,
				initialized: true,
				queue:       kv.queue,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStorage{
				kvstorage:   tt.fields.kvstorage,
				queue:       tt.fields.queue,
				initialized: tt.fields.initialized,
			}
			if got := kv.Set(tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("KVStorage.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStorage_Get(t *testing.T) {
	var kv KVStorage
	kv.Init()
	kv.Set("key11", "new key11 value")

	type fields struct {
		kvstorage   map[string]*element.Element
		mux         sync.Mutex
		queue       list.List
		initialized bool
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		// Test cases.
		{
			name:   "Storage is not initialized",
			fields: fields{initialized: true},
			want:   "",
			want1:  false,
		},
		{
			name:  "Empty key",
			args:  args{""},
			want:  "",
			want1: false,
		},
		{
			name: "Key is not in storage",
			args: args{"unknown key"},
			fields: fields{
				kvstorage:   make(map[string]*element.Element),
				initialized: true,
			},
			want:  "",
			want1: false,
		},
		{
			name: "Normal key and in storage",
			args: args{"key11"},
			fields: fields{
				kvstorage:   kv.kvstorage,
				queue:       kv.queue,
				initialized: true,
			},
			want:  "new key11 value",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStorage{
				kvstorage:   tt.fields.kvstorage,
				queue:       tt.fields.queue,
				initialized: tt.fields.initialized,
			}
			got, got1 := kv.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("KVStorage.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("KVStorage.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestKVStorage_OldestElementTime(t *testing.T) {
	var kv KVStorage
	kv.Init()
	kv.Set("key11", "new key11 value")
	elementTime := kv.kvstorage["key11"].Timestamp

	type fields struct {
		kvstorage   map[string]*element.Element
		mux         sync.Mutex
		queue       list.List
		initialized bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    time.Time
		wantErr bool
	}{
		// Test cases.
		{
			name:    "Storage is not initialized",
			fields:  fields{initialized: false},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "Empty storage",
			fields: fields{
				kvstorage:   make(map[string]*element.Element),
				initialized: true,
				queue:       *list.New(),
			},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name: "Normal operation",
			fields: fields{
				kvstorage:   kv.kvstorage,
				queue:       kv.queue,
				initialized: true,
			},
			want:    elementTime,
			wantErr: false,
		},
		{
			name: "INNORMAL OPERATION, NOT POSSIBLE!!!",
			fields: fields{
				kvstorage:   make(map[string]*element.Element),
				initialized: true,
				queue:       kv.queue,
			},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStorage{
				kvstorage:   tt.fields.kvstorage,
				queue:       tt.fields.queue,
				initialized: tt.fields.initialized,
			}
			got, err := kv.OldestElementTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("KVStorage.OldestElementTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KVStorage.OldestElementTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStorage_Delete(t *testing.T) {
	var kv KVStorage
	kv.Init()
	kv.Set("key11", "new key11 value")

	type fields struct {
		kvstorage   map[string]*element.Element
		mux         sync.Mutex
		queue       list.List
		initialized bool
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// Test cases.
		{
			name:   "Storage is not initialized",
			fields: fields{initialized: false},
			args:   args{"key1"},
			want:   false,
		},
		{
			name:   "Empty key",
			fields: fields{initialized: true},
			args:   args{""},
			want:   false,
		},
		{
			name: "Good key, empty storage",
			fields: fields{
				kvstorage:   make(map[string]*element.Element),
				initialized: true,
			},
			args: args{"key1"},
			want: false,
		},
		{
			name: "Not found in storage",
			fields: fields{
				kvstorage:   kv.kvstorage,
				queue:       kv.queue,
				initialized: true,
			},
			args: args{"not found"},
			want: false,
		},
		{
			name: "Normal operation",
			fields: fields{
				kvstorage:   kv.kvstorage,
				queue:       kv.queue,
				initialized: true,
			},
			args: args{"key11"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStorage{
				kvstorage:   tt.fields.kvstorage,
				queue:       tt.fields.queue,
				initialized: tt.fields.initialized,
			}
			if got := kv.Delete(tt.args.key); got != tt.want {
				t.Errorf("KVStorage.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStorage_DeleteFrontIfOlder(t *testing.T) {
	var kv KVStorage
	kv.Init()
	kv.Set("key11", "new key11 value")

	type fields struct {
		kvstorage   map[string]*element.Element
		mux         sync.Mutex
		queue       list.List
		initialized bool
	}
	type args struct {
		ctxTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// Test cases.
		{
			name:   "Storage is not initialized",
			fields: fields{initialized: false},
			want:   false,
		},
		{
			name: "Empty storage",
			fields: fields{
				kvstorage:   make(map[string]*element.Element),
				initialized: true,
				queue:       *list.New(),
			},
			args: args{time.Now()},
			want: false,
		},
		{
			name: "Element younger than NOW - 60 secs",
			fields: fields{
				kvstorage:   kv.kvstorage,
				initialized: true,
				queue:       kv.queue,
			},
			args: args{time.Now().Add(-60 * time.Second)},
			want: false,
		},
		{
			name: "Element older than NOW",
			fields: fields{
				kvstorage:   kv.kvstorage,
				initialized: true,
				queue:       kv.queue,
			},
			args: args{time.Now()},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStorage{
				kvstorage:   tt.fields.kvstorage,
				queue:       tt.fields.queue,
				initialized: tt.fields.initialized,
			}
			if got := kv.DeleteFrontIfOlder(tt.args.ctxTime); got != tt.want {
				t.Errorf("KVStorage.DeleteFrontIfOlder() = %v, want %v", got, tt.want)
			}
		})
	}
}
