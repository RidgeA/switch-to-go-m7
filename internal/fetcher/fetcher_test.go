package fetcher

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func Test_load(t *testing.T) {
	type args struct {
		ctx  context.Context
		path string
	}
	tests := []struct {
		name    string
		server  *httptest.Server
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Load page",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = fmt.Fprintf(w, "<html></html>")
			})),
			args: args{
				ctx:  context.Background(),
				path: "/path",
			},
			want: []byte("<html></html>"),
		},

		{
			name: "Timeout",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(20 * time.Millisecond)
				_, _ = fmt.Fprintf(w, "<html></html>")
			})),
			args: args{
				ctx: func() context.Context {
					ctx, _ := context.WithTimeout(context.Background(), 10*time.Microsecond)
					return ctx
				}(),
				path: "/path",
			},
			wantErr: true,
		},
		{
			name: "HTTP error",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
			})),
			args: args{
				ctx:  context.Background(),
				path: "/path",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			url := tt.server.URL + tt.args.path

			got, err := load(tt.args.ctx, url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parse(t *testing.T) {
	type args struct {
		baseUrl string
		reader  string
	}
	tests := []struct {
		name      string
		args      args
		wantTitle string
		wantLinks []string
	}{
		{
			name: "A single anchor tag",
			args: args{
				baseUrl: "https://example.com",
				reader:  `<a href="https://example.com">Link<\a>`,
			},
			wantTitle: "",
			wantLinks: []string{"https://example.com"},
		},
		{
			name: "Anchor tag without 'href' attribute",
			args: args{
				baseUrl: "https://example.com",
				reader:  `<a>Link<\a>`,
			},
			wantTitle: "",
			wantLinks: []string{},
		},
		{
			name: "Self-closing anchor tag",
			args: args{
				baseUrl: "https://example.com",
				reader:  `<a href="https://example.com" \>`,
			},
			wantTitle: "",
			wantLinks: []string{"https://example.com"},
		},
		{
			name: "Several anchor tags",
			args: args{
				baseUrl: "https://example.com",
				reader: `<a href="https://example.com/1">Link<\a>
						 <a href="https://example.com/2">Link<\a>
						 <a href="https://example.com/3">Link<\a>
						 <a href="https://example.com/4">Link<\a>`,
			},
			wantTitle: "",
			wantLinks: []string{
				"https://example.com/1",
				"https://example.com/2",
				"https://example.com/3",
				"https://example.com/4",
			},
		},
		{
			name: "Several anchor tags with short links",
			args: args{
				baseUrl: "https://example.com",
				reader: `<a href="/1">Link<\a>
						 <a href="/2">Link<\a>
						 <a href="/3">Link<\a>
						 <a href="/4">Link<\a>`,
			},
			wantTitle: "",
			wantLinks: []string{
				"https://example.com/1",
				"https://example.com/2",
				"https://example.com/3",
				"https://example.com/4",
			},
		},
		{
			name:
			"Should select an article header",
			args: args{
				baseUrl: "https://example.com",
				reader:  `<h1 id="firstHeading">Heading</h1>`,
			},
			wantTitle: "Heading",
			wantLinks: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTitle, gotLinks := parse(tt.args.baseUrl, []byte(tt.args.reader))
			if gotTitle != tt.wantTitle {
				t.Errorf("Got title = %v, want %v", gotTitle, tt.wantTitle)
			}
			if !reflect.DeepEqual(gotLinks, tt.wantLinks) {
				t.Errorf("Got links = %v, want %v", gotLinks, tt.wantLinks)
			}
		})
	}
}
