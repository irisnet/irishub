package iservice

import (
	"github.com/irisnet/irishub/modules/iservice"
)

type ServiceOutput struct {
	iservice.MsgSvcDef
	Methods []iservice.MethodProperty `json:"methods"`
}
