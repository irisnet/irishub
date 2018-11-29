package stateverification

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/context"
	"github.com/tendermint/tendermint/libs/log"
)

func verifyAccountsState(logger log.Logger, cliCtx context.CLIContext, accountsState []app.GenesisAccount) error {
	logger.Info("-------------------------------------------------------------------------------")
	logger.Info("Verifying account state")
	decoder := authcmd.GetAccountDecoder(cliCtx.Codec)
	for _, acc := range accountsState {
		res, err := cliCtx.QueryStore(auth.AddressStoreKey(acc.Address), accountStore)
		if err != nil {
			return err
		}
		if len(res) == 0 {
			return fmt.Errorf("account %s doesn't exist", acc.Address.String())
		}
		account, err := decoder(res)
		if err != nil {
			return fmt.Errorf("account %s: failed to decode account info", acc.Address.String())
		}
		if !acc.Coins.IsEqual(account.GetCoins()) {
			return fmt.Errorf("account %s: token amount doesn't match, expect %s, got %s",
				acc.Address.String(), acc.Coins.String(), account.GetCoins().String())
		}
		if acc.AccountNumber != account.GetAccountNumber() {
			return fmt.Errorf("account %s: account number doesn't match, expect %d, got %d",
				acc.Address.String(), acc.AccountNumber, account.GetAccountNumber())
		}
		if acc.Sequence != account.GetSequence() {
			return fmt.Errorf("account %s: account sequence doesn't match, expect %d, got %d",
				acc.Address.String(), acc.Sequence, account.GetSequence())
		}
	}
	logger.Info("Account state is verified")
	return nil
}
