package upgrade

import (
	promutil "github.com/irisnet/irishub/tools/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	cfg "github.com/tendermint/tendermint/config"
	"strconv"
)

const (
	MetricsSubsystem = "module_upgrade"
	ValidatorLabel   = "validator_address"
	VersionLabel     = "version"
)

type Metrics struct {
	Signal  *stdprometheus.GaugeVec
	Version *stdprometheus.GaugeVec
}

// PrometheusMetrics returns Metrics build using Prometheus client library.
func PrometheusMetrics(config *cfg.InstrumentationConfig) *Metrics {
	if !config.Prometheus {
		return NopMetrics()
	}

	signalVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "signal",
		Help:      "validator runs the new software",
	}, []string{ValidatorLabel, VersionLabel})

	upgradeVec := stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{
		Namespace: config.Namespace,
		Subsystem: MetricsSubsystem,
		Name:      "upgrade",
		Help:      "alert the validators to install new version",
	}, []string{})

	promutil.RegisterMetrics(signalVec, upgradeVec)

	return &Metrics{
		Signal:  signalVec,
		Version: upgradeVec,
	}
}

func (m *Metrics) SetSignal(valAddr string, version uint64) {
	promutil.SafeExec(func() {
		m.Signal.With(stdprometheus.Labels{
			ValidatorLabel: valAddr,
			VersionLabel:   strconv.FormatUint(version, 10),
		}).Set(1)
	})
}

func (m *Metrics) DeleteSignal(valAddr string, version uint64) {
	promutil.SafeExec(func() {
		m.Signal.Delete(stdprometheus.Labels{
			ValidatorLabel: valAddr,
			VersionLabel:   strconv.FormatUint(version, 10),
		})
	})
}

func (m *Metrics) SetVersion(version uint64) {
	promutil.SafeExec(func() {
		m.Signal.WithLabelValues().Set(float64(version))
	})
}

func (m *Metrics) DeleteVersion() {
	promutil.SafeExec(func() {
		m.Signal.DeleteLabelValues()
	})
}

func NopMetrics() *Metrics {
	return &Metrics{
		Signal:  promutil.EmptyGaugeVec(),
		Version: promutil.EmptyGaugeVec(),
	}
}
