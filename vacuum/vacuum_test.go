package vacuum

import (
	"kvserver/kvstorage"
	"testing"
)

func createChanString() *chan string {
	out := make(chan string, 10)
	return &out
}

func TestLifo_Init(t *testing.T) {
	type args struct {
		stor *kvstorage.KVStorage
		in   *chan string
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
				in:   createChanString(),
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

// func TestLifo_Run(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		q    *Lifo
// 		want bool
// 	}{
// 		{
// 			name: "Канал = хранилище = nil, ttl = 0",
// 			q: &Lifo{
// 				inpElemChan: nil,
// 				storage:     nil,
// 				ttl:         0,
// 			},
// 			want: false,
// 		},
// 		// {
// 		// 	name: "Канал = хранилище = nil, ttl = 0",
// 		// 	q: &Lifo{
// 		// 		inpElemChan: createChanString(),
// 		// 		storage:     &kvstorage.KVStorage{},
// 		// 		ttl:         2,
// 		// 	},
// 		// 	want: false,
// 		// },
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.q.Run()
// 		})
// 	}
// }

// func TestLifo_cleanUp(t *testing.T) {
// 	type args struct {
// 		elem *list.Element
// 	}
// 	tests := []struct {
// 		name string
// 		q    *Lifo
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.q.cleanUp(tt.args.elem)
// 		})
// 	}
// }
