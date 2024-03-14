package rpc

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/irisnet/irishub/v3/app/keepers"
)

var overrideModules = map[string]overrideHandler{
	authtypes.ModuleName: overrideAuthServices,
}

type overrideHandler func(configurator module.Configurator, appKeepers keepers.AppKeepers)

// RegisterService allows a module to register services.
func RegisterService(mod module.AppModule, configurator module.Configurator, appKeepers keepers.AppKeepers) {
	handler, has := overrideModules[mod.Name()]
	if has {
		handler(configurator, appKeepers)
		return
	}

	m, ok := mod.(module.HasServices)
	if ok {
		m.RegisterServices(configurator)
	}
}
