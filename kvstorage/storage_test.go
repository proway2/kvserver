package kvstorage

import (
	"container/list"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/proway2/kvserver/element"
)

var KEYNAME string = "key1"
var KEYVALUE string = "key1 value"

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name string
		want *KVStorage
	}{
		{
			name: "Normal run",
			want: &KVStorage{
				kvstorage:   make(map[string]*element.Element),
				mux:         &sync.Mutex{},
				initialized: true,
				queue:       list.New(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStorage_Set(t *testing.T) {
	// because we need to test the case when key-value pair already in the storage - one storage will be in use by all testcases.
	storage := NewStorage()

	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name   string
		fields *KVStorage
		args   args
		want   bool
	}{
		{
			name:   "Storage is not initialized",
			fields: &KVStorage{initialized: false},
			args:   args{key: "test_key", value: "test_value"},
			want:   false,
		},
		{
			name:   "Empty key",
			fields: storage,
			args:   args{key: "", value: KEYVALUE},
			want:   false,
		},
		{
			name:   "Good key, empty value",
			fields: storage,
			args:   args{key: KEYNAME, value: ""},
			want:   true,
		},
		{
			// this key-value pair will be in use by next testcase "Key already in storage (update key-value pair)"
			name:   "Good key, good value",
			fields: storage,
			args:   args{key: "key2", value: "value for key2"},
			want:   true,
		},
		{
			// this testcase relies on the result of the previous testcase "Good key, good value"
			name:   "Key already in storage (update key-value pair)",
			fields: storage,
			args:   args{key: "key2", value: "new value for key2"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &tt.fields
			if got := (*kv).Set(tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("KVStorage.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStorage_Get(t *testing.T) {
	// because we need to test the case when key-value pair already in the storage - one storage will be in use by all testcases.
	storage := NewStorage()
	storage.Set(KEYNAME, KEYVALUE)

	// this storage is not initialized
	badStorage := NewStorage()
	badStorage.initialized = false

	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields *KVStorage
		args   args
		want   string
		want1  bool
	}{
		{
			name:   "Storage is not initialized",
			fields: badStorage,
			args:   args{"key1"},
			want:   "",
			want1:  false,
		},
		{
			name:   "Empty key",
			fields: storage,
			args:   args{""},
			want:   "",
			want1:  false,
		},
		{
			name:   "Key is in storage",
			fields: storage,
			args:   args{KEYNAME},
			want:   KEYVALUE,
			want1:  true,
		},
		{
			name:   "Key is not in storage",
			fields: storage,
			args:   args{KEYVALUE + "xxx"},
			want:   "",
			want1:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := tt.fields
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
	// because we need to test the case when key-value pair already in the storage - one storage will be in use by all testcases.
	goodStorage := NewStorage()
	goodStorage.Set(KEYNAME, KEYVALUE)
	elementTime := goodStorage.kvstorage[KEYNAME].Timestamp

	emptyStorage := NewStorage()

	badStorage := NewStorage()
	badStorage.initialized = false

	tests := []struct {
		name    string
		fields  *KVStorage
		want    time.Time
		wantErr bool
	}{
		{
			name:    "Storage is not initialized",
			fields:  badStorage,
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "Empty storage",
			fields:  emptyStorage,
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "Normal operation",
			fields:  goodStorage,
			want:    elementTime,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := tt.fields
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
	// because we need to test the case when key-value pair already in the storage - one storage will be in use by all testcases.
	storage := NewStorage()
	storage.Set(KEYNAME, KEYVALUE)

	// empty storage
	emptyStorage := NewStorage()

	// uninitialized storage
	badStorage := NewStorage()
	badStorage.initialized = false

	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields *KVStorage
		args   args
		want   bool
		want2  int // number of elements in the storage
	}{
		{
			name:   "Storage is not initialized",
			fields: badStorage,
			args:   args{KEYNAME},
			want:   false,
			want2:  0,
		},
		{
			name:   "Empty key, normal storage",
			fields: storage,
			args:   args{""},
			want:   false,
			want2:  1,
		},
		{
			name:   "Good key, empty storage",
			fields: emptyStorage,
			args:   args{KEYNAME},
			want:   false,
			want2:  0,
		},
		{
			name:   "Key is not found in storage (storage is not empty)",
			fields: storage,
			args:   args{KEYNAME + "xxx"},
			want:   false,
			want2:  1,
		},
		{
			name:   "Key is found in storage",
			fields: storage,
			args:   args{KEYNAME},
			want:   true,
			want2:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := tt.fields
			if got := kv.Delete(tt.args.key); got != tt.want {
				t.Errorf("KVStorage.Delete() = %v, want %v", got, tt.want)
			}
			if len := kv.queue.Len(); len != tt.want2 {
				t.Errorf("Queue length  = %v, want %v", len, tt.want2)
			}
		})
	}
}

func TestKVStorage_DeleteFrontIfOlder(t *testing.T) {
	// because we need to test the case when key-value pair already in the storage - one storage will be in use by all testcases.
	goodStorage := NewStorage()
	goodStorage.Set(KEYNAME, KEYVALUE)

	emptyStorage := NewStorage()

	badStorage := NewStorage()
	badStorage.initialized = false

	type args struct {
		ctxTime time.Time
	}
	tests := []struct {
		name   string
		fields *KVStorage
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{
			name:   "Storage is not initialized",
			fields: badStorage,
			want:   false,
		},
		{
			name:   "Empty storage",
			fields: emptyStorage,
			want:   false,
		},
		{
			name:   "Element younger than NOW-60 secs",
			fields: goodStorage,
			args:   args{time.Now().Add(-60 * time.Second)},
			want:   false,
		},
		{
			name:   "Element older than NOW",
			fields: goodStorage,
			args:   args{time.Now()},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := tt.fields
			if got := kv.DeleteFrontIfOlder(tt.args.ctxTime); got != tt.want {
				t.Errorf("KVStorage.DeleteFrontIfOlder() = %v, want %v", got, tt.want)
			}
		})
	}
}
