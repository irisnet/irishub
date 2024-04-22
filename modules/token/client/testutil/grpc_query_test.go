package testutil_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokencli "github.com/irisnet/irismod/modules/token/client/cli"
	tokentestutil "github.com/irisnet/irismod/modules/token/client/testutil"
	tokentypes "github.com/irisnet/irismod/modules/token/types"
	v1 "github.com/irisnet/irismod/modules/token/types/v1"
	"github.com/irisnet/irismod/simapp"
)

type IntegrationTestSuite struct {
	suite.Suite

	network simapp.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = simapp.SetupNetwork(s.T())
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
	symbol := "kitty"
	name := "Kitty Token"
	minUnit := "kitty"
	scale := 0
	initialSupply := int64(100000000)
	maxSupply := int64(200000000)
	mintable := true
	baseURL := val.APIAddress

	//------test GetCmdIssueToken()-------------
	args := []string{
		fmt.Sprintf("--%s=%s", tokencli.FlagSymbol, symbol),
		fmt.Sprintf("--%s=%s", tokencli.FlagName, name),
		fmt.Sprintf("--%s=%s", tokencli.FlagMinUnit, minUnit),
		fmt.Sprintf("--%s=%d", tokencli.FlagScale, scale),
		fmt.Sprintf("--%s=%d", tokencli.FlagInitialSupply, initialSupply),
		fmt.Sprintf("--%s=%d", tokencli.FlagMaxSupply, maxSupply),
		fmt.Sprintf("--%s=%t", tokencli.FlagMintable, mintable),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
		),
	}
	txResult := tokentestutil.IssueTokenExec(s.T(), s.network, clientCtx, from.String(), args...)

	tokenSymbol := s.network.GetAttribute(
		tokentypes.EventTypeIssueToken,
		tokentypes.AttributeKeySymbol,
		txResult.Events,
	)

	//------test GetCmdQueryTokens()-------------
	url := fmt.Sprintf("%s/irismod/token/v1/tokens", baseURL)
	resp, err := testutil.GetRequest(url)
	respType := proto.Message(&v1.QueryTokensResponse{})
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
	tokensResp := respType.(*v1.QueryTokensResponse)
	s.Require().Equal(2, len(tokensResp.Tokens))

	//------test GetCmdQueryToken()-------------
	url = fmt.Sprintf("%s/irismod/token/v1/tokens/%s", baseURL, tokenSymbol)
	resp, err = testutil.GetRequest(url)
	respType = proto.Message(&v1.QueryTokenResponse{})
	var token v1.TokenI
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
	tokenResp := respType.(*v1.QueryTokenResponse)
	err = clientCtx.InterfaceRegistry.UnpackAny(tokenResp.Token, &token)
	s.Require().NoError(err)
	s.Require().Equal(name, token.GetName())
	s.Require().Equal(symbol, token.GetSymbol())
	s.Require().Equal(uint64(initialSupply), token.GetInitialSupply())

	//------test GetCmdQueryFee()-------------
	url = fmt.Sprintf("%s/irismod/token/v1/tokens/%s/fees", baseURL, tokenSymbol)
	resp, err = testutil.GetRequest(url)
	respType = proto.Message(&v1.QueryFeesResponse{})
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
	feeResp := respType.(*v1.QueryFeesResponse)
	expectedFeeResp := "{\"exist\":true,\"issue_fee\":{\"denom\":\"stake\",\"amount\":\"13015\"},\"mint_fee\":{\"denom\":\"stake\",\"amount\":\"1301\"}}"
	result, _ := json.Marshal(feeResp)
	s.Require().Equal(expectedFeeResp, string(result))

	//------test GetCmdQueryParams()-------------
	url = fmt.Sprintf("%s/irismod/token/v1/params", baseURL)
	resp, err = testutil.GetRequest(url)
	respType = proto.Message(&v1.QueryParamsResponse{})
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
	paramsResp := respType.(*v1.QueryParamsResponse)
	s.Require().NoError(err)
	expectedParams := "{\"token_tax_rate\":\"0.400000000000000000\",\"issue_token_base_fee\":{\"denom\":\"stake\",\"amount\":\"60000\"},\"mint_token_fee_ratio\":\"0.100000000000000000\",\"enable_erc20\":true}"
	result, _ = json.Marshal(paramsResp.Params)
	s.Require().Equal(expectedParams, string(result))
}
