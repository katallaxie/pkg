package debug

import (
	"context"
	"net/http"
	"net/http/pprof"

	o "github.com/katallaxie/pkg/opts"
	"github.com/katallaxie/pkg/server"
)

var _ server.Listener = (*debug)(nil)

type debug struct {
	opts    Opts
	mux     *http.ServeMux
	handler *http.Server
}

const (
	// Addr ...
	Addr o.Opt = iota
	// Routes ...
	Routes
)

// DebugOptions are the options for the debug listener.
type Options struct {
	o.Options
}

// Opts are the options for the debug listener
type Opts interface {
	// Addr ...
	Addr() string
	// Routes ...
	Routes() map[string]http.Handler

	o.Opts
}

// NewOpts returns a new instance of the debug options.
func NewOpts(opts ...o.OptFunc) Opts {
	opts = append([]o.OptFunc{func(opts o.Opts) {
		opts.Set(Addr, ":8443")
		opts.Set(Routes, make(map[string]http.Handler))
	}}, opts...)

	oo := new(Options)
	oo.Configure(opts...)

	return oo
}

// Addr ...
func (o *Options) Addr() string {
	v, _ := o.Get(Addr)

	return v.(string)
}

// Addr ...
func (o *Options) Routes() map[string]http.Handler {
	v, _ := o.Get(Routes)

	return v.(map[string]http.Handler)
}

// New ...
func New(opts ...o.OptFunc) server.Listener {
	options := NewOpts(opts...)

	d := new(debug)
	d.opts = options

	// create the mux
	d.mux = http.NewServeMux()

	configureMux(d)

	d.handler = new(http.Server)
	d.handler.Addr = d.opts.Addr()
	d.handler.Handler = d.mux

	return d
}

// Start ...
func (d *debug) Start(ctx context.Context, ready server.ReadyFunc) func() error {
	return func() error {
		// noop, call to be ready
		ready()

		if err := d.handler.ListenAndServe(); err != nil {
			return err
		}

		return nil
	}
}

// WithStatusAddr is adding this status addr as an option.
func WithStatusAddr(addr string) o.OptFunc {
	return func(opts o.Opts) {
		opts.Set(Addr, addr)
	}
}

// WithPprof ...
func WithPprof() o.OptFunc {
	return func(opts o.Opts) {
		v, _ := opts.Get(Routes)
		vv := v.(map[string]http.Handler)

		vv["/debug/pprof/trace"] = http.HandlerFunc(pprof.Trace)
		vv["/debug/pprof/"] = http.HandlerFunc(pprof.Index)
		vv["/debug/pprof/cmdline"] = http.HandlerFunc(pprof.Cmdline)
		vv["/debug/pprof/profile"] = http.HandlerFunc(pprof.Profile)
		vv["/debug/pprof/symbol"] = http.HandlerFunc(pprof.Symbol)

		opts.Set(Routes, vv)
	}
}

// WithPrometheusHandler is adding this prometheus http handler as an option.
func WithPrometheusHandler(handler http.Handler) o.OptFunc {
	return func(opts o.Opts) {
		v, _ := opts.Get(Routes)
		vv := v.(map[string]http.Handler)

		vv["/debug/metrics"] = handler

		opts.Set(Routes, vv)
	}
}

func configureMux(d *debug) {
	for route, handler := range d.opts.Routes() {
		d.mux.Handle(route, handler)
	}
}
