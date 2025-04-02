package fsmx

import (
	"sync"

	"github.com/katallaxie/pkg/slices"
)

// Subscription is the type of the subscription of the FSM.
type Subscription interface {
	// ID gets the ID of the subscription.
	ID() int
	// Listen listens to the subscription.
	Listen() <-chan State
	// Cancel cancels the subscription.
	Cancel()
}

var _ Subscription = (*subscription)(nil)

type subscription struct {
	id      int
	store   *store
	subOnce sync.Once
}

// Cancel cancels the subscription.
func (s *subscription) Cancel() {
	s.store.Lock()
	defer s.store.Unlock()

	for i, sub := range s.store.subscribers {
		if i == s.id {
			slices.Delete(i, s.store.subscribers...)
			close(sub)
		}
	}
}

// Listen listens to the subscription.
func (s *subscription) Listen() <-chan State {
	l := make(chan State)

	s.subOnce.Do(func() {
		s.store.Lock()
		defer s.store.Unlock()

		s.store.subscribers = append(s.store.subscribers, l)
	})

	return l
}

// ID gets the ID of the subscription.
func (s *subscription) ID() int {
	return s.id
}

// ActionPayload is the type of the payload of the action.
type ActionPayload interface{}

// Action is the type of the action of the FSM.
type ActionType int

// State is the type of the state of the FSM.
type State interface{}

// Actionable is the interface that wraps the basic Action methods.
type Action interface {
	// Paylog gets the payload of the action
	Payload(p ...ActionPayload) ActionPayload
	// Type gets the type of the action
	Type(a ...ActionType) ActionType
}

// NewAction creates a new action.
func NewAction(actionType ActionType, payload ActionPayload) Action {
	return &action{
		actionType: actionType,
		payload:    payload,
	}
}

type action struct {
	actionType ActionType
	payload    ActionPayload
}

// Payload gets the payload of the action
func (a *action) Payload(payloads ...ActionPayload) ActionPayload {
	if slices.Len(payloads) > 0 {
		a.payload = slices.First(payloads...)
	}

	return a.payload
}

// Type gets the type of the action
func (a *action) Type(types ...ActionType) ActionType {
	if slices.Len(types) > 0 {
		a.actionType = slices.First(types...)
	}

	return a.actionType
}

// Reducer is the type of the reducer of the FSM.
type Reducer func(prev State, action Action) State

// Store is the type of the store of the FSM.
type Store interface {
	// Dispatch dispatches an event to the FSM.
	Dispatch(actions ...Action)
	// Subscribe subscribes to the store.
	Subscribe() Subscription
	// State gets the current state of the store.
	State(s ...State) State
	// Drain drains the store.
	Drain()
}

type store struct {
	id          int
	state       State
	reducers    []Reducer
	subscribers []chan<- State

	sync.RWMutex
}

// New creates a new store.
func New(initialState State, reducers ...Reducer) Store {
	s := new(store)
	s.state = initialState
	s.reducers = slices.Append(reducers, s.reducers...)

	return s
}

// Dispatch dispatches an event to the FSM.
func (s *store) Dispatch(actions ...Action) {
	s.Lock()
	defer s.Unlock()

	for _, action := range actions {
		for _, reducer := range s.reducers {
			s.state = reducer(s.state, action)

			for _, sub := range s.subscribers {
				go func(sub chan<- State) { // background
					sub <- s.state
				}(sub)
			}
		}
	}
}

// State gets the current state of the store.
func (s *store) State(states ...State) State {
	s.Lock()
	defer s.Unlock()

	if slices.Len(states) > 0 {
		s.state = slices.First(states...)
	}

	return s.state
}

// Subscribe subscribes to the store.
func (s *store) Subscribe() Subscription {
	s.Lock()
	defer s.Unlock()

	sub := new(subscription)
	s.id++
	sub.id = s.id
	sub.store = s

	return sub
}

// Drain drains the store.
func (s *store) Drain() {
	s.Lock()
	defer s.Unlock()

	for _, sub := range s.subscribers {
		close(sub)
	}
}
