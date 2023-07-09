package lru

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// ErrNoPositiveSize indicates that the size for the LRU cache
// must be a positive number.
var ErrNoPositiveSize = errors.New("lru: must provide a positive size")

// Cache is the interface for the LRU cache
type Cache interface {
	// Adds a value to the cache, or updates an item in the cache.
	// It returns true if an item needed to be removed for storing the new item.
	Add(key interface{}, value interface{}, ttl int64) bool

	// Returns the value of the provided key, and updates status of the item
	// in the cache.
	Get(key interface{}) (value interface{}, ok bool)

	// Check if a key exsists in the cache.
	Contains(key interface{}) (ok bool)

	// Expires returns the time of expiration.
	Expires(key interface{}) (expires time.Time, ok bool)

	// Fetches a value which has expired, or does not exits and fills the cache.
	Fetch(key interface{}, ttl int64, call func() (interface{}, error)) (value interface{}, ok bool, err error)

	// Removes a key from the cache.
	Remove(key interface{}) bool

	// Removes the oldest entry from cache.
	RemoveOldest() (interface{}, interface{}, bool)

	// Returns the oldest entry from the cache.
	GetOldest() (interface{}, interface{}, bool)

	// Returns a slice of the keys in the cache, from oldest to newest.
	Keys() []interface{}

	// Returns the number of items in the cache.
	Len() int

	// Purge is purging the full cache.
	Purge()
}

// cache is a thread-safe, fixed size LRU cache with TTL.
type SimpleCache struct {
	sync.RWMutex

	lru Cache
}

// Sized is the interface to a cache size calculation.
type Sized interface {
	Size() int64
}

// LRUCache represents the instance of an LRU cache.
type LRUCache struct {
	size  int
	items map[interface{}]*list.Element
	list  *list.List
}

// Item represents the internal presentation of a cache entry
type Item struct {
	key       interface{}
	ttl       int64
	timestamp int64
	size      int64
	value     interface{}
}

var _ Cache = (*LRUCache)(nil)

// NewLRU returns a new instance of the LRU cache with
// a certain size of elements that can be stored in time.
func NewLRU(size int) (Cache, error) {
	if size <= 0 {
		return nil, ErrNoPositiveSize
	}

	c := &LRUCache{
		size:  size,
		list:  list.New(),
		items: make(map[interface{}]*list.Element),
	}

	return c, nil
}

// Add is adding a key and value with a TTL to the store.
// Setting the TTL to 0 signales that this key will not expire.
func (l *LRUCache) Add(key interface{}, value interface{}, ttl int64) bool {
	if e, ok := l.items[key]; ok {
		l.list.MoveToFront(e)
		e.Value.(*Item).value = value
		e.Value.(*Item).ttl = ttl
		e.Value.(*Item).timestamp = time.Now().UnixNano()
	}

	e := newItem(key, value, ttl)
	h := l.list.PushFront(e)
	l.items[key] = h

	rm := l.list.Len() > l.size
	if rm {
		l.remove()
	}

	return rm
}

// Get is returning the value of a provided key.
func (l *LRUCache) Get(key interface{}) (value interface{}, ok bool) {
	e, ok := l.items[key]
	if !ok {
		return
	}

	if e.Value.(*Item).Expired() {
		l.removeElement(e)

		return nil, false
	}

	l.list.MoveToFront(e)
	e.Value.(*Item).timestamp = time.Now().UnixNano()

	return e.Value.(*Item).value, true
}

// Contains is checking if a provided key exists in the cache
func (l *LRUCache) Contains(key interface{}) (ok bool) {
	_, ok = l.items[key]
	return ok
}

// Fetch is fetching a value for key that does not exists or has expired.
// The fetching is done by a provided function.
func (l *LRUCache) Fetch(key interface{}, ttl int64, call func() (interface{}, error)) (value interface{}, ok bool, err error) {
	if e, ok := l.Get(key); ok {
		return e, ok, nil
	}

	v, err := call()
	if err != nil {
		return nil, false, err
	}

	e := newItem(key, v, ttl)
	h := l.list.PushFront(e)
	l.items[key] = h

	rm := l.list.Len() > l.size
	if rm {
		l.remove()
	}

	return e.Value(), rm, nil
}

// GetOldest returns the oldest item of the cache.
func (l *LRUCache) GetOldest() (key interface{}, value interface{}, ok bool) {
	e := l.list.Back()
	if e != nil {
		kv := e.Value.(*Item)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

// RemoveOldest removes the oldest item in the cache.
func (l *LRUCache) RemoveOldest() (key interface{}, value interface{}, ok bool) {
	e := l.list.Back()
	if e != nil {
		l.removeElement(e)
		kv := e.Value.(*Item)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

// Remove is removing a provided key from the cache.
func (l *LRUCache) Remove(key interface{}) (ok bool) {
	if e, ok := l.items[key]; ok {
		l.removeElement(e)
		return true
	}
	return false
}

// Expires returns the time.Time when the provided key will expire.
func (l *LRUCache) Expires(key interface{}) (expires time.Time, ok bool) {
	if e, ok := l.items[key]; ok {
		return e.Value.(*Item).Expires(), true
	}
	return
}

// Keys returning the keys of the current cache.
func (l *LRUCache) Keys() []interface{} {
	keys := make([]interface{}, len(l.items))
	i := 0
	for e := l.list.Back(); e != nil; e = e.Prev() {
		keys[i] = e.Value.(*Item).key
		i++
	}
	return keys
}

// Len returns the length/number of elements that are in the cache.
func (l *LRUCache) Len() int {
	return l.list.Len()
}

// Purge is purging the cache.
func (l *LRUCache) Purge() {
	for k := range l.items {
		delete(l.items, k)
	}
	l.list.Init()
}

func (l *LRUCache) remove() {
	e := l.list.Back()
	if e != nil {
		l.removeElement(e)
	}
}

func (l *LRUCache) removeElement(e *list.Element) {
	l.list.Remove(e)
	kv := e.Value.(*Item)
	delete(l.items, kv.key)
}
