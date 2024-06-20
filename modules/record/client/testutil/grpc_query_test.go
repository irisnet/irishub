package testutil_test

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/cosmos/gogoproto/proto"
// 	"github.com/stretchr/testify/suite"

// 	"github.com/cosmos/cosmos-sdk/client/flags"
// 	"github.com/cosmos/cosmos-sdk/testutil"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/irisnet/irismod/simapp"
// 	recordcli "irismod.io/record/client/cli"
// 	recordtestutil "irismod.io/record/client/testutil"
// 	recordtypes "irismod.io/record/types"
// )

// type IntegrationTestSuite struct {
// 	suite.Suite

// 	network simapp.Network
// }

// func (s *IntegrationTestSuite) SetupSuite() {
// 	s.T().Log("setting up integration test suite")

// 	s.network = simapp.SetupNetwork(s.T())
// }

// func (s *IntegrationTestSuite) TearDownSuite() {
// 	s.T().Log("tearing down integration test suite")
// 	s.network.Cleanup()
// }

// func TestIntegrationTestSuite(t *testing.T) {
// 	suite.Run(t, new(IntegrationTestSuite))
// }

// func (s *IntegrationTestSuite) TestQueryRecordGRPC() {
// 	val := s.network.Validators[0]
// 	clientCtx := val.ClientCtx

// 	// ---------------------------------------------------------------------------

// 	from := val.Address
// 	digest := "digest"
// 	digestAlgo := "digest-algo"
// 	uri := "https://example.abc"
// 	meta := "meta data"

// 	args := []string{
// 		fmt.Sprintf("--%s=%s", recordcli.FlagURI, uri),
// 		fmt.Sprintf("--%s=%s", recordcli.FlagMeta, meta),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf(
// 			"--%s=%s",
// 			flags.FlagFees,
// 			sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String(),
// 		),
// 	}

// 	expectedCode := uint32(0)

// 	txResult := recordtestutil.CreateRecordExec(s.T(),
// 		s.network,
// 		clientCtx, from.String(), digest, digestAlgo, args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	recordID := s.network.GetAttribute(
// 		recordtypes.EventTypeCreateRecord,
// 		recordtypes.AttributeKeyRecordID,
// 		txResult.Events,
// 	)
// 	// ---------------------------------------------------------------------------

// 	baseURL := val.APIAddress
// 	url := fmt.Sprintf("%s/irismod/record/records/%s", baseURL, recordID)

// 	respType := proto.Message(&recordtypes.QueryRecordResponse{})
// 	expectedContents := []recordtypes.Content{{
// 		Digest:     digest,
// 		DigestAlgo: digestAlgo,
// 		URI:        uri,
// 		Meta:       meta,
// 	}}

// 	resp, err := testutil.GetRequest(url)
// 	s.Require().NoError(err)
// 	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
// 	record := respType.(*recordtypes.QueryRecordResponse).Record
// 	s.Require().Equal(expectedContents, record.Contents)
// }
