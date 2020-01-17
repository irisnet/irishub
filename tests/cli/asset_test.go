package cli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/irisnet/irishub/app/v3/asset"
	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestIrisCLIToken(t *testing.T) {
	t.Parallel()
	chainID, servAddr, port, irisHome, iriscliHome, p2pAddr := initializeFixtures(t)

	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v --output=json", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v --p2p.laddr=%v", irisHome, servAddr, p2pAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))
	barAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show bar --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	//barAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	barCoin := sdk.NewCoin("btc-min", sdk.ZeroInt())
	require.Equal(t, "50iris", fooCoin)

	symbol := "btc"
	name := "Bitcoin"
	minUnit := "satoshi"
	initialSupply := 10000000000
	maxSupply := 20000000000
	decimal := 18
	mintable := true

	tokenID := asset.GetTokenID(symbol)

	// issue a token
	issueCmd := fmt.Sprintf("iriscli asset token issue %v", flags)
	issueCmd += fmt.Sprintf(" --from=%s", "foo")
	issueCmd += fmt.Sprintf(" --symbol=%s", symbol)
	issueCmd += fmt.Sprintf(" --name=%s", name)
	issueCmd += fmt.Sprintf(" --min-unit=%s", minUnit)
	issueCmd += fmt.Sprintf(" --initial-supply=%d", initialSupply)
	issueCmd += fmt.Sprintf(" --max-supply=%d", maxSupply)
	issueCmd += fmt.Sprintf(" --scale=%d", decimal)
	issueCmd += fmt.Sprintf(" --mintable=%v", mintable)
	issueCmd += fmt.Sprintf(" --fee=%s", "0.4iris")

	require.True(t, executeWrite(t, issueCmd, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	amt := getAmountFromCoinStr(fooCoin)

	// 30iris is used to issue tokens
	if !(amt > 19 && amt < 20) {
		t.Error("Test Failed: (19, 20) expected, received:", amt)
	}

	query := fmt.Sprintf("--token-id=%s ", tokenID)
	token := executeGetToken(t, fmt.Sprintf("iriscli asset token tokens %s %v", query, flags))
	require.Equal(t, strings.ToLower(strings.TrimSpace(symbol)), token.Symbol)
	require.Equal(t, strings.ToLower(strings.TrimSpace(minUnit)), token.MinUnit)
	require.Equal(t, strings.TrimSpace(name), token.Name)
	require.Equal(t, sdk.NewIntWithDecimal(int64(initialSupply), decimal), token.InitialSupply)
	require.Equal(t, sdk.NewIntWithDecimal(int64(maxSupply), decimal), token.MaxSupply)
	require.Equal(t, uint8(decimal), token.Scale)
	require.Equal(t, mintable, token.Mintable)

	// edit a token
	name = "BTC_Token"
	maxSupply = 30000000000
	mintable = true

	editCmd := fmt.Sprintf("iriscli asset token edit %s", tokenID)
	editCmd += fmt.Sprintf(" --from=%s", "foo")
	editCmd += fmt.Sprintf(" --name=%s", name)
	editCmd += fmt.Sprintf(" --max-supply=%d", maxSupply)
	editCmd += fmt.Sprintf(" --mintable=%v", mintable)
	editCmd += fmt.Sprintf(" --fee=%s %v", "0.4iris", flags)
	require.True(t, executeWrite(t, editCmd, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	query = fmt.Sprintf("--token-id=%s ", tokenID)
	token = executeGetToken(t, fmt.Sprintf("iriscli asset token tokens %s %v", query, flags))

	require.Equal(t, name, token.Name)
	require.Equal(t, sdk.NewIntWithDecimal(int64(maxSupply), decimal), token.MaxSupply)
	require.Equal(t, mintable, token.Mintable)

	//mint a token
	amount := 1000
	mintCmd := fmt.Sprintf("iriscli asset token mint %s", tokenID)
	mintCmd += fmt.Sprintf(" --from=%s", "foo")
	mintCmd += fmt.Sprintf(" --to=%s", barAddr.String())
	mintCmd += fmt.Sprintf(" --amount=%d", amount)
	mintCmd += fmt.Sprintf(" --fee=%s %v", "0.4iris", flags)
	require.True(t, executeWrite(t, mintCmd, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	barAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", barAddr, flags))
	balanceExp := barCoin.Add(sdk.NewCoin("btc-min", sdk.NewIntWithDecimal(1000, decimal)))
	require.Equal(t, balanceExp.String(), barAcc.GetCoins().String())

	//transfer a token
	transferCmd := fmt.Sprintf("iriscli asset token transfer %s", tokenID)
	transferCmd += fmt.Sprintf(" --from=%s", "foo")
	transferCmd += fmt.Sprintf(" --to=%s", barAddr.String())
	transferCmd += fmt.Sprintf(" --fee=%s %v", "0.4iris", flags)
	require.True(t, executeWrite(t, transferCmd, sdk.DefaultKeyPass))
	tests.WaitForNextNBlocksTM(2, port)

	query = fmt.Sprintf("--owner=%s ", barAddr.String())
	token = executeGetToken(t, fmt.Sprintf("iriscli asset token tokens %s %v", query, flags))
	require.Equal(t, barAddr.String(), token.Owner.String())
}
