package router

import (
	"fmt"
	"kvserver/kvstorage"
	"net/http"
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
		reqHandler, isHandlerExists := creator(r.Method)
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

// Фабричная функция
func creator(method string) (func(*kvstorage.KVStorage, string, *http.Request) (string, int), bool) {
	if method == "GET" {
		return methodGET, true
	}
	if method == "POST" {
		return methodPOST, true
	}
	return nil, false
}

// methodGET - функция обработчика метода GET
func methodGET(stor *kvstorage.KVStorage, key string, r *http.Request) (string, int) {
	code := 200
	// Получаем значение по ключу
	val, res := stor.Get(key)
	if !res {
		// значение ключа не найдено - ошибка 404
		val = "404 There is no record in storage for this key.\n"
		code = 404
	}
	return val, code
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
			val = "404 There is no record in storage for this key.\n"
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
