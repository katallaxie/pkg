package opts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	Verbose Opt = iota
)

func BenchmarkOpts_Get(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		o := New[Opt, any]()
		o.Set(Verbose, true)

		for pb.Next() {
			_, _ = o.Get(Verbose)
		}
	})
}

func TestConfig_NewDefaultOpts(t *testing.T) {
	Verbose := Opt(0)

	cond := []struct {
		desc string
		in   Opt
		out  interface{}
	}{
		{desc: "", in: Verbose, out: Opt(0)},
	}

	for _, tt := range cond {
		t.Run(tt.desc, func(t *testing.T) {
			o := New[int, any]()
			o.Set(tt.in, tt.out)

			v, err := o.Get(tt.in)
			require.NoError(t, err)
			assert.Equal(t, tt.out, v)
		})
	}
}

func TestConfig_WithLogger(t *testing.T) {
	logger, err := zap.NewProduction()
	defer func() { _ = logger.Sync() }()
	require.NoError(t, err)

	Logger := Opt(0)

	o := New[int, any]()
	o.Set(Logger, logger)

	v, err := o.Get(0)
	require.NoError(t, err)
	assert.NotNil(t, v)
}
