package async_test

import (
	"testing"

	"github.com/katallaxie/pkg/async"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkAll(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	b.Run("resolve all", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			future := async.All(
				async.New(func() (string, error) {
					return "hello", nil
				}),
				async.New(func() (string, error) {
					return "world", nil
				}),
				async.New(func() (string, error) {
					return "world", nil
				}),
			)

			v, err := future.Await()

			assert.NotNil(b, v)
			require.NoError(b, err)
			assert.Equal(b, []string{"hello", "world", "world"}, v)
		}
	})
}
