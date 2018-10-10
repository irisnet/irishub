package clitest

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/bank"
	"github.com/irisnet/irishub/client/context"
	govcli "github.com/irisnet/irishub/client/gov"
	"github.com/irisnet/irishub/client/keys"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
)

var (
	irisHome    = ""
	iriscliHome = ""
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

//___________________________________________________________________________________
// helper methods

func getTestingHomeDirs() (string, string) {
	tmpDir := os.TempDir()
	irisHome := fmt.Sprintf("%s%s.test_iris", tmpDir, string(os.PathSeparator))
	iriscliHome := fmt.Sprintf("%s%s.test_iriscli", tmpDir, string(os.PathSeparator))
	return irisHome, iriscliHome
}

//___________________________________________________________________________________
// executors

func executeWrite(t *testing.T, cmdStr string, writes ...string) bool {
	proc := tests.GoExecuteT(t, cmdStr)

	for _, write := range writes {
		_, err := proc.StdinPipe.Write([]byte(write + "\n"))
		require.NoError(t, err)
	}
	stdout, stderr, err := proc.ReadAll()
	if err != nil {
		fmt.Println("Err on proc.ReadAll()", err, cmdStr)
	}
	// Log output.
	if len(stdout) > 0 {
		t.Log("Stdout:", cmn.Green(string(stdout)))
	}
	if len(stderr) > 0 {
		t.Log("Stderr:", cmn.Red(string(stderr)))
	}

	proc.Wait()
	return proc.ExitState.Success()
	//	bz := proc.StdoutBuffer.Bytes()
	//	fmt.Println("EXEC WRITE", string(bz))
}

func executeInit(t *testing.T, cmdStr string) (chainID string) {
	out := tests.ExecuteT(t, cmdStr, app.DefaultKeyPass)

	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(t, err)

	err = json.Unmarshal(initRes["chain_id"], &chainID)
	require.NoError(t, err)

	return
}

func executeGetAddrPK(t *testing.T, cmdStr string) (sdk.AccAddress, crypto.PubKey) {
	out := tests.ExecuteT(t, cmdStr, "")
	var ko keys.KeyOutput
	keys.UnmarshalJSON([]byte(out), &ko)

	pk, err := sdk.GetAccPubKeyBech32(ko.PubKey)
	require.NoError(t, err)

	return ko.Address, pk
}

func executeGetAccount(t *testing.T, cmdStr string) (acc *bank.BaseAccount) {
	out := tests.ExecuteT(t, cmdStr, "")
	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(t, err, "out %v, err %v", out, err)

	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)

	err = cdc.UnmarshalJSON([]byte(out), &acc)
	require.NoError(t, err, "acc %v, err %v", string(out), err)

	return acc
}

func executeGetValidator(t *testing.T, cmdStr string) stake.Validator {
	out := tests.ExecuteT(t, cmdStr, "")
	var validator stake.Validator
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &validator)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return validator
}

func executeGetProposal(t *testing.T, cmdStr string) govcli.ProposalOutput {
	out := tests.ExecuteT(t, cmdStr, "")
	var proposal govcli.ProposalOutput
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &proposal)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return proposal
}

func executeGetVote(t *testing.T, cmdStr string) gov.Vote {
	out := tests.ExecuteT(t, cmdStr, "")
	var vote gov.Vote
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &vote)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return vote
}

func executeGetVotes(t *testing.T, cmdStr string) []gov.Vote {
	out := tests.ExecuteT(t, cmdStr, "")
	var votes []gov.Vote
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &votes)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return votes
}
