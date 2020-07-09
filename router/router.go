package router

import (
	"fmt"
	"net/http"
)

type writer interface {
	Set(key, value string) error
	Delete(key string) (bool, error)
}

type reader interface {
	Get(key string) ([]byte, error)
}

type readerWriter interface {
	reader
	writer
}

// POST form field name (contains data for storing the key)
const valueFormFieldName = "value"

var httpStatusCodeMessages = map[int]string{
	200: "",
	400: "400 Malformed request.\n",
	404: "404 There is no record in the storage for key '%v'.\n",
	500: "500 Internal storage error.\n",
}

// GetURLrouter - возвращает функцию маршрутизатор HTTP запросов в зависимости от типа.
func GetURLrouter(stor readerWriter) func(
	w http.ResponseWriter, r *http.Request,
) {
	// замыкание необходимо для оборачивания локальных переменных в обработчик URL
	return func(w http.ResponseWriter, r *http.Request) {
		keyName, ok := getKeyFromURL(r.URL.Path)
		if !ok {
			w.WriteHeader(400) // Bad request
			fmt.Fprint(w, httpStatusCodeMessages[400])
			return
		}
		reqHandler, isHandlerExists := requestFactory(r.Method)
		if !isHandlerExists {
			w.WriteHeader(400) // Bad request
			fmt.Fprint(w, httpStatusCodeMessages[400])
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
func requestFactory(method string) (func(readerWriter, string, *http.Request) (string, int), bool) {
	if method == http.MethodGet {
		return methodGET, true
	}
	if method == http.MethodPost {
		return methodPOST, true
	}
	return nil, false
}

// methodGET returns value and the HTTP code for the key.
func methodGET(stor readerWriter, key string, r *http.Request) (string, int) {
	code := 200
	// get the value by its key
	val, err := stor.Get(key)
	if err != nil {
		return httpStatusCodeMessages[500], 500
	}
	if val == nil {
		// either error occured or key is not found in the storage (code 404)
		code = 404
		val = []byte(fmt.Sprintf(httpStatusCodeMessages[code], key))
	}
	return string(val), code
}

// methodPOST - функция обработчика метода POST
func methodPOST(stor readerWriter, key string, r *http.Request) (string, int) {
	// требуется для извлечения значений метода POST - заполняется r.Form
	r.ParseForm()
	value := r.FormValue(valueFormFieldName)
	postProcessingMethod := postMethodFactory(len(r.Form))
	httpCode := postProcessingMethod(stor, key, value)
	return httpStatusCodeMessages[httpCode], httpCode
}

func postMethodFactory(formLen int) func(storage readerWriter, key, value string) int {
	if formLen == 0 {
		// deleting the element
		return deleteElementRequest
	}
	// setting the element
	return setElementRequest
}

// deleteElementRequest processes delete HTTP request and returns HTTP code.
func deleteElementRequest(storage readerWriter, key, value string) int {
	// deleting element by its key
	delStatus, err := storage.Delete(key)
	if err != nil {
		// something went wrong with the storage
		return 500
	}
	if delStatus {
		// element deleted successfully
		return 200
	}
	// element was not found and is not deleted
	return 404
}

func setElementRequest(storage readerWriter, key, value string) int {
	// setting (updating) the value by its key
	err := storage.Set(key, value)
	if err != nil {
		// something went wrong with the storage
		return 500
	}
	if value != "" {
		return 200
	}
	return 400
}
