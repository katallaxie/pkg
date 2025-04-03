package fsmx_test

import (
	"context"
	"testing"

	"github.com/katallaxie/pkg/fsmx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHook(t *testing.T) {
	t.Parallel()

	type noopState struct {
		Text string
	}

	tests := []struct {
		name     string
		state    noopState
		expected fsmx.State
		reducer  fsmx.Reducer[noopState]
		hook     fsmx.EffectFunc[noopState]
	}{
		{
			name: "non nil state",
			state: noopState{
				Text: "foo",
			},
			expected: noopState{
				Text: "bar",
			},
			reducer: func(prev noopState, action fsmx.Action) noopState {
				return noopState{
					Text: action.Payload().(string),
				}
			},
			hook: func(ctx context.Context) (fsmx.Action, error) {
				return fsmx.NewAction(1, "bar"), nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := fsmx.New(tt.state, tt.reducer)
			err := fsmx.Effect(store, tt.hook)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, store.State())
		})
	}
}
