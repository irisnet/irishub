package service

import (
	"github.com/irisnet/irishub/modules/service"
)

type ServiceOutput struct {
	service.SvcDef
	Methods []service.MethodProperty `json:"methods"`
}
