package stateverification

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub/client/context"
	"github.com/tendermint/tendermint/libs/log"
)

func verifyStakeState(logger log.Logger, cliCtx context.CLIContext, stakeState stake.GenesisState) error {
	logger.Info("-------------------------------------------------------------------------------")
	logger.Info("Verifying auth state")
	collectedFeesKey := []byte("collectedFees")
	res, err := cliCtx.QueryStore(collectedFeesKey, authStore)
	if err != nil {
		return err
	}
	var CollectedFees sdk.Coins
	err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(res, &CollectedFees)
	if err != nil {
		return err
	}
	logger.Info("Auth state is verified")
	return nil
}
