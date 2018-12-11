package keeper


import "github.com/irisnet/irishub/types/common"

type UpgradeConfig struct {
	ProposalID   uint64
	Definition   common.ProtocolDefinition
}