package cli_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	minttestutil "github.com/irisnet/irishub/v2/modules/mint/client/testutil"
	minttypes "github.com/irisnet/irishub/v2/modules/mint/types"
	"github.com/irisnet/irishub/v2/simapp"
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

	var err error
	s.cfg = cfg
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
