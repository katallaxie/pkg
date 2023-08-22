package ulid_test

import (
	"fmt"
	"testing"

	"github.com/katallaxie/pkg/ulid"

	"github.com/stretchr/testify/assert"
)

func ExampleULID() {
	fmt.Println(ulid.MustNew())
}

func TestNew(t *testing.T) {
	t.Parallel()

	u, err := ulid.New()
	assert.NoError(t, err)
	assert.NotEmpty(t, u)
}

func TestNewReverse(t *testing.T) {
	t.Parallel()

	u, err := ulid.NewReverse()
	assert.NoError(t, err)
	assert.NotEmpty(t, u)
}

func BenchmarkNew(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = ulid.New()
		}
	})
}

func BenchmarkNewReverse(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = ulid.NewReverse()
		}
	})
}

func BenchmarkMax(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = ulid.Max()
			assert.NotEmpty(b, ulid.Max())
		}
	})
}

func BenchmarkParse(b *testing.B) {
	u := ulid.MustNew()
	bb := u.Bytes()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = ulid.Parse(bb)
		}
	})
}

func BenchmarkParseString(b *testing.B) {
	u := ulid.MustNew()
	s := u.String()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = ulid.ParseString(s)
		}
	})
}
