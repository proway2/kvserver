package vacuum

import (
	"kvserver/kvstorage"
	"log"
	"time"
)

// Lifo - структура, описывающая очередь элементов, тип LIFO
type Lifo struct {
	storage     *kvstorage.KVStorage
	ttl         uint64
	ttlDelim    uint
	initialized bool
}

// Init - функция инициализации структуры Lifo
func (q *Lifo) Init(stor *kvstorage.KVStorage, ttl uint64) bool {
	if stor == nil || q.initialized {
		return false
	}
	q.storage = stor
	q.ttl = ttl
	q.ttlDelim = 2 // hard coded делитель времени !
	q.initialized = true
	return true
}

// Run - infinite storage cleaner
func (q *Lifo) Run() {
	if !q.initialized {
		log.Fatalln("Cleaner is not properly initialized.")
	}
	// we need to hit the oldest element periodically
	for {
		sleepPeriod := q.getSleepPeriod()
		select {
		case <-time.After(sleepPeriod):
			testTime := time.Now().Add(
				time.Duration(
					-q.ttl * uint64(time.Second),
				),
			)
			q.storage.DeleteFrontIfOlder(testTime)
		}
	}
}

func (q *Lifo) getSleepPeriod() time.Duration {
	// to calculate sleeping time the oldest's element
	// in queue time must be known
	oldestElementTime, err := q.storage.OldestElementTime()
	if err != nil {
		return time.Duration(
			float64(q.ttl) * float64(time.Second) / float64(q.ttlDelim),
		)
	}
	oldestElementFinalTime := oldestElementTime.Add(
		time.Duration(
			int64(q.ttl) * int64(time.Second),
		),
	)

	timeDiffNS := float64(
		oldestElementFinalTime.Sub(time.Now()).Nanoseconds(),
	) / float64(q.ttlDelim)

	// to handle already expired elements must check for negative numbers
	if timeDiffNS < 0.0 {
		return time.Duration(0 * time.Nanosecond)
	}

	sleepDuration := time.Duration(timeDiffNS * float64(time.Nanosecond))
	if sleepDuration.Nanoseconds() < 1.0 {
		return time.Duration(0 * time.Nanosecond)
	}
	return sleepDuration
}
