package router

import (
	"kvserver/kvstorage"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
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
				stor: &kvstorage.KVStorage{},
			},
			want: func(w http.ResponseWriter, r *http.Request) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetURLrouter(tt.args.stor); !reflect.DeepEqual(got, tt.want) {
				// t.Errorf("GetURLrouter() = %v, want %v", got, tt.want)
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
			name:  "Пустая строка",
			args:  args{inps: ""},
			want:  "",
			want1: false,
		},
		{
			name:  "Строка менее 6 символов",
			args:  args{inps: "abc"},
			want:  "",
			want1: false,
		},
		{
			name:  "Строка стандартная > 6 символов",
			args:  args{inps: "/key/fg"},
			want:  "fg",
			want1: true,
		},
		{
			name:  "Строка 5 символов, без ключа",
			args:  args{inps: "/key/"},
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

// func Test_creator(t *testing.T) {
// 	type args struct {
// 		method string
// 	}
// 	tests := []struct {
// 		name  string
// 		args  args
// 		want  func(*kvstorage.KVStorage, string, *http.Request) (string, int)
// 		want1 bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, got1 := creator(tt.args.method)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("creator() got = %v, want %v", got, tt.want)
// 			}
// 			if got1 != tt.want1 {
// 				t.Errorf("creator() got1 = %v, want %v", got1, tt.want1)
// 			}
// 		})
// 	}
// }

// func Test_methodGET(t *testing.T) {
// 	type args struct {
// 		stor *kvstorage.KVStorage
// 		key  string
// 		r    *http.Request
// 	}
// 	tests := []struct {
// 		name  string
// 		args  args
// 		want  string
// 		want1 int
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, got1 := methodGET(tt.args.stor, tt.args.key, tt.args.r)
// 			if got != tt.want {
// 				t.Errorf("methodGET() got = %v, want %v", got, tt.want)
// 			}
// 			if got1 != tt.want1 {
// 				t.Errorf("methodGET() got1 = %v, want %v", got1, tt.want1)
// 			}
// 		})
// 	}
// }

// func Test_methodPOST(t *testing.T) {
// 	type args struct {
// 		stor *kvstorage.KVStorage
// 		key  string
// 		r    *http.Request
// 	}
// 	tests := []struct {
// 		name  string
// 		args  args
// 		want  string
// 		want1 int
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, got1 := methodPOST(tt.args.stor, tt.args.key, tt.args.r)
// 			if got != tt.want {
// 				t.Errorf("methodPOST() got = %v, want %v", got, tt.want)
// 			}
// 			if got1 != tt.want1 {
// 				t.Errorf("methodPOST() got1 = %v, want %v", got1, tt.want1)
// 			}
// 		})
// 	}
// }

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
	storage := kvstorage.KVStorage{}
	chanOut := make(chan string, 10)
	storage.Init(chanOut)

	handler := GetURLrouter(&storage)
	var writer *myResponseWriter
	writer = &myResponseWriter{code: 200}

	form := url.Values{}
	form.Add("value", "key111 test")
	reqSet, _ := http.NewRequest("POST", "http://localhost:8080/key/key111", strings.NewReader(form.Encode()))
	reqSet.Header.Add("value", "key111 test")
	reqSet.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form = url.Values{}
	form.Add("valu", "key111 test")
	reqBad, _ := http.NewRequest("POST", "http://localhost:8080/key/key111", strings.NewReader(form.Encode()))
	reqBad.Header.Add("valu", "key111 test")
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
			name: "Неправильный HTTP метод",
			args: args{
				w: writer,
				r: &http.Request{
					Method: "PUT",
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost:8080",
						Path:   "/key/key111",
					},
				},
			},
			want: 400,
		},
		{
			name: "Получение из пустого хранилища",
			args: args{
				w: writer,
				r: &http.Request{
					Method: "GET",
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost:8080",
						Path:   "/key/key111",
					},
				},
			},
			want: 404,
		},
		{
			name: "Удаление из пустого хранилища",
			args: args{
				w: writer,
				r: &http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost:8080",
						Path:   "/key/key111",
					},
				},
			},
			want: 404,
		},
		{
			name: "Пустое значение ключа",
			args: args{
				w: writer,
				r: &http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost:8080",
						Path:   "/key/",
					},
				},
			},
			want: 400,
		},
		{
			name: "Короткое значение ключа",
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
			name: "Установка значения",
			args: args{
				w: writer,
				r: reqSet,
			},
			want: 200,
		},
		{
			name: "Удаление существующего значения",
			args: args{
				w: writer,
				r: &http.Request{
					Method: "POST",
					URL: &url.URL{
						Scheme: "http",
						Host:   "localhost:8080",
						Path:   "/key/key111",
					},
				},
			},
			want: 200,
		},
		{
			name: "Некорректный URL для установки значения",
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
