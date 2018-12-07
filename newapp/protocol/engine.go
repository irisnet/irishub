package protocol

import (
	"fmt"
	sdk "github.com/irisnet/irishub/types"
)

/*
func (pb *ProtocolBase) GetEngine() *ProtocolEngine {
	return pb.engine
}
*/
type ProtocolEngine struct {
	protocols map[uint64]Protocol
	next      uint64
	current   uint64
	//	app			*app.IrisApp
}

func NewProtocolEngine() ProtocolEngine {
	engine := ProtocolEngine{
		make(map[uint64]Protocol),
		0,
		0,
		//		irisApp,
	}
	return engine
}

func (pe *ProtocolEngine) LoadCurrentProtocol() {
	//find the current version From DB( EngineKeeper?)
	current := uint64(0)
	next := uint64(1)
	p, flag := pe.protocols[current]
	if flag == true {
		p.Load()
		pe.current = current
		pe.next = next
	}
}

// To be used for Protocol with version > 0
func (pe *ProtocolEngine) Activate(version uint64) bool {
	p, flag := pe.protocols[version]
	if flag == true {
		p.Load()
		p.Init()
		pe.current = version
	}
	return flag
}

func (pe *ProtocolEngine) GetCurrent() Protocol {
	return pe.protocols[pe.current]
}

func (pe *ProtocolEngine) Add(p Protocol) Protocol {
	if pe.next != p.GetDefinition().GetVersion() {
		panic(fmt.Errorf("The next protocol expected is %d", pe.next))
	}
	pe.protocols[pe.next] = p
	pe.next++
	return p
}

func (pe *ProtocolEngine) GetByVersion(v uint64) (Protocol, bool) {
	p, flag := pe.protocols[v]
	return p, flag
}

func (pe *ProtocolEngine) GetKVStoreKeys() []*sdk.KVStoreKey {
	return []*sdk.KVStoreKey{
		KeyMain,
		KeyAccount,
		KeyStake,
		KeyMint,
		KeyDistr,
		KeySlashing,
		KeyGov,
		KeyRecord,
		KeyFeeCollection,
		KeyParams,
		//keyUpgrade,
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
