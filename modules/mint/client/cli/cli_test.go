package cli_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	minttestutil "github.com/irisnet/irishub/modules/mint/client/testutil"
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

func (s *IntegrationTestSuite) TestMint() {
	val := s.network.Validators[0]

	//------test GetCmdQueryParams()-------------
	respType := proto.Message(&minttypes.Params{})
	bz, err := minttestutil.QueryParamsExec(val.ClientCtx)
	s.Require().NoError(err)
	s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	params := respType.(*minttypes.Params)
	s.Require().Equal("stake", params.MintDenom)
	s.Require().Equal("0.040000000000000000", params.Inflation.String())
}
