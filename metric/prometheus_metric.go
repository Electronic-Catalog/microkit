package metric

import "time"

import (
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type prometheusConfig struct {
	labels      []string
	errorLabels []string
	buckets     []float64
}

type promethuesOption func(*prometheusConfig)

func Labels(labels ...string) promethuesOption {
	return func(pc *prometheusConfig) {
		pc.labels = labels
		if len(pc.errorLabels) == 0 {
			pc.errorLabels = labels
		}
	}
}

func Buckets(buckets ...float64) promethuesOption {
	return func(pc *prometheusConfig) {
		pc.buckets = buckets
	}
}

func ErrorLabels(errorLabels ...string) promethuesOption {
	return func(pc *prometheusConfig) {
		pc.errorLabels = errorLabels
	}
}

type prometheusMetric struct {
	totalCounter          *prom.CounterVec
	errorCounter          *prom.CounterVec
	responseTimeHistogram *prom.HistogramVec
}

var _ Metric = prometheusMetric{}

func RegisterMetric(
	namespace string,
	subsystem string,
	name string,
	options ...promethuesOption,
) prometheusMetric {
	var conf prometheusConfig
	conf.buckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 15, 20, 25, 100}
	for _, po := range options {
		po(&conf)
	}

	return prometheusMetric{
		totalCounter: promauto.NewCounterVec(prom.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      name + "_total",
		}, conf.labels),
		errorCounter: promauto.NewCounterVec(prom.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      name + "_error",
		}, conf.errorLabels),
		responseTimeHistogram: promauto.NewHistogramVec(prom.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      name + "_response_time",
			Buckets:   conf.buckets,
		}, conf.labels),
	}

}

func (p prometheusMetric) IncrementTotal(labelValues ...string) {
	p.totalCounter.WithLabelValues(labelValues...).Inc()
}

func (p prometheusMetric) IncrementError(errorLabelValues ...string) {
	p.errorCounter.WithLabelValues(errorLabelValues...).Inc()
}

func (p prometheusMetric) ObserveResponseTime(duration time.Duration, labelValues ...string) {
	p.responseTimeHistogram.WithLabelValues(labelValues...).Observe(duration.Seconds())
}

func Handler() http.Handler {
	return promhttp.Handler()
}
