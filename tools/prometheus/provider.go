package prometheus

import (
	"github.com/irisnet/irishub/tools/prometheus/p2p"
	cs "github.com/irisnet/irishub/tools/prometheus/consensus"
	mempl "github.com/irisnet/irishub/tools/prometheus/mempool"
)

// DefaultMetricsProvider returns consensus, p2p and mempool Metrics build
// using Prometheus client library.
func DefaultMetricsProvider() (*cs.Metrics, *p2p.Metrics, *mempl.Metrics) {
	return cs.PrometheusMetrics(), p2p.PrometheusMetrics(), mempl.PrometheusMetrics()
}