package utils

import (
	"testing"
)

func TestB2S(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "default",
			args: args{
				[]byte{72, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 33},
			},
			want: "Hello world!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := B2S(tt.args.b); got != tt.want {
				t.Errorf("B2S() = %v, want %v", got, tt.want)
			}
		})
	}
}
