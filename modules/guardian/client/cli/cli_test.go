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
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/cli"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	guardiancli "github.com/irisnet/irishub/v2/modules/guardian/client/cli"
	guardiantestutil "github.com/irisnet/irishub/v2/modules/guardian/client/testutil"
	guardiantypes "github.com/irisnet/irishub/v2/modules/guardian/types"
	"github.com/irisnet/irishub/v2/simapp"
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

func (s *IntegrationTestSuite) TestGuardian() {
	val := s.network.Validators[0]
	from := val.Address
	description := "test"
	clientCtx := val.ClientCtx
	privKeyStr := cosmoscrypto.EncryptArmorPrivKey(privKey, "", "")
	clientCtx.Keyring.ImportPrivKey(addr.String(), privKeyStr, "")
	pubKeyStr := cosmoscrypto.ArmorPubKeyBytes(pubKey.Bytes(), "")
	clientCtx.Keyring.ImportPubKey(addr.String(), pubKeyStr)
	expectedCode := uint32(0)

	amount := sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(100000000))
	args := []string{
		from.String(),
		addr.String(),
		amount.String(),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	result := simapp.ExecTxCmdWithResult(s.T(), s.network, clientCtx, bankcli.NewSendTxCmd(), args)
	s.Require().Equal(uint32(0), result.TxResult.Code, result.TxResult.Log)

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
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	result = guardiantestutil.CreateSuperExec(
		s.T(),
		s.network,
		val.ClientCtx, addr.String(), args...,
	)
	s.Require().Equal(expectedCode, result.TxResult.Code, result.TxResult.Log)

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
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	result = guardiantestutil.DeleteSuperExec(
		s.T(),
		s.network,
		val.ClientCtx, addr.String(), args...)
	s.Require().Equal(expectedCode, result.TxResult.Code, result.TxResult.Log)

	respType = proto.Message(&guardiantypes.QuerySupersResponse{})
	bz, err = guardiantestutil.QuerySupersExec(clientCtx)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(bz.Bytes(), respType))
	supersResp = respType.(*guardiantypes.QuerySupersResponse)
	s.Require().Equal(1, len(supersResp.Supers))
}
