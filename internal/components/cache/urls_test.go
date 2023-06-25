package cache

import (
	"testing"
	"yaGoShortURL/internal/entity"
)

func TestUrls_ReadURLFromCash(t *testing.T) {
	type fields struct {
		URLs    URLsAllInfo
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
				URLs: struct {
					IDKey   map[string]entity.JSONAllInfo
					URLKey  map[string]string
					Counter int
				}{
					IDKey: map[string]entity.JSONAllInfo{
						"0": {FullURL: "https://golang-blog.blogspot.com/2020/01/map-golang.html"},
					},
				},
			},
			args: args{
				id: "0",
			},
			want:    "https://golang-blog.blogspot.com/2020/01/map-golang.html",
			wantErr: false,
		},
		{
			name: "positiveRead1",
			fields: fields{
				URLs: struct {
					IDKey   map[string]entity.JSONAllInfo
					URLKey  map[string]string
					Counter int
				}{
					IDKey: map[string]entity.JSONAllInfo{
						"0":  {FullURL: "https://golang-blog.blogspot.com/2020/01/map-golang.html"},
						"30": {FullURL: "https://go-traps.appspot.com/#slice-traversal"},
					},
				},
			},
			args: args{

				id: "30",
			},
			want:    "https://go-traps.appspot.com/#slice-traversal",
			wantErr: false,
		},
		{
			name: "negativeRead2",
			fields: fields{
				URLs: struct {
					IDKey   map[string]entity.JSONAllInfo
					URLKey  map[string]string
					Counter int
				}{
					IDKey: map[string]entity.JSONAllInfo{
						"0":  {FullURL: "https://golang-blog.blogspot.com/2020/01/map-golang.html"},
						"30": {FullURL: "https://go-traps.appspot.com/#slice-traversal"},
					},
				},
			},
			args: args{
				id: "11",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &Urls{
				URLs: tt.fields.URLs,
			}
			got, err := u.FullURL(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FullURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FullURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrls_WriteURLInCash(t *testing.T) {
	type fields struct {
		URLs      URLsAllInfo
		usersUrls map[uint64][]string
		baseURL   string
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
				URLs: struct {
					IDKey   map[string]entity.JSONAllInfo
					URLKey  map[string]string
					Counter int
				}{
					IDKey:   make(map[string]entity.JSONAllInfo, defaultURLsNumber),
					URLKey:  make(map[string]string, defaultURLsNumber),
					Counter: 0,
				},
				usersUrls: map[uint64][]string{},
				baseURL:   "http://localhost:8080",
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
				URLs: struct {
					IDKey   map[string]entity.JSONAllInfo
					URLKey  map[string]string
					Counter int
				}{
					IDKey: map[string]entity.JSONAllInfo{
						"0": {FullURL: "https://golang-blog.blogspot.com/2020/01/map-golang.html"},
					},
					URLKey: make(map[string]string, defaultURLsNumber),
				},
				usersUrls: map[uint64][]string{},
				baseURL:   "http://localhost:8080",
			},

			args: args{
				fullURL: "https://blog.mozilla.org/en/",
				userID:  15,
			},
			want:    "0",
			wantErr: false,
		},
		{
			name: "negativeWrite2",
			fields: fields{
				URLs: struct {
					IDKey   map[string]entity.JSONAllInfo
					URLKey  map[string]string
					Counter int
				}{
					IDKey: map[string]entity.JSONAllInfo{
						"0": {FullURL: "https://golang-blog.blogspot.com/2020/01/map-golang.html"},
					},
					URLKey: map[string]string{
						"https://golang-blog.blogspot.com/2020/01/map-golang.html": "0",
					},
				},
				usersUrls: map[uint64][]string{},
				baseURL:   "http://localhost:8080",
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
				URLs: struct {
					IDKey   map[string]entity.JSONAllInfo
					URLKey  map[string]string
					Counter int
				}{
					IDKey: map[string]entity.JSONAllInfo{
						"0": {FullURL: "https://golang-blog.blogspot.com/2020/01/map-golang.html"},
					},
					URLKey: make(map[string]string, defaultURLsNumber),
				},
				usersUrls: map[uint64][]string{},
				baseURL:   "http://localhost:8080",
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
				URLs: struct {
					IDKey   map[string]entity.JSONAllInfo
					URLKey  map[string]string
					Counter int
				}{
					IDKey: map[string]entity.JSONAllInfo{
						"0": {FullURL: "https://golang-blog.blogspot.com/2020/01/map-golang.html"},
					},
					URLKey: make(map[string]string, defaultURLsNumber),
				},
				usersUrls: map[uint64][]string{},
				baseURL:   "http://localhost:8080",
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
				URLs:      tt.fields.URLs,
				usersUrls: tt.fields.usersUrls,
				baseURL:   tt.fields.baseURL,
			}
			userIDB := tt.args.userID
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
