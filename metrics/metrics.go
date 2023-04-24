package metrics

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// DefaultRegistry is a default prometheus registry.
	DefaultRegistry = NewRegistry()

	// DefaultRegisterer is a default prometheus registerer.
	DefaultRegisterer prometheus.Registerer = DefaultRegistry

	// DefaultGatherer is a default prometheus gatherer.
	DefaultGatherer prometheus.Gatherer = DefaultRegistry
)

// Registry is a prometheus registry.
type Registry struct {
	*prometheus.Registry
}

// NewRegistry is a constructor for Registry.
func NewRegistry() *Registry {
	r := prometheus.NewRegistry()

	r.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewBuildInfoCollector(),
	)

	return &Registry{Registry: r}
}

// Handler returns a HTTP handler for this registry. Should be registered at "/metrics".
func (r *Registry) Handler() http.Handler {
	return promhttp.InstrumentMetricHandler(r, promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
}

// Collector is a type that can collect metrics.
type Collector interface {
	// Collect ...
	Collect(ch chan<- Metric)
}

// Metric is a type that can be collected.
type Metric interface {
	// Write ...
	Write(m Monitor) error
}

// Monitor is a type that can monitor metrics.
type Monitor interface {
	// Write ...
	Write(m Metric) error
}

// Probe is a type that can probe metrics.
type Probe[K, V any] interface {
	// Do ...
	Do(ctx context.Context, monitor Monitor) error

	Collector
}

// Gatherer is a type that can gather metrics.
type Gatherer interface {
	// Gather ...
	Gather(collector Collector)
}
