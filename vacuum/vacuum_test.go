package vacuum

import (
	"reflect"
	"testing"
	"time"

	"github.com/proway2/kvserver/kvstorage"
)

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
			got := getSleepPeriod(tt.args.elementTime, tt.args.err, tt.args.ttl, tt.args.ttlDelim)
			if got.Nanoseconds() < tt.want.Nanoseconds()-toleranceNS || got.Nanoseconds() > tt.want.Nanoseconds()+toleranceNS {
				t.Errorf("getSleepPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCleaner(t *testing.T) {
	type args struct {
		w   writer
		ttl uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *Vacuum
		wantErr bool
	}{
		{
			name: "No storage provided.",
			args: args{
				w:   nil,
				ttl: 10,
			},
			want:    &Vacuum{},
			wantErr: true,
		},
		{
			name: "Wrong TTL < 1",
			args: args{
				w:   kvstorage.NewStorage(),
				ttl: 0,
			},
			want:    &Vacuum{},
			wantErr: true,
		},
		{
			name: "Normal operation",
			args: args{
				w:   kvstorage.NewStorage(),
				ttl: 20,
			},
			want: &Vacuum{
				storage:     kvstorage.NewStorage(),
				ttl:         20,
				ttlDelim:    2,
				initialized: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCleaner(tt.args.w, tt.args.ttl)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCleaner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCleaner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSleepPeriodEmptyQueue(t *testing.T) {
	type args struct {
		ttl      uint64
		ttlDelim uint
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "Normal operation",
			args: args{
				ttl:      60,
				ttlDelim: 2,
			},
			want: time.Duration(30 * time.Second),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSleepPeriodEmptyQueue(tt.args.ttl, tt.args.ttlDelim); got != tt.want {
				t.Errorf("getSleepPeriodEmptyQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}
