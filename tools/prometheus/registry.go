package prometheus

import (
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func RegisterMetrics(cs ...stdprometheus.Collector) {
	for _, c := range cs {
		if err := stdprometheus.Register(c); err != nil {
			CheckRegisterError(err)
		}
	}
}

func CheckRegisterError(err error) {
	if register, ok := err.(stdprometheus.AlreadyRegisteredError); ok {
		stdprometheus.Unregister(register.ExistingCollector)
		stdprometheus.MustRegister(register.NewCollector)
	} else {
		panic(err)
	}
}

func SafeExec(fn func(), enable bool) {
	defer func() {
		if r := recover(); r != nil {
			//TODO
		}
	}()
	if enable {
		fn()
	}
}

func EmptyGaugeVec() *stdprometheus.GaugeVec {
	return stdprometheus.NewGaugeVec(stdprometheus.GaugeOpts{}, []string{})
}

func EmptyCounterVec() *stdprometheus.CounterVec {
	return stdprometheus.NewCounterVec(stdprometheus.CounterOpts{}, []string{})
}
