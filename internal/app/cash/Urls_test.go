package cash

import (
	"reflect"
	"testing"
)

func TestNewCash(t *testing.T) {
	tests := []struct {
		name string
		want *Cash
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCash(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUrls(t *testing.T) {
	tests := []struct {
		name string
		want *Urls
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUrls(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUrls() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrls_ReadURLFromCash(t *testing.T) {
	type fields struct {
		urlsMap map[string]string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "typicalRead",
			fields: fields{
				urlsMap: map[string]string{
					"id:0": "https://golang-blog.blogspot.com/2020/01/map-golang.html",
				},
			},
			args: args{

				id: "id:0",
			},
			want:    "https://golang-blog.blogspot.com/2020/01/map-golang.html",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Urls{
				urlsMap: tt.fields.urlsMap,
			}
			got, err := u.ReadURLFromCash(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadURLFromCash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReadURLFromCash() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrls_WriteURLInCash(t *testing.T) {
	type fields struct {
		urlsMap map[string]string
	}
	type args struct {
		fullURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "typicalWrite0",
			fields: fields{
				urlsMap: map[string]string{},
			},
			args: args{
				fullURL: "https://golang-blog.blogspot.com/2020/01/map-golang.html",
			},
			want:    "0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Urls{
				urlsMap: tt.fields.urlsMap,
			}
			got, err := u.WriteURLInCash(tt.args.fullURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteURLInCash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WriteURLInCash() got = %v, want %v", got, tt.want)
			}
		})
	}
}
