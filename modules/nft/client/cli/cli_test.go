package cli_test

import (
	"fmt"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/tidwall/gjson"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	nftcli "github.com/irisnet/irismod/modules/nft/client/cli"
	nfttestutil "github.com/irisnet/irismod/modules/nft/client/testutil"
	"github.com/irisnet/irismod/simapp"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	cfg := network.DefaultConfig()
	cfg.AppConstructor = simapp.SimAppConstructor
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

func (s *IntegrationTestSuite) TestNft() {
	val := s.network.Validators[0]
	clientCtx := val.ClientCtx

	// ---------------------------------------------------------------------------

	from := val.Address
	//tokenName := "name"
	//tokenURI  := "uri"
	//tokenData := "data"
	//recipient := "recipient"
	//owner     := "owner"
	denomName := "name"
	denom     := "denom"
	schema    := "schema"

	args := []string{
		fmt.Sprintf("--%s=%s", nftcli.FlagDenomName, denomName),
		fmt.Sprintf("--%s=%s", nftcli.FlagSchema, schema),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	respType := proto.Message(&sdk.TxResponse{})
	expectedCode := uint32(0)

	bz, err := nfttestutil.IssueDenomExec(clientCtx, from.String(), denom, args...)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType), bz.String())
	txResp := respType.(*sdk.TxResponse)
	s.Require().Equal(expectedCode, txResp.Code)

	denomID := gjson.Get(txResp.RawLog, "0.events.0.attributes.0.value").String()


	respType = proto.Message(&nfttypes.Denom{})
	bz, err = nfttestutil.QueryDenomExec(val.ClientCtx, denomID)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(bz.Bytes(), respType))
	denomItem:=respType.(*nfttypes.Denom)
	s.Require().Equal(denomName, denomItem.Name)
	s.Require().Equal(schema, denomItem.Schema)



	// ---------------------------------------------------------------------------
	//args = []string{
	//	fmt.Sprintf("--%s=%s", nftcli.FlagTokenURI, tokenURI),
	//	fmt.Sprintf("--%s=%s", nftcli.FlagRecipient, recipient),
	//	fmt.Sprintf("--%s=%s", nftcli.FlagTokenData, tokenData),
	//	fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, tokenName),
	//
	//	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	//	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	//	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	//}
	//
	//
	//// ---------------------------------------------------------------------------
	//args = []string{
	//	fmt.Sprintf("--%s=%s", nftcli.FlagTokenURI, tokenURI),
	//	fmt.Sprintf("--%s=%s", nftcli.FlagTokenData, tokenData),
	//	fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, tokenName),
	//
	//	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	//	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	//	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	//}
	//
	//args = []string{
	//	fmt.Sprintf("--%s=%s", nftcli.FlagTokenURI, tokenURI),
	//	fmt.Sprintf("--%s=%s", nftcli.FlagTokenData, tokenData),
	//	fmt.Sprintf("--%s=%s", nftcli.FlagTokenName, tokenName),
	//
	//	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	//	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	//	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	//}
	//
	//args = []string{
	//	fmt.Sprintf("--%s=%s", nftcli.FlagOwner, owner),
	//
	//	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	//	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	//	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	//}
	//
	//args = []string{
	//	fmt.Sprintf("--%s=%s", nftcli.FlagDenom, denom),
	//
	//	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	//	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	//	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	//}
}
