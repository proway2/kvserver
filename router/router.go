package router

import (
	"fmt"
	"net/http"

	"github.com/proway2/kvserver/kvstorage"
)

// GetURLrouter - возвращает функцию маршрутизатор HTTP запросов в зависимости от типа.
func GetURLrouter(stor *kvstorage.KVStorage) func(
	w http.ResponseWriter, r *http.Request,
) {
	// замыкание необходимо для оборачивания локальных переменных в обработчик URL
	return func(w http.ResponseWriter, r *http.Request) {
		keyName, ok := getKeyFromURL(r.URL.Path)
		if !ok {
			w.WriteHeader(400) // Bad request
			fmt.Fprint(w, "400 Malformed request.\n")
			return
		}
		reqHandler, isHandlerExists := requestFactory(r.Method)
		if !isHandlerExists {
			w.WriteHeader(400) // Bad request
			fmt.Fprint(w, "400 Malformed request.\n")
			return
		}
		val, code := reqHandler(stor, keyName, r)
		w.WriteHeader(code)
		fmt.Fprint(w, val)
	}
}

func getKeyFromURL(inps string) (string, bool) {
	if len(inps) < 6 {
		return "", false
	}
	return inps[5:], true
}

// requestFactory returns function which can be use to handle different types of HTTP request (GET or POST)
func requestFactory(method string) (func(*kvstorage.KVStorage, string, *http.Request) (string, int), bool) {
	if method == "GET" {
		return methodGET, true
	}
	if method == "POST" {
		return methodPOST, true
	}
	return nil, false
}

// methodGET returns value and the HTTP code for the key.
func methodGET(stor *kvstorage.KVStorage, key string, r *http.Request) (string, int) {
	code := 200
	// get the value by its key
	val, err := stor.Get(key)
	if err != nil {
		return "500 Internal storage error.\n", 500
	}
	if val == nil {
		// either error occured or key is not found in the storage (code 404)
		val = []byte(fmt.Sprintf("404 There is no record in the storage for key '%v'.\n", key))
		code = 404
	}
	return string(val), code
}

// methodPOST - функция обработчика метода POST
func methodPOST(stor *kvstorage.KVStorage, key string, r *http.Request) (string, int) {
	// требуется для извлечения значений метода POST - заполняется r.Form
	val := ""
	code := 200
	r.ParseForm()
	if len(r.Form) == 0 {
		// удаление значения по ключу
		if !stor.Delete(key) {
			// значение не удалено
			code = 404
			val = fmt.Sprintf("404 There is no record in the storage for key '%v'.\n", key)
		}
	} else {
		if value, ok := r.Form["value"]; ok {
			// требуется установка значения
			stor.Set(key, value[0])
		} else {
			// Bad request
			code = 400
			val = "400 Malformed request.\n"
		}
	}
	return val, code
}
