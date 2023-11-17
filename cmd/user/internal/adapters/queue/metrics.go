package queue

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	labelTopic = "topic"
)

// Metrics contains general metrics for DAL methods.
type Metrics struct {
	errorsTotal *prometheus.CounterVec
}

// NewMetrics registers and returns common queue metrics.
func NewMetrics(reg *prometheus.Registry, namespace, subsystem string, topics []string) (metric Metrics) {
	metric.errorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "errors_total",
			Help:      "Amount of subscription errors.",
		},
		[]string{labelTopic},
	)
	reg.MustRegister(metric.errorsTotal)

	for _, topic := range topics {
		l := prometheus.Labels{
			labelTopic: topic,
		}
		metric.errorsTotal.With(l)
	}

	return metric
}
