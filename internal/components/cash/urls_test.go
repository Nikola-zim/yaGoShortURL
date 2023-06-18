package cash

import (
	"encoding/binary"
	"testing"
	"yaGoShortURL/internal/entity"
)

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
			name: "positiveRead0",
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
		{
			name: "positiveRead1",
			fields: fields{
				urlsMap: map[string]string{
					"id:0":  "https://golang-blog.blogspot.com/2020/01/map-golang.html",
					"id:30": "https://go-traps.appspot.com/#slice-traversal",
				},
			},
			args: args{

				id: "id:30",
			},
			want:    "https://go-traps.appspot.com/#slice-traversal",
			wantErr: false,
		},
		{
			name: "negativeRead2",
			fields: fields{
				urlsMap: map[string]string{
					"id:0":  "https://golang-blog.blogspot.com/2020/01/map-golang.html",
					"id:30": "https://go-traps.appspot.com/#slice-traversal",
					"id:11": "",
				},
			},
			args: args{

				id: "id:11",
			},
			want:    "",
			wantErr: true,
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
		urlsMap     map[string]string
		usersUrls   map[uint64][]string
		URLsAllInfo map[string]entity.JSONAllInfo
		baseURL     string
	}
	type args struct {
		fullURL string
		userID  uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "positiveWrite0",
			fields: fields{
				urlsMap:     map[string]string{},
				usersUrls:   map[uint64][]string{},
				URLsAllInfo: map[string]entity.JSONAllInfo{},
				baseURL:     "http://localhost:8080",
			},
			args: args{
				fullURL: "https://golang-blog.blogspot.com/2020/01/map-golang.html",
				userID:  15,
			},
			want:    "0",
			wantErr: false,
		},
		{
			name: "positiveWrite1",
			fields: fields{
				urlsMap: map[string]string{
					"id:0": "https://golang-blog.blogspot.com/2020/01/map-golang.html",
					"url:https://golang-blog.blogspot.com/2020/01/map-golang.html": "https://golang-blog.blogspot.com/2020/01/map-golang.html",
				},
				usersUrls:   map[uint64][]string{},
				URLsAllInfo: map[string]entity.JSONAllInfo{},
				baseURL:     "http://localhost:8080",
			},
			args: args{
				fullURL: "https://blog.mozilla.org/en/",
				userID:  15,
			},
			want:    "1",
			wantErr: false,
		},
		{
			name: "negativeWrite2",
			fields: fields{
				urlsMap: map[string]string{
					"id:0": "https://golang-blog.blogspot.com/2020/01/map-golang.html",
					"url:https://golang-blog.blogspot.com/2020/01/map-golang.html": "https://golang-blog.blogspot.com/2020/01/map-golang.html",
				},
				usersUrls:   map[uint64][]string{},
				URLsAllInfo: map[string]entity.JSONAllInfo{},
				baseURL:     "http://localhost:8080",
			},
			args: args{
				fullURL: "https://golang-blog.blogspot.com/2020/01/map-golang.html",
				userID:  15,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "negativeWrite3",
			fields: fields{
				urlsMap: map[string]string{
					"id:0": "https://golang-blog.blogspot.com/2020/01/map-golang.html",
					"url:https://golang-blog.blogspot.com/2020/01/map-golang.html": "https://golang-blog.blogspot.com/2020/01/map-golang.html",
				},
				usersUrls:   map[uint64][]string{},
				URLsAllInfo: map[string]entity.JSONAllInfo{},
				baseURL:     "http://localhost:8080",
			},
			args: args{
				fullURL: "",
				userID:  15,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "negativeWrite4",
			fields: fields{
				urlsMap: map[string]string{
					"id:0": "https://golang-blog.blogspot.com/2020/01/map-golang.html",
					"url:https://golang-blog.blogspot.com/2020/01/map-golang.html": "https://golang-blog.blogspot.com/2020/01/map-golang.html",
				},
				usersUrls:   map[uint64][]string{},
				URLsAllInfo: map[string]entity.JSONAllInfo{},
				baseURL:     "http://localhost:8080",
			},
			args: args{
				fullURL: "qwer",
				userID:  15,
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Urls{
				urlsMap:     tt.fields.urlsMap,
				usersUrls:   tt.fields.usersUrls,
				URLsAllInfo: tt.fields.URLsAllInfo,
				baseURL:     tt.fields.baseURL,
			}
			userIDB := make([]byte, 8)
			binary.LittleEndian.PutUint64(userIDB, tt.args.userID)
			got, err := u.WriteURL(tt.args.fullURL, userIDB)
			if (err != nil) != tt.wantErr {
				t.Errorf("WriteURLInFS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("WriteURLInFS() got = %v, want %v", got, tt.want)
			}
		})
	}
}
