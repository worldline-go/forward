package server_test

import (
	"reflect"
	"testing"

	"github.com/worldline-go/forward/internal/server"
)

func TestParse(t *testing.T) {
	type args struct {
		hosts   []string
		sockets []string
	}

	tests := []struct {
		name string
		args args
		want []server.Holder
	}{
		{
			name: "default",
			args: args{
				hosts:   []string{"0.0.0.0:8080"},
				sockets: []string{"/docker/:*,-POST,-PUT,-DELETE"},
			},
			want: []server.Holder{
				{
					Name:   "default",
					Host:   "0.0.0.0:8080",
					Socket: []string{"/docker/:*,-POST,-PUT,-DELETE"},
				},
			},
		},
		{
			name: "multiple",
			args: args{
				hosts:   []string{"default@0.0.0.0:8080", "test@0.0.0.0:8081"},
				sockets: []string{"/docker/:*,-POST,-PUT,-DELETE", "test@/docker/:*,-POST,-PUT,-DELETE"},
			},
			want: []server.Holder{
				{
					Name:   "default",
					Host:   "0.0.0.0:8080",
					Socket: []string{"/docker/:*,-POST,-PUT,-DELETE"},
				},
				{
					Name:   "test",
					Host:   "0.0.0.0:8081",
					Socket: []string{"/docker/:*,-POST,-PUT,-DELETE"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := server.Parse(tt.args.hosts, tt.args.sockets); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
