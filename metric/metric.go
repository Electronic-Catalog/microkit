package metric

import (
	"time"
)

type Metric interface {
	IncrementTotal(labelValues ...string)
	IncrementError(errorLabelValues ...string)
	ObserveResponseTime(duration time.Duration, labelValues ...string)
}
