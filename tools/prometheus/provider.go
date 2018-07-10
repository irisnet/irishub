package prometheus

import (
	cs "github.com/irisnet/irishub/tools/prometheus/consensus"
	mempl "github.com/irisnet/irishub/tools/prometheus/mempool"
	"github.com/irisnet/irishub/tools/prometheus/p2p"
	sys "github.com/irisnet/irishub/tools/prometheus/system"
)

// DefaultMetricsProvider returns consensus, p2p and mempool Metrics build
// using Prometheus client library.
func DefaultMetricsProvider() (*cs.Metrics, *p2p.Metrics, *mempl.Metrics, *sys.Metrics) {
	return cs.PrometheusMetrics(), p2p.PrometheusMetrics(), mempl.PrometheusMetrics(), sys.PrometheusMetrics()
}
