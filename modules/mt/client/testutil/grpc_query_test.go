package testutil_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	mtcli "github.com/irisnet/irismod/modules/mt/client/cli"
	mttestutil "github.com/irisnet/irismod/modules/mt/client/testutil"
	mttypes "github.com/irisnet/irismod/modules/mt/types"

	"github.com/irisnet/irismod/simapp"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	network simapp.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.network = simapp.SetupNetwork(s.T())
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestMT() {
	denomName := "name"
	data := "data"
	mintAmt := "10"
	mintAmtUint := uint64(10)

	denomID := ""
	mtID := ""

	val := s.network.Validators[0]
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
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
	}
	txResult := mttestutil.IssueDenomExec(
		s.T(),
		s.network,
		clientCtx,
		from.String(),
		args...,
	)
	s.Require().Equal(expectedCode, txResult.Code)
	denomID = s.network.GetAttribute(mttypes.EventTypeIssueDenom, mttypes.AttributeKeyDenomID, txResult.Events)

	// Mint
	args = []string{
		fmt.Sprintf("--%s=%s", mtcli.FlagRecipient, from.String()),
		fmt.Sprintf("--%s=%s", mtcli.FlagAmount, mintAmt),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(100))).String()),
	}

	txResult = mttestutil.MintMTExec(s.T(),
		s.network,
		clientCtx, from.String(), denomID, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	mtID = s.network.GetAttribute(mttypes.EventTypeMintMT, mttypes.AttributeKeyMTID, txResult.Events)

	//Denom
	respType := proto.Message(&mttypes.QueryDenomResponse{})
	url := fmt.Sprintf("%s/irismod/mt/denoms/%s", baseURL, denomID)
	resp, err := testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	denomItem := respType.(*mttypes.QueryDenomResponse)
	s.Require().Equal(denomID, denomItem.Denom.Id)
	s.Require().Equal([]byte(data), denomItem.Denom.Data)
	s.Require().Equal(val.Address.String(), denomItem.Denom.Owner)

	//Denoms
	respType = proto.Message(&mttypes.QueryDenomsResponse{})
	url = fmt.Sprintf("%s/irismod/mt/denoms", baseURL)
	resp, err = testutil.GetRequest(url)

	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	denomsItem := respType.(*mttypes.QueryDenomsResponse)
	s.Require().Equal(1, len(denomsItem.Denoms))
	s.Require().Equal(denomID, denomsItem.Denoms[0].Id)

	//MTSupply
	respType = proto.Message(&mttypes.QueryMTSupplyResponse{})
	url = fmt.Sprintf("%s/irismod/mt/mts/%s/%s/supply", baseURL, denomID, mtID)
	resp, err = testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	mtSupplyItem := respType.(*mttypes.QueryMTSupplyResponse)
	s.Require().Equal(mintAmtUint, mtSupplyItem.Amount)

	//MT
	respType = proto.Message(&mttypes.QueryMTResponse{})
	url = fmt.Sprintf("%s/irismod/mt/mts/%s/%s", baseURL, denomID, mtID)
	resp, err = testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	mtItem := respType.(*mttypes.QueryMTResponse)
	s.Require().Equal(mtID, mtItem.Mt.Id)

	//MTs
	respType = proto.Message(&mttypes.QueryMTsResponse{})
	url = fmt.Sprintf("%s/irismod/mt/mts/%s", baseURL, denomID)
	resp, err = testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	mtsItem := respType.(*mttypes.QueryMTsResponse)
	s.Require().Equal(1, len(mtsItem.Mts))
}
