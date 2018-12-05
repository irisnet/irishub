package protocol

import "github.com/irisnet/irishub/types/common"

type ProtocolVersion0 struct {
	*ProtocolBase
}

func NewProtocolVersion0(engine *ProtocolEngine) ProtocolVersion0 {
	base := ProtocolBase{
		definition: common.ProtocolDefinition{
			0,
			"https://github.com/irisnet/irishub/releases/tag/v0.7.0",
			1,
		},
//		engine: engine,
	}
	p0 := ProtocolVersion0{
		&base,
	}
	return p0
}

// load the configuration of this Protocol
func (p *ProtocolVersion0) Load() {
	ConfigKeepers()
	ConfigHooks()
	ConfigRouters()
	ConfigStores()
}

// create all Keepers
func ConfigKeepers() {
}

// wire all Keepers together in a loosely coupled manner
func ConfigHooks() {
}

// configure all Routers
func ConfigRouters() {

}

// configure all Stores
func ConfigStores() {

}

// initialize State for this Protocol
func (p *ProtocolVersion0) Init() {

}

