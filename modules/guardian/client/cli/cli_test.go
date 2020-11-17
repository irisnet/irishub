package cli_test

import (
	"fmt"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"

	sdk "github.com/cosmos/cosmos-sdk/types"


	guardiancli "github.com/irisnet/irishub/modules/guardian/client/cli"
	guardiantestutil "github.com/irisnet/irishub/modules/guardian/client/testutil"
	guardiantypes "github.com/irisnet/irishub/modules/guardian/types"
	"github.com/irisnet/irismod/simapp"
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

func (s *IntegrationTestSuite) TestGuardian() {
	val := s.network.Validators[0]
	from := val.Address
	address := sdk.AccAddress(crypto.AddressHash([]byte("dgsbl")))
	description := "test"
	clientCtx := val.ClientCtx

	//------test GetCmdCreateSuper()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", guardiancli.FlagAddress, address.String()),
		fmt.Sprintf("--%s=%s", guardiancli.FlagDescription, description),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := guardiantestutil.CreateSuperExec(val.ClientCtx, from.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	println(bz.String())
	s.Require().Equal(expectedCode, txResp.Code)


	//------test GetCmdDeleteSuper()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", guardiancli.FlagAddress, address.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}


	respType = proto.Message(&sdk.TxResponse{})

	bz, err = guardiantestutil.DeleteSuperExec(val.ClientCtx, from.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	println(bz.String())
	s.Require().Equal(expectedCode, txResp.Code)



	//------test GetCmdQuerySupers()-------------
	respType = proto.Message(&guardiantypes.QuerySupersResponse{})
	bz, err = guardiantestutil.QuerySupersExec(clientCtx)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	supersResp := respType.(*guardiantypes.QuerySupersResponse)
	println(supersResp.String())
	s.Require().NoError(err)
}