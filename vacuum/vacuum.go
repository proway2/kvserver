package vacuum

import (
	"container/list"
	"kvserver/kvstorage"
)

// Lifo - структура, описывающая очередь элементов, тип LIFO
type Lifo struct {
	storage     *kvstorage.KVStorage
	queue       list.List
	inpElemChan *chan string
	ttl         uint64
}

// Init - функция инициализации структуры Lifo
func (q *Lifo) Init(
	stor *kvstorage.KVStorage,
	in *chan string,
	ttl uint64) bool {

	if in == nil || stor == nil || ttl == 0 {
		return false
	}
	q.storage = stor
	q.inpElemChan = in
	q.queue.Init()
	q.ttl = ttl
	return true
}

// Run - метод, запускающий бесконечный цикл очистки
func (q *Lifo) Run() {
	if q.inpElemChan == nil || q.storage == nil || q.ttl == 0 {
		return
	}
	// цикл на прием сообщений из канала и запись их в конец очереди
	// при блокировании канала выполняется очистка очереди и удаление
	// элементов с истекшим TTL
	for {
		select {
		case elem := <-*q.inpElemChan:
			q.queue.PushBack(elem)
			q.cleanUp(q.queue.Front())
		default:
			if elem := q.queue.Front(); elem != nil {
				q.cleanUp(elem)
			}
		}
	}
}

func (q *Lifo) cleanUp(elem *list.Element) {
	key := elem.Value.(string)
	if !q.storage.IsInStorage(key) {
		// элемента нет в хранилище - очистка очереди
		q.queue.Remove(elem)
		return
	}
	// проверяем обновление элемента и его TTL
	// при необходимости очищаем из хранилища и очереди
	if q.storage.IsElemUpdated(key) {
		// элемент был обновлен, просто удаляем его из очереди
		q.queue.Remove(elem)
		// отправляем обратно значение ключа, которое надо обновить
		q.storage.ResetUpdated(key)
	} else {
		// перед удалением из хранилища надо проверить TTL
		if q.storage.IsElemTTLOver(key, q.ttl) {
			q.storage.Delete(key) // удаление из хранилища
			q.queue.Remove(elem)  // удаление из очереди
		}
	}
}
