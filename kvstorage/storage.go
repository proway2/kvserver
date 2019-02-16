package kvstorage

import (
	"kvserver/element"
	"sync"
	"time"
)

// KVStorage - Структура с методами, описывающая хранилище
type KVStorage struct {
	kvstorage  map[string]*element.Element
	mux        sync.Mutex
	outElmChan *chan string
}

// Init - Функция инициализации хранилища
func (kv *KVStorage) Init(out *chan string) bool {
	if out == nil {
		return false
	}
	kv.kvstorage = make(map[string]*element.Element)
	kv.outElmChan = out
	return true
}

// Set - установка ключ-значение
func (kv *KVStorage) Set(key, value string) bool {
	if len(key) == 0 {
		return false
	}
	kv.mux.Lock()
	// проверяем есть ли у нас такой ключ в карте
	_, found := kv.kvstorage[key]
	// ключ есть, надо обновить его значение и установить признак updated
	elem := &element.Element{
		Val: value, Timestamp: time.Now().Unix(), Updated: found,
	}
	kv.kvstorage[key] = elem
	kv.mux.Unlock()

	if kv.outElmChan != nil {
		*kv.outElmChan <- key
	}
	return true
}

// Get - получение значения по ключу,
// второй параметр указывает на успешность получения значения по ключу
func (kv *KVStorage) Get(key string) (string, bool) {
	if len(key) == 0 {
		return "", false
	}
	kv.mux.Lock()
	defer kv.mux.Unlock()
	elem, ok := kv.kvstorage[key]
	if ok {
		return elem.Val, true
	}
	return "", false
}

// Delete - удаление значения по ключу
func (kv *KVStorage) Delete(key string) bool {
	if len(key) == 0 {
		return false
	}
	kv.mux.Lock()
	defer kv.mux.Unlock()
	_, ok := kv.kvstorage[key]
	if ok {
		delete(kv.kvstorage, key)
	}
	return ok
}

// ResetUpdated - сброс признака обновления элемента в false
func (kv *KVStorage) ResetUpdated(key string) bool {
	if len(key) == 0 {
		return false
	}
	kv.mux.Lock()
	defer kv.mux.Unlock()
	if elem, ok := kv.kvstorage[key]; ok {
		elem.Updated = false
		return true
	}
	return false
}

// IsElemUpdated - получение признака был ли элемент обновлен
func (kv *KVStorage) IsElemUpdated(key string) bool {
	if len(key) == 0 {
		return false
	}
	kv.mux.Lock()
	defer kv.mux.Unlock()
	elem, ok := kv.kvstorage[key]
	if ok {
		return elem.Updated
	}
	return ok
}

// IsElemTTLOver - функция проверяет должен ли элемент быть удален
func (kv *KVStorage) IsElemTTLOver(key string, ttl uint64) bool {
	if len(key) == 0 {
		return false
	}
	kv.mux.Lock()
	defer kv.mux.Unlock()
	if elem, ok := kv.kvstorage[key]; ok {
		return elem.IsTTLOver(ttl)
	}
	return true
}

// IsInStorage - проверка нахождения элемента с данным ключем в хранилище
func (kv *KVStorage) IsInStorage(key string) bool {
	if len(key) == 0 {
		return false
	}
	kv.mux.Lock()
	defer kv.mux.Unlock()
	_, ok := kv.kvstorage[key]
	return ok
}
