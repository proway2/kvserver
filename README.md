# kvserver
Простой сервер хранения "ключ-значение" на основе протокола HTTP. Данные хранятся только в оперативной памяти.
Реализованы три пользовательские операции:

- установка/обновление значения для ключа
- получение значения для ключа
- удаление значения для ключа

Периодически сервер проводит очистку хранилища от старых элементов, тех, у которых превышен TTL.

# Установка
В папке проекта запустить ```go install```

# Запуск сервера
Параметры командной строки:
```bash
$ kvserver -h
Usage of kvserver:
  -addr string
      IP адрес для подключения сервера (default "127.0.0.1")
  -port int
      номер порта для подключения (default 8080)
  -ttl int
      время жизни элемента (ключ-значение) в хранилище, сек. (default 60)
```
# API
Базовый URL ```http://<host>:<port>/key/<key_name>```, где ```<key_name>``` - имя ключа для работы с хранилищем. Ключ и его значение могут быть только текстовыми.
## Установка значения для ключа
_HTTP метод_: ```POST```    
_Имя параметра запроса для передачи данных_: ```value```    
_Код состояния в случае успеха_: ```200```    
_Код состояния в случае ошибки_: кода нет, такая ситуация воспринимается как крах сервера.    
_Примечание_: при повторной установке значения существующего ключа, время существования записи продлевается.

## Получение значения для ключа
_HTTP метод_: ```GET```    
_Имя параметра запроса для передачи данных_: для получения из хранилища параметры не передаются.    
_Код состояния в случае успеха_: ```200```, в теле ответа передается текстовое значение для заданного ключа.    
_Код состояния в случае ошибки_: ```404```

## Удаление значения для ключа
_HTTP метод_: ```POST```    
_Имя параметра запроса для передачи данных_: для удаления из хранилища параметры не передаются.    
_Код состояния в случае успеха_: ```200```    
_Код состояния в случае ошибки_: ```404```

При ошибочном запросе, возвращается код ошибки ```400```.

# Тесты
В папке проекта запустить ```run_test.sh```. Работает при наличии ```bash```.

# Лицензия
GPL v3
