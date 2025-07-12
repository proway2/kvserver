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
		name    string
		fields  *KVStorage
		args    args
		wantErr bool
	}{
		{
			name:    "Storage is not initialized",
			fields:  &KVStorage{initialized: false},
			args:    args{key: KEYNAME, value: KEYVALUE},
			wantErr: true,
		},
		{
			name:    "Empty key",
			fields:  storage,
			args:    args{key: "", value: KEYVALUE},
			wantErr: true,
		},
		{
			name:    "Good key, empty value",
			fields:  storage,
			args:    args{key: KEYNAME, value: ""},
			wantErr: false,
		},
		{
			// this key-value pair will be in use by next testcase "Key already in storage (update key-value pair)"
			name:    "Good key, good value",
			fields:  storage,
			args:    args{key: "key2", value: "value for key2"},
			wantErr: false,
		},
		{
			// this testcase relies on the result of the previous testcase "Good key, good value"
			name:    "Key already in storage (update key-value pair)",
			fields:  storage,
			args:    args{key: "key2", value: "new value for key2"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &tt.fields
			if err := (*kv).Set(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("KVStorage.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKVStorage_Get(t *testing.T) {
	// because we need to test the case when key-value pair already in the storage - one storage will be in use by all testcases.
	goodStorage := NewStorage()
	err := goodStorage.Set(KEYNAME, KEYVALUE)
	check(err, t)
	// empty value for the key
	keyForEmptyValue := "empty key"
	err = goodStorage.Set(keyForEmptyValue, "")
	check(err, t)

	// this storage is not initialized
	badStorage := NewStorage()
	badStorage.initialized = false

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  *KVStorage
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "Storage is not initialized",
			fields:  badStorage,
			args:    args{KEYNAME},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty key",
			fields:  goodStorage,
			args:    args{""},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Key is in the storage, value is not empty",
			fields:  goodStorage,
			args:    args{KEYNAME},
			want:    []byte(KEYVALUE),
			wantErr: false,
		},
		{
			name:    "Key is in the storage, value is empty",
			fields:  goodStorage,
			args:    args{keyForEmptyValue},
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "Key is not in the storage",
			fields:  goodStorage,
			args:    args{KEYVALUE + "xxx"},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &tt.fields
			got, err := (*kv).Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("KVStorage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KVStorage.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStorage_OldestElementTime(t *testing.T) {
	// because we need to test the case when key-value pair already in the storage - one storage will be in use by all testcases.
	goodStorage := NewStorage()
	err := goodStorage.Set(KEYNAME, KEYVALUE)
	check(err, t)
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
	goodStorage := NewStorage()
	err := goodStorage.Set(KEYNAME, KEYVALUE)
	check(err, t)

	// empty storage
	emptyStorage := NewStorage()

	// uninitialized storage
	badStorage := NewStorage()
	badStorage.initialized = false

	tests := []struct {
		name    string
		fields  *KVStorage
		key     string
		want    bool
		wantErr bool
		wantLen int // number of elements in the storage
	}{
		{
			name:    "Storage is not initialized",
			fields:  badStorage,
			key:     KEYNAME,
			want:    false,
			wantErr: true,
			wantLen: 0,
		},
		{
			name:    "Empty key, normal storage",
			fields:  goodStorage,
			key:     "",
			want:    false,
			wantErr: true,
			wantLen: 0,
		},
		{
			name:    "Good key, empty storage",
			fields:  emptyStorage,
			key:     KEYNAME,
			want:    false,
			wantErr: false,
			wantLen: 0,
		},
		{
			name:    "Key is not found in storage (storage is not empty)",
			fields:  goodStorage,
			key:     KEYNAME + "xxx",
			want:    false,
			wantErr: false,
			wantLen: 1,
		},
		{
			name:    "Key is found in storage",
			fields:  goodStorage,
			key:     KEYNAME,
			want:    true,
			wantErr: false,
			wantLen: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := tt.fields
			got, err := kv.Delete(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("KVStorage.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len := kv.queue.Len(); got != tt.want && len != tt.wantLen {
				t.Errorf("KVStorage.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStorage_DeleteFrontIfOlder(t *testing.T) {
	// because we need to test the case when key-value pair already in the storage - one storage will be in use by all testcases.
	goodStorage := NewStorage()
	err := goodStorage.Set(KEYNAME, KEYVALUE)
	check(err, t)

	emptyStorage := NewStorage()

	badStorage := NewStorage()
	badStorage.initialized = false

	type args struct {
		ctxTime time.Time
	}
	tests := []struct {
		name    string
		fields  *KVStorage
		args    args
		want    bool
		wantErr bool
		wantLen int
	}{
		{
			name:    "Storage is not initialized",
			fields:  badStorage,
			want:    false,
			wantErr: true,
			wantLen: 0,
		},
		{
			name:    "Empty storage",
			fields:  emptyStorage,
			want:    false,
			wantErr: false,
			wantLen: 0,
		},
		{
			name:    "Element younger than NOW-60 secs (no delete)",
			fields:  goodStorage,
			args:    args{time.Now().Add(-60 * time.Second)},
			want:    false,
			wantErr: false,
			wantLen: 1,
		},
		{
			name:    "Element older than NOW (delete)",
			fields:  goodStorage,
			args:    args{time.Now()},
			want:    true,
			wantErr: false,
			wantLen: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := tt.fields
			got, err := kv.DeleteFrontIfOlder(tt.args.ctxTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("KVStorage.DeleteFrontIfOlder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len := kv.queue.Len(); got != tt.want && len != tt.wantLen {
				t.Errorf("KVStorage.DeleteFrontIfOlder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func check(e error, t *testing.T) {
	if e != nil {
		t.Error("Something is wrong with tests")
	}
}
