package lru

import (
	"time"
)

// New creates a new Cache of the given size
func New(size int) (Cache, error) {
	lru, err := NewLRU(size)
	if err != nil {
		return nil, err
	}

	c := &SimpleCache{
		lru: lru,
	}

	return c, nil
}

// Adds a value to the cache, or updates an item in the cache.
// It returns true if an item needed to be removed for storing the new item.
func (c *SimpleCache) Add(key interface{}, value interface{}, ttl int64) bool {
	c.Lock()
	defer c.Unlock()

	return c.lru.Add(key, value, ttl)
}

// Returns the value of the provided key, and updates status of the item
// in the cache.
func (c *SimpleCache) Get(key interface{}) (value interface{}, ok bool) {
	c.Lock()
	defer c.Unlock()

	return c.lru.Get(key)
}

// Check if a key exsists in the cache.
func (c *SimpleCache) Contains(key interface{}) (ok bool) {
	c.Lock()
	defer c.Unlock()

	return c.lru.Contains(key)
}

// Expires returns the time of expiration.
func (c *SimpleCache) Expires(key interface{}) (expires time.Time, ok bool) {
	c.RLock()
	defer c.RUnlock()

	return c.lru.Expires(key)
}

// Fetches a value which has expired, or does not exits and fills the cache.
func (c *SimpleCache) Fetch(key interface{}, ttl int64, call func() (interface{}, error)) (value interface{}, ok bool, err error) {
	c.Lock()
	defer c.Unlock()

	return c.lru.Fetch(key, ttl, call)
}

// Removes a key from the cache.
func (c *SimpleCache) Remove(key interface{}) bool {
	c.Lock()
	defer c.Unlock()

	return c.lru.Remove(key)
}

// Removes the oldest entry from cache.
func (c *SimpleCache) RemoveOldest() (interface{}, interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	return c.lru.RemoveOldest()
}

// Returns the oldest entry from the cache.
func (c *SimpleCache) GetOldest() (interface{}, interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	return c.lru.GetOldest()
}

// Returns a slice of the keys in the cache, from oldest to newest.
func (c *SimpleCache) Keys() []interface{} {
	c.RLock()
	defer c.RUnlock()

	return c.lru.Keys()
}

// Returns the number of items in the cache.
func (c *SimpleCache) Len() int {
	c.RLock()
	defer c.RUnlock()

	return c.lru.Len()
}

// Purge is purging the full cache.
func (c *SimpleCache) Purge() {
	c.Lock()
	defer c.Unlock()

	c.lru.Purge()
}
