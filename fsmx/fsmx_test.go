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

	tests := []struct {
		name     string
		state    fsmx.State
		expected fsmx.State
		reducers []fsmx.Reducer
	}{
		{
			name: "non nil state",
			state: struct {
				Text string
			}{
				Text: "foo",
			},
			expected: struct {
				Text string
			}{
				Text: "bar",
			},
			reducers: []fsmx.Reducer{
				func(prev fsmx.State, action fsmx.Action) fsmx.State {
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
			listener := sub.Listen()
			s.Dispatch(nil)

			output := <-listener
			defer sub.Cancel()
			require.NotNil(t, output)
			require.Equal(t, tt.expected, output)

			s.Drain()
		})
	}
}

func TestDrain(t *testing.T) {
	t.Parallel()

	noopState := struct{}{}

	s := fsmx.New(noopState)
	require.NotNil(t, s)

	sub := s.Subscribe()
	listener := sub.Listen()
	s.Drain()

	_, ok := <-listener
	require.False(t, ok)
}

func TestSubscriptionID(t *testing.T) {
	t.Parallel()

	noopState := struct{}{}

	s := fsmx.New(noopState)
	require.NotNil(t, s)

	sub := s.Subscribe()
	_ = sub.Listen()
	defer sub.Cancel()

	assert.NotEmpty(t, sub.ID())
	assert.NotEmpty(t, 1, sub.ID())

	sub2 := s.Subscribe()
	_ = sub2.Listen()
	defer sub2.Cancel()

	assert.NotEmpty(t, sub2.ID())
	assert.NotEmpty(t, 2, sub2.ID())
}
