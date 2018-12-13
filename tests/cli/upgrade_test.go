package cli

import (
	"fmt"
	"testing"

	"github.com/irisnet/irishub/server"
	"github.com/irisnet/irishub/tests"
	"github.com/stretchr/testify/require"

	//sdk "github.com/irisnet/irishub/types"
	"sync"

	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/app/v0"
)

var (
	lastSwitchHeight int64
	wg               sync.WaitGroup
)

func TestIrisCLISoftwareUpgrade(t *testing.T) {
	t.SkipNow()
	chainID, servAddr, port := initializeFixtures(t)
	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v", irisHome, servAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	require.Equal(t, "50iris", fooCoin)

	// check the upgrade info
	upgradeInfo := executeGetUpgradeInfo(t, fmt.Sprintf("iriscli upgrade info --output=json %v", flags))
	require.Equal(t, uint64(0), upgradeInfo.CurrentProposalId)
	require.Equal(t, int64(0), upgradeInfo.Verion.Id)

	/////////////////// Upgrade Proposal /////////////////////////////////
	// submit a upgrade proposal
	spStr := fmt.Sprintf("iriscli gov submit-proposal %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --deposit=%s", "10iris")
	spStr += fmt.Sprintf(" --type=%s", "SoftwareUpgrade")
	spStr += fmt.Sprintf(" --title=%s", "Upgrade")
	spStr += fmt.Sprintf(" --description=%s", "test")
	spStr += fmt.Sprintf(" --fee=%s", "0.004iris")

	executeWrite(t, spStr, v0.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	proposal1 := executeGetProposal(t, fmt.Sprintf("iriscli gov query-proposal --proposal-id=1 --output=json %v", flags))
	require.Equal(t, uint64(1), proposal1.ProposalID)
	require.Equal(t, gov.StatusVotingPeriod, proposal1.Status)

	voteStr := fmt.Sprintf("iriscli gov vote %v", flags)
	voteStr += fmt.Sprintf(" --from=%s", "foo")
	voteStr += fmt.Sprintf(" --proposal-id=%s", "1")
	voteStr += fmt.Sprintf(" --option=%s", "Yes")
	voteStr += fmt.Sprintf(" --fee=%s", "0.004iris")

	executeWrite(t, voteStr, v0.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	votes := executeGetVotes(t, fmt.Sprintf("iriscli gov query-votes --proposal-id=1 --output=json %v", flags))
	require.Len(t, votes, 1)
	require.Equal(t, uint64(1), votes[0].ProposalID)
	require.Equal(t, gov.OptionYes, votes[0].Option)

	tests.WaitForNextNBlocksTM(12, port)
	proposal1 = executeGetProposal(t, fmt.Sprintf("iriscli gov query-proposal --proposal-id=1 --output=json %v", flags))
	require.Equal(t, uint64(1), proposal1.ProposalID)
	require.Equal(t, gov.StatusPassed, proposal1.Status)

	/////////////// Stop and Run new version Software ////////////////////
	// kill iris
	proc.Stop(true)

	// start iris1 server
	proc1 := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris1 start --home=%s --rpc.laddr=%v", irisHome, servAddr))
	defer proc1.Stop(false)

	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	// check the upgrade info
	upgradeInfo = executeGetUpgradeInfo(t, fmt.Sprintf("iriscli1 upgrade info --output=json %v", flags))
	require.Equal(t, uint64(1), upgradeInfo.CurrentProposalId)
	//require.Equal(t, votingStartBlock1+10, upgradeInfo.CurrentProposalAcceptHeight)
	require.Equal(t, int64(0), upgradeInfo.Verion.Id)

	// submit switch msg
	switchStr := fmt.Sprintf("iriscli1 upgrade submit-switch %v", flags)
	switchStr += fmt.Sprintf(" --from=%s", "foo")
	switchStr += fmt.Sprintf(" --proposal-id=%s", "1")
	switchStr += fmt.Sprintf(" --title=%s", "Upgrade")
	switchStr += fmt.Sprintf(" --fee=%s", "0.004iris")

	executeWrite(t, switchStr, v0.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	// check switch msg
	switchMsg := executeGetSwitch(t, fmt.Sprintf("iriscli1 upgrade query-switch --proposal-id=1 --voter=%v --output=json %v", fooAddr.String(), flags))
	require.Equal(t, uint64(1), switchMsg.ProposalID)
	require.Equal(t, "Upgrade", switchMsg.Title)

	// check whether switched to the new version
	lastSwitchHeight = upgradeInfo.CurrentProposalAcceptHeight + 17
	tests.WaitForHeightTM(lastSwitchHeight, port)

	upgradeInfo = executeGetUpgradeInfo(t, fmt.Sprintf("iriscli1 upgrade info --output=json %v", flags))
	require.Equal(t, uint64(0), upgradeInfo.CurrentProposalId)
	//require.Equal(t, votingStartBlock1+10, upgradeInfo.CurrentProposalAcceptHeight)
	require.Equal(t, int64(1), upgradeInfo.Verion.Id)

	//////////////////// replay from version 0 for new coming node /////////////////////////////
	/// start a old node with old version and then use a new version to start
	startOldNodeBToReplay(t, chainID)

	//////////////////////////////// Bugfix Software Upgrade ////////////////////////////////

	/////////////////// Upgrade Proposal /////////////////////////////////
	// submit a upgrade proposal
	spStr = fmt.Sprintf("iriscli1 gov submit-proposal %v", flags)
	spStr += fmt.Sprintf(" --from=%s", "foo")
	spStr += fmt.Sprintf(" --deposit=%s", "10iris")
	spStr += fmt.Sprintf(" --type=%s", "SoftwareUpgrade")
	spStr += fmt.Sprintf(" --title=%s", "Upgrade")
	spStr += fmt.Sprintf(" --description=%s", "test")
	spStr += fmt.Sprintf(" --fee=%s", "0.004iris")

	executeWrite(t, spStr, v0.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	proposal2 := executeGetProposal(t, fmt.Sprintf("iriscli1 gov query-proposal --proposal-id=2 --output=json %v", flags))
	require.Equal(t, uint64(2), proposal2.ProposalID)
	require.Equal(t, gov.StatusVotingPeriod, proposal2.Status)

	//votingStartBlock2 := proposal2.VotingStartBlock

	voteStr = fmt.Sprintf("iriscli1 gov vote %v", flags)
	voteStr += fmt.Sprintf(" --from=%s", "foo")
	voteStr += fmt.Sprintf(" --proposal-id=%s", "2")
	voteStr += fmt.Sprintf(" --option=%s", "Yes")
	voteStr += fmt.Sprintf(" --fee=%s", "0.004iris")

	executeWrite(t, voteStr, v0.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	votes = executeGetVotes(t, fmt.Sprintf("iriscli1 gov query-votes --proposal-id=2 --output=json %v", flags))
	require.Len(t, votes, 1)
	require.Equal(t, uint64(2), votes[0].ProposalID)
	require.Equal(t, gov.OptionYes, votes[0].Option)

	tests.WaitForNextNBlocksTM(12, port)
	proposal2 = executeGetProposal(t, fmt.Sprintf("iriscli1 gov query-proposal --proposal-id=2 --output=json %v", flags))
	require.Equal(t, uint64(2), proposal2.ProposalID)
	require.Equal(t, gov.StatusPassed, proposal2.Status)

	/////////////// Stop and Run new version Software ////////////////////
	// kill iris
	proc1.Stop(true)

	// start iris1 server
	proc2 := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris2-bugfix start --home=%s --rpc.laddr=%v", irisHome, servAddr))
	defer proc2.Stop(false)

	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	// check the upgrade info
	upgradeInfo = executeGetUpgradeInfo(t, fmt.Sprintf("iriscli2-bugfix upgrade info --output=json %v", flags))
	require.Equal(t, uint64(2), upgradeInfo.CurrentProposalId)
	//require.Equal(t, votingStartBlock2+10, upgradeInfo.CurrentProposalAcceptHeight)
	require.Equal(t, int64(1), upgradeInfo.Verion.Id)

	// submit switch msg
	switchStr = fmt.Sprintf("iriscli2-bugfix upgrade submit-switch %v", flags)
	switchStr += fmt.Sprintf(" --from=%s", "foo")
	switchStr += fmt.Sprintf(" --proposal-id=%s", "2")
	switchStr += fmt.Sprintf(" --title=%s", "Upgrade")
	switchStr += fmt.Sprintf(" --fee=%s", "0.004iris")

	executeWrite(t, switchStr, v0.DefaultKeyPass)
	tests.WaitForNextNBlocksTM(2, port)

	// check switch msg
	switchMsg = executeGetSwitch(t, fmt.Sprintf("iriscli2-bugfix upgrade query-switch --proposal-id=2 --voter=%v --output=json %v", fooAddr.String(), flags))
	require.Equal(t, uint64(2), switchMsg.ProposalID)
	require.Equal(t, "Upgrade", switchMsg.Title)

	// check whether switched to the new version
	lastSwitchHeight = upgradeInfo.CurrentProposalAcceptHeight + 17
	tests.WaitForHeightTM(lastSwitchHeight, port)

	upgradeInfo = executeGetUpgradeInfo(t, fmt.Sprintf("iriscli2-bugfix upgrade info --output=json %v", flags))
	require.Equal(t, uint64(0), upgradeInfo.CurrentProposalId)
	//require.Equal(t, votingStartBlock2+10, upgradeInfo.CurrentProposalAcceptHeight)
	require.Equal(t, int64(2), upgradeInfo.Verion.Id)

	//////////////////// replay from version 0 for new coming node /////////////////////////////
	/// start a new node

	go startNodeBToReplay(t, chainID)

	wg.Add(1)
	wg.Wait()
	proc2.Stop(true)
}

func startOldNodeBToReplay(t *testing.T, chainID string) {
	irisBHome, iriscliBHome := getTestingHomeDirsB()
	require.True(t, irisBHome != irisHome)
	require.True(t, iriscliBHome != iriscliHome)

	tests.ExecuteT(t, fmt.Sprintf("iris --home=%s unsafe-reset-all", irisBHome), "")
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s foo", iriscliBHome), v0.DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s bar", iriscliBHome), v0.DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("iriscli keys add --home=%s foo", iriscliBHome), v0.DefaultKeyPass)
	executeInit(t, fmt.Sprintf("iris init -o --moniker=foo --home=%s", irisBHome))

	err := setupGenesisAndConfig(irisHome, irisBHome)
	require.NoError(t, err)

	// get a free port, also setup some common flags
	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)
	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliBHome, servAddr, chainID)

	// start old iris server
	tests.ExecuteT(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v", irisBHome, servAddr), "")

	// start new iris2-bugfix server
	proc3 := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris1 start --replay --home=%s --rpc.laddr=%v", irisBHome, servAddr))
	defer proc3.Stop(false)

	tests.WaitForTMStart(port)
	tests.WaitForHeightTM(lastSwitchHeight+10, port)

	// check the upgrade info
	upgradeInfo := executeGetUpgradeInfo(t, fmt.Sprintf("iriscli1 upgrade info --output=json %v", flags))
	require.Equal(t, uint64(0), upgradeInfo.CurrentProposalId)
	require.Equal(t, int64(1), upgradeInfo.Verion.Id)

}

func startNodeBToReplay(t *testing.T, chainID string) {
	irisBHome, iriscliBHome := getTestingHomeDirsB()
	require.True(t, irisBHome != irisHome)
	require.True(t, iriscliBHome != iriscliHome)

	tests.ExecuteT(t, fmt.Sprintf("iris2-bugfix --home=%s unsafe-reset-all", irisBHome), "")
	executeWrite(t, fmt.Sprintf("iriscli2-bugfix keys delete --home=%s foo", iriscliBHome), v0.DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("iriscli2-bugfix keys delete --home=%s bar", iriscliBHome), v0.DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("iriscli2-bugfix keys add --home=%s foo", iriscliBHome), v0.DefaultKeyPass)
	executeInit(t, fmt.Sprintf("iris2-bugfix init -o --moniker=foo --home=%s", irisBHome))

	err := setupGenesisAndConfig(irisHome, irisBHome)
	require.NoError(t, err)

	// get a free port, also setup some common flags
	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)
	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliBHome, servAddr, chainID)

	// start new iris2-bugfix server
	proc3 := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris2-bugfix start --home=%s --rpc.laddr=%v", irisBHome, servAddr))
	defer proc3.Stop(false)

	tests.WaitForTMStart(port)
	tests.WaitForHeightTM(lastSwitchHeight+10, port)

	// check the upgrade info
	upgradeInfo := executeGetUpgradeInfo(t, fmt.Sprintf("iriscli2-bugfix upgrade info --output=json %v", flags))
	require.Equal(t, uint64(0), upgradeInfo.CurrentProposalId)
	require.Equal(t, int64(2), upgradeInfo.Verion.Id)

	wg.Done()
}

func TestIrisStartTwoNodesToSyncBlocks(t *testing.T) {
	t.SkipNow()

	chainID, servAddr, port := initializeFixtures(t)
	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliHome, servAddr, chainID)

	// start iris server
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v", irisHome, servAddr))

	defer proc.Stop(false)
	tests.WaitForTMStart(port)
	tests.WaitForNextNBlocksTM(2, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	require.Equal(t, "50iris", fooCoin)

	//////////////////////// start node B ////////////////////////////

	go irisStartNodeB(t, chainID)

	wg.Add(1)
	wg.Wait()
	proc.Stop(true)

}

func irisStartNodeB(t *testing.T, chainID string) {
	irisBHome, iriscliBHome := getTestingHomeDirsB()
	require.True(t, irisBHome != irisHome)
	require.True(t, iriscliBHome != iriscliHome)

	tests.ExecuteT(t, fmt.Sprintf("iris --home=%s unsafe-reset-all", irisBHome), "")
	executeWrite(t, fmt.Sprintf("iriscli keys delete --home=%s foo", iriscliBHome), v0.DefaultKeyPass)
	executeWrite(t, fmt.Sprintf("iriscli keys add --home=%s foo", iriscliBHome), v0.DefaultKeyPass)
	executeInit(t, fmt.Sprintf("iris init -o --moniker=foo --home=%s", irisBHome))
	err := setupGenesisAndConfig(irisHome, irisBHome)
	require.NoError(t, err)

	// get a free port, also setup some common flags
	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)
	flags := fmt.Sprintf("--home=%s --node=%v --chain-id=%v", iriscliHome, servAddr, chainID)

	// start new iris2-bugfix server
	proc3 := tests.GoExecuteTWithStdout(t, fmt.Sprintf("iris start --home=%s --rpc.laddr=%v", irisBHome, servAddr))
	defer proc3.Stop(false)

	tests.WaitForTMStart(port)
	tests.WaitForHeightTM(10, port)

	fooAddr, _ := executeGetAddrPK(t, fmt.Sprintf("iriscli keys show foo --output=json --home=%s", iriscliHome))

	fooAcc := executeGetAccount(t, fmt.Sprintf("iriscli bank account %s %v", fooAddr, flags))
	fooCoin := convertToIrisBaseAccount(t, fooAcc)
	require.Equal(t, "50iris", fooCoin)

	wg.Done()
}
