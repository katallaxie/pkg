package redux_test

import (
	"context"
	"testing"

	"github.com/katallaxie/pkg/redux"
	"github.com/stretchr/testify/assert"
)

func TestBackgroundOnce(t *testing.T) {
	t.Parallel()

	type noopState struct {
		Text string
	}

	tests := []struct {
		name     string
		state    noopState
		expected redux.State
		reducer  redux.Reducer[noopState]
		effect   redux.EffectFunc
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
			effect: func(ctx context.Context) redux.Action {
				return redux.NewAction(1, "bar")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := redux.New(tt.state, tt.reducer)
			sub := store.Subscribe()
			defer store.CancelSubscription(sub)

			changeString := redux.BackgroundOnce(context.Background(), store, tt.effect)
			changeString()

			state := <-sub
			assert.Equal(t, tt.expected, state.Curr())
		})
	}
}
