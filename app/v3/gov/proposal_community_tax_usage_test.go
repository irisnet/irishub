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

func TestCommunityTaxUsageProposal_Execute(t *testing.T) {
	ctx, k, accs := createTestInput(t, sdk.NewInt(100), 2)

	taxAmount := sdk.NewCoins(sdk.NewInt64Coin(sdk.IrisAtto, 10))
	spendAmount := sdk.NewCoins(sdk.NewInt64Coin(sdk.IrisAtto, 10))

	_, _, err := k.ck.AddCoins(ctx, auth.CommunityTaxCoinsAccAddr, taxAmount)
	k.ck.IncreaseLoosenToken(ctx, taxAmount)
	require.NoError(t, err)

	err = k.guardianKeeper.AddTrustee(ctx, guardian.NewGuardian("", guardian.Ordinary, accs[0].GetAddress(), accs[0].GetAddress()))
	require.NoError(t, err)

	proposals := []CommunityTaxUsageProposal{
		{BasicProposal{ProposalType: ProposalTypeCommunityTaxUsage}, TaxUsage{UsageTypeGrant, accs[0].GetAddress(), sdk.NewDecWithPrec(2, 1), sdk.NewCoins(sdk.NewInt64Coin(sdk.IrisAtto, 11))}},
		{BasicProposal{ProposalType: ProposalTypeCommunityTaxUsage}, TaxUsage{UsageTypeDistribute, accs[1].GetAddress(), sdk.NewDecWithPrec(2, 1), spendAmount}},
	}

	errTests := []struct {
		proposal CommunityTaxUsageProposal
	}{
		{proposals[0]}, // greater than tax amount
		{proposals[1]}, // not trustee account for distribute usage
	}

	for i, tc := range errTests {
		require.NotNil(t, tc.proposal.Execute(ctx, k), "test: %d", i)
	}

	// burn 2iris-atto to empty destination address
	burnAmount := sdk.NewCoins(sdk.NewInt64Coin(sdk.IrisAtto, 2))
	burnProposal := CommunityTaxUsageProposal{BasicProposal{ProposalType: ProposalTypeCommunityTaxUsage},
		TaxUsage{UsageTypeBurn, nil, sdk.ZeroDec(), burnAmount}}
	err = burnProposal.Execute(ctx, k)
	require.NoError(t, err)
	remainingTax := k.ck.GetCoins(ctx, auth.CommunityTaxCoinsAccAddr)
	require.Equal(t, "8iris-atto", remainingTax.String())
	BurnedAccCoins := k.ck.GetCoins(ctx, auth.BurnedCoinsAccAddr)
	require.Equal(t, burnAmount.String(), BurnedAccCoins.String())

	// grant 3iris-atto to destination address(not trustee address)
	grantAddr := sdk.AccAddress([]byte("grantAddr"))
	grantAmount := sdk.NewCoins(sdk.NewInt64Coin(sdk.IrisAtto, 3))
	grantProposal := CommunityTaxUsageProposal{BasicProposal{ProposalType: ProposalTypeCommunityTaxUsage},
		TaxUsage{UsageTypeGrant, grantAddr, sdk.ZeroDec(), grantAmount}}
	err = grantProposal.Execute(ctx, k)
	require.NoError(t, err)
	remainingTax = k.ck.GetCoins(ctx, auth.CommunityTaxCoinsAccAddr)
	require.Equal(t, "5iris-atto", remainingTax.String())
	require.Equal(t, grantAmount.String(), k.ck.GetCoins(ctx, grantAddr).String())

	// distribute 2iris-atto to destination address(not trustee address)
	distributeAddr := sdk.AccAddress([]byte("distributeAddr"))
	distributeAmount := sdk.NewCoins(sdk.NewInt64Coin(sdk.IrisAtto, 2))
	distributeProposal := CommunityTaxUsageProposal{BasicProposal{ProposalType: ProposalTypeCommunityTaxUsage},
		TaxUsage{UsageTypeDistribute, distributeAddr, sdk.ZeroDec(), distributeAmount}}
	err = distributeProposal.Execute(ctx, k)
	require.Error(t, err)
	remainingTax = k.ck.GetCoins(ctx, auth.CommunityTaxCoinsAccAddr)
	require.Equal(t, "5iris-atto", remainingTax.String())
	require.Equal(t, sdk.Coins{}.String(), k.ck.GetCoins(ctx, distributeAddr).String())

	// distribute 5iris-atto to trustee address(not trustee address)
	distributeAddr = accs[0].GetAddress()
	genesisAmount := k.ck.GetCoins(ctx, distributeAddr)
	distributeAmount = sdk.NewCoins(sdk.NewInt64Coin(sdk.IrisAtto, 5))
	distributeProposal = CommunityTaxUsageProposal{BasicProposal{ProposalType: ProposalTypeCommunityTaxUsage},
		TaxUsage{UsageTypeDistribute, distributeAddr, sdk.ZeroDec(), distributeAmount}}
	err = distributeProposal.Execute(ctx, k)
	require.NoError(t, err)
	remainingTax = k.ck.GetCoins(ctx, auth.CommunityTaxCoinsAccAddr)
	require.Equal(t, sdk.NewCoins().String(), remainingTax.String())
	require.Equal(t, distributeAmount.String(), k.ck.GetCoins(ctx, distributeAddr).Sub(genesisAmount).String())
}
