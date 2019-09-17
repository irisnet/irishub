package cli

import (
	"fmt"
	"github.com/irisnet/irishub/app/v1/gov"
	"github.com/irisnet/irishub/tests"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIrisCLISubmitProposal(t *testing.T) {
	t.Parallel()
	chainID, servAddr, port, irisHome, iriscliHome, p2pAddr := initializeFixtures(t)

	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v --output=json", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v --p2p.laddr=%v", irisHome, servAddr, p2pAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	require.Equal(t, "50iris", fooCoin)

	proposalsQuery, _ := tests.ExecuteT(t, fmt.Sprintf("iriscli gov query-proposals %v", flags), "")
	require.Equal(t, "null", proposalsQuery)

	// submit a test proposal
	spStr := fmt.Sprintf("iriscli gov submit-proposal %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --deposit=%s", "5iris")
	spStr += fmt.Sprintf(" --type=%s", "Parameter")
	spStr += fmt.Sprintf(" --title=%s", "Test")
	spStr += fmt.Sprintf(" --description=%s", "test")
	spStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	spStr += fmt.Sprintf(" --param=mint/Inflation=%s", "0.04")

	executeWrite(t, spStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num := getAmountFromCoinStr(fooCoin)

	if !(num > 44 && num < 45) {
		t.Error("Test Failed: (44, 45) expected, received:", num)
	}

	proposal1 := executeGetProposal(t, fmt.Sprintf("iriscli gov query-proposal --proposal-id=1 --output=json %v", flags))
	require.Equal(t, uint64(1), proposal1.GetProposalID())
	require.Equal(t, gov.StatusDepositPeriod, proposal1.GetStatus())

	proposals := executeGetProposals(t, fmt.Sprintf("iriscli gov query-proposals %v", flags))
	require.Equal(t, 1, len(proposals))
	require.Equal(t, "Test", proposals[0].GetTitle())

	depositStr := fmt.Sprintf("iriscli gov deposit %v", flags)
	depositStr += fmt.Sprintf(" --from=%s", "foo")
	depositStr += fmt.Sprintf(" --deposit=%s", "5iris")
	depositStr += fmt.Sprintf(" --proposal-id=%s", "1")
	depositStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	executeWrite(t, depositStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin = convertToIrisBaseAccount(t, fooAcc)
	num = getAmountFromCoinStr(fooCoin)

	if !(num > 39 && num < 40) {
		t.Error("Test Failed: (39, 40) expected, received: ", num)
	}

	proposal1 = executeGetProposal(t, fmt.Sprintf("iriscli gov query-proposal --proposal-id=1 --output=json %v", flags))
	require.Equal(t, uint64(1), proposal1.GetProposalID())
	require.Equal(t, gov.StatusVotingPeriod, proposal1.GetStatus())

	voteStr := fmt.Sprintf("iriscli gov vote %v", flags)
	voteStr += fmt.Sprintf(" --from=%s", "foo")
	voteStr += fmt.Sprintf(" --proposal-id=%s", "1")
	voteStr += fmt.Sprintf(" --option=%s", "Yes")
	voteStr += fmt.Sprintf(" --fee=%s", "0.4iris")

	executeWrite(t, voteStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	vote := executeGetVote(t, fmt.Sprintf("iriscli gov query-vote --proposal-id=1 --voter=%s --output=json %v", fooAddr, flags))
	require.Equal(t, uint64(1), vote.ProposalID)
	require.Equal(t, gov.OptionYes, vote.Option)

	votes := executeGetVotes(t, fmt.Sprintf("iriscli gov query-votes --proposal-id=1 --output=json %v", flags))
	require.Len(t, votes, 1)
	require.Equal(t, uint64(1), votes[0].ProposalID)
	require.Equal(t, gov.OptionYes, votes[0].Option)

	proposalsQuery, _ = tests.ExecuteT(t, fmt.Sprintf("iriscli gov query-proposals --status=DepositPeriod %v", flags), "")
	require.Equal(t, "null", proposalsQuery)

	proposals = executeGetProposals(t, fmt.Sprintf("iriscli gov query-proposals %v", flags))
	require.Equal(t, 1, len(proposals))
	require.Equal(t, "Test", proposals[0].GetTitle())

	tests.WaitForNextNBlocksTM(5, port)
	proposal1 = executeGetProposal(t, fmt.Sprintf("iriscli gov query-proposal --proposal-id=1 --output=json %v", flags))
	require.Equal(t, uint64(1), proposal1.GetProposalID())
	require.Equal(t, gov.StatusPassed, proposal1.GetStatus())

	// submit a second test proposal
	spStr = fmt.Sprintf("iriscli gov submit-proposal %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --deposit=%s", "5iris")
	spStr += fmt.Sprintf(" --type=%s", "Parameter")
	spStr += fmt.Sprintf(" --title=%s", "Apples")
	spStr += fmt.Sprintf(" --description=%s", "test")
	spStr += fmt.Sprintf(" --fee=%s", "0.4iris")
	spStr += fmt.Sprintf(" --param=mint/Inflation=%s", "0.05")

	executeWrite(t, spStr, sdk.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	proposals = executeGetProposals(t, fmt.Sprintf("iriscli gov query-proposals --limit=1 %v", flags))
	require.Equal(t, 1, len(proposals))
	require.Equal(t, "Apples", proposals[0].GetTitle())
}
