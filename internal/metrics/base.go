package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Type aliases.
type (
	GaugeOpts     prometheus.GaugeOpts
	CounterOpts   prometheus.CounterOpts
	SummaryOpts   prometheus.SummaryOpts
	HistogramOpts prometheus.HistogramOpts

	Gauge     prometheus.Gauge
	Counter   prometheus.Counter
	Summary   prometheus.Summary
	Histogram prometheus.Histogram

	GaugeVec     = prometheus.GaugeVec
	CounterVec   = prometheus.CounterVec
	SummaryVec   = prometheus.SummaryVec
	HistogramVec = prometheus.HistogramVec
)

// Default prometheus implementations for Registerer and Gatherer.
var (
	DefaultRegisterer = prometheus.DefaultRegisterer
)

// Register registers collector in the default metrics registerer.
func Register(c prometheus.Collector) error {
	return DefaultRegisterer.Register(c)
}

// Metrics contains registered collectors.
type Metrics struct {
	Namespace   string
	Subsystem   string
	Registerer  prometheus.Registerer
	ConstLabels prometheus.Labels

	mu      sync.Mutex
	metrics []prometheus.Collector
}

// NewMetrics returns a Metricer instance.
func NewMetrics() *Metrics {
	return &Metrics{
		Registerer: DefaultRegisterer,
	}
}

// WithNamespace sets namespace for metrics.
func (m *Metrics) WithNamespace(namespace string) *Metrics {
	m.Namespace = namespace

	return m
}

// WithSubsystem sets subsystem for metrics.
func (m *Metrics) WithSubsystem(subsystem string) *Metrics {
	m.Subsystem = subsystem

	return m
}

func (m *Metrics) WithConstLabels(labels prometheus.Labels) *Metrics {
	m.ConstLabels = labels

	return m
}

// NewGauge creates a new Gauge based on the provided GaugeOpts.
func (m *Metrics) NewGauge(o GaugeOpts) Gauge {
	c := m.register().NewGauge(
		prometheus.GaugeOpts{
			Namespace:   m.Namespace,
			Subsystem:   m.Subsystem,
			ConstLabels: m.ConstLabels,
			Name:        o.Name,
			Help:        o.Help,
		},
	)
	m.collect(c)

	return c
}

// NewCounter creates a new Counter based on the provided CounterOpts.
func (m *Metrics) NewCounter(o CounterOpts) Counter {
	c := m.register().NewCounter(
		prometheus.CounterOpts{
			Namespace:   m.Namespace,
			Subsystem:   m.Subsystem,
			ConstLabels: m.ConstLabels,
			Name:        o.Name,
			Help:        o.Help,
		},
	)
	m.collect(c)

	return c
}

// NewHistogram creates a new Histogram based on the provided HistogramOpts. It
// panics if the buckets in HistogramOpts are not in strictly increasing order.
func (m *Metrics) NewHistogram(o HistogramOpts) Histogram {
	c := m.register().NewHistogram(
		prometheus.HistogramOpts{
			Namespace:   m.Namespace,
			Subsystem:   m.Subsystem,
			ConstLabels: m.ConstLabels,
			Name:        o.Name,
			Help:        o.Help,
			Buckets:     o.Buckets,
		},
	)
	m.collect(c)

	return c
}

// NewGaugeVec creates a new GaugeVec based on the provided GaugeOpts and
// partitioned by the given label names.
func (m *Metrics) NewGaugeVec(o GaugeOpts, labelNames []string) *GaugeVec {
	c := m.register().NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace:   m.Namespace,
			Subsystem:   m.Subsystem,
			ConstLabels: m.ConstLabels,
			Name:        o.Name,
			Help:        o.Help,
		},
		labelNames,
	)
	m.collect(c)

	return c
}

// NewCounterVec creates a new CounterVec based on the provided CounterOpts and
// partitioned by the given label names.
func (m *Metrics) NewCounterVec(o CounterOpts, labelNames []string) *CounterVec {
	c := m.register().NewCounterVec(
		prometheus.CounterOpts{
			Namespace:   m.Namespace,
			Subsystem:   m.Subsystem,
			ConstLabels: m.ConstLabels,
			Name:        o.Name,
			Help:        o.Help,
		},
		labelNames,
	)
	m.collect(c)

	return c
}

// NewHistogramVec creates a new HistogramVec based on the provided HistogramOpts and
// partitioned by the given label names.
func (m *Metrics) NewHistogramVec(o HistogramOpts, labelNames []string) *HistogramVec {
	c := m.register().NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace:   m.Namespace,
			Subsystem:   m.Subsystem,
			ConstLabels: m.ConstLabels,
			Name:        o.Name,
			Help:        o.Help,
			Buckets:     o.Buckets,
		},
		labelNames,
	)
	m.collect(c)

	return c
}

func (m *Metrics) register() promauto.Factory {
	return promauto.With(m.Registerer)
}

func (m *Metrics) collect(c prometheus.Collector) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metrics = append(m.metrics, c)
}
