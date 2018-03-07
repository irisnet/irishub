package iservice

import (
	"github.com/cosmos/cosmos-sdk/state"
	"github.com/tendermint/go-wire"
)

// nolint
var (
	ParamKey                   = []byte{0x01} // key for global parameters relating to iservice
	ServiceDefinitionKeyPrefix = []byte{0x02} // prefix for each key to a service definition
)

//---------------------------------------------------------------------

// load/save the global iservice params
func loadParams(store state.SimpleDB) (params Params) {
	b := store.Get(ParamKey)
	if b == nil {
		return defaultParams()
	}

	err := wire.ReadBinaryBytes(b, &params)
	if err != nil {
		panic(err) // This error should never occure big problem if does
	}

	return
}
func saveParams(store state.SimpleDB, params Params) {
	b := wire.BinaryBytes(params)
	store.Set(ParamKey, b)
}

func saveService(store state.SimpleDB, service *ServiceDefinition) {

	b := wire.BinaryBytes(*service)
	store.Set(GetServiceDefinitionKey(service.Name), b)
}

func loadService(store state.SimpleDB, name string) *ServiceDefinition {
	serviceBytes := store.Get(GetServiceDefinitionKey(name))
	if serviceBytes == nil {
		return nil
	}

	service := new(ServiceDefinition)
	err := wire.ReadBinaryBytes(serviceBytes, service)
	if err != nil {
		panic(err)
	}
	return service
}

// GetServiceDefinitionKey - get the key for a service definition
func GetServiceDefinitionKey(name string) []byte {
	return append(ServiceDefinitionKeyPrefix, wire.BinaryBytes(&name)...)
}
