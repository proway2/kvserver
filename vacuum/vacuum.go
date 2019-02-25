package vacuum

import (
	"container/list"
	"kvserver/kvstorage"
	"time"
)

// Lifo - структура, описывающая очередь элементов, тип LIFO
type Lifo struct {
	storage     *kvstorage.KVStorage
	queue       list.List
	inpElemChan chan string
	ttl         uint64
	ttlDelim    uint
}

// Init - функция инициализации структуры Lifo
func (q *Lifo) Init(
	stor *kvstorage.KVStorage,
	in chan string,
	ttl uint64) bool {

	if in == nil || stor == nil || ttl == 0 {
		return false
	}
	q.storage = stor
	q.inpElemChan = in
	q.queue.Init()
	q.ttl = ttl
	q.ttlDelim = 2 // hard coded делитель времени !
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
		sleepPeriod := q.getSleepPeriod()
		select {
		case elem := <-q.inpElemChan:
			q.queue.PushBack(elem)
			q.cleanUp(q.queue.Front())
		case <-time.After(sleepPeriod):
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
		_ = q.storage.ResetUpdated(key)
	} else {
		// перед удалением из хранилища надо проверить TTL
		if q.storage.IsElemTTLOver(key, q.ttl) {
			_ = q.storage.Delete(key) // удаление из хранилища
			q.queue.Remove(elem)      // удаление из очереди
		}
	}
}

func (q *Lifo) getSleepPeriod() time.Duration {

	element := q.queue.Front()
	if element == nil {
		return time.Duration(
			float64(q.ttl) * float64(time.Second) / float64(q.ttlDelim),
		)
	}

	elemUnixTime := q.storage.GetTimestamp(element.Value.(string))
	// надо проверить правильность ответа времени элемента
	if elemUnixTime == 0 {
		// элемент в очереди, но не в хранилище - удалить из очереди
		return time.Duration(0)
	}

	elemTime := time.Unix(
		elemUnixTime,
		0,
	)
	ttlDuration := time.Duration(q.ttl * uint64(time.Second))
	elemExpireTime := elemTime.Add(ttlDuration)
	// надо проверить закончился TTL или нет
	if elemTime.Add(ttlDuration).Before(time.Now()) {
		// надо срочно удалить этот элемент
		return time.Duration(0)
	}

	sleepPeriodNS := float64(elemExpireTime.Sub(time.Now()).Nanoseconds()) / float64(q.ttlDelim)
	sleepDuration := time.Duration(
		sleepPeriodNS * float64(time.Nanosecond),
	)
	if sleepDuration.Nanoseconds() < 1.0 {
		return time.Duration(1 * time.Nanosecond)
	}

	return sleepDuration
}
