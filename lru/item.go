package lru

import (
	"sync/atomic"
	"time"
)

// Expired is used to check if the TTL of the item has expired
func (i *Item) Expired() bool {
	ttl := atomic.LoadInt64(&i.ttl)
	t := atomic.LoadInt64(&i.timestamp)
	return ttl > 0 && (time.Now().UnixNano()-t) > ttl
}

// Expires returns the time.Time when the item expires
func (i *Item) Expires() time.Time {
	ttl := atomic.LoadInt64(&i.ttl)
	t := atomic.LoadInt64(&i.timestamp)
	return time.Unix(0, t+ttl)
}

// Value returns the value of the item
func (i *Item) Value() interface{} {
	return i.value
}

// newItem is creating a new item in the LRU
func newItem(key interface{}, value interface{}, ttl int64) *Item {
	s := int64(1)
	if sized, ok := value.(Sized); ok {
		s = sized.Size()
	}

	return &Item{
		key:       key,
		value:     value,
		size:      s,
		ttl:       ttl,
		timestamp: time.Now().UnixNano(),
	}
}
