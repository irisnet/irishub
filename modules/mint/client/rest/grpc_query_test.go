package rest_test

import (
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/cosmos/cosmos-sdk/types/rest"

	minttypes "github.com/irisnet/irishub/modules/mint/types"
	"github.com/irisnet/irishub/simapp"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := simapp.NewConfig()
	cfg.NumValidators = 1

	s.cfg = cfg
	s.network = network.New(s.T(), cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) TestParams() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress
	//------test GetCmdQueryParams()-------------
	url := fmt.Sprintf("%s/irishub/mint/params", baseURL)
	resp, err := rest.GetRequest(url)
	respType := proto.Message(&minttypes.QueryParamsResponse{})
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, respType))
	paramsResp := respType.(*minttypes.QueryParamsResponse)
	s.Require().Equal("stake", paramsResp.Params.MintDenom)
	s.Require().Equal("0.040000000000000000", paramsResp.Params.Inflation.String())
}
