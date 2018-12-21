package protocol

import (
	protocolKeeper "github.com/irisnet/irishub/app/protocol/keeper"
	"github.com/irisnet/irishub/codec"
	sdk "github.com/irisnet/irishub/types"
)

/*
func (pb *ProtocolBase) GetEngine() *ProtocolEngine {
	return pb.engine
}
*/
type ProtocolEngine struct {
	protocols map[uint64]Protocol
	current   uint64
	pk        protocolKeeper.Keeper
}

func NewProtocolEngine(cdc *codec.Codec) ProtocolEngine {
	engine := ProtocolEngine{
		make(map[uint64]Protocol),
		0,
		protocolKeeper.NewKeeper(cdc, KeyProtocol),
		//		irisApp,
	}
	return engine
}

func (pe *ProtocolEngine) LoadCurrentProtocol(kvStore sdk.KVStore) {
	//find the current version From DB( EngineKeeper?)
	current := pe.pk.GetCurrentProtocolVersionByStore(kvStore)
	p, flag := pe.protocols[current]
	if flag == true {
		p.Load(pe.pk)
		pe.current = current
	}
}

func (pe *ProtocolEngine) GetCurrentProtocolVersionByStore(kvStore sdk.KVStore) uint64 {
	return pe.pk.GetCurrentProtocolVersionByStore(kvStore)
}

func (pe *ProtocolEngine) GetUpgradeConfigByStore(kvStore sdk.KVStore) (protocolKeeper.UpgradeConfig, bool) {
	return pe.pk.GetUpgradeConfigByStore(kvStore)
}

// To be used for Protocol with version > 0
func (pe *ProtocolEngine) Activate(version uint64) bool {
	p, flag := pe.protocols[version]
	if flag == true {
		p.Load(pe.pk)
		p.Init()
		pe.current = version
	}
	return flag
}

func (pe *ProtocolEngine) GetCurrentProtocol() Protocol {
	return pe.protocols[pe.current]
}

func (pe *ProtocolEngine) GetCurrentVersion() uint64 {
	return pe.current
}

func (pe *ProtocolEngine) Add(p Protocol) Protocol {
	pe.protocols[p.GetDefinition().GetVersion()] = p
	return p
}

func (pe *ProtocolEngine) GetByVersion(v uint64) (Protocol, bool) {
	p, flag := pe.protocols[v]
	return p, flag
}

func (pe *ProtocolEngine) GetKVStoreKeys() []*sdk.KVStoreKey {
	return []*sdk.KVStoreKey{
		KeyMain,
		KeyProtocol,
		KeyAccount,
		KeyStake,
		KeyMint,
		KeyDistr,
		KeySlashing,
		KeyGov,
		KeyRecord,
		KeyFeeCollection,
		KeyParams,
		KeyUpgrade,
		KeyService,
		KeyGuardian}
}

func (pe *ProtocolEngine) GetTransientStoreKeys() []*sdk.TransientStoreKey {
	return []*sdk.TransientStoreKey{
		TkeyStake,
		TkeyDistr,
		TkeyParams}
}

func (pe *ProtocolEngine) GetKeyMain() *sdk.KVStoreKey {
	return KeyMain
}
