package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/proway2/kvserver/kvstorage"
	"github.com/proway2/kvserver/router"
	"github.com/proway2/kvserver/vacuum"
)

func getCLIargs() (string, int, uint64) {
	ttlP := flag.Uint64(
		"ttl",
		60,
		"element's (key-value) lifetime in the storage, secs.",
	)
	addr := flag.String(
		"addr",
		"127.0.0.1",
		"IP address to bind to",
	)
	port := flag.Int(
		"port",
		8080,
		"port to listen to",
	)
	flag.Parse()
	return *addr, *port, *ttlP
}

func main() {
	// для дальнейшей работы надо или получить аргументы
	// из командной строки или установить значения по умолчанию
	addr, port, ttl := getCLIargs()

	// инициализация хранилища
	storage := kvstorage.NewStorage()
	if storage == nil {
		log.Fatal("Cannot initialize storage!")
	}

	// инициализация очистки
	cleaner := vacuum.Vacuum{}
	if initRes := cleaner.Init(storage, ttl); !initRes {
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
