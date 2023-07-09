package lru_test

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/katallaxie/pkg/lru"
	"github.com/stretchr/testify/assert"
)

const (
	size = 64
)

func TestCache(t *testing.T) {
	tests := []struct {
		desc string
		ttl  int64
	}{
		{
			desc: "with ttl of 1min",
			ttl:  int64((5 * time.Second) / time.Nanosecond),
		},
		{
			desc: "with no ttl set",
			ttl:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			l, err := lru.New(size)

			assert.NoError(t, err)

			for i := 0; i < size; i++ {
				l.Add(i, i, tt.ttl)
			}

			assert.Equal(t, l.Len(), size)

			for i, k := range l.Keys() {
				v, ok := l.Get(k)

				assert.True(t, ok)
				assert.Equal(t, v, k)
				assert.NotEqual(t, v, i+size)
			}

			for i := 0; i < size; i++ {
				_, ok := l.Get(i)

				assert.True(t, ok)
			}

			for i := size; i < size*2; i++ {
				_, ok := l.Get(i)

				assert.False(t, ok)
			}

			for i := size; i < (size/2)+size; i++ {
				ok := l.Remove(i)
				assert.False(t, ok)

				ok = l.Remove(i)
				assert.False(t, ok)

				_, ok = l.Get(i)
				assert.False(t, ok)
			}

			l.Purge()
			assert.Equal(t, 0, l.Len())

			for i := 0; i < size; i++ {
				l.Add(i, i, tt.ttl)
			}

			if tt.ttl > 0 {
				time.Sleep(time.Duration(tt.ttl))

				for i := 0; i < size; i++ {
					_, ok := l.Get(i)

					assert.False(t, ok)
				}
			}
		})
	}
}

func TestCacheFetch(t *testing.T) {
	tests := []struct {
		desc  string
		value interface{}
		key   interface{}
		ttl   int64
		call  func() (interface{}, error)
	}{
		{
			desc:  "should fetch new value without err and no ttl",
			key:   1,
			value: 1,
			ttl:   0,
			call: func() (interface{}, error) {
				return 1, nil
			},
		},
		{
			desc:  "should fetch new value with err and no ttl",
			key:   1,
			value: 1,
			ttl:   0,
			call: func() (interface{}, error) {
				return nil, errors.New("no value")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			l, err := lru.New(1)

			assert.NoError(t, err)

			v, ok, err := l.Fetch(tt.key, tt.ttl, tt.call)
			if err != nil {
				assert.Error(t, err)
				assert.Nil(t, v)
				assert.False(t, ok)

				return
			}

			assert.Equal(t, tt.value, v)
			assert.False(t, ok)
			assert.NoError(t, err)
		})
	}
}

func TestCacheAdd(t *testing.T) {
	type item struct {
		key   interface{}
		value interface{}
		ttl   int64
	}

	tests := []struct {
		desc         string
		items        []item
		expectRemove bool
	}{
		{
			desc: "should not be removed",
			items: []item{
				{
					key:   1,
					value: 1,
					ttl:   0,
				},
			},
			expectRemove: false,
		},
		{
			desc: "should be removed",
			items: []item{
				{
					key:   1,
					value: 1,
					ttl:   0,
				},
				{
					key:   2,
					value: 2,
					ttl:   0,
				},
			},
			expectRemove: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			l, err := lru.NewLRU(1)

			assert.NoError(t, err)

			var ok bool
			for _, item := range tt.items {
				ok = l.Add(item.key, item.value, item.ttl)
			}

			assert.Equal(t, tt.expectRemove, ok)
		})
	}
}

func BenchmarkCache_Rand(b *testing.B) {
	tests := []struct {
		desc string
		size int
		ttl  int64
	}{
		{
			desc: "with size of 4096 and no ttl set",
			size: 4096,
			ttl:  0,
		},
		{
			desc: "with size of 8092 items and no ttl set",
			size: 8092,
			ttl:  0,
		},
	}

	for _, tt := range tests {
		b.Run(tt.desc, func(b *testing.B) {
			l, err := lru.New(tt.size)
			assert.NoError(b, err)

			trace := make([]int64, b.N*2)
			for i := 0; i < b.N*2; i++ {
				trace[i] = rand.Int63() % 32768
			}

			b.ResetTimer()

			var hit, miss int
			for i := 0; i < 2*b.N; i++ {
				if i%2 == 0 {
					l.Add(trace[i], trace[i], tt.ttl)
				} else {
					_, ok := l.Get(trace[i])
					if ok {
						hit++
					} else {
						miss++
					}
				}
			}
			b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
		})
	}
}

func BenchmarkCache_Freq(b *testing.B) {
	tests := []struct {
		desc string
		size int
		ttl  int64
	}{
		{
			desc: "with size of 4096 and no ttl set",
			size: 4096,
			ttl:  0,
		},
		{
			desc: "with size of 8092 items and no ttl set",
			size: 8092,
			ttl:  0,
		},
	}

	for _, tt := range tests {
		b.Run(tt.desc, func(b *testing.B) {
			l, err := lru.New(8192)
			assert.NoError(b, err)

			trace := make([]int64, b.N*2)
			for i := 0; i < b.N*2; i++ {
				if i%2 == 0 {
					trace[i] = rand.Int63() % 16384
				} else {
					trace[i] = rand.Int63() % 32768
				}
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				l.Add(trace[i], trace[i], tt.ttl)
			}
			var hit, miss int
			for i := 0; i < b.N; i++ {
				_, ok := l.Get(trace[i])
				if ok {
					hit++
				} else {
					miss++
				}
			}
			b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
		})
	}
}
