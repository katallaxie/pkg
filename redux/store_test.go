package redux_test

import (
	"fmt"
	"testing"

	"github.com/katallaxie/pkg/redux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleNew() {
	type noopState struct {
		Name string
	}

	r := func(prev noopState, action redux.Action) noopState {
		return noopState{Name: "bar"}
	}
	a := redux.NewAction(1, "foo")

	s := redux.New(noopState{Name: "foo"}, r)
	defer s.Dispose()

	sub := s.Subscribe()
	s.Dispatch(a)

	out := <-sub
	fmt.Println(out)
	// Output: {bar}
}

func TestNew(t *testing.T) {
	t.Parallel()

	noopState := struct{}{}

	s := redux.New(noopState)
	require.NotNil(t, s)
}

func TestNewAction(t *testing.T) {
	t.Parallel()

	a := redux.NewAction(1, "foo")
	require.NotNil(t, a)
	assert.Equal(t, redux.ActionType(1), a.Type())
	assert.Equal(t, "foo", a.Payload())

	a.Payload("bar")
	assert.Equal(t, "bar", a.Payload())

	a.Type(2)
	assert.Equal(t, redux.ActionType(2), a.Type())
}

func TestDispatch(t *testing.T) {
	t.Parallel()

	type noopState struct {
		Text string
	}

	tests := []struct {
		name     string
		state    noopState
		expected redux.State
		reducers []redux.Reducer[noopState]
	}{
		{
			name: "non nil state",
			state: noopState{
				Text: "foo",
			},
			expected: noopState{
				Text: "bar",
			},
			reducers: []redux.Reducer[noopState]{
				func(prev noopState, action redux.Action) noopState {
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

			s := redux.New(tt.state, tt.reducers...)
			defer s.Dispose()
			require.NotNil(t, s)

			sub := s.Subscribe()
			s.Dispatch(nil)

			output := <-sub
			defer s.CancelSubscription(sub)
			require.NotNil(t, output)
			require.Equal(t, tt.expected, output)
		})
	}
}

func TestDispose(t *testing.T) {
	t.Parallel()

	noopState := struct{}{}

	s := redux.New(noopState)
	require.NotNil(t, s)

	sub := s.Subscribe()
	s.Dispose()

	_, ok := <-sub
	require.False(t, ok)
}
