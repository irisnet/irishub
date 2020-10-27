package cli_test

import (
	"fmt"
	"testing"

	"github.com/tidwall/gjson"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	tokencli "github.com/irisnet/irismod/modules/token/client/cli"
	tokentestutil "github.com/irisnet/irismod/modules/token/client/testutil"
	tokentypes "github.com/irisnet/irismod/modules/token/types"
	"github.com/stretchr/testify/suite"

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

func (s *IntegrationTestSuite) TestToken() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx
	// ---------------------------------------------------------------------------

	from := val.Address
	symbol := "Kitty"
	name := "Kitty Token"
	minUnit := "kitty"
	scale := 0
	initialSupply := 100000000
	maxSupply := 100000000
	mintable := true
	//mintAmount := 50000000

	args := []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagSymbol, symbol),
		fmt.Sprintf("--%s=%s", tokencli.FlagName, name),
		fmt.Sprintf("--%s=%s", tokencli.FlagMinUnit, minUnit),
		fmt.Sprintf("--%s=%d", tokencli.FlagScale, scale),
		fmt.Sprintf("--%s=%d", tokencli.FlagInitialSupply, initialSupply),
		fmt.Sprintf("--%s=%d", tokencli.FlagMaxSupply, maxSupply),
		fmt.Sprintf("--%s=%t", tokencli.FlagMintable, mintable),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}
	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)
	bz, err := tokentestutil.IssueTokenExec(clientCtx, from.String(), args...)

	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)
	tokenSymbol := gjson.Get(txResp.RawLog, "0.events.0.attributes.0.value").String()

	tokens := &[]tokentypes.TokenI{}

	tokentypes.RegisterLegacyAminoCodec(clientCtx.LegacyAmino)
	tokentypes.RegisterInterfaces(clientCtx.InterfaceRegistry)
	bz, err = tokentestutil.QueryTokensExec(clientCtx, from.String())
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.LegacyAmino.UnmarshalJSON(bz.Bytes(), tokens))

	var token tokentypes.TokenI
	respType = proto.Message(&types.Any{})
	bz, err = tokentestutil.QueryTokenExec(clientCtx, tokenSymbol)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	err = clientCtx.InterfaceRegistry.UnpackAny(respType.(*types.Any), &token)
	s.Require().NoError(err)

	//
	//args = []string{
	//	fmt.Sprintf("--%s=%s", tokencli.FlagName, name),
	//	fmt.Sprintf("--%s=%d", tokencli.FlagMaxSupply, maxSupply),
	//	fmt.Sprintf("--%s=%t", tokencli.FlagMintable, mintable),
	//
	//	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	//	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	//	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	//}
	//respType = proto.Message(&sdk.TxResponse{})
	//bz, err = tokentestutil.EditTokenExec(clientCtx, from.String(), symbol, args...)
	//
	//s.Require().NoError(err)
	//s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	//txResp = respType.(*sdk.TxResponse)
	//s.Require().Equal(expectedCode, txResp.Code)
	//
	//args = []string{
	//	fmt.Sprintf("--%s=%s", tokencli.FlagTo, to),
	//	fmt.Sprintf("--%s=%d", tokencli.FlagAmount, mintAmount),
	//
	//	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	//	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	//	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	//}
	//respType = proto.Message(&sdk.TxResponse{})
	//bz, err = tokentestutil.MintTokenExec(clientCtx, from.String(), symbol, args...)
	//
	//s.Require().NoError(err)
	//s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	//txResp = respType.(*sdk.TxResponse)
	//s.Require().Equal(expectedCode, txResp.Code)
	//
	//args = []string{
	//	fmt.Sprintf("--%s=%s", tokencli.FlagTo, to),
	//
	//	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	//	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	//	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	//}
	//respType = proto.Message(&sdk.TxResponse{})
	//bz, err = tokentestutil.TransferTokenOwnerExec(clientCtx, from.String(), symbol, args...)
	//
	//s.Require().NoError(err)
	//s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	//txResp = respType.(*sdk.TxResponse)
	//s.Require().Equal(expectedCode, txResp.Code)
	//
	//// ---------------------------------------------------------------------------
	//

	//
	//respType = proto.Message(&sdk.TxResponse{})
	//bz, err = tokentestutil.QueryTokensExec(clientCtx, from.String())
	//
	//s.Require().NoError(err)
	//s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	//txResp = respType.(*sdk.TxResponse)
	//s.Require().Equal(expectedCode, txResp.Code)
	//
	//respType = proto.Message(&sdk.TxResponse{})
	//bz, err = tokentestutil.QueryFeeExec(clientCtx, from.String())
	//
	//s.Require().NoError(err)
	//s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	//txResp = respType.(*sdk.TxResponse)
	//s.Require().Equal(expectedCode, txResp.Code)
	//
	//respType = proto.Message(&sdk.TxResponse{})
	//bz, err = tokentestutil.QueryParamsExec(clientCtx, from.String())
	//
	//s.Require().NoError(err)
	//s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	//txResp = respType.(*sdk.TxResponse)
	//s.Require().Equal(expectedCode, txResp.Code)
}
