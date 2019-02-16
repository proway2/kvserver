package element

import "time"

// Element - структура описывающая один элемент хранилища
type Element struct {
	Val       string
	Timestamp int64
	Updated   bool
}

// IsTTLOver - принимает TTL и возвращает признак
// должен ли элемент быть удален
func (elem Element) IsTTLOver(ttl uint64) bool {
	if ttl < 1 {
		return true
	}
	res := time.Now().Unix() - elem.Timestamp
	return uint64(res) >= ttl
}
