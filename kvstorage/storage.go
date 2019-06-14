package kvstorage

import (
	"container/list"
	"errors"
	"sync"
	"time"

	"github.com/proway2/kvserver/element"
)

// KVStorage - Структура с методами, описывающая хранилище
type KVStorage struct {
	kvstorage   map[string]*element.Element
	mux         sync.Mutex
	queue       list.List // LIFO - the oldest element is always at the front!!!
	initialized bool
}

// Init - Функция инициализации хранилища
func (kv *KVStorage) Init() bool {
	if kv.initialized {
		return false
	}
	kv.kvstorage = make(map[string]*element.Element)
	kv.initialized = true
	kv.queue.Init()
	return true
}

// Set - установка ключ-значение
func (kv *KVStorage) Set(key, value string) bool {
	if !kv.initialized || len(key) == 0 {
		return false
	}
	kv.mux.Lock()
	defer kv.mux.Unlock()

	// проверяем есть ли у нас такой ключ в карте
	if elem, found := kv.kvstorage[key]; found {
		// для поддержания порядка очереди LIFO,
		// надо удалить найденный элемент из очереди
		// вместо него будет новый с таким же ключом
		kv.queue.Remove(elem.QueueElement)
	}
	// in order to maintain LIFO new elements push back
	elem := &element.Element{
		Val:          value,
		Timestamp:    time.Now(),
		QueueElement: kv.queue.PushBack(key),
	}
	kv.kvstorage[key] = elem

	return true
}

// Get - получение значения по ключу,
// второе возвращаемое значение указывает на успешность получения значения по ключу
func (kv *KVStorage) Get(key string) (string, bool) {
	if !kv.initialized || len(key) == 0 {
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

// OldestElementTime - получить метку времени старейшего элемента
func (kv *KVStorage) OldestElementTime() (time.Time, error) {
	if !kv.initialized {
		return time.Time{}, errors.New("storage is not initialized")
	}
	kv.mux.Lock()
	defer kv.mux.Unlock()

	oldestElemInQueue := kv.queue.Front()
	if oldestElemInQueue == nil {
		return time.Time{}, errors.New("not found in storage")
	}

	key := oldestElemInQueue.Value.(string)
	elem, ok := kv.kvstorage[key]
	if ok {
		return elem.Timestamp, nil
	}
	// WTF? Element is in queue, but not in the map?!
	return time.Time{}, errors.New("shit happens")
}

// Delete - delete element from storage by its key
func (kv *KVStorage) Delete(key string) bool {
	if !kv.initialized || len(key) == 0 {
		return false
	}
	kv.mux.Lock()
	defer kv.mux.Unlock()
	_, ok := kv.kvstorage[key]
	if ok {
		kv.purgeElement(key)
	}
	return ok
}

// DeleteFrontIfOlder - removes front element if its older than ctxTime
func (kv *KVStorage) DeleteFrontIfOlder(ctxTime time.Time) bool {
	if !kv.initialized {
		return false
	}
	kv.mux.Lock()
	defer kv.mux.Unlock()

	// first need to check if there is something in the queue
	oldestElemInQueue := kv.queue.Front()
	if oldestElemInQueue == nil {
		return false
	}

	// map key is needed to test the element against input time
	key := oldestElemInQueue.Value.(string)
	// need to check if it's in storage.
	// it MUST be in storage at this point
	elem, _ := kv.kvstorage[key]
	if elem.Timestamp.Before(ctxTime) {
		kv.purgeElement(key)
		return true
	}
	return false
}

func (kv *KVStorage) purgeElement(key string) {
	// THIS IS NOT THREAD SAFE FUNCTION !!!
	// INTERNAL USE ONLY !!!
	// CALL THIS FUNCTION WITHIN CRITICAL SECTION
	// WHEN THREAD IS LOCKED !!!
	// NOT INTENDED FOR SEPARATE USE !!!
	kv.queue.Remove(kv.kvstorage[key].QueueElement)
	delete(kv.kvstorage, key)
}
