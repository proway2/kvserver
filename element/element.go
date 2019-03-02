package element

import (
	"container/list"
	"time"
)

// Element - структура описывающая один элемент хранилища
type Element struct {
	Val          string
	Timestamp    time.Time
	QueueElement *list.Element
}
