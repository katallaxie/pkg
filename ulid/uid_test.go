package ulid_test

import (
	"testing"

	"github.com/katallaxie/pkg/ulid"
)

func BenchmarkMax(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ulid.Max()
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
