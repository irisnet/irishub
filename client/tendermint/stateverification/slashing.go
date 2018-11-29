package stateverification

import (
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/irisnet/irishub/client/context"
)

func verifySlashingState(logger log.Logger, cliCtx context.CLIContext, slashingState slashing.GenesisState) error {
	return nil
}