# kvserver
Simple yet fully functional key-value server based on HTTP protocol. All objects stored in RAM only. All operations run for constant time.    
Server provides three types of operations:

- storing/updating value by its key
- getting value by its key
- deleting value by its key

Server cleans up the storage periodically from elements with expired TTL.

# Features

- all operations have time complexity of ```O(1)```, i.e. always run for constant time.    
- hits TTL as much accurate as it's possible.    
- lower CPU cycles consumption during approximation and idle.    
- делитель оставшегося времени всегда 2.

# Installation
В папке проекта запустить ```go install```

# Usage
Параметры командной строки:
```bash
$ kvserver -h
Usage of kvserver:
  -addr string
      IP адрес для подключения сервера (default "127.0.0.1")
  -port int
      номер порта для подключения (default 8080)
  -ttl uint
      время жизни элемента (ключ-значение) в хранилище, сек. (default 60)
```
# API
Базовый URL ```http://<host>:<port>/key/<key_name>```, где ```<key_name>``` - имя ключа для работы с хранилищем. Ключ и его значение могут быть только текстовыми.
## Storing/Updating value by its key
_HTTP метод_: ```POST```    
_Имя параметра запроса для передачи данных_: ```value```    
_Код состояния в случае успеха_: ```200```    
_Код состояния в случае ошибки_: кода нет, такая ситуация воспринимается как крах сервера.    
_Примечание_: при повторной установке значения существующего ключа, время существования записи продлевается.

## Getting value by its key
_HTTP метод_: ```GET```    
_Имя параметра запроса для передачи данных_: для получения из хранилища параметры не передаются.    
_Код состояния в случае успеха_: ```200```, в теле ответа передается текстовое значение для заданного ключа.    
_Код состояния в случае ошибки_: ```404```

## Deleting value by its key
_HTTP метод_: ```POST```    
_Имя параметра запроса для передачи данных_: для удаления из хранилища параметры не передаются.    
_Код состояния в случае успеха_: ```200```    
_Код состояния в случае ошибки_: ```404```

При ошибочном запросе, возвращается код ошибки ```400```.

# Tests
В папке проекта запустить ```run_test.sh```. Работает при наличии ```bash```.

# License
GPL v3
