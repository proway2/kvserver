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
func getHandler(stor *kvstorage.KVStorage) func(
	w http.ResponseWriter, r *http.Request,
) {
	return func(w http.ResponseWriter, r *http.Request) {
		keyName, ok := getKeyFromURL(r.URL.Path)
		if !ok {
			w.WriteHeader(400) // Bad request
			return
		}
		switch r.Method {
		case "GET":
			// Получаем значение по ключу
			val, res := stor.Get(keyName)
			if !res {
				// значение ключа не найдено - возвращаем ошибку
				w.WriteHeader(404)
			}
			fmt.Fprint(w, val)
		case "POST":
			// требуется для извлечения значений метода POST - заполняется r.Form
			// fmt.Println("Значение value ", r.PostFormValue("value"))
			r.ParseForm()
			if len(r.Form) == 0 {
				// удаление значения по ключу
				if !stor.Delete(keyName) {
					// значение не удалено
					w.WriteHeader(404)
				}
			} else {
				if value, ok := r.Form["value"]; ok {
					// требуется установка значения
					stor.Set(keyName, value[0])
				} else {
					w.WriteHeader(400) // Bad request
				}
			}
		}
	}
}

func getKeyFromURL(inps string) (string, bool) {
	if len(inps) < 6 {
		return "", false
	}
	return inps[5:], true
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
	urlHandler := getHandler(storage)

	// для работы веб-сервера требуется определить обработчик URL
	http.HandleFunc("/key/", urlHandler)
	log.Fatal(server.ListenAndServe())
}
