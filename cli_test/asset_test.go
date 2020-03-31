package clitest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/app"
	token "github.com/irisnet/irishub/modules/asset/01-token"
)

func TestIrisCLIIssueToken(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)

	// start iris server
	proc := f.GDStart()
	defer proc.Stop(false)

	fooAddr := f.KeyAddress(keyFoo)
	barAddr := f.KeyAddress(keyBar)

	fooAcc := f.QueryAccount(fooAddr)
	startTokens := sdk.TokensFromConsensusPower(50)
	require.Equal(t, startTokens, fooAcc.GetCoins().AmountOf(sdk.DefaultBondDenom))

	initTokenNum := len(token.DefaultTokens()) + 1

	tokensQuery := f.QueryAssetTokens()
	require.Len(t, tokensQuery, initTokenNum)

	symbol := "abcdefgf"
	name := "Bitcoin"
	initialSupply := int64(100000000)
	scale := 18
	minUnit := "Satoshi"

	// Test --dry-run
	success, _, _ := f.TxAssetIssueToken(keyFoo, symbol, name, minUnit,
		initialSupply, scale, "--dry-run")
	require.True(t, success)

	// issue token
	f.TxAssetIssueToken(keyFoo, symbol, name, minUnit,
		initialSupply, scale, "-y")
	tests.WaitForNextNBlocksTM(1, f.Port)

	// Ensure transaction tags can be queried
	searchResult := f.QueryTxs(1, 50, "message.action=issue_token", fmt.Sprintf("message.sender=%s", fooAddr))
	require.Len(t, searchResult.Txs, 1)

	// Ensure token is directly queryable
	tokensQuery = f.QueryAssetTokens()
	require.Equal(t, initTokenNum+1, len(tokensQuery))

	token := f.QueryAssetToken(symbol)
	require.Equal(t, name, token.Name)
	require.Equal(t, strings.ToLower(symbol), token.Symbol)

	// check total supply
	totalSupply := f.QueryTotalSupplyOf(token.MinUnit)
	require.Equal(t, sdk.NewIntWithDecimal(initialSupply, scale).String(), totalSupply.String())

	// check foo account
	fooAmount := f.QueryAccount(fooAddr).Coins.AmountOf(token.MinUnit)
	require.Equal(t, sdk.NewIntWithDecimal(initialSupply, scale).String(), fooAmount.String())

	name1 := "BTC_Token"
	maxSupply1 := int64(200000000)
	mintable := true

	// Test --dry-run
	success, _, _ = f.TxAssetEditToken(keyFoo, symbol, name1, maxSupply1,
		mintable, "--dry-run")
	require.True(t, success)

	// edit token
	f.TxAssetEditToken(keyFoo, symbol, name1, maxSupply1,
		mintable, "-y")
	tests.WaitForNextNBlocksTM(1, f.Port)

	// Ensure transaction tags can be queried
	searchResult = f.QueryTxs(1, 50, "message.action=edit_token", fmt.Sprintf("message.sender=%s", fooAddr))
	require.Len(t, searchResult.Txs, 1)

	// Ensure token has been edited
	token1 := f.QueryAssetToken(symbol)
	require.Equal(t, strings.ToLower(token.Symbol), token1.Symbol)
	require.Equal(t, name1, token1.Name)
	require.Equal(t, sdk.NewInt(maxSupply1).String(), token1.MaxSupply.String())
	require.Equal(t, mintable, token1.Mintable)

	mintAmount := int64(50000000)
	// Test --dry-run
	success, _, _ = f.TxAssetMintToken(keyFoo, symbol, mintAmount, barAddr, "--dry-run")
	require.True(t, success)

	// mint token
	f.TxAssetMintToken(keyFoo, symbol, mintAmount, barAddr, "-y")
	tests.WaitForNextNBlocksTM(1, f.Port)

	// Ensure transaction tags can be queried
	searchResult = f.QueryTxs(1, 50, "message.action=mint_token", fmt.Sprintf("message.sender=%s", fooAddr))
	require.Len(t, searchResult.Txs, 1)

	// Ensure token has been minted
	totalSupply1 := f.QueryTotalSupplyOf(token.MinUnit)
	require.Equal(t, sdk.NewIntWithDecimal(initialSupply+mintAmount, scale).String(), totalSupply1.String())

	// check bar account
	barAmount := f.QueryAccount(barAddr).Coins.AmountOf(token.MinUnit)
	require.Equal(t, sdk.NewIntWithDecimal(mintAmount, scale).String(), barAmount.String())

	// Test --dry-run
	success, _, _ = f.TxAssetTransferToken(keyFoo, symbol, barAddr, "--dry-run")
	require.True(t, success)

	// transfer token owner
	f.TxAssetTransferToken(keyFoo, symbol, barAddr, "-y")
	tests.WaitForNextNBlocksTM(1, f.Port)

	// Ensure transaction tags can be queried
	searchResult = f.QueryTxs(1, 50, "message.action=transfer_token", fmt.Sprintf("message.sender=%s", fooAddr))
	require.Len(t, searchResult.Txs, 1)

	// Ensure token owner has been modified
	token2 := f.QueryAssetToken(symbol)
	require.Equal(t, barAddr.String(), token2.Owner.String())

	//
	burnAmount := int64(10000000)
	// Test --dry-run
	amount := sdk.NewCoin(token.MinUnit, sdk.NewIntWithDecimal(burnAmount, scale))
	success, _, _ = f.TxAssetBurnToken(keyFoo, amount.String(), "--dry-run")
	require.True(t, success)

	// burn token
	f.TxAssetBurnToken(keyBar, amount.String(), "-y")
	tests.WaitForNextNBlocksTM(1, f.Port)

	// Ensure transaction tags can be queried
	searchResult = f.QueryTxs(1, 50, "message.action=burn_token", fmt.Sprintf("message.sender=%s", barAddr))
	require.Len(t, searchResult.Txs, 1)

	// Ensure token has been burn
	totalSupply2 := f.QueryTotalSupplyOf(amount.Denom)
	require.Equal(t, sdk.NewIntWithDecimal(initialSupply+mintAmount-burnAmount, scale).String(), totalSupply2.String())

	// check bar account
	barAmount = f.QueryAccount(barAddr).Coins.AmountOf(token.MinUnit)
	require.Equal(t, sdk.NewIntWithDecimal(mintAmount-burnAmount, scale).String(), barAmount.String())

	f.Cleanup()
}

// QueryAssetTokens is iriscli query asset token tokens
func (f *Fixtures) QueryAssetTokens(flags ...string) token.Tokens {
	cmd := fmt.Sprintf("%s query asset token tokens %v", f.IriscliBinary, f.Flags())
	stdout, stderr := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	if strings.Contains(stderr, "no matching tokens found") {
		return token.Tokens{}
	}
	require.Empty(f.T, stderr)
	var out token.Tokens
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(stdout), &out)
	require.NoError(f.T, err)
	return out
}

// QueryAssetToken is iriscli query asset token tokens
func (f *Fixtures) QueryAssetToken(symbol string, flags ...string) token.FungibleToken {
	cmd := fmt.Sprintf("%s query asset token tokens %s ", f.IriscliBinary, f.Flags())
	cmd += fmt.Sprintf("--symbol=%s ", symbol)
	stdout, stderr := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	if strings.Contains(stderr, "no matching token found") {
		return token.FungibleToken{}
	}
	require.Empty(f.T, stderr)
	var out token.Tokens
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(stdout), &out)
	require.NoError(f.T, err)
	return out[0]
}

// TxAssetIssueToken is iriscli tx asset token issue
func (f *Fixtures) TxAssetIssueToken(from, symbol, name, minUnit string,
	initialSupply int64, scale int, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx asset token issue %v --keyring-backend=test --from=%s", f.IriscliBinary, f.Flags(), from)
	cmd += fmt.Sprintf(" --symbol=%s --name=%s --scale=%d --min-unit=%s --initial-supply=%d ",
		symbol, name, scale, minUnit, initialSupply)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxAssetEditToken is iriscli tx asset token edit
func (f *Fixtures) TxAssetEditToken(from, symbol, name string, maxSupply int64,
	mintable bool, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx asset token edit %s %v --keyring-backend=test --from=%s", f.IriscliBinary, symbol, f.Flags(), from)
	cmd += fmt.Sprintf(" --name=%s --max-supply=%d --mintable=%v",
		name, maxSupply, mintable)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxAssetMintToken is iriscli tx asset token mint
func (f *Fixtures) TxAssetMintToken(from, symbol string, amount int64, to sdk.AccAddress, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx asset token mint %v --keyring-backend=test --from=%s", f.IriscliBinary, f.Flags(), from)
	cmd += fmt.Sprintf(" %s --recipient=%s --amount=%d", symbol, to, amount)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxAssetTransferToken is iriscli tx asset token transfer
func (f *Fixtures) TxAssetTransferToken(from, symbol string, to sdk.AccAddress, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx asset token transfer %v --keyring-backend=test --from=%s", f.IriscliBinary, f.Flags(), from)
	cmd += fmt.Sprintf(" %s --recipient=%s", symbol, to)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}

// TxAssetTransferToken is iriscli tx asset token transfer
func (f *Fixtures) TxAssetBurnToken(from, amount string, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx asset token burn %s %v --keyring-backend=test --from=%s", f.IriscliBinary, amount, f.Flags(), from)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), keys.DefaultKeyPass)
}
