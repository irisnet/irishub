package stateverification

import (
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/irisnet/irishub/client/context"
	"github.com/tendermint/tendermint/libs/log"
)

func verifyMintState(logger log.Logger, cliCtx context.CLIContext, mintState mint.GenesisState) error {
	return nil
}