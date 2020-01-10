package gov

import (
	"testing"

	"github.com/irisnet/irishub/app/v1/auth"
	"github.com/irisnet/irishub/modules/guardian"
	sdk "github.com/irisnet/irishub/types"
	"github.com/stretchr/testify/require"
)

func TestCommunityTaxUsageProposal_Validate(t *testing.T) {
	ctx, k, accs := createTestInput(t, sdk.NewInt(100), 2)

	taxAmount := sdk.NewCoins(sdk.NewInt64Coin(sdk.IrisAtto, 10))
	spendAmount := sdk.NewCoins(sdk.NewInt64Coin(sdk.IrisAtto, 10))

	_, _, err := k.ck.AddCoins(ctx, auth.CommunityTaxCoinsAccAddr, taxAmount)
	require.NoError(t, err)

	err = k.guardianKeeper.AddTrustee(ctx, guardian.NewGuardian("", guardian.Ordinary, accs[0].GetAddress(), accs[0].GetAddress()))
	require.NoError(t, err)

	proposals := []CommunityTaxUsageProposal{
		{BasicProposal{ProposalType: ProposalTypeCommunityTaxUsage}, TaxUsage{UsageTypeGrant, accs[0].GetAddress(), sdk.NewDecWithPrec(2, 1), sdk.Coins{sdk.Coin{Denom: sdk.IrisAtto, Amount: sdk.ZeroInt()}}}},
		{BasicProposal{ProposalType: ProposalTypeCommunityTaxUsage}, TaxUsage{UsageTypeGrant, accs[0].GetAddress(), sdk.NewDecWithPrec(2, 1), sdk.NewCoins(sdk.NewInt64Coin(sdk.IrisAtto, 11))}},
		{BasicProposal{ProposalType: ProposalTypeCommunityTaxUsage}, TaxUsage{UsageTypeDistribute, accs[1].GetAddress(), sdk.NewDecWithPrec(2, 1), spendAmount}},
		{BasicProposal{ProposalType: ProposalTypeCommunityTaxUsage}, TaxUsage{UsageTypeGrant, accs[0].GetAddress(), sdk.NewDecWithPrec(2, 1), spendAmount}},
	}

	tests := []struct {
		expectPass bool
		proposal   CommunityTaxUsageProposal
	}{
		{true, proposals[0]},  // zero amount
		{false, proposals[1]}, // greater than tax amount
		{false, proposals[2]}, // not trustee account for distribute usage
		{true, proposals[3]},  // success
	}

	for i, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.proposal.Validate(ctx, k, true), "test: %d", i)
		} else {
			require.NotNil(t, tc.proposal.Validate(ctx, k, true), "test: %d", i)
		}
	}
}
