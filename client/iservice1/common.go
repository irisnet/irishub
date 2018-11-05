package iservice

import (
	"github.com/irisnet/irishub/modules/iservice1"
)

type ServiceOutput struct {
	iservice.SvcDef
	Methods []iservice.MethodProperty `json:"methods"`
}
