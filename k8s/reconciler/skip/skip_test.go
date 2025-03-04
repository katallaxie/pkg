package skip_test

import (
	"context"
	"testing"

	"github.com/katallaxie/pkg/k8s/reconciler/skip"

	"github.com/stretchr/testify/require"
)

func TestSkipEnableSkip(t *testing.T) {
	ctx := context.Background()
	ctx = skip.EnableSkip(ctx)
	require.True(t, skip.Skip(ctx))
}
