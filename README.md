[![Build Status](https://travis-ci.org/proway2/kvserver.svg?branch=master)](https://travis-ci.org/proway2/kvserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/proway2/kvserver)](https://goreportcard.com/report/github.com/proway2/kvserver)

# kvserver
Simple yet fully functional key-value storage server based on HTTP protocol. All objects stored in RAM only. All operations run for constant time.    
Server provides three types of operations:

- storing/updating value by its key
- getting value by its key
- deleting value by its key

Key-value element is cleaned up when expired.

# Features

- all operations have time complexity of ```O(1)```, i.e. always run for constant time.    
- hits TTL as much accurate as it's possible.    
- lower CPU cycles consumption during approximation and idle.    
- TTL approximation's divider is always 2.

# Installation
Clone and run ```go install``` in project folder.

# Usage
Command line arguments:
```bash
$ kvserver -h
Usage of kvserver:
  -addr string
    	IP address to bind to (default "127.0.0.1")
  -port int
    	port to listen to (default 8080)
  -ttl uint
    	element's (key-value) lifetime in the storage, secs. (default 60)
```
# API
Base URL ```http://<host>:<port>/key/<key_name>```, where ```<key_name>``` - is the name of the key to be stored. Key and its value are always string.
## Storing/Updating value by its key
_HTTP method_: ```POST```    
_Имя параметра запроса для передачи данных_: ```value```    
_Success code_: ```200```    
_Error code_: кода нет, такая ситуация воспринимается как крах сервера.    
_Примечание_: при повторной установке значения существующего ключа, время существования записи продлевается.

## Getting value by its key
_HTTP method_: ```GET```    
_Имя параметра запроса для передачи данных_: для получения из хранилища параметры не передаются.    
_Success code_: ```200```, в теле ответа передается текстовое значение для заданного ключа.    
_Error code_: ```404```

## Deleting value by its key
_HTTP method_: ```POST```    
_Имя параметра запроса для передачи данных_: для удаления из хранилища параметры не передаются.    
_Success code_: ```200```    
_Error code_: ```404```

When error is occured code ```400``` is returned by server.

# Tests
Run ```run_test.sh``` in project folder.

# License
GPL v3
