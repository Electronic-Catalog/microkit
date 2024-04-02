package metric

import "time"

type nopMetric struct {
}

func NewNop() *nopMetric {
	return &nopMetric{}
}

func (n *nopMetric) IncrementTotal(labelValues ...string) {
}

func (n *nopMetric) IncrementError(errorLabelValues ...string) {
}

func (n *nopMetric) ObserveResponseTime(duration time.Duration, labelValues ...string) {
}
