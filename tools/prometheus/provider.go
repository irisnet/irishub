package prometheus

import (
	"github.com/irisnet/irishub/tools/prometheus/p2p"
	cs "github.com/irisnet/irishub/tools/prometheus/consensus"
	sys "github.com/programokey/irishub/tools/prometheus/system"
	mempl "github.com/irisnet/irishub/tools/prometheus/mempool"
)

// DefaultMetricsProvider returns consensus, p2p and mempool Metrics build
// using Prometheus client library.
func DefaultMetricsProvider() (*cs.Metrics, *p2p.Metrics, *mempl.Metrics, *sys.Metrics) {
	return cs.PrometheusMetrics(), p2p.PrometheusMetrics(), mempl.PrometheusMetrics(), sys.PrometheusMetrics()
}