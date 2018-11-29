package stateverification

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/irisnet/irishub/client/context"
	"github.com/tendermint/tendermint/libs/log"
)

func verifyAuthState(logger log.Logger, cliCtx context.CLIContext, authState auth.GenesisState) error {
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
	if !authState.CollectedFees.IsEqual(CollectedFees) {
		return fmt.Errorf("CollectedFees doesn't match")
	}
	logger.Info("Auth state is verified")
	return nil
}
