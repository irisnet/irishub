package clitest

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/modules/asset"
	"github.com/stretchr/testify/require"
	tmtypes "github.com/tendermint/tendermint/types"
)

func TestIrisCLIEditToken(t *testing.T) {
	t.Parallel()
	f := InitFixtures(t)
	cdc := app.MakeCodec()

	// Update asset params for test
	genesisState := f.GenesisState()
	var assetData asset.GenesisState
	err := cdc.UnmarshalJSON(genesisState[asset.ModuleName], &assetData)
	require.NoError(t, err)
	assetData.Params.IssueTokenBaseFee = sdk.NewInt(30)
	assetDataBz, err := cdc.MarshalJSON(assetData)
	require.NoError(t, err)
	genesisState[asset.ModuleName] = assetDataBz

	genFile := filepath.Join(f.IrisdHome, "config", "genesis.json")
	genDoc, err := tmtypes.GenesisDocFromFile(genFile)
	require.NoError(t, err)
	genDoc.AppState, err = cdc.MarshalJSON(genesisState)
	require.NoError(t, genDoc.SaveAs(genFile))

	// start iris server
	proc := f.GDStart()
	defer proc.Stop(false)

	fooAddr := f.KeyAddress(keyFoo)

	fooAcc := f.QueryAccount(fooAddr)
	startTokens := sdk.TokensFromConsensusPower(50)
	require.Equal(t, startTokens, fooAcc.GetCoins().AmountOf(sdk.DefaultBondDenom))

	tokensQuery := f.QueryAssetTokens()
	require.Empty(t, tokensQuery)

	family := "fungible"
	source := "native"
	symbol := "AbcdefgH"
	name := "Bitcoin"
	initialSupply := int64(100000000)
	decimal := 18
	canonicalSymbol := "Btc"
	minUnitAlias := "Satoshi"

	// Test --dry-run
	success, _, _ := f.TxAssetIssueToken(keyFoo, source, family, name, symbol, canonicalSymbol, minUnitAlias,
		initialSupply, decimal, "--dry-run")
	require.True(t, success)

	// issue token
	f.TxAssetIssueToken(keyFoo, source, family, name, symbol, canonicalSymbol, minUnitAlias,
		initialSupply, decimal, "-y")
	tests.WaitForNextNBlocksTM(1, f.Port)

	// Ensure transaction tags can be queried
	searchResult := f.QueryTxs(1, 50, "message.action:issue_token", fmt.Sprintf("message.sender:%s", fooAddr))
	require.Len(t, searchResult.Txs, 1)

	// Ensure token is directly queryable
	tokensQuery = f.QueryAssetTokens()
	require.Equal(t, 1, len(tokensQuery))

	token := f.QueryAssetToken(symbol)
	require.Equal(t, name, token.Name)
	require.Equal(t, strings.ToLower(symbol), token.Symbol)

	// Test --dry-run
	name1 := "BTC_Token"
	maxSupply1 := int64(200000000)
	mintable := true
	success, _, _ = f.TxAssetEditToken(keyFoo, symbol, name1, maxSupply1,
		mintable, "--dry-run")
	require.True(t, success)

	// edit token
	f.TxAssetEditToken(keyFoo, symbol, name1, maxSupply1,
		mintable, "-y")
	tests.WaitForNextNBlocksTM(1, f.Port)

	// Ensure transaction tags can be queried
	searchResult = f.QueryTxs(1, 50, "message.action:edit_token", fmt.Sprintf("message.sender:%s", fooAddr))
	require.Len(t, searchResult.Txs, 1)

	// Ensure token is directly queryable
	token1 := f.QueryAssetToken(symbol)
	require.Equal(t, strings.ToLower(token.Symbol), token1.Symbol)
	require.Equal(t, name1, token1.Name)
	require.Equal(t, sdk.NewIntWithDecimal(maxSupply1, decimal).String(), token1.MaxSupply.String())
	require.Equal(t, mintable, token1.Mintable)

	// TODO: mint token and transfer owner test
}

// QueryAssetTokens is iriscli query asset tokens
func (f *Fixtures) QueryAssetTokens(flags ...string) asset.Tokens {
	cmd := fmt.Sprintf("%s query asset tokens %v", f.IriscliBinary, f.Flags())
	stdout, stderr := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	if strings.Contains(stderr, "no matching proposals found") {
		return asset.Tokens{}
	}
	require.Empty(f.T, stderr)
	var out asset.Tokens
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(stdout), &out)
	require.NoError(f.T, err)
	return out
}

// QueryAssetToken is iriscli query asset token
func (f *Fixtures) QueryAssetToken(symbol string, flags ...string) asset.FungibleToken {
	cmd := fmt.Sprintf("%s query asset token %s %v", f.IriscliBinary, symbol, f.Flags())
	stdout, stderr := tests.ExecuteT(f.T, addFlags(cmd, flags), "")
	if strings.Contains(stderr, "no matching proposals found") {
		return asset.FungibleToken{}
	}
	require.Empty(f.T, stderr)
	var out asset.FungibleToken
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(stdout), &out)
	require.NoError(f.T, err)
	return out
}

// TxAssetIssueToken is iriscli tx asset issue-token
func (f *Fixtures) TxAssetIssueToken(from, source, family, name, symbol, canonicalSymbol, minUnitAlias string,
	initialSupply int64, decimal int, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx asset issue-token %v --from=%s", f.IriscliBinary, f.Flags(), from)
	cmd += fmt.Sprintf(" --source=%s --family=%s --symbol=%s --name=%s --initial-supply=%d --decimal=%d"+
		" --canonical-symbol=%s --min-unit-alias=%s",
		source, family, symbol, name, initialSupply, decimal, canonicalSymbol, minUnitAlias)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}

// TxAssetEditToken is iriscli tx asset edit-token
func (f *Fixtures) TxAssetEditToken(from, symbol, name string, maxSupply int64,
	mintable bool, flags ...string) (bool, string, string) {
	cmd := fmt.Sprintf("%s tx asset edit-token %v --from=%s", f.IriscliBinary, f.Flags(), from)
	cmd += fmt.Sprintf(" %s --name=%s --max-supply=%d --mintable=%v",
		symbol, name, maxSupply, mintable)
	return executeWriteRetStdStreams(f.T, addFlags(cmd, flags), client.DefaultKeyPass)
}
