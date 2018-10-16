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

func TestIrisCLISubmitProposal(t *testing.T) {
	tests.ExecuteT(t, fmt.Sprintf("iris --home=%s unsafe_reset_all", irisHome), "")
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s foo", iriscliHome), app.DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s bar", iriscliHome), app.DefaultKeyPass)
	chainID, nodeID = executeInit(t, fmt.Sprintf("iris init -o --name=foo --home=%s --home-client=%s", irisHome, iriscliHome))
	executeWrite(t, fmt.Sprintf("iriscli keys add --home=%s bar", iriscliHome), app.DefaultKeyPass)

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
	spStr += fmt.Sprintf(" --deposit=%s", "5iris")
	spStr += fmt.Sprintf(" --type=%s", "Text")
	spStr += fmt.Sprintf(" --title=%s", "Test")
	spStr += fmt.Sprintf(" --description=%s", "test")
	spStr += fmt.Sprintf(" --fee=%s", "0.004iris")

	executeWrite(t, spStr, app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num := getAmuntFromCoinStr(fooCoin)

	if !(num > 94 && num < 95) {
		t.Error("Test Failed: (94, 95) expected, recieved: {}", num)
	}

	proposal1 := executeGetProposal(t, fmt.Sprintf("iriscli gov query-proposal --proposal-id=1 --output=json %v", flags))
	require.Equal(t, int64(1), proposal1.ProposalID)
	require.Equal(t, gov.StatusDepositPeriod, proposal1.Status)

	proposalsQuery = tests.ExecuteT(t, fmt.Sprintf("iriscli gov query-proposals %v", flags), "")
	require.Equal(t, "  1 - Test", proposalsQuery)

	depositStr := fmt.Sprintf("iriscli gov deposit %v", flags)
	depositStr += fmt.Sprintf(" --from=%s", "foo")
	depositStr += fmt.Sprintf(" --deposit=%s", "5iris")
	depositStr += fmt.Sprintf(" --proposal-id=%s", "1")
	depositStr += fmt.Sprintf(" --fee=%s", "0.004iris")

	executeWrite(t, depositStr, app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num = getAmuntFromCoinStr(fooCoin)

	if !(num > 89 && num < 90) {
		t.Error("Test Failed: (89, 90) expected, recieved: {}", num)
	}

	proposal1 = executeGetProposal(t, fmt.Sprintf("iriscli gov query-proposal --proposal-id=1 --output=json %v", flags))
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

	votes := executeGetVotes(t, fmt.Sprintf("iriscli gov query-votes --proposal-id=1 --output=json %v", flags))
	require.Len(t, votes, 1)
	require.Equal(t, int64(1), votes[0].ProposalID)
	require.Equal(t, gov.OptionYes, votes[0].Option)

	proposalsQuery = tests.ExecuteT(t, fmt.Sprintf("iriscli gov query-proposals --status=DepositPeriod %v", flags), "")
	require.Equal(t, "No matching proposals found", proposalsQuery)

	proposalsQuery = tests.ExecuteT(t, fmt.Sprintf("iriscli gov query-proposals --status=VotingPeriod %v", flags), "")
	require.Equal(t, "  1 - Test", proposalsQuery)

	tests.WaitForHeightTM(votingStartBlock1+20, port)
	proposal1 = executeGetProposal(t, fmt.Sprintf("iriscli gov query-proposal --proposal-id=1 --output=json %v", flags))
	require.Equal(t, int64(1), proposal1.ProposalID)
	require.Equal(t, gov.StatusPassed, proposal1.Status)

	// submit a second test proposal
	spStr = fmt.Sprintf("iriscli gov submit-proposal %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --deposit=%s", "5iris")
	spStr += fmt.Sprintf(" --type=%s", "Text")
	spStr += fmt.Sprintf(" --title=%s", "Apples")
	spStr += fmt.Sprintf(" --description=%s", "test")
	spStr += fmt.Sprintf(" --fee=%s", "0.004iris")

	executeWrite(t, spStr, app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	proposalsQuery = tests.ExecuteT(t, fmt.Sprintf("iriscli gov query-proposals --latest=1 %v", flags), "")
	require.Equal(t, "  2 - Apples", proposalsQuery)
}
