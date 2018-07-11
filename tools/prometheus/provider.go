package prometheus

import (
	"github.com/irisnet/irishub/tools"
	cs "github.com/irisnet/irishub/tools/prometheus/consensus"
	mempl "github.com/irisnet/irishub/tools/prometheus/mempool"
	"github.com/irisnet/irishub/tools/prometheus/p2p"
	sys "github.com/irisnet/irishub/tools/prometheus/system"
)

type Monitor struct {
	providers []MetricsProvider
	ctx       tools.Context
}

func DefaultMonitor(ctx tools.Context) *Monitor {
	var providers []MetricsProvider
	monitor := &Monitor{
		providers: providers,
		ctx:ctx,
	}
	monitor.AddMetricsProvider(cs.PrometheusMetrics()).
		AddMetricsProvider(p2p.PrometheusMetrics()).
		AddMetricsProvider(mempl.PrometheusMetrics()).
		AddMetricsProvider(sys.PrometheusMetrics())
	return monitor
}

func (m *Monitor) AddMetricsProvider(provider MetricsProvider) *Monitor {
	m.providers = append(m.providers, provider)
	return m
}

func (m *Monitor) Start() {
	for _, provider := range m.providers {
		provider.Start(m.ctx)
	}
}

type MetricsProvider interface {
	Start(ctx tools.Context)
}
