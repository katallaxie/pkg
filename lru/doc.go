// Package lru is a simple implementation of an LRU cache.
// LRU is enhanced by allowing to set a TTL for cached items.
// It also supports to fetch the values for non-existing keys,
// by providing a function to fetch these.
package lru
