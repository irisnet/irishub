package iservice

import (
	"github.com/irisnet/irishub/modules/iservice"
)

type ServiceOutput struct {
	iservice.SvcDef
	Methods []iservice.MethodProperty `json:"methods"`
}
