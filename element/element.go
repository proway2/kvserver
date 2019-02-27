package element

import "time"

// Element - структура описывающая один элемент хранилища
type Element struct {
	Val       string
	Timestamp time.Time
}

// IsTTLOver - принимает TTL и возвращает признак
// старше элемент переданного времени или нет
func (elem Element) IsOlder(testTime time.Time) bool {
	return elem.Timestamp.Before(testTime)
}
