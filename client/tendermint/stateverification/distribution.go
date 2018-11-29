package stateverification

import (
	"github.com/irisnet/irishub/client/context"
	"github.com/tendermint/tendermint/libs/log"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
)

func verifyDistrState(logger log.Logger, cliCtx context.CLIContext, distrState distr.GenesisState) error {
	return nil
}
