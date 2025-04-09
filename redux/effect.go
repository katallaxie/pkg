package redux

import "context"

// Hook is a function that runs any side-effect.
type Hook func()

// EffectFunc is a function that takes a context and returns an action
type EffectFunc func(ctx context.Context) Action

// BackgroundOnce is an effect that is run once an action is created once.
func BackgroundOnce[S State](ctx context.Context, store Store[S], effect EffectFunc) Hook {
	return func() {
		go store.Dispatch(effect(ctx))
	}
}
