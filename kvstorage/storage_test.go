package kvstorage

import (
	"kvserver/element"
	"sync"
	"testing"
	"time"
)

func createChanString() *chan string {
	out := make(chan string, 10)
	return &out
}

func createStorageMap() map[string]*element.Element {
	// res := make(map[string]*element.Element)
	// return res
	return make(map[string]*element.Element)
}

func fillStorageMap(storage *map[string]*element.Element) {
	kvals := []struct {
		elem *element.Element
		key  string
	}{
		{
			key: "key111",
			elem: &element.Element{
				Val:       "key111 value",
				Updated:   false,
				Timestamp: time.Now().Unix() - 100,
			},
		},
		{
			key: "key222",
			elem: &element.Element{
				Val:       "key222 value 123456",
				Updated:   true,
				Timestamp: time.Now().Unix() + 100,
			},
		},
		{
			key: "empty key",
			elem: &element.Element{
				Val:     "",
				Updated: true,
			},
		},
	}
	for _, kv := range kvals {
		(*storage)[kv.key] = kv.elem
	}
}

func TestKVStorage_Init(t *testing.T) {
	type fields struct {
		kvstorage  map[string]*element.Element
		mux        sync.Mutex
		outElmChan *chan string
	}
	type args struct {
		out *chan string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Канал nil",
			fields: fields{outElmChan: nil},
			args:   args{nil},
			want:   false,
		},
		{
			name:   "Канал не nil",
			fields: fields{},
			args:   args{createChanString()},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStorage{
				kvstorage:  tt.fields.kvstorage,
				mux:        tt.fields.mux,
				outElmChan: tt.fields.outElmChan,
			}
			if got := kv.Init(tt.args.out); got != tt.want {
				t.Errorf("KVStorage.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStorage_Set(t *testing.T) {
	type fields struct {
		kvstorage  map[string]*element.Element
		mux        sync.Mutex
		outElmChan *chan string
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
		{
			name: "Пустой ключ",
			fields: fields{
				kvstorage:  make(map[string]*element.Element),
				outElmChan: createChanString(),
			},
			args: args{"", "empty key"},
			want: false,
		},
		{
			name: "Ключ длиной > 0, канал nil",
			fields: fields{
				kvstorage:  make(map[string]*element.Element),
				outElmChan: nil,
			},
			args: args{"key1", "key1 value"},
			want: true,
		},
		{
			name: "Ключ длиной > 0, канал действующий",
			fields: fields{
				kvstorage:  make(map[string]*element.Element),
				outElmChan: createChanString(),
			},
			args: args{"key2", "key2 value"},
			want: true,
		},
		{
			name: "Ключ длиной > 0, пустое значение, канал действующий",
			fields: fields{
				kvstorage:  make(map[string]*element.Element),
				outElmChan: createChanString(),
			},
			args: args{"key3", ""},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStorage{
				kvstorage:  tt.fields.kvstorage,
				mux:        tt.fields.mux,
				outElmChan: tt.fields.outElmChan,
			}
			if got := kv.Set(tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("KVStorage.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStorage_Get(t *testing.T) {
	storage := createStorageMap()
	fillStorageMap(&storage)
	type fields struct {
		kvstorage  map[string]*element.Element
		mux        sync.Mutex
		outElmChan *chan string
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
		{
			name:   "Пустой ключ",
			fields: fields{},
			args:   args{""},
			want:   "",
			want1:  false,
		},
		{
			name: "Ключ длиной > 0, хранилище пустое",
			fields: fields{
				kvstorage: make(map[string]*element.Element),
			},
			args:  args{"key1"},
			want:  "",
			want1: false,
		},
		{
			name: "Ключ длиной > 0, в хранилище отсутствует",
			fields: fields{
				kvstorage:  storage,
				outElmChan: createChanString(),
			},
			args:  args{"keyNNN"},
			want:  "",
			want1: false,
		},
		{
			name: "Ключ длиной > 0, в хранилище",
			fields: fields{
				kvstorage:  storage,
				outElmChan: createChanString(),
			},
			args:  args{"key111"},
			want:  "key111 value",
			want1: true,
		},
		{
			name: "Ключ длиной > 0, в хранилище, значение пустое",
			fields: fields{
				kvstorage:  storage,
				outElmChan: createChanString(),
			},
			args:  args{"empty key"},
			want:  "",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStorage{
				kvstorage:  tt.fields.kvstorage,
				mux:        tt.fields.mux,
				outElmChan: tt.fields.outElmChan,
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

func TestKVStorage_Delete(t *testing.T) {
	storage := createStorageMap()
	fillStorageMap(&storage)
	type fields struct {
		kvstorage  map[string]*element.Element
		mux        sync.Mutex
		outElmChan *chan string
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
		{
			name:   "Пустой ключ",
			fields: fields{},
			args:   args{""},
			want:   false,
		},
		{
			name: "Ключ длиной > 0, хранилище пустое",
			fields: fields{
				kvstorage: make(map[string]*element.Element),
			},
			args: args{"key1"},
			want: false,
		},
		{
			name: "Ключ длиной > 0, в хранилище отсутствует",
			fields: fields{
				kvstorage: storage,
			},
			args: args{"keyNNN"},
			want: false,
		},
		{
			name: "Ключ длиной > 0, в хранилище",
			fields: fields{
				kvstorage: storage,
			},
			args: args{"key111"},
			want: true,
		},
		{
			name: "Ключ длиной > 0, в хранилище, значение пустое",
			fields: fields{
				kvstorage: storage,
			},
			args: args{"empty key"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStorage{
				kvstorage:  tt.fields.kvstorage,
				mux:        tt.fields.mux,
				outElmChan: tt.fields.outElmChan,
			}
			if got := kv.Delete(tt.args.key); got != tt.want {
				t.Errorf("KVStorage.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStorage_ResetUpdated(t *testing.T) {
	storage := createStorageMap()
	fillStorageMap(&storage)
	type fields struct {
		kvstorage  map[string]*element.Element
		mux        sync.Mutex
		outElmChan *chan string
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
		{
			name:   "Пустой ключ",
			fields: fields{},
			args:   args{""},
			want:   false,
		},
		{
			name: "Ключ длиной > 0, хранилище пустое",
			fields: fields{
				kvstorage: make(map[string]*element.Element),
			},
			args: args{"key1"},
			want: false,
		},
		{
			name: "Ключ длиной > 0, в хранилище отсутствует",
			fields: fields{
				kvstorage: storage,
			},
			args: args{"keyNNN"},
			want: false,
		},
		{
			name: "Ключ длиной > 0, в хранилище",
			fields: fields{
				kvstorage: storage,
			},
			args: args{"key111"},
			want: true,
		},
		{
			name: "Ключ длиной > 0, в хранилище, значение пустое",
			fields: fields{
				kvstorage: storage,
			},
			args: args{"empty key"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStorage{
				kvstorage:  tt.fields.kvstorage,
				mux:        tt.fields.mux,
				outElmChan: tt.fields.outElmChan,
			}
			if got := kv.ResetUpdated(tt.args.key); got != tt.want {
				t.Errorf(
					"KVStorage.ResetUpdated() = %v, want %v", got, tt.want,
				)
			}

		})
	}
}

func TestKVStorage_IsElemUpdated(t *testing.T) {
	storage := createStorageMap()
	fillStorageMap(&storage)
	type fields struct {
		kvstorage  map[string]*element.Element
		mux        sync.Mutex
		outElmChan *chan string
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
		{
			name:   "Пустой ключ",
			fields: fields{},
			args:   args{""},
			want:   false,
		},
		{
			name: "Ключ длиной > 0, хранилище пустое",
			fields: fields{
				kvstorage: make(map[string]*element.Element),
			},
			args: args{"key1"},
			want: false,
		},
		{
			name: "Ключ длиной > 0, в хранилище отсутствует",
			fields: fields{
				kvstorage: storage,
			},
			args: args{"keyNNN"},
			want: false,
		},
		{
			name: "Ключ длиной > 0, в хранилище, элемент не обновлен",
			fields: fields{
				kvstorage: storage,
			},
			args: args{"key111"},
			want: false,
		},
		{
			name: "Ключ длиной > 0, в хранилище, элемент обновлен",
			fields: fields{
				kvstorage: storage,
			},
			args: args{"empty key"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStorage{
				kvstorage:  tt.fields.kvstorage,
				mux:        tt.fields.mux,
				outElmChan: tt.fields.outElmChan,
			}
			if got := kv.IsElemUpdated(tt.args.key); got != tt.want {
				t.Errorf("KVStorage.IsElemUpdated() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStorage_IsElemTTLOver(t *testing.T) {
	storage := createStorageMap()
	fillStorageMap(&storage)
	type fields struct {
		kvstorage  map[string]*element.Element
		mux        sync.Mutex
		outElmChan *chan string
	}
	type args struct {
		key string
		ttl uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Пустой ключ",
			fields: fields{},
			args:   args{key: "", ttl: 10},
			want:   false,
		},
		{
			name: "Ключ длиной > 0, хранилище пустое",
			fields: fields{
				kvstorage: make(map[string]*element.Element),
			},
			args: args{key: "key1", ttl: 10},
			want: true,
		},
		{
			name: "Ключ длиной > 0, в хранилище отсутствует",
			fields: fields{
				kvstorage: storage,
			},
			args: args{key: "keyNNN", ttl: 10},
			want: true,
		},
		{
			name: "Устаревший элемент",
			fields: fields{
				kvstorage: storage,
			},
			args: args{key: "key111", ttl: 10},
			want: true,
		},
		{
			name: "Новый элемент",
			fields: fields{
				kvstorage: storage,
			},
			args: args{key: "key222", ttl: 10},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStorage{
				kvstorage:  tt.fields.kvstorage,
				mux:        tt.fields.mux,
				outElmChan: tt.fields.outElmChan,
			}
			if got := kv.IsElemTTLOver(tt.args.key, tt.args.ttl); got != tt.want {
				t.Errorf("KVStorage.IsElemTTLOver() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKVStorage_IsInStorage(t *testing.T) {
	storage := createStorageMap()
	fillStorageMap(&storage)
	type fields struct {
		kvstorage  map[string]*element.Element
		mux        sync.Mutex
		outElmChan *chan string
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
		{
			name:   "Пустой ключ",
			fields: fields{},
			args:   args{""},
			want:   false,
		},
		{
			name: "Ключ длиной > 0, хранилище пустое",
			fields: fields{
				kvstorage: make(map[string]*element.Element),
			},
			args: args{"key1"},
			want: false,
		},
		{
			name: "Ключ длиной > 0, хранилище пустое",
			fields: fields{
				kvstorage: make(map[string]*element.Element),
			},
			args: args{"key1"},
			want: false,
		},
		{
			name: "Ключ длиной > 0, в хранилище отсутствует",
			fields: fields{
				kvstorage: storage,
			},
			args: args{"keyNNN"},
			want: false,
		},
		{
			name: "Ключ длиной > 0, в хранилище",
			fields: fields{
				kvstorage: storage,
			},
			args: args{"key111"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kv := &KVStorage{
				kvstorage:  tt.fields.kvstorage,
				mux:        tt.fields.mux,
				outElmChan: tt.fields.outElmChan,
			}
			if got := kv.IsInStorage(tt.args.key); got != tt.want {
				t.Errorf("KVStorage.IsInStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}
