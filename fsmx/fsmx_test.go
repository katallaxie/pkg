package fsmx_test

import (
	"testing"

	"github.com/katallaxie/pkg/fsmx"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	s := fsmx.New(nil)
	require.NotNil(t, s)
}

func TestDispatch(t *testing.T) {
	tests := []struct {
		name     string
		state    fsmx.State
		expected fsmx.State
		reducers []fsmx.Reducable
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
			reducers: []fsmx.Reducable{
				func(prev fsmx.State, action fsmx.Actionable) fsmx.State {
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
			require.NotNil(t, output)
			require.Equal(t, tt.expected, output)

			s.Drain()
		})
	}
}

func TestDrain(t *testing.T) {
	s := fsmx.New(nil)
	require.NotNil(t, s)

	sub := s.Subscribe()
	s.Drain()

	_, ok := <-sub
	require.False(t, ok)
}
