package server

import (
	"net/http"
	"reflect"
	"testing"
	"yaGoShortURL/internal/app/cash"
)

func TestNewServer(t *testing.T) {
	type args struct {
		memory cash.Cash
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(tt.args.memory); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_StartServer(t *testing.T) {
	type fields struct {
		memory cash.Cash
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				memory: tt.fields.memory,
			}
			s.StartServer()
		})
	}
}

func TestServer_shorterServer(t *testing.T) {
	type fields struct {
		memory cash.Cash
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				memory: tt.fields.memory,
			}
			s.shorterServer(tt.args.w, tt.args.r)
		})
	}
}
