package redux_test

import (
	"context"
	"testing"

	"github.com/katallaxie/pkg/redux"
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
		expected redux.State
		reducer  redux.Reducer[noopState]
		hook     redux.EffectFunc
	}{
		{
			name: "non nil state",
			state: noopState{
				Text: "foo",
			},
			expected: noopState{
				Text: "bar",
			},
			reducer: func(prev noopState, action redux.Action) noopState {
				return noopState{
					Text: action.Payload().(string),
				}
			},
			hook: func(ctx context.Context) (redux.Action, error) {
				return redux.NewAction(1, "bar"), nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := redux.New(tt.state, tt.reducer)
			err := redux.Effect(store, tt.hook)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, store.State())
		})
	}
}
