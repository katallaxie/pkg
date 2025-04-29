package redux_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/katallaxie/pkg/redux"
	"github.com/stretchr/testify/require"
)

func ExampleNew() {
	type noopState struct {
		Name string
	}

	r := func(prev noopState, action redux.Update) noopState {
		return noopState{Name: "bar"}
	}

	a := func() redux.Update {
		return "foo"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := redux.New(ctx, noopState{Name: "foo"}, r)
	defer s.Dispose()

	sub := s.Subscribe()
	s.Dispatch(a)

	out := <-sub
	fmt.Println(out.Curr(), out.Prev())
	// Output: {bar} {foo}
}

type fooMsg string

func fooAction() redux.Action {
	return func() redux.Update {
		return fooMsg("foo")
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	noopState := struct{}{}

	s := redux.New(t.Context(), noopState)
	require.NotNil(t, s)
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
				func(prev noopState, action redux.Update) noopState {
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
			s := redux.New(t.Context(), tt.state, tt.reducers...)
			defer s.Dispose()
			require.NotNil(t, s)

			sub := s.Subscribe()
			s.Dispatch(fooAction())

			output := <-sub
			defer s.CancelSubscription(sub)
			require.NotNil(t, output)
			require.Equal(t, tt.expected, output.Curr())
		})
	}
}

func BenchmarkDispatch(b *testing.B) {
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
				func(prev noopState, action redux.Update) noopState {
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
		b.Run(tt.name, func(b *testing.B) {
			s := redux.New(b.Context(), tt.state, tt.reducers...)
			defer s.Dispose()
			require.NotNil(b, s)

			sub := s.Subscribe()
			for i := 0; i < b.N; i++ {
				s.Dispatch(fooAction())
				output := <-sub
				defer s.CancelSubscription(sub)
				require.NotNil(b, output)
				require.Equal(b, tt.expected, output.Curr())
			}
		})
	}
}

func TestDispose(t *testing.T) {
	t.Parallel()

	noopState := struct{}{}

	s := redux.New(t.Context(), noopState)
	require.NotNil(t, s)

	sub := s.Subscribe()
	s.Dispose()

	_, ok := <-sub
	require.False(t, ok)
}

func TestNewStateChange(t *testing.T) {
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
				func(prev noopState, action redux.Update) noopState {
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

			s := redux.New(t.Context(), tt.state, tt.reducers...)
			defer s.Dispose()
			require.NotNil(t, s)

			sub := s.Subscribe()
			s.Dispatch(fooAction())

			output := <-sub
			defer s.CancelSubscription(sub)

			require.NotNil(t, output)
			require.Equal(t, tt.state, output.Prev())
			require.Equal(t, tt.expected, output.Curr())
		})
	}
}
