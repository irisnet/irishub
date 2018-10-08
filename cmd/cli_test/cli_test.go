// +build cli_test

package clitest

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/tests"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/keys"
	"github.com/irisnet/irishub/modules/gov"
)

var (
	irisHome    = ""
	iriscliHome = ""
)

func init() {
	irisHome, iriscliHome = getTestingHomeDirs()
}

func TestIrisCLISubmitProposal(t *testing.T) {
	tests.ExecuteT(t, fmt.Sprintf("iris --home=%s unsafe_reset_all", irisHome), "")
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s foo", iriscliHome), app.DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s bar", iriscliHome), app.DefaultKeyPass)
	chainID := executeInit(t, fmt.Sprintf("iris init -o --name=foo --home=%s --home-client=%s", irisHome, iriscliHome))
	executeWrite(t, fmt.Sprintf("iriscli keys add --home=%s bar", iriscliHome), app.DefaultKeyPass)

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

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli account %s %v", fooAddr, flags))
	require.Equal(t, int64(50), fooAcc.GetCoins().AmountOf("steak").Int64())

	proposalsQuery := tests.ExecuteT(t, fmt.Sprintf("iriscli gov query-proposals %v", flags), "")
	require.Equal(t, "No matching proposals found", proposalsQuery)

	// submit a test proposal
	spStr := fmt.Sprintf("iriscli gov submit-proposal %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --deposit=%s", "5steak")
	spStr += fmt.Sprintf(" --type=%s", "Text")
	spStr += fmt.Sprintf(" --title=%s", "Test")
	spStr += fmt.Sprintf(" --description=%s", "test")

	executeWrite(t, spStr, app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli account %s %v", fooAddr, flags))
	require.Equal(t, int64(45), fooAcc.GetCoins().AmountOf("steak").Int64())

	proposal1 := executeGetProposal(t, fmt.Sprintf("iriscli gov query-proposal --proposal-id=1 --output=json %v", flags))
	require.Equal(t, int64(1), proposal1.GetProposalID())
	require.Equal(t, gov.StatusDepositPeriod, proposal1.GetStatus())

	proposalsQuery = tests.ExecuteT(t, fmt.Sprintf("iriscli gov query-proposals %v", flags), "")
	require.Equal(t, "  1 - Test", proposalsQuery)

	depositStr := fmt.Sprintf("iriscli gov deposit %v", flags)
	depositStr += fmt.Sprintf(" --from=%s", "foo")
	depositStr += fmt.Sprintf(" --deposit=%s", "10steak")
	depositStr += fmt.Sprintf(" --proposal-id=%s", "1")

	executeWrite(t, depositStr, app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	fooAcc = executeGetAccount(t, fmt.Sprintf("iriscli account %s %v", fooAddr, flags))
	require.Equal(t, int64(35), fooAcc.GetCoins().AmountOf("steak").Int64())
	proposal1 = executeGetProposal(t, fmt.Sprintf("iriscli gov query-proposal --proposal-id=1 --output=json %v", flags))
	require.Equal(t, int64(1), proposal1.GetProposalID())
	require.Equal(t, gov.StatusVotingPeriod, proposal1.GetStatus())

	voteStr := fmt.Sprintf("iriscli gov vote %v", flags)
	voteStr += fmt.Sprintf(" --from=%s", "foo")
	voteStr += fmt.Sprintf(" --proposal-id=%s", "1")
	voteStr += fmt.Sprintf(" --option=%s", "Yes")

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

	// submit a second test proposal
	spStr = fmt.Sprintf("iriscli gov submit-proposal %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --deposit=%s", "5steak")
	spStr += fmt.Sprintf(" --type=%s", "Text")
	spStr += fmt.Sprintf(" --title=%s", "Apples")
	spStr += fmt.Sprintf(" --description=%s", "test")

	executeWrite(t, spStr, app.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	proposalsQuery = tests.ExecuteT(t, fmt.Sprintf("iriscli gov query-proposals --latest=1 %v", flags), "")
	require.Equal(t, "  2 - Apples", proposalsQuery)
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

func executeGetAccount(t *testing.T, cmdStr string) auth.BaseAccount {
	out := tests.ExecuteT(t, cmdStr, "")
	var initRes map[string]json.RawMessage
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(t, err, "out %v, err %v", out, err)
	value := initRes["value"]
	var acc auth.BaseAccount
	cdc := wire.NewCodec()
	wire.RegisterCrypto(cdc)
	err = cdc.UnmarshalJSON(value, &acc)
	require.NoError(t, err, "value %v, err %v", string(value), err)
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

func executeGetProposal(t *testing.T, cmdStr string) gov.Proposal {
	out := tests.ExecuteT(t, cmdStr, "")
	var proposal gov.Proposal
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
