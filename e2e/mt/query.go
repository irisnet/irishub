package mt

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	"mods.irisnet.org/e2e"
	mtcli "mods.irisnet.org/modules/mt/client/cli"
	mttypes "mods.irisnet.org/modules/mt/types"
)

// QueryTestSuite is a suite of end-to-end tests for the mt module
type QueryTestSuite struct {
	e2e.TestSuite
}

// TestQueryCmd tests all query command in the mt module
func (s *QueryTestSuite) TestQueryCmd() {
	denomName := "name"
	data := "data"
	mintAmt := "10"
	mintAmtUint := uint64(10)

	var (
		denomID string
		mtID    string
	)

	val := s.Validators[0]
	from := val.Address
	baseURL := val.APIAddress

	expectedCode := uint32(0)
	clientCtx := val.ClientCtx

	// Issue
	args := []string{
		fmt.Sprintf("--%s=%s", mtcli.FlagName, denomName),
		fmt.Sprintf("--%s=%s", mtcli.FlagData, data),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.BondDenom, sdk.NewInt(10))).String(),
		),
	}
	txResult := IssueDenomExec(
		s.T(),
		s.Network,
		clientCtx,
		from.String(),
		args...,
	)
	s.Require().Equal(expectedCode, txResult.Code)
	denomID = s.GetAttribute(
		mttypes.EventTypeIssueDenom,
		mttypes.AttributeKeyDenomID,
		txResult.Events,
	)

	// Mint
	args = []string{
		fmt.Sprintf("--%s=%s", mtcli.FlagRecipient, from.String()),
		fmt.Sprintf("--%s=%s", mtcli.FlagAmount, mintAmt),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.BondDenom, sdk.NewInt(100))).String(),
		),
	}

	txResult = MintMTExec(s.T(),
		s.Network,
		clientCtx, from.String(), denomID, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	mtID = s.GetAttribute(
		mttypes.EventTypeMintMT,
		mttypes.AttributeKeyMTID,
		txResult.Events,
	)

	// Denom
	respType := proto.Message(&mttypes.QueryDenomResponse{})
	url := fmt.Sprintf("%s/irismod/mt/denoms/%s", baseURL, denomID)
	resp, err := testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	denomItem := respType.(*mttypes.QueryDenomResponse)
	s.Require().Equal(denomID, denomItem.Denom.Id)
	s.Require().Equal([]byte(data), denomItem.Denom.Data)
	s.Require().Equal(val.Address.String(), denomItem.Denom.Owner)

	// Denoms
	respType = proto.Message(&mttypes.QueryDenomsResponse{})
	url = fmt.Sprintf("%s/irismod/mt/denoms", baseURL)
	resp, err = testutil.GetRequest(url)

	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	denomsItem := respType.(*mttypes.QueryDenomsResponse)
	s.Require().Equal(1, len(denomsItem.Denoms))
	s.Require().Equal(denomID, denomsItem.Denoms[0].Id)

	// MTSupply
	respType = proto.Message(&mttypes.QueryMTSupplyResponse{})
	url = fmt.Sprintf("%s/irismod/mt/mts/%s/%s/supply", baseURL, denomID, mtID)
	resp, err = testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	mtSupplyItem := respType.(*mttypes.QueryMTSupplyResponse)
	s.Require().Equal(mintAmtUint, mtSupplyItem.Amount)

	// MT
	respType = proto.Message(&mttypes.QueryMTResponse{})
	url = fmt.Sprintf("%s/irismod/mt/mts/%s/%s", baseURL, denomID, mtID)
	resp, err = testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	mtItem := respType.(*mttypes.QueryMTResponse)
	s.Require().Equal(mtID, mtItem.Mt.Id)

	// MTs
	respType = proto.Message(&mttypes.QueryMTsResponse{})
	url = fmt.Sprintf("%s/irismod/mt/mts/%s", baseURL, denomID)
	resp, err = testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	mtsItem := respType.(*mttypes.QueryMTsResponse)
	s.Require().Equal(1, len(mtsItem.Mts))
}
