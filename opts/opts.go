package opts

import (
	"fmt"
	"sync"

	"github.com/katallaxie/pkg/utils"
)

// ErrNotFound signals that this option is not set.
var ErrNotFound = fmt.Errorf("options: option not found")

// Opt is the identifier for the option.
type Opt int

// Opts ...
type Opts[K comparable, V any] interface {
	// Get ...
	Get(Opt) (V, error)
	// Set ...
	Set(Opt, V)
	// Configure ...
	Configure(...OptFunc[K, V])
}

// OptFunc is an option
type OptFunc[K comparable, V any] func(Opts[K, V])

// Options is default options structure.
type Options[K comparable, V any] struct {
	opts map[Opt]V

	sync.RWMutex
}

// New returns a new instance of the options.
func New[K comparable, V any](opts ...OptFunc[K, V]) Opts[K, V] {
	o := new(Options[K, V])
	o.Configure(opts...)

	return o
}

// Get is returning the value of the option.
func (o *Options[K, V]) Get(opt Opt) (V, error) {
	o.RLock()
	defer o.RUnlock()

	v, ok := o.opts[opt]
	if !ok {
		return utils.Zero[V](), ErrNotFound
	}

	return v, nil
}

// Set is setting the value of the option.
func (o *Options[K, V]) Set(opt Opt, v V) {
	o.Lock()
	defer o.Unlock()

	o.opts[opt] = v
}

// Configure os configuring the options.
func (o *Options[K, V]) Configure(opts ...OptFunc[K, V]) {
	if o.opts == nil {
		o.opts = make(map[Opt]V)
	}

	for _, opt := range opts {
		opt(o)
	}
}
