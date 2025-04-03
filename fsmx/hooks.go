package fsmx

import "context"

// Effect is a function that takes a context and returns an action
func Effect[S State](store Store[S], hook EffectFunc[S]) error {
	action, err := hook(context.Background())
	if err != nil {
		return err
	}

	store.Dispatch(action)

	return nil
}

// EffectFunc is a function that takes a context and returns an action
type EffectFunc[S State] func(ctx context.Context) (Action, error)
