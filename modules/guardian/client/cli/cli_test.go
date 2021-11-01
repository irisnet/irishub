package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	cosmoscrypto "github.com/cosmos/cosmos-sdk/crypto"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	guardiancli "github.com/irisnet/irishub/modules/guardian/client/cli"
	guardiantestutil "github.com/irisnet/irishub/modules/guardian/client/testutil"
	guardiantypes "github.com/irisnet/irishub/modules/guardian/types"
	"github.com/irisnet/irishub/simapp"
)

var privKey cryptotypes.PrivKey
var pubKey cryptotypes.PubKey
var addr sdk.AccAddress

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := simapp.NewConfig()
	cfg.NumValidators = 1

	privKey, pubKey, addr = testdata.KeyTestPubAddr()
	guardian := guardiantypes.NewSuper("test", guardiantypes.Genesis, addr, addr)

	var guardianGenState guardiantypes.GenesisState
	cfg.Codec.MustUnmarshalJSON(cfg.GenesisState[guardiantypes.ModuleName], &guardianGenState)
	guardianGenState.Supers = append(guardianGenState.Supers, guardian)

	cfg.GenesisState[guardiantypes.ModuleName] = cfg.Codec.MustMarshalJSON(&guardianGenState)

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
	description := "test"
	clientCtx := val.ClientCtx
	privKeyStr := cosmoscrypto.EncryptArmorPrivKey(privKey, "", "")
	clientCtx.Keyring.ImportPrivKey(addr.String(), privKeyStr, "")
	pubKeyStr := cosmoscrypto.ArmorPubKeyBytes(pubKey.Bytes(), "")
	clientCtx.Keyring.ImportPubKey(addr.String(), pubKeyStr)

	amount := sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100000000))
	args := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	_, err := banktestutil.MsgSendExec(clientCtx, from, addr, amount, args...)
	s.Require().NoError(err)

	//------test GetCmdQuerySupers()-------------
	respType := proto.Message(&guardiantypes.QuerySupersResponse{})
	bz, err := guardiantestutil.QuerySupersExec(clientCtx)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	supersResp := respType.(*guardiantypes.QuerySupersResponse)
	s.Require().Equal(1, len(supersResp.Supers))

	//------test GetCmdCreateSuper()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", guardiancli.FlagAddress, from.String()),
		fmt.Sprintf("--%s=%s", guardiancli.FlagDescription, description),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err = guardiantestutil.CreateSuperExec(val.ClientCtx, addr.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&guardiantypes.QuerySupersResponse{})
	bz, err = guardiantestutil.QuerySupersExec(clientCtx)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	supersResp = respType.(*guardiantypes.QuerySupersResponse)
	s.Require().Equal(2, len(supersResp.Supers))

	//------test GetCmdDeleteSuper()-------------
	args = []string{
		fmt.Sprintf("--%s=%s", guardiancli.FlagAddress, from.String()),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType = proto.Message(&sdk.TxResponse{})

	bz, err = guardiantestutil.DeleteSuperExec(val.ClientCtx, addr.String(), args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp = respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	respType = proto.Message(&guardiantypes.QuerySupersResponse{})
	bz, err = guardiantestutil.QuerySupersExec(clientCtx)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	supersResp = respType.(*guardiantypes.QuerySupersResponse)
	s.Require().Equal(1, len(supersResp.Supers))
}
