package vacuum

import (
	"container/list"
	"kvserver/element"
	"kvserver/kvstorage"
	"testing"
	"time"
)

func fillStorage(kvs *kvstorage.KVStorage) {
	kvals := []struct {
		elem *element.Element
		key  string
	}{
		{
			key: "key111",
			elem: &element.Element{
				Val: "key111 value",
			},
		},
		{
			key: "key222",
			elem: &element.Element{
				Val: "key222 value 123456",
			},
		},
		{
			key: "key333",
			elem: &element.Element{
				Val: "key333 value 123456",
			},
		},
		{
			key: "empty key",
			elem: &element.Element{
				Val: "",
			},
		},
	}
	for _, kv := range kvals {
		kvs.Set(kv.key, kv.elem.Val)
	}
}

func TestLifo_Init(t *testing.T) {
	type args struct {
		stor *kvstorage.KVStorage
		in   chan string
		ttl  uint64
	}
	tests := []struct {
		name string
		q    *Lifo
		args args
		want bool
	}{
		{
			name: "Хранилище = канал = nil, ttl = 0",
			q:    &Lifo{},
			args: args{stor: nil, in: nil, ttl: 0},
			want: false,
		},
		{
			name: "Хранилище, канал != nil, ttl != 0",
			q:    &Lifo{},
			args: args{
				stor: &kvstorage.KVStorage{},
				in:   make(chan string, 10),
				ttl:  10,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Init(tt.args.stor, tt.args.in, tt.args.ttl); got != tt.want {
				t.Errorf("Lifo.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLifo_Run(t *testing.T) {
	tests := []struct {
		name string
		q    *Lifo
		want bool
	}{
		{
			name: "Канал = хранилище = nil, ttl = 0",
			q: &Lifo{
				inpElemChan: nil,
				storage:     nil,
				ttl:         0,
			},
			want: false,
		},
		// {
		// 	name: "Канал = хранилище = nil, ttl = 0",
		// 	q: &Lifo{
		// 		inpElemChan: createChanString(),
		// 		storage:     &kvstorage.KVStorage{},
		// 		ttl:         2,
		// 	},
		// 	want: false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Run()
		})
	}
}

func TestLifo_cleanUp(t *testing.T) {
	chanElm := make(chan string, 10)

	emptyKVS := kvstorage.KVStorage{}

	fullKVS := kvstorage.KVStorage{}
	fullKVS.Init(chanElm)
	fillStorage(&fullKVS)

	elementNO := list.Element{Value: "test123"}
	elementUPD := list.Element{Value: "key222"}
	elementUPDTTLover := list.Element{Value: "key333"}

	type args struct {
		elem *list.Element
	}
	tests := []struct {
		name string
		q    *Lifo
		args args
	}{
		{
			name: "Элемента нет в хранилище, хранилище пустое",
			q: &Lifo{
				storage:     &emptyKVS,
				queue:       list.List{},
				inpElemChan: chanElm,
				ttl:         1,
			},
			args: args{elem: &elementNO},
		},
		{
			name: "Элемента нет в хранилище, хранилище не пустое",
			q: &Lifo{
				storage:     &fullKVS,
				queue:       list.List{},
				inpElemChan: chanElm,
				ttl:         1,
			},
			args: args{elem: &elementNO},
		},
		{
			name: "Элемента в хранилище и обновлен, но TTL не вышел",
			q: &Lifo{
				storage:     &fullKVS,
				queue:       list.List{},
				inpElemChan: chanElm,
				ttl:         1,
			},
			args: args{elem: &elementUPD},
		},
		{
			name: "Элемента в хранилище и обновлен, но TTL вышел",
			q: &Lifo{
				storage:     &fullKVS,
				queue:       list.List{},
				inpElemChan: chanElm,
				ttl:         1,
			},
			args: args{elem: &elementUPDTTLover},
		},
	}
	// иначе не будет работать TTL
	time.Sleep(2 * time.Second)
	// обновляем метку элемента
	_ = fullKVS.Set(elementUPD.Value.(string), "new value 123")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.cleanUp(tt.args.elem)
		})
	}
}
