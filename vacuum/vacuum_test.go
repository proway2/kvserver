package vacuum

import (
	"errors"
	"testing"
	"time"

	"github.com/proway2/kvserver/kvstorage"
)

func TestVacuum_Init(t *testing.T) {
	type fields struct {
		storage     *kvstorage.KVStorage
		ttl         uint64
		ttlDelim    uint
		initialized bool
	}
	type args struct {
		stor *kvstorage.KVStorage
		ttl  uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// Test cases.
		{
			name: "No storage provided.",
			args: args{
				ttl: 2,
			},
			want: false,
		},
		{
			name: "Wrong TTL < 1",
			args: args{
				stor: &kvstorage.KVStorage{},
				ttl:  0,
			},
			want: false,
		},
		{
			name: "Already initialized",
			args: args{
				stor: &kvstorage.KVStorage{},
				ttl:  10,
			},
			fields: fields{initialized: true},
			want:   false,
		},
		{
			name: "Normal operation",
			args: args{
				stor: &kvstorage.KVStorage{},
				ttl:  10,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Vacuum{
				storage:     tt.fields.storage,
				ttl:         tt.fields.ttl,
				ttlDelim:    tt.fields.ttlDelim,
				initialized: tt.fields.initialized,
			}
			if got := q.Init(tt.args.stor, tt.args.ttl); got != tt.want {
				t.Errorf("Vacuum.Init() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestVacuum_Run(t *testing.T) {
// 	type fields struct {
// 		storage     *kvstorage.KVStorage
// 		ttl         uint64
// 		ttlDelim    uint
// 		initialized bool
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 	}{
// 		// Test cases.
// 		{
// 			name:   "Already initialized",
// 			fields: fields{initialized: false},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			q := &Vacuum{
// 				storage:     tt.fields.storage,
// 				ttl:         tt.fields.ttl,
// 				ttlDelim:    tt.fields.ttlDelim,
// 				initialized: tt.fields.initialized,
// 			}
// 			q.Run()
// 		})
// 	}
// }

func Test_getSleepPeriod(t *testing.T) {
	type args struct {
		elementTime time.Time
		err         error
		ttl         uint64
		ttlDelim    uint
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		// Test cases.
		{
			name: "No elements in queue",
			args: args{
				err:      errors.New("Testing: No elements in queue"),
				ttl:      60,
				ttlDelim: 2,
			},
			want: time.Duration(30 * time.Second),
		},
		{
			name: "No TTL provided",
			args: args{
				elementTime: time.Now(),
				err:         nil,
				ttlDelim:    3,
			},
			want: time.Duration(1 * time.Second),
		},
		{
			name: "No TTL delimiter provided",
			args: args{
				elementTime: time.Now(),
				err:         nil,
				ttl:         60,
			},
			want: time.Duration(1 * time.Second),
		},
		{
			name: "Expired element",
			args: args{
				elementTime: time.Now().Add((-100 * time.Second)),
				err:         nil,
				ttl:         60,
				ttlDelim:    2,
			},
			want: time.Duration(0 * time.Nanosecond),
		},
		{
			name: "30 secs to now",
			args: args{
				elementTime: time.Now().Add((-30 * time.Second)),
				err:         nil,
				ttl:         60,
				ttlDelim:    2,
			},
			want: time.Duration(15 * time.Second),
		},
	}
	var toleranceNS int64 = 1000000
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := getSleepPeriod(tt.args.elementTime, tt.args.err, tt.args.ttl, tt.args.ttlDelim); got != tt.want {
			got := getSleepPeriod(tt.args.elementTime, tt.args.err, tt.args.ttl, tt.args.ttlDelim)
			if got.Nanoseconds() < tt.want.Nanoseconds()-toleranceNS || got.Nanoseconds() > tt.want.Nanoseconds()+toleranceNS {
				t.Errorf("getSleepPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}
