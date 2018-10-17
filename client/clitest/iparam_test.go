package clitest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/tests"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/modules/gov"
)

func init() {
	irisHome, iriscliHome = getTestingHomeDirs()
}

func TestIrisCLIParameterChangeProposal(t *testing.T) {
	tests.ExecuteT(t, fmt.Sprintf("iris --home=%s unsafe_reset_all", irisHome), "")
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s foo", iriscliHome), app.DefaultKeyPass)
	chainID, nodeID = executeInit(t, fmt.Sprintf("iris init -o --name=foo --home=%s --home-client=%s", irisHome, iriscliHome))

	err := modifyGenesisFile(irisHome)
	require.NoError(t, err)

	// get a free port, also setup some common flags
	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)
	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v", irisHome, servAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	require.Equal(t, "100iris", fooCoin)

	proposalsQuery := tests.ExecuteT(t, fmt.Sprintf("iriscli gov query-proposals %v", flags), "")
	require.Equal(t, "No matching proposals found", proposalsQuery)

	// submit a test proposal
	spStr := fmt.Sprintf("iriscli gov submit-proposal %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --deposit=%s", "10iris")
	spStr += fmt.Sprintf(" --type=%s", "ParameterChange")
	spStr += fmt.Sprintf(" --title=%s", "Test")
	spStr += fmt.Sprintf(" --description=%s", "test")
	spStr += fmt.Sprintf(" --fee=%s", "0.004iris")
	spStr += fmt.Sprintf(" --param=%s", "{\"key\":\"Gov/gov/DepositProcedure\",\"value\":\"{\\\"min_deposit\\\":[{\\\"denom\\\":\\\"iris-atto\\\",\\\"amount\\\":\\\"10000000000000000000\\\"}],\\\"max_deposit_period\\\":1440}\",\"op\":\"update\"}")

	executeWrite(t, spStr, app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num := getAmountFromCoinStr(fooCoin)

	if !(num > 89 && num < 90) {
		t.Error("Test Failed: (89, 90) expected, recieved: {}", num)
	}

	proposal1 := executeGetProposal(t, fmt.Sprintf("iriscli gov query-proposal --proposal-id=1 --output=json %v", flags))
	require.Equal(t, int64(1), proposal1.ProposalID)
	require.Equal(t, gov.StatusVotingPeriod, proposal1.Status)

	votingStartBlock1 := proposal1.VotingStartBlock

	voteStr := fmt.Sprintf("iriscli gov vote %v", flags)
	voteStr += fmt.Sprintf(" --from=%s", "foo")
	voteStr += fmt.Sprintf(" --proposal-id=%s", "1")
	voteStr += fmt.Sprintf(" --option=%s", "Yes")
	voteStr += fmt.Sprintf(" --fee=%s", "0.004iris")

	executeWrite(t, voteStr, app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	vote := executeGetVote(t, fmt.Sprintf("iriscli gov query-vote --proposal-id=1 --voter=%s --output=json %v", fooAddr, flags))
	require.Equal(t, int64(1), vote.ProposalID)
	require.Equal(t, gov.OptionYes, vote.Option)

	tests.WaitForHeightTM(votingStartBlock1+20, port)
	proposal1 = executeGetProposal(t, fmt.Sprintf("iriscli gov query-proposal --proposal-id=1 --output=json %v", flags))
	require.Equal(t, int64(1), proposal1.ProposalID)
	require.Equal(t, gov.StatusPassed, proposal1.Status)

	param := executeGetParam(t, fmt.Sprintf("iriscli gov query-params --key=%v --output=json %v ","Gov/gov/DepositProcedure",flags))
	require.Equal(t,param,gov.Param{Key:"Gov/gov/DepositProcedure",Value:"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":1440}",Op:""})
}


func TestIrisCLIQueryParams(t *testing.T) {
	tests.ExecuteT(t, fmt.Sprintf("iris --home=%s unsafe_reset_all", irisHome), "")
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s foo", iriscliHome), app.DefaultKeyPass)
	chainID, nodeID = executeInit(t, fmt.Sprintf("iris init -o --name=foo --home=%s --home-client=%s", irisHome, iriscliHome))

	err := modifyGenesisFile(irisHome)
	require.NoError(t, err)

	// get a free port, also setup some common flags
	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)
	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v", irisHome, servAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	require.Equal(t, "100iris", fooCoin)

	param := executeGetParam(t, fmt.Sprintf("iriscli gov query-params --key=%v --output=json %v ","Gov/gov/DepositProcedure",flags))
	require.Equal(t,param,gov.Param{Key:"Gov/gov/DepositProcedure",Value:"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":10}",Op:""})
}

func TestIrisCLIPullParams(t *testing.T) {
	tests.ExecuteT(t, fmt.Sprintf("iris --home=%s unsafe_reset_all", irisHome), "")
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s foo", iriscliHome), app.DefaultKeyPass)
	chainID, nodeID = executeInit(t, fmt.Sprintf("iris init -o --name=foo --home=%s --home-client=%s", irisHome, iriscliHome))

	err := modifyGenesisFile(irisHome)
	require.NoError(t, err)

	// get a free port, also setup some common flags
	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)
	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v", irisHome, servAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	require.Equal(t, "100iris", fooCoin)

	tests.ExecuteT(t, fmt.Sprintf("iriscli gov pull-params --path=%v --output=json %v ",irisHome,flags), "")
}