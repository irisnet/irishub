package cli

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/irishub/app/v1/bank"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/irisnet/irishub/app"
	v1 "github.com/irisnet/irishub/app/v1"
	"github.com/irisnet/irishub/app/v1/asset"
	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/app/v1/gov"
	"github.com/irisnet/irishub/app/v1/service"
	"github.com/irisnet/irishub/app/v1/upgrade"
	"github.com/irisnet/irishub/client/context"
	"github.com/irisnet/irishub/client/keys"
	servicecli "github.com/irisnet/irishub/client/service"
	"github.com/irisnet/irishub/client/stake"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/guardian"
	"github.com/irisnet/irishub/server"
	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/types"
)

//___________________________________________________________________________________
// irisnet helper methods

func convertToIrisBaseAccount(t *testing.T, acc auth.BaseAccount) string {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	cliCtx := context.NewCLIContext().
		WithCodec(cdc)

	coinstr := acc.Coins.String()
	coins, err := cliCtx.ConvertToMainUnit(coinstr)
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

func modifyGenesisState(genesisState v1.GenesisFileState) v1.GenesisFileState {
	genesisState.GovData = gov.DefaultGenesisStateForCliTest()
	genesisState.UpgradeData = upgrade.DefaultGenesisStateForTest()
	genesisState.ServiceData = service.DefaultGenesisStateForTest()
	genesisState.GuardianData = guardian.DefaultGenesisStateForTest()
	genesisState.AssetData = asset.DefaultGenesisStateForTest()

	// genesis add a profiler
	if len(genesisState.Accounts) > 0 {
		gd := guardian.Guardian{
			Description: "genesis",
			AccountType: guardian.Genesis,
			Address:     genesisState.Accounts[0].Address,
			AddedBy:     genesisState.Accounts[0].Address,
		}
		genesisState.GuardianData.Profilers[0] = gd
		genesisState.GuardianData.Trustees[0] = gd
	}

	return genesisState
}

func getTestingHomeDirs(name string) (string, string) {
	tmpDir := os.TempDir()
	irisHome := fmt.Sprintf("%s%s%s%s.test_iris", tmpDir, string(os.PathSeparator), name, string(os.PathSeparator))
	iriscliHome := fmt.Sprintf("%s%s%s%s.test_iriscli", tmpDir, string(os.PathSeparator), name, string(os.PathSeparator))
	return irisHome, iriscliHome
}

func getTestingHomeDirsB() (string, string) {
	tmpDir := os.TempDir()
	irisHome := fmt.Sprintf("%s%s.test_iris_b", tmpDir, string(os.PathSeparator))
	iriscliHome := fmt.Sprintf("%s%s.test_iriscli_b", tmpDir, string(os.PathSeparator))
	return irisHome, iriscliHome
}

//___________________________________________________________________________________
// helper methods

func initializeFixtures(t *testing.T) (chainID, servAddr, port, irisHome, iriscliHome, p2pAddr string) {
	irisHome, iriscliHome = getTestingHomeDirs(t.Name())
	tests.ExecuteT(t, fmt.Sprintf("rm -rf %s ", irisHome), "")
	//tests.ExecuteT(t, fmt.Sprintf("iris --home=%s unsafe-reset-all", irisHome), "")
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s foo", iriscliHome), sdk.DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s bar", iriscliHome), sdk.DefaultKeyPass)
	executeWriteCheckErr(t, fmt.Sprintf("iriscli keys add --home=%s foo", iriscliHome), sdk.DefaultKeyPass)
	executeWriteCheckErr(t, fmt.Sprintf("iriscli keys add --home=%s bar", iriscliHome), sdk.DefaultKeyPass)
	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf(
		"iriscli keys show foo --output=json --home=%s", iriscliHome))
	chainID = executeInit(t, fmt.Sprintf("iris init -o --moniker=foo --home=%s", irisHome))
	genFile := filepath.Join(irisHome, "config", "genesis.json")
	genDoc := readGenesisFile(t, genFile)
	var appState v1.GenesisFileState
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON(genDoc.AppState, &appState)
	require.NoError(t, err)
	appState.Accounts = []v1.GenesisFileAccount{v1.NewDefaultGenesisFileAccount(fooAddr)}
	appState = modifyGenesisState(appState)
	appStateJSON, err := codec.Cdc.MarshalJSON(appState)
	require.NoError(t, err)
	genDoc.AppState = appStateJSON
	genDoc.SaveAs(genFile)
	executeWriteCheckErr(t, fmt.Sprintf(
		"iris gentx --name=foo --home=%s --home-client=%s", irisHome, iriscliHome),
		sdk.DefaultKeyPass)
	executeWriteCheckErr(t, fmt.Sprintf("iris collect-gentxs --home=%s", irisHome), sdk.DefaultKeyPass)
	// get a free port, also setup some common flags
	servAddr, port, err = server.FreeTCPAddr()
	require.NoError(t, err)
	p2pAddr, _, err = server.FreeTCPAddr()
	require.NoError(t, err)
	return
}

func unmarshalStdTx(t *testing.T, s string) (stdTx auth.StdTx) {
	cdc := app.MakeLatestCodec()
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
	// broadcast transaction and return after the transaction is included by a block
	if strings.Contains(cmdStr, "--from") && strings.Contains(cmdStr, "--fee") {
		cmdStr = cmdStr + " --commit"
	}

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
	_, stderr := tests.ExecuteT(t, cmdStr, sdk.DefaultKeyPass)

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

func executeGetAccount(t *testing.T, cmdStr string) (acc auth.BaseAccount) {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(t, err, "out %v, err %v", out, err)

	cdc := app.MakeLatestCodec()

	err = cdc.UnmarshalJSON([]byte(out), &acc)
	require.NoError(t, err, "acc %v, err %v", string(out), err)

	return acc
}

func executeGetTokenStatsForAsset(t *testing.T, cmdStr string) (tokenStats bank.TokenStats) {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(t, err, "out %v, err %v", out, err)

	cdc := app.MakeLatestCodec()

	err = cdc.UnmarshalJSON([]byte(out), &tokenStats)
	require.NoError(t, err, "token-stats %v, err %v", string(out), err)

	return tokenStats
}

func executeGetValidatorPK(t *testing.T, cmdStr string) string {
	out, errMsg := tests.ExecuteT(t, cmdStr, "")
	require.Empty(t, errMsg)

	return out
}

func executeGetValidator(t *testing.T, cmdStr string) stake.ValidatorOutput {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var validator stake.ValidatorOutput
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &validator)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return validator
}

func executeGetProposal(t *testing.T, cmdStr string) gov.Proposal {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var proposal gov.Proposal
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &proposal)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return proposal
}

func executeGetProposals(t *testing.T, cmdStr string) []gov.Proposal {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var proposals []gov.Proposal
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &proposals)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return proposals
}

func executeGetVote(t *testing.T, cmdStr string) gov.Vote {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var vote gov.Vote
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &vote)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return vote
}

func executeGetVotes(t *testing.T, cmdStr string) []gov.Vote {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var votes []gov.Vote
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &votes)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return votes
}

func executeGetServiceDefinition(t *testing.T, cmdStr string) servicecli.DefOutput {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var serviceDef servicecli.DefOutput
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &serviceDef)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return serviceDef
}

func executeGetServiceBinding(t *testing.T, cmdStr string) service.SvcBinding {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var serviceBinding service.SvcBinding
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &serviceBinding)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return serviceBinding
}

func executeGetServiceBindings(t *testing.T, cmdStr string) []service.SvcBinding {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var serviceBindings []service.SvcBinding
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &serviceBindings)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return serviceBindings
}

func executeGetProfilers(t *testing.T, cmdStr string) []guardian.Guardian {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var profilers []guardian.Guardian
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &profilers)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return profilers
}

func executeGetTrustees(t *testing.T, cmdStr string) []guardian.Guardian {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var trustees []guardian.Guardian
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &trustees)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return trustees
}

func executeGetServiceRequests(t *testing.T, cmdStr string) []service.SvcRequest {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var svcRequests []service.SvcRequest
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &svcRequests)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return svcRequests
}

func executeGetServiceFees(t *testing.T, cmdStr string) servicecli.FeesOutput {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var feesOutput servicecli.FeesOutput
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &feesOutput)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return feesOutput
}

func executeGetToken(t *testing.T, cmdStr string) asset.FungibleToken {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var token asset.FungibleToken
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &token)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return token
}

func executeGetGateway(t *testing.T, cmdStr string) asset.Gateway {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var gateway asset.Gateway
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &gateway)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return gateway
}

func executeGetGateways(t *testing.T, cmdStr string) []asset.Gateway {
	out, _ := tests.ExecuteT(t, cmdStr, "")
	var gateways []asset.Gateway
	cdc := app.MakeLatestCodec()
	err := cdc.UnmarshalJSON([]byte(out), &gateways)
	require.NoError(t, err, "out %v\n, err %v", out, err)
	return gateways
}

func executeWriteCheckErr(t *testing.T, cmdStr string, writes ...string) {
	require.True(t, executeWrite(t, cmdStr, writes...))
}
