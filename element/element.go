package element

import (
	"container/list"
	"time"
)

// Element - структура описывающая один элемент хранилища
type Element struct {
	Val          string        // the actual value of the element
	Timestamp    time.Time     // time when element is created or updated
	QueueElement *list.Element // pointer to the position in the queue (LIFO stack)
}
