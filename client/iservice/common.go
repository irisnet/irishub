package iservice

import (
	iservice1 "github.com/irisnet/irishub/modules/iservice1"
	"github.com/irisnet/irishub/modules/iservice"
)

type Service1Output struct {
	iservice1.SvcDef
	Methods []iservice1.MethodProperty `json:"methods"`
}

type ServiceOutput struct {
	iservice.MsgSvcDef
	Methods []iservice.MethodProperty `json:"methods"`
}
