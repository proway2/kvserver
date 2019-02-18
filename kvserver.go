package main

import (
	"flag"
	"fmt"
	"kvserver/kvstorage"
	"kvserver/vacuum"
	"log"
	"net/http"
	"strconv"
)

// замыкание необходимо для оборачивания локальных переменных в обработчик URL
func getURLHandler(stor *kvstorage.KVStorage) func(
	w http.ResponseWriter, r *http.Request,
) {
	return func(w http.ResponseWriter, r *http.Request) {
		keyName, ok := getKeyFromURL(r.URL.Path)
		if !ok {
			w.WriteHeader(400) // Bad request
			fmt.Fprint(w, "400 Malformed request.\n")
			return
		}
		reqHandler, isHandlerExists := creator(r.Method)
		if !isHandlerExists {
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
		return GEThandler, true
	}
	if method == "POST" {
		return POSThandler, true
	}
	return nil, false
}

// GEThandler - функция обработчика метода GET
func GEThandler(stor *kvstorage.KVStorage, key string, r *http.Request) (string, int) {
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

// POSThandler - функция обработчика метода POST
func POSThandler(stor *kvstorage.KVStorage, key string, r *http.Request) (string, int) {
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

func getCLIargs() (string, int, uint64) {
	ttlP := flag.Uint64(
		"ttl",
		60,
		"время жизни элемента (ключ-значение) в хранилище, сек.",
	)
	addr := flag.String(
		"addr",
		"127.0.0.1",
		"IP адрес для подключения сервера",
	)
	port := flag.Int(
		"port",
		8080,
		"номер порта для подключения",
	)
	flag.Parse()
	return *addr, *port, *ttlP
}

func main() {
	// для дальнейшей работы надо или получить аргументы
	// из командной строки или установить значения по умолчанию
	addr, port, ttl := getCLIargs()

	// канал необходим для постановки в очередь элементов для очистителя
	outElmChan := make(chan string)

	// инициализация хранилища
	storage := &kvstorage.KVStorage{}
	// storage.Init(nil)
	if initRes := storage.Init(&outElmChan); !initRes {
		log.Fatal("Cannot initialize storage!")
	}

	// инициализация очистки
	cleaner := vacuum.Lifo{}
	if initRes := cleaner.Init(storage, &outElmChan, ttl); !initRes {
		log.Fatal("Cannot initialize cleaner!")
	}
	// для очистки хранилища от старых элементов используем отдельный поток
	go cleaner.Run()

	server := &http.Server{
		Addr: addr + ":" + strconv.Itoa(port),
	}
	urlHandler := getURLHandler(storage)

	// для работы веб-сервера требуется определить обработчик URL
	http.HandleFunc("/key/", urlHandler)
	log.Fatal(server.ListenAndServe())
}
