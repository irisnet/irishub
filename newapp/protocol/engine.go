package protocol

import (
	"github.com/irisnet/irishub/types/common"
	"fmt"
)

type Protocol interface {
	GetDefinition() common.ProtocolDefinition
	Load()
	Init()
}

type ProtocolBase struct {
	definition	common.ProtocolDefinition
//	engine 		*ProtocolEngine
}

func (pb ProtocolBase) GetDefinition() common.ProtocolDefinition {
	return pb.definition
}
/*
func (pb *ProtocolBase) GetEngine() *ProtocolEngine {
	return pb.engine
}
*/
type ProtocolEngine struct {
	protocols	map[uint]Protocol
	next		uint
	current		uint
//	app			*app.IrisApp
}

func NewProtocolEngine() ProtocolEngine {
	engine := ProtocolEngine{
		make(map[uint]Protocol),
		0,
		0,
//		irisApp,
	}
	return engine
}

// To be used for Protocol with version > 0
func (pe *ProtocolEngine) Activate(version uint) bool {
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

func (pe *ProtocolEngine) GetByVersion(v uint) (Protocol, bool) {
	p, flag := pe.protocols[v]
	return p, flag
}