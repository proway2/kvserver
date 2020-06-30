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
	mux         *sync.Mutex
	queue       *list.List // LIFO - the oldest element is always at the front!!!
	initialized bool
}

// NewStorage returns an initialized key-value storage
func NewStorage() *KVStorage {
	return &KVStorage{
		kvstorage:   make(map[string]*element.Element),
		mux:         &sync.Mutex{},
		initialized: true,
		queue:       list.New(),
	}
}

// Set adds new or updates existing element into the storage
func (kv *KVStorage) Set(key, value string) error {
	if !kv.initialized || len(key) == 0 {
		return errors.New("set: Storage is not initialized or key is empty")
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
	// in order to maintain LIFO new elements pushed back
	elem := &element.Element{
		Val:          value,
		Timestamp:    time.Now(),
		QueueElement: kv.queue.PushBack(key),
	}
	kv.kvstorage[key] = elem

	return nil
}

// Get returns value by it's key
func (kv *KVStorage) Get(key string) ([]byte, error) {
	if !kv.initialized || len(key) == 0 {
		return nil, errors.New("get: Storage is not initialized or key is empty")
	}
	kv.mux.Lock()
	defer kv.mux.Unlock()
	elem, ok := kv.kvstorage[key]
	if ok {
		return []byte(elem.Val), nil
	}
	// element with the key is not found, but this is not an error
	return nil, nil
}

// OldestElementTime - получить метку времени старейшего элемента
func (kv *KVStorage) OldestElementTime() (time.Time, error) {
	if !kv.initialized {
		return time.Time{}, errors.New("oldestelementtime: Storage is not initialized")
	}
	kv.mux.Lock()
	defer kv.mux.Unlock()

	oldestElemInQueue := kv.queue.Front()
	if oldestElemInQueue == nil {
		return time.Time{}, errors.New("oldestelementtime: Element is not found in storage")
	}

	key := oldestElemInQueue.Value.(string)
	elem, ok := kv.kvstorage[key]
	if ok {
		return elem.Timestamp, nil
	}
	// WTF? Element is in queue, but not in the map?!
	panic("element is in the queue, but not in the map")
}

// Delete removes element from storage by its key
func (kv *KVStorage) Delete(key string) (bool, error) {
	if !kv.initialized || len(key) == 0 {
		return false, errors.New("delete: Storage is not initialized or key is empty")
	}
	kv.mux.Lock()
	defer kv.mux.Unlock()
	_, ok := kv.kvstorage[key]
	if ok {
		kv.purgeElement(key)
	}
	return ok, nil
}

// DeleteFrontIfOlder - removes front element if its older than ctxTime
func (kv *KVStorage) DeleteFrontIfOlder(ctxTime time.Time) (bool, error) {
	if !kv.initialized {
		return false, errors.New("deletefrontifolder: Storage is not initialized")
	}
	kv.mux.Lock()
	defer kv.mux.Unlock()

	// first need to check if there is something in the queue
	oldestElemInQueue := kv.queue.Front()
	if oldestElemInQueue == nil {
		return false, nil
	}

	// map key is needed to test the element against input time
	key := oldestElemInQueue.Value.(string)
	// need to check if it's in storage.
	// it MUST be in storage at this point
	elem, _ := kv.kvstorage[key]
	if elem.Timestamp.Before(ctxTime) {
		kv.purgeElement(key)
		return true, nil
	}
	// the element is not in the storage right now
	return false, nil
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
