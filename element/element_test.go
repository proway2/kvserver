package element

import (
	"testing"
	"time"
)

func TestElement_IsTTLOver(t *testing.T) {
	type fields struct {
		Val       string
		Timestamp int64
		Updated   bool
	}
	type args struct {
		ttl uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{
			name:   "Значение TTL = 0",
			fields: fields{Val: "test1", Timestamp: time.Now().Unix()},
			args:   args{ttl: 0},
			want:   true,
		},
		{
			name:   "Значение TTL = 1, элемент устаревший",
			fields: fields{Val: "test2", Timestamp: time.Now().Unix() - 100},
			args:   args{ttl: 1},
			want:   true,
		},
		{
			name:   "Значение TTL = 1, элемент новый",
			fields: fields{Val: "test3", Timestamp: time.Now().Unix()},
			args:   args{ttl: 1},
			want:   false,
		},
		{
			name:   "Значение TTL = 10, элемент устаревший",
			fields: fields{Val: "test4", Timestamp: time.Now().Unix() - 100},
			args:   args{ttl: 10},
			want:   true,
		},
		{
			name:   "Значение TTL = 10, элемент новый",
			fields: fields{Val: "test5", Timestamp: time.Now().Unix() + 100},
			args:   args{ttl: 10},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elem := Element{
				Val:       tt.fields.Val,
				Timestamp: tt.fields.Timestamp,
				Updated:   tt.fields.Updated,
			}
			if got := elem.IsTTLOver(tt.args.ttl); got != tt.want {
				t.Errorf("Element.IsTTLOver() = %v, want %v", got, tt.want)
			}
		})
	}
}
