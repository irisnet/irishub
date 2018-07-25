package p2p

import (
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	tools "github.com/irisnet/irishub/tools"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"time"
)

// Metrics contains metrics exposed by this package.
type Metrics struct {
	// Number of peers.
	Peers metrics.Gauge
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics() *Metrics {
	return &Metrics{
		Peers: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Subsystem: "p2p",
			Name:      "peers",
			Help:      "Number of peers.",
		}, []string{}),
	}
}

func (m *Metrics) Start(ctx tools.Context) {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			result := ctx.NetInfo()
			m.Peers.Set(float64(result.NPeers))
		}
	}()
}
