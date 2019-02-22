package main

import (
	"flag"
	"kvserver/kvstorage"
	"kvserver/router"
	"kvserver/vacuum"
	"log"
	"net/http"
	"strconv"
)

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
	if initRes := storage.Init(outElmChan); !initRes {
		log.Fatal("Cannot initialize storage!")
	}

	// инициализация очистки
	cleaner := vacuum.Lifo{}
	if initRes := cleaner.Init(storage, outElmChan, ttl); !initRes {
		log.Fatal("Cannot initialize cleaner!")
	}
	// для очистки хранилища от старых элементов используем отдельный поток
	go cleaner.Run()

	server := &http.Server{
		Addr: addr + ":" + strconv.Itoa(port),
	}
	urlHandler := router.GetURLrouter(storage)

	// для работы веб-сервера требуется определить обработчик URL
	http.HandleFunc("/key/", urlHandler)
	log.Fatal(server.ListenAndServe())
}
