package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"bufio"
	"io"

	irisInit "github.com/irisnet/irishub/server/init"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/server"
	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	distributiontypes "github.com/irisnet/irishub/modules/distribution/types"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/bank"
	"github.com/irisnet/irishub/client/context"
	distributionclient "github.com/irisnet/irishub/client/distribution"
	"github.com/irisnet/irishub/client/keys"
	recordCli "github.com/irisnet/irishub/client/record"
	servicecli "github.com/irisnet/irishub/client/service"
	stakecli "github.com/irisnet/irishub/client/stake"
	"github.com/irisnet/irishub/client/tendermint/tx"
	upgcli "github.com/irisnet/irishub/client/upgrade"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/record"
	"github.com/irisnet/irishub/modules/upgrade"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
	"github.com/irisnet/irishub/modules/service"
	"github.com/irisnet/irishub/modules/auth"
	"path/filepath"
	"io/ioutil"
	"github.com/irisnet/irishub/modules/arbitration"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/app/v0"
)

var (
	irisHome    = ""
	iriscliHome = ""
	chainID     = ""
	nodeID      = ""
)

func init() {
	irisHome, iriscliHome = getTestingHomeDirs()
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(irisInit.Bech32PrefixAccAddr, irisInit.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(irisInit.Bech32PrefixValAddr, irisInit.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(irisInit.Bech32PrefixConsAddr, irisInit.Bech32PrefixConsPub)
	config.Seal()
}

//___________________________________________________________________________________
// irisnet helper methods

func convertToIrisBaseAccount(t *testing.T, acc bank.BaseAccount) string {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
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

func modifyGenesisState(genesisState v0.GenesisFileState) v0.GenesisFileState {
	genesisState.GovData = gov.DefaultGenesisStateForCliTest()
	genesisState.UpgradeData = upgrade.DefaultGenesisStateForTest()
	genesisState.ServiceData = service.DefaultGenesisStateForTest()
	genesisState.GuardianData = guardian.DefaultGenesisStateForTest()
	genesisState.ArbitrationData = arbitration.DefaultGenesisStateForTest()

	// genesis add a profiler
	if len(genesisState.Accounts) > 0 {
		profiler := guardian.Profiler{
			Name:      "genesis",
			Addr:      genesisState.Accounts[0].Address,
			AddedAddr: genesisState.Accounts[0].Address,
		}
		genesisState.GuardianData.Profilers[0] = profiler
		genesisState.GuardianData.Trustees[0].Addr = genesisState.Accounts[0].Address
	}

	return genesisState
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
// helper methods

func initializeFixtures(t *testing.T) (chainID, servAddr, port string) {
	tests.ExecuteT(t, fmt.Sprintf("rm -rf %s ", irisHome), "")
	//tests.ExecuteT(t, fmt.Sprintf("iris --home=%s unsafe-reset-all", irisHome), "")
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s foo", iriscliHome), v0.DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s bar", iriscliHome), v0.DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("iriscli keys add --home=%s foo", iriscliHome), v0.DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("iriscli keys add --home=%s bar", iriscliHome), v0.DefaultKeyPass)
	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf(
		"iriscli keys show foo --output=json --home=%s", iriscliHome))
	chainID = executeInit(t, fmt.Sprintf("iris init -o --moniker=foo --home=%s", irisHome))
	nodeID, _ = tests.ExecuteT(t, fmt.Sprintf("iris tendermint show-node-id --home=%s ", irisHome), "")
	genFile := filepath.Join(irisHome, "config", "genesis.json")
	genDoc := readGenesisFile(t, genFile)
	var appState v0.GenesisFileState
	err := codec.Cdc.UnmarshalJSON(genDoc.AppState, &appState)
	require.NoError(t, err)
	appState.Accounts = []v0.GenesisFileAccount{v0.NewDefaultGenesisFileAccount(fooAddr)}
	appState = modifyGenesisState(appState)
	appStateJSON, err := codec.Cdc.MarshalJSON(appState)
	require.NoError(t, err)
	genDoc.AppState = appStateJSON
	genDoc.SaveAs(genFile)
	executeWrite(t, fmt.Sprintf(
		"iris gentx --name=foo --home=%s --home-client=%s", irisHome, iriscliHome),
		v0.DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("iris collect-gentxs --home=%s", irisHome), v0.DefaultKeyPass)
	// get a free port, also setup some common flags
	servAddr, port, err = server.FreeTCPAddr()
	require.NoError(t, err)
	return
}

func unmarshalStdTx(t *testing.T, s string) (stdTx auth.StdTx) {
	cdc := app.MakeCodec()
	require.Nil(t, cdc.UnmarshalJSON([]byte(s), &stdTx))
	return
}

func readGenesisFile(t *testing.T, genFile string) types.GenesisDoc {
	var genDoc types.GenesisDoc
	fp, err := os.Open(genFile)
	require.NoError(t, err)
	fileContents, err := ioutil.ReadAll(fp)
	require.NoError(t, err)
	defer fp.Close()
	err = codec.Cdc.UnmarshalJSON(fileContents, &genDoc)
	require.NoError(t, err)
	return genDoc
}

//___________________________________________________________________________________
// executors

func executeWrite(t *testing.T, cmdStr string, writes ...string) (exitSuccess bool) {
	exitSuccess, _, _ = executeWriteRetStdStreams(t, cmdStr, writes...)
	return
}

func executeWriteRetStdStreams(t *testing.T, cmdStr string, writes ...string) (bool, string, string) {
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
	return proc.ExitState.Success(), string(stdout), string(stderr)
}

func executeInit(t *testing.T, cmdStr string) (chainID string) {
	_, stderr := tests.ExecuteT(t, cmdStr, v0.DefaultKeyPass)

	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(stderr), &initRes)
	require.NoError(t, err)

	err = json.Unmarshal(initRes["chain_id"], &chainID)
	require.NoError(t, err)

	return
}

func executeGetAddrPK(t *testing.T, cmdStr string) (sdk.AccAddress, crypto.PubKey) {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var ko keys.KeyOutput
	keys.UnmarshalJSON([]byte(out), &ko)

	pk, err := sdk.GetAccPubKeyBech32(ko.PubKey)
	require.NoError(t, err)

	accAddr, err := sdk.AccAddressFromBech32(ko.Address)
	require.NoError(t, err)

	return accAddr, pk
}

// irisnet-module-helper function

func executeGetAccount(t *testing.T, cmdStr string) (acc bank.BaseAccount) {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(t, err, "out %v, err %v", out, err)

	cdc := codec.New()
	codec.RegisterCrypto(cdc)

	err = cdc.UnmarshalJSON([]byte(out), &acc)
	require.NoError(t, err, "acc %v, err %v", string(out), err)

	return acc
}

func executeGetValidatorPK(t *testing.T, cmdStr string) string {
	out, errMsg := tests.ExecuteT(t, cmdStr, "")
	require.Empty(t, errMsg)

	return out
}

func executeGetDelegatorDistrInfo(t *testing.T, cmdStr string) []distributiontypes.DelegationDistInfo {
	out, errMsg := tests.ExecuteT(t, cmdStr, "")
	require.Empty(t, errMsg)

	cdc := app.MakeCodec()
	var ddiList []distributiontypes.DelegationDistInfo
	err := cdc.UnmarshalJSON([]byte(out), &ddiList)

	require.Empty(t, err)
	return ddiList
}

func executeGetDelegationDistrInfo(t *testing.T, cmdStr string) distributiontypes.DelegationDistInfo {
	out, errMsg := tests.ExecuteT(t, cmdStr, "")
	require.Empty(t, errMsg)

	cdc := app.MakeCodec()
	var ddi distributiontypes.DelegationDistInfo
	err := cdc.UnmarshalJSON([]byte(out), &ddi)

	require.Empty(t, err)
	return ddi
}

func executeGetValidatorDistrInfo(t *testing.T, cmdStr string) distributionclient.ValidatorDistInfoOutput {
	out, errMsg := tests.ExecuteT(t, cmdStr, "")
	require.Empty(t, errMsg)

	cdc := app.MakeCodec()
	var vdi distributionclient.ValidatorDistInfoOutput
	err := cdc.UnmarshalJSON([]byte(out), &vdi)

	require.Empty(t, err)
	return vdi
}

func executeGetValidator(t *testing.T, cmdStr string) stakecli.ValidatorOutput {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var validator stakecli.ValidatorOutput
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &validator)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return validator
}

func executeGetProposal(t *testing.T, cmdStr string) gov.ProposalOutput {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var proposal gov.ProposalOutput
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &proposal)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return proposal
}

func executeGetVote(t *testing.T, cmdStr string) gov.Vote {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var vote gov.Vote
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &vote)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return vote
}

func executeGetVotes(t *testing.T, cmdStr string) []gov.Vote {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var votes []gov.Vote
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &votes)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return votes
}

func executeGetParam(t *testing.T, cmdStr string) gov.Param {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var param gov.Param
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &param)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return param
}

func executeGetUpgradeInfo(t *testing.T, cmdStr string) upgcli.UpgradeInfoOutput {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var info upgcli.UpgradeInfoOutput
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &info)

	require.NoError(t, err, "out %v\n, err %v", out, err)
	return info
}

func executeGetSwitch(t *testing.T, cmdStr string) upgrade.MsgSwitch {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var switchMsg upgrade.MsgSwitch
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &switchMsg)

	require.NoError(t, err, "out %v\n, err %v", out, err)
	return switchMsg
}

func executeGetServiceDefinition(t *testing.T, cmdStr string) servicecli.DefOutput {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var serviceDef servicecli.DefOutput
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &serviceDef)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return serviceDef
}

func executeGetServiceBinding(t *testing.T, cmdStr string) service.SvcBinding {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var serviceBinding service.SvcBinding
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &serviceBinding)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return serviceBinding
}

func executeGetServiceBindings(t *testing.T, cmdStr string) []service.SvcBinding {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var serviceBindings []service.SvcBinding
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &serviceBindings)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return serviceBindings
}

func executeGetProfilers(t *testing.T, cmdStr string) []guardian.Profiler {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var profilers []guardian.Profiler
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &profilers)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return profilers
}

func executeGetTrustees(t *testing.T, cmdStr string) []guardian.Trustee {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var trustees []guardian.Trustee
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &trustees)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return trustees
}

func executeGetServiceRequests(t *testing.T, cmdStr string) []service.SvcRequest {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var svcRequests []service.SvcRequest
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &svcRequests)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return svcRequests
}

func executeGetServiceFees(t *testing.T, cmdStr string) servicecli.FeesOutput {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var feesOutput servicecli.FeesOutput
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &feesOutput)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return feesOutput
}

func executeSubmitRecordAndGetTxHash(t *testing.T, cmdStr string, writes ...string) string {
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

	type toJSON struct {
		Height int64  `json:"Height"`
		TxHash string `json:"TxHash"`
		//Response string `json:"Response"`
	}
	var res toJSON
	cdc := app.MakeCodec()
	err = cdc.UnmarshalJSON([]byte(stdout), &res)
	require.NoError(t, err, "out %v\n, err %v", stdout, err)

	return res.TxHash
}

func executeGetRecordID(t *testing.T, cmdStr string) string {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var info tx.Info
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &info)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	recordMsg, ok := info.Tx.GetMsgs()[0].(record.MsgSubmitRecord)
	if !ok {
		fmt.Println("Err MsgSubmitRecord type assertion failed")
		return ""
	}
	return recordMsg.RecordID
}

func executeGetRecord(t *testing.T, cmdStr string) recordCli.RecordOutput {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var record recordCli.RecordOutput
	cdc := app.MakeCodec()
	err := cdc.UnmarshalJSON([]byte(out), &record)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return record
}

func executeDownloadRecord(t *testing.T, cmdStr string, filePath string, force bool) bool {

	if force {
		os.Remove(filePath)
	}

	proc := tests.GoExecuteT(t, cmdStr)
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

	if !proc.ExitState.Success() {
		return false
	}

	// Check whether download file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true

}
