package router

import (
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/proway2/kvserver/kvstorage"
)

const (
	correctValueName   = "value"
	incorrectValueName = "valu"
	correctKey         = "key1"
	emptyKey           = ""
	correctValue       = "key1 test value"
)

func TestGetURLrouter(t *testing.T) {
	type args struct {
		stor *kvstorage.KVStorage
	}
	tests := []struct {
		name string
		args args
		want func(w http.ResponseWriter, r *http.Request)
	}{
		{
			name: "Создание обработчика URL",
			args: args{
				stor: kvstorage.NewStorage(),
			},
			want: func(w http.ResponseWriter, r *http.Request) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedType := reflect.TypeOf(func(w http.ResponseWriter, r *http.Request) {})
			if got := GetURLrouter(tt.args.stor); reflect.TypeOf(got) != expectedType {
				t.Errorf("GetURLrouter() = %T, want %T", got, tt.want)
			}
		})
	}
}

func Test_getKeyFromURL(t *testing.T) {
	type args struct {
		inps string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name:  "Empty URL path",
			args:  args{inps: ""},
			want:  "",
			want1: false,
		},
		{
			name:  "Incorrect path",
			args:  args{inps: "abc"},
			want:  "",
			want1: false,
		},
		{
			name:  "Good path",
			args:  args{inps: "/key/fg"},
			want:  "fg",
			want1: true,
		},
		{
			name:  "No key path",
			args:  args{inps: "/key/"},
			want:  "",
			want1: false,
		},
		{
			name:  "Incorrect path, no key",
			args:  args{inps: "/key2/"},
			want:  "",
			want1: false,
		},
		{
			name:  "Incorrect path with a key",
			args:  args{inps: "/key2/abc"},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getKeyFromURL(tt.args.inps)
			if got != tt.want {
				t.Errorf("getKeyFromURL() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getKeyFromURL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

type myResponseWriter struct {
	code int
}

func (resp *myResponseWriter) Header() http.Header {
	res := make(map[string][]string)
	return res
}

func (resp *myResponseWriter) WriteHeader(code int) {
	resp.code = code
}

func (resp *myResponseWriter) Write(inp []byte) (int, error) {
	var er error
	return int(1), er
}

func (resp *myResponseWriter) getCode() int {
	return resp.code
}

func Test_closure(t *testing.T) {
	storage := kvstorage.NewStorage()

	handler := GetURLrouter(storage)
	var writer = &myResponseWriter{code: 200}

	form := url.Values{}
	form.Add(correctValueName, correctValue)
	reqSet, _ := http.NewRequest("POST", "http://localhost:8080/key/"+correctKey, strings.NewReader(form.Encode()))
	reqSet.Header.Add(correctValueName, correctValue)
	reqSet.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form = url.Values{}
	form.Add(incorrectValueName, correctValue)
	reqBad, _ := http.NewRequest("POST", "http://localhost:8080/key/"+correctKey, strings.NewReader(form.Encode()))
	reqBad.Header.Add(incorrectValueName, correctValue)
	reqBad.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name string
		args args
		want int // код ошибки
		text string
	}{
		{
			name: "Incorrect HTTP method (verb)",
			args: args{
				w: writer,
				r: &http.Request{
					Method: "PUT",
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost:8080",
						Path:   "/key/" + correctKey,
					},
				},
			},
			want: 400,
		},
		{
			name: "Getting value from the empty storage",
			args: args{
				w: writer,
				r: &http.Request{
					Method: "GET",
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost:8080",
						Path:   "/key/" + correctKey,
					},
				},
			},
			want: 404,
		},
		{
			name: "Deleting from the empty storage",
			args: args{
				w: writer,
				r: &http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost:8080",
						Path:   "/key/" + correctKey,
					},
				},
			},
			want: 404,
		},
		{
			name: "Setting value with empty key",
			args: args{
				w: writer,
				r: &http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost:8080",
						Path:   "/key/" + emptyKey,
					},
				},
			},
			want: 400,
		},
		{
			name: "URL keyword is incorrect",
			args: args{
				w: writer,
				r: &http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost:8080",
						Path:   "/ke/",
					},
				},
			},
			want: 400,
		},
		{
			name: "Correct setting the value by its key",
			args: args{
				w: writer,
				r: reqSet,
			},
			want: 200,
		},
		{
			name: "Deleting existing value by its key",
			args: args{
				w: writer,
				r: &http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost:8080",
						Path:   "/key/" + correctKey,
					},
				},
			},
			want: 200,
		},
		{
			name: "Incorrect URL",
			args: args{
				w: writer,
				r: reqBad,
			},
			want: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// сброс кода ответа
			tt.args.w.WriteHeader(200)

			handler(tt.args.w, tt.args.r)
			got := writer.getCode()
			if got != tt.want {
				t.Errorf("urlHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}
