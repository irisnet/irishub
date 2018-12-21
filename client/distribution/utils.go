package distribution

import (
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/distribution"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/utils"
	stakeTypes "github.com/irisnet/irishub/modules/stake/types"
)

// distribution info for a particular validator
type ValidatorDistInfoOutput struct {
	OperatorAddr            sdk.ValAddress          `json:"operator_addr"`
	FeePoolWithdrawalHeight int64                   `json:"fee_pool_withdrawal_height"`
	DelAccum                distribution.TotalAccum `json:"del_accum"`
	DelPool                 string                  `json:"del_pool"`
	ValCommission           string                  `json:"val_commission"`
}

func ConvertToValidatorDistInfoOutput(cliCtx context.CLIContext, vdi distribution.ValidatorDistInfo) ValidatorDistInfoOutput {
	exRate := utils.ExRateFromStakeTokenToMainUnit(cliCtx)
	delPool := utils.ConvertDecToRat(vdi.DelPool.AmountOf(stakeTypes.StakeDenom)).Mul(exRate).FloatString() + stakeTypes.StakeDenomName
	valCommission := utils.ConvertDecToRat(vdi.ValCommission.AmountOf(stakeTypes.StakeDenom)).Mul(exRate).FloatString() + stakeTypes.StakeDenomName
	return ValidatorDistInfoOutput{
		OperatorAddr:            vdi.OperatorAddr,
		FeePoolWithdrawalHeight: vdi.FeePoolWithdrawalHeight,
		DelAccum:                vdi.DelAccum,
		DelPool:                 delPool,
		ValCommission:           valCommission,
	}
}
