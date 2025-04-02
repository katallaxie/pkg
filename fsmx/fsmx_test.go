package fsmx_test

import (
	"testing"

	"github.com/katallaxie/pkg/fsmx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	noopState := struct{}{}

	s := fsmx.New(noopState)
	require.NotNil(t, s)
}

func TestNewAction(t *testing.T) {
	t.Parallel()

	a := fsmx.NewAction(1, "foo")
	require.NotNil(t, a)
	assert.Equal(t, fsmx.ActionType(1), a.Type())
	assert.Equal(t, "foo", a.Payload())

	a.Payload("bar")
	assert.Equal(t, "bar", a.Payload())

	a.Type(2)
	assert.Equal(t, fsmx.ActionType(2), a.Type())
}

func TestDispatch(t *testing.T) {
	t.Parallel()

	type noopState struct {
		Text string
	}

	tests := []struct {
		name     string
		state    noopState
		expected fsmx.State
		reducers []fsmx.Reducer[noopState]
	}{
		{
			name: "non nil state",
			state: noopState{
				Text: "foo",
			},
			expected: noopState{
				Text: "bar",
			},
			reducers: []fsmx.Reducer[noopState]{
				func(prev noopState, action fsmx.Action) noopState {
					return struct {
						Text string
					}{
						Text: "bar",
					}
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := fsmx.New(tt.state, tt.reducers...)
			require.NotNil(t, s)

			sub := s.Subscribe()
			s.Dispatch(nil)

			output := <-sub
			defer s.CancelSubscription(sub)
			require.NotNil(t, output)
			require.Equal(t, tt.expected, output)

			s.Cancel()
		})
	}
}

func TestDrain(t *testing.T) {
	t.Parallel()

	noopState := struct{}{}

	s := fsmx.New(noopState)
	require.NotNil(t, s)

	sub := s.Subscribe()
	s.Cancel()

	_, ok := <-sub
	require.False(t, ok)
}
