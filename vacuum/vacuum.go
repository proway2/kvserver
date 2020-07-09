package vacuum

import (
	"errors"
	"log"
	"time"
)

type writer interface {
	OldestElementTime() (time.Time, error)
	DeleteFrontIfOlder(time.Time) (bool, error)
}

// Vacuum - struct for cleaner
type Vacuum struct {
	storage     writer
	ttl         uint64
	ttlDelim    uint
	initialized bool
}

// NewCleaner returns an initialized cleaner for storage 'w' with TTL of 'ttl'
func NewCleaner(w writer, ttl uint64) (*Vacuum, error) {
	if w == nil || ttl == 0 {
		return &Vacuum{}, errors.New("newCleaner: no storage provided or TTL = 0")
	}
	return &Vacuum{
		storage:     w,
		ttl:         ttl,
		ttlDelim:    2,
		initialized: true,
	}, nil
}

// Run - infinite storage cleaner
func (q *Vacuum) Run() {
	if !q.initialized {
		log.Fatalln("Cleaner is not properly initialized.")
	}
	// we need to hit the oldest element periodically
	for {
		elementTime, err := q.storage.OldestElementTime()
		var sleepPeriod time.Duration
		if err != nil {
			sleepPeriod = getSleepPeriodEmptyQueue(q.ttl, q.ttlDelim)
		} else {
			sleepPeriod = getSleepPeriod(elementTime, nil, q.ttl, q.ttlDelim)
		}
		select {
		case <-time.After(sleepPeriod):
			testTime := time.Now().Add(
				time.Duration(-q.ttl * uint64(time.Second)),
			)
			q.storage.DeleteFrontIfOlder(testTime)
		}
	}
}

func getSleepPeriodEmptyQueue(ttl uint64, ttlDelim uint) time.Duration {
	return time.Duration(
		float64(ttl) * float64(time.Second) / float64(ttlDelim),
	)
}

func getSleepPeriod(elementTime time.Time, err error, ttl uint64, ttlDelim uint) time.Duration {
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
