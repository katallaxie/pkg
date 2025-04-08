package redux

import "context"

// Effect is a function that takes a context and returns an action
func Effect[S State](store Store[S], hook EffectFunc) error {
	action, err := hook(context.Background())
	if err != nil {
		return err
	}

	store.Dispatch(action)

	return nil
}

// EffectFunc is a function that takes a context and returns an action
type EffectFunc func(ctx context.Context) (Action, error)
