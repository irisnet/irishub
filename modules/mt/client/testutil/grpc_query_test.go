package testutil_test

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/testutil/rest"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	mtcli "github.com/irisnet/irismod/modules/mt/client/cli"
	mttestutil "github.com/irisnet/irismod/modules/mt/client/testutil"
	mttypes "github.com/irisnet/irismod/modules/mt/types"
	"github.com/tidwall/gjson"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/irisnet/irismod/simapp"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := simapp.NewConfig()
	cfg.NumValidators = 2

	s.cfg = cfg

	var err error
	s.network, err = network.New(s.T(), s.T().TempDir(), cfg)
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
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

	// Issue
	args := []string{
		fmt.Sprintf("--%s=%s", mtcli.FlagName, denomName),
		fmt.Sprintf("--%s=%s", mtcli.FlagData, data),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType := proto.Message(&sdk.TxResponse{})
	bz, err := mttestutil.IssueDenomExec(val.ClientCtx, from.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
	denomID = gjson.Get(txResp.RawLog, "0.events.0.attributes.0.value").String()

	// Mint
	args = []string{
		fmt.Sprintf("--%s=%s", mtcli.FlagRecipient, from.String()),
		fmt.Sprintf("--%s=%s", mtcli.FlagAmount, mintAmt),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})
	bz, err = mttestutil.MintMTExec(val.ClientCtx, from.String(), denomID, args...)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
	mtID = gjson.Get(txResp.RawLog, "0.events.1.attributes.0.value").String()

	//Denom
	respType = proto.Message(&mttypes.QueryDenomResponse{})
	url := fmt.Sprintf("%s/irismod/mt/denoms/%s", baseURL, denomID)
	resp, err := rest.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	denomItem := respType.(*mttypes.QueryDenomResponse)
	s.Require().Equal(denomID, denomItem.Denom.Id)
	s.Require().Equal([]byte(data), denomItem.Denom.Data)
	s.Require().Equal(val.Address.String(), denomItem.Denom.Owner)

	//Denoms
	respType = proto.Message(&mttypes.QueryDenomsResponse{})
	url = fmt.Sprintf("%s/irismod/mt/denoms", baseURL)
	resp, err = rest.GetRequest(url)

	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	denomsItem := respType.(*mttypes.QueryDenomsResponse)
	s.Require().Equal(1, len(denomsItem.Denoms))
	s.Require().Equal(denomID, denomsItem.Denoms[0].Id)

	//MTSupply
	respType = proto.Message(&mttypes.QueryMTSupplyResponse{})
	url = fmt.Sprintf("%s/irismod/mt/mts/%s/%s/supply", baseURL, denomID, mtID)
	resp, err = rest.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	mtSupplyItem := respType.(*mttypes.QueryMTSupplyResponse)
	s.Require().Equal(mintAmtUint, mtSupplyItem.Amount)

	//MT
	respType = proto.Message(&mttypes.QueryMTResponse{})
	url = fmt.Sprintf("%s/irismod/mt/mts/%s/%s", baseURL, denomID, mtID)
	resp, err = rest.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	mtItem := respType.(*mttypes.QueryMTResponse)
	s.Require().Equal(mtID, mtItem.Mt.Id)

	//MTs
	respType = proto.Message(&mttypes.QueryMTsResponse{})
	url = fmt.Sprintf("%s/irismod/mt/mts/%s", baseURL, denomID)
	resp, err = rest.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))

	mtsItem := respType.(*mttypes.QueryMTsResponse)
	s.Require().Equal(1, len(mtsItem.Mts))
}
