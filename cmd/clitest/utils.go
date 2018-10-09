package clitest

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/client/bank"
	"github.com/irisnet/irishub/client/context"
	"github.com/stretchr/testify/require"
	"strconv"
	"strings"
	"testing"
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

func getAmuntFromCoinStr(t *testing.T, coinStr string) int {
	index := strings.Index(coinStr, "iris")
	if index <= 0 {
		return -1
	}

	numStr := coinStr[:index]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return -1
	}

	return num
}
