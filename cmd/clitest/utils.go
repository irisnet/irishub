package clitest

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client/bank"
	"github.com/irisnet/irishub/client/context"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/types"
	"strconv"
	"strings"
	"testing"
	"fmt"
	"os"
	"github.com/irisnet/irishub/app"
)

func convertToIrisBaseAccount(t *testing.T, acc *bank.BaseAccount) string {
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	cliCtx := context.NewCLIContext().
		WithCodec(cdc)

	coinstr := acc.Coins[0]
	for i := 1; i < len(acc.Coins); i++ {
		coinstr += ("," + acc.Coins[i])
	}

	coins, err := cliCtx.ConvertCoinToMainUnit(coinstr)
	require.NoError(t, err, "coins %v, err %v", coinstr, err)

	return coins[0]
}

func getAmuntFromCoinStr(t *testing.T, coinStr string) float64 {
	index := strings.Index(coinStr, "iris")
	if index <= 0 {
		return -1
	}

	numStr := coinStr[:index]
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return -1
	}

	return num
}

func modifyGenesisFile(t *testing.T, irisHome string) error {
	genesisFilePath := fmt.Sprintf("%s%sconfig%sgenesis.json", irisHome, string(os.PathSeparator), string(os.PathSeparator))

	genesisDoc, err := types.GenesisDocFromFile(genesisFilePath)
	if err != nil {
		return err
	}

	var genesisState app.GenesisState

	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	cliCtx := context.NewCLIContext().
		WithCodec(cdc)

	err = cdc.UnmarshalJSON(genesisDoc.AppState, &genesisState)
	if err != nil {
		return err
	}

	coin, err := cliCtx.ParseCoin("1000000000000iris")
	if err != nil {
		return err
	}

	genesisState.Accounts[0].Coins[0] = coin
	bz, err := cdc.MarshalJSON(genesisState)
	if err != nil {
		return err
	}

	genesisDoc.AppState = bz
	return genesisDoc.SaveAs(genesisFilePath)
}
