package vacuum

import (
	"kvserver/kvstorage"
	"log"
	"time"
)

// Vacuum - struct for cleaner
type Vacuum struct {
	storage     *kvstorage.KVStorage
	ttl         uint64
	ttlDelim    uint
	initialized bool
}

// Init - функция инициализации структуры Lifo
func (q *Vacuum) Init(stor *kvstorage.KVStorage, ttl uint64) bool {
	if stor == nil || q.initialized || ttl == 0 {
		return false
	}
	q.storage = stor
	q.ttl = ttl
	q.ttlDelim = 2 // hard coded time delimiter !
	q.initialized = true
	return true
}

// Run - infinite storage cleaner
func (q *Vacuum) Run() {
	if !q.initialized {
		log.Fatalln("Cleaner is not properly initialized.")
	}
	// we need to hit the oldest element periodically
	for {
		element, err := q.storage.OldestElementTime()
		sleepPeriod := getSleepPeriod(element, err, q.ttl, q.ttlDelim)
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

func getSleepPeriod(elementTime time.Time, err error, ttl uint64, ttlDelim uint) time.Duration {
	if err != nil {
		return time.Duration(
			float64(ttl) * float64(time.Second) / float64(ttlDelim),
		)
	}
	// need to handle special case scenario when
	// either no ttl or ttlDelim provided or these are wrong
	if ttl < 1 || ttlDelim < 2 {
		return time.Duration(1 * time.Second)
	}
	oldestElementFinalTime := elementTime.Add(
		time.Duration(
			int64(ttl) * int64(time.Second),
		),
	)

	timeDiffNS := float64(
		oldestElementFinalTime.Sub(time.Now()).Nanoseconds(),
	) / float64(ttlDelim)

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
