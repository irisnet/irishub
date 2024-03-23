package e2e

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"golang.org/x/sync/errgroup"

	"github.com/docker/docker/client"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	fromVersion     = "v2.1.0"
	toVersion       = "v3.0.0-rc1"
	repository      = "zhiqiangz/irishub"
	upgradePlanName = "v3"
)

const (
	haltHeightDelta    = uint64(10) // The number of blocks after which to apply upgrade after creation of proposal.
	blocksAfterUpgrade = uint64(10) // The number of blocks to wait for after the upgrade has been applied.
	votingPeriod       = "10s"      // Reducing voting period for testing
	maxDepositPeriod   = "10s"      // Reducing max deposit period for testing
	depositDenom       = "uiris"    // The bond denom to be used to deposit for propsals
)

func TestChainUpgrade(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}
	t.Parallel()

	irisChain, dockerClient, ctx := startChain(t)
	chainUser := fundChainUser(t, ctx, irisChain)
	haltHeight := submitUpgradeProposalAndVote(t, ctx, irisChain, chainUser)

	height, err := irisChain.Height(ctx)
	require.NoError(t, err, "error fetching height before upgrade")

	timeoutCtx, timeoutCtxCancel := context.WithTimeout(ctx, time.Second*45)
	defer timeoutCtxCancel()

	// This should timeout due to chain halt at upgrade height.
	_ = testutil.WaitForBlocks(timeoutCtx, int(haltHeight-height)+1, irisChain)

	height, err = irisChain.Height(ctx)
	require.NoError(t, err, "error fetching height after chain should have halted")

	// Make sure that chain is halted
	require.Equal(t, haltHeight, height, "height is not equal to halt height")

	// Bring down nodes to prepare for upgrade
	err = irisChain.StopAllNodes(ctx)
	require.NoError(t, err, "error stopping node(s)")

	// Upgrade version on all nodes - We are passing in the local image for the upgrade build using `make build-docker`
	irisChain.UpgradeVersion(ctx, dockerClient, repository, toVersion)

	// Start all nodes back up.
	// Validators reach consensus on first block after upgrade height
	// And chain block production resumes 🎉
	err = irisChain.StartAllNodes(ctx)
	require.NoError(t, err, "error starting upgraded node(s)")

	timeoutCtx, timeoutCtxCancel = context.WithTimeout(ctx, time.Second*45)
	defer timeoutCtxCancel()

	err = testutil.WaitForBlocks(timeoutCtx, int(blocksAfterUpgrade), irisChain)
	require.NoError(t, err, "chain did not produce blocks after upgrade")
}

func submitUpgradeProposalAndVote(t *testing.T, ctx context.Context, irisChain *cosmos.CosmosChain, chainUser ibc.Wallet) uint64 {
	height, err := irisChain.Height(ctx) // The current chain height
	require.NoError(t, err, "error fetching height before submit upgrade proposal")

	haltHeight := height + haltHeightDelta // The height at which upgrade should be applied
	upgradeMsg := &upgradetypes.MsgSoftwareUpgrade{
		Authority: "iaa10d07y265gmmuvt4z0w9aw880jnsr700j0qn55z",
		Plan: upgradetypes.Plan{
			Name:   upgradePlanName,
			Height: int64(haltHeight),
		},
	}

	upgradeMsgJSON, err := irisChain.Config().EncodingConfig.Codec.MarshalInterfaceJSON(upgradeMsg)
	proposal := cosmos.TxProposalv1{
		Messages: []json.RawMessage{
			upgradeMsgJSON,
		},
		Deposit: "10000000000" + irisChain.Config().Denom,
		Title:   "Chain Upgrade 1",
		Summary: "First chain software upgrade",
	}

	upgradeTx, err := irisChain.SubmitProposal(ctx, chainUser.KeyName(), proposal) // Submitting the software upgrade proposal
	require.NoError(t, err, "error submitting software upgrade proposal tx")

	// err = irisChain.VoteOnProposalAllValidators(ctx, upgradeTx.ProposalID, cosmos.ProposalVoteYes)
	err = voteOnProposalAllValidators(ctx, irisChain, upgradeTx.ProposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit votes")

	_, err = cosmos.PollForProposalStatus(ctx, irisChain, height, height+haltHeightDelta, upgradeTx.ProposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")
	return haltHeight
}

func fundChainUser(t *testing.T, ctx context.Context, irisChain *cosmos.CosmosChain) ibc.Wallet {
	const userFunds = int64(10_000_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, irisChain)
	return users[0]
}

func startChain(t *testing.T) (*cosmos.CosmosChain, *client.Client, context.Context) {
	// Configuring the chain factory. We are building Stargaze chain with the version that matches the `fromVersion` value
	numOfVals := 5
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:      "irisnet",
			ChainName: "irishub-1",
			Version:   fromVersion,
			ChainConfig: ibc.ChainConfig{
				ModifyGenesis: cosmos.ModifyGenesis(getTestGenesis()), // Modifying genesis to have test-friendly gov params
				Images: []ibc.DockerImage{
					{
						Repository: "irisnet/irishub",
						Version:    fromVersion,
					},
				},
			},
			NumValidators: &numOfVals,
		},
	})
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)
	irisChain := chains[0].(*cosmos.CosmosChain)

	ic := interchaintest.NewInterchain().AddChain(irisChain)
	dockerClient, network := interchaintest.DockerSetup(t)
	ctx := context.Background()
	require.NoError(t, ic.Build(ctx, nil, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           dockerClient,
		NetworkID:        network,
		SkipPathCreation: true,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})
	return irisChain, dockerClient, ctx
}

func getTestGenesis() []cosmos.GenesisKV {
	return []cosmos.GenesisKV{
		{
			Key:   "app_state.gov.params.voting_period",
			Value: votingPeriod,
		},
		{
			Key:   "app_state.gov.params.max_deposit_period",
			Value: maxDepositPeriod,
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.denom",
			Value: depositDenom,
		},
	}
}

func voteOnProposalAllValidators(ctx context.Context, c *cosmos.CosmosChain, proposalID string, vote string) error {
	VoteOnProposal := func(ctx context.Context, tn *cosmos.ChainNode, keyName string, proposalID string, vote string) error {
		_, err := tn.ExecTx(ctx, keyName,
			"gov", "vote",
			proposalID, vote,
		)
		return err
	}

	var eg errgroup.Group
	for _, n := range c.Nodes() {
		if n.Validator {
			n := n
			eg.Go(func() error {
				return VoteOnProposal(ctx, n, "validator", proposalID, vote)
			})
		}
	}
	return eg.Wait()
}
