package clitest

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"bufio"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/bank"
	"github.com/irisnet/irishub/client/context"
	govcli "github.com/irisnet/irishub/client/gov"
	"github.com/irisnet/irishub/client/keys"
	stakecli "github.com/irisnet/irishub/client/stake"
	iservicecli "github.com/irisnet/irishub/client/iservice"
	upgcli "github.com/irisnet/irishub/client/upgrade"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
	"io"
)

var (
	irisHome    = ""
	iriscliHome = ""
	chainID     = ""
	nodeID      = ""
)

//___________________________________________________________________________________
// helper methods

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

func getAmountFromCoinStr(coinStr string) float64 {
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

func setupGenesisAndConfig(srcHome, dstHome string) error {
	genesisSrcFilePath := fmt.Sprintf("%s%sconfig%sgenesis.json", srcHome, string(os.PathSeparator), string(os.PathSeparator))
	configSrcFilePath := fmt.Sprintf("%s%sconfig%sconfig.toml", srcHome, string(os.PathSeparator), string(os.PathSeparator))

	genesisDstFilePath := fmt.Sprintf("%s%sconfig%sgenesis.json", dstHome, string(os.PathSeparator), string(os.PathSeparator))
	configDstFilePath := fmt.Sprintf("%s%sconfig%sconfig.toml", dstHome, string(os.PathSeparator), string(os.PathSeparator))

	err := os.Remove(genesisDstFilePath)
	if err != nil {
		return err
	}
	err = os.Remove(configDstFilePath)
	if err != nil {
		return err
	}

	err = copyFile(genesisDstFilePath, genesisSrcFilePath)
	if err != nil {
		return err
	}
	err = modifyConfigFile(configSrcFilePath, configDstFilePath)
	if err != nil {
		return err
	}
	return nil
}

func modifyGenesisFile(irisHome string) error {
	genesisFilePath := fmt.Sprintf("%s%sconfig%sgenesis.json", irisHome, string(os.PathSeparator), string(os.PathSeparator))

	genesisDoc, err := types.GenesisDocFromFile(genesisFilePath)
	if err != nil {
		return err
	}

	var genesisState app.GenesisState

	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)

	err = cdc.UnmarshalJSON(genesisDoc.AppState, &genesisState)
	if err != nil {
		return err
	}

	genesisState.GovData = gov.DefaultGenesisStateForCliTest()
	genesisState.UpgradeData = upgrade.DefaultGenesisStateForTest()

	bz, err := cdc.MarshalJSON(genesisState)
	if err != nil {
		return err
	}

	genesisDoc.AppState = bz
	return genesisDoc.SaveAs(genesisFilePath)
}

func modifyConfigFile(configSrcPath, configDstPath string) error {
	fsrc, err := os.Open(configSrcPath)
	if err != nil {
		return err
	}
	defer fsrc.Close()

	fdst, err := os.Create(configDstPath)
	if err != nil {
		return err
	}
	defer fdst.Close()

	w := bufio.NewWriter(fdst)
	br := bufio.NewReader(fsrc)

	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}

		newline := strings.Replace(string(line), "266", "366", -1)

		if strings.Index(newline, "persistent_peers") != -1 {
			newline = fmt.Sprintf("persistent_peers = \"%s@127.0.0.1:26656\"", nodeID)
		}
		fmt.Fprintln(w, newline)
	}

	return w.Flush()
}

func getTestingHomeDirs() (string, string) {
	tmpDir := os.TempDir()
	irisHome := fmt.Sprintf("%s%s.test_iris", tmpDir, string(os.PathSeparator))
	iriscliHome := fmt.Sprintf("%s%s.test_iriscli", tmpDir, string(os.PathSeparator))
	return irisHome, iriscliHome
}

func getTestingHomeDirsB() (string, string) {
	tmpDir := os.TempDir()
	irisHome := fmt.Sprintf("%s%s.test_iris_b", tmpDir, string(os.PathSeparator))
	iriscliHome := fmt.Sprintf("%s%s.test_iriscli_b", tmpDir, string(os.PathSeparator))
	return irisHome, iriscliHome
}

func copyFile(dstFile, srcFile string) error {
	src, err := os.Open(srcFile)
	if err != nil {
		return err
	}

	defer src.Close()
	dst, err := os.OpenFile(dstFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
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

func executeInit(t *testing.T, cmdStr string) (chainID, nodeID string) {
	out := tests.ExecuteT(t, cmdStr, app.DefaultKeyPass)

	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(t, err)

	err = json.Unmarshal(initRes["chain_id"], &chainID)
	require.NoError(t, err)

	err = json.Unmarshal(initRes["node_id"], &nodeID)
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

func executeGetValidator(t *testing.T, cmdStr string) stakecli.ValidatorOutput {
	out := tests.ExecuteT(t, cmdStr, "")
	var validator stakecli.ValidatorOutput
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

func executeGetParam(t *testing.T, cmdStr string) gov.Param {
	out := tests.ExecuteT(t, cmdStr, "")
	var param gov.Param
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &param)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return param
}

func executeGetUpgradeInfo(t *testing.T, cmdStr string) upgcli.UpgradeInfoOutput {
	out := tests.ExecuteT(t, cmdStr, "")
	var info upgcli.UpgradeInfoOutput
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &info)

	require.NoError(t, err, "out %v\n, err %v", out, err)
	return info
}

func executeGetSwitch(t *testing.T, cmdStr string) upgrade.MsgSwitch {
	out := tests.ExecuteT(t, cmdStr, "")
	var switchMsg upgrade.MsgSwitch
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &switchMsg)

	require.NoError(t, err, "out %v\n, err %v", out, err)
	return switchMsg
}

func executeGetServiceDefinition(t *testing.T, cmdStr string) iservicecli.ServiceOutput {
	out := tests.ExecuteT(t, cmdStr, "")
	var serviceDef iservicecli.ServiceOutput
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &serviceDef)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return serviceDef
}