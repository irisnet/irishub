package cli_test

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/stretchr/testify/suite"

// 	"github.com/cosmos/cosmos-sdk/client/flags"
// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"mods.irisnet.org/simapp"
// 	recordcli "mods.irisnet.org/record/client/cli"
// 	recordtestutil "mods.irisnet.org/record/client/testutil"
// 	recordtypes "mods.irisnet.org/record/types"
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

// func (s *IntegrationTestSuite) TestRecord() {
// 	val := s.network.Validators[0]
// 	clientCtx := val.ClientCtx

// 	// ---------------------------------------------------------------------------

// 	from := val.Address
// 	digest := "digest"
// 	digestAlgo := "digest-algo"
// 	uri := "uri"
// 	meta := "meta"

// 	args := []string{
// 		fmt.Sprintf("--%s=%s", recordcli.FlagURI, uri),
// 		fmt.Sprintf("--%s=%s", recordcli.FlagMeta, meta),

// 		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
// 		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.network.BondDenom, sdk.NewInt(10))).String()),
// 	}

// 	expectedCode := uint32(0)

// 	txResult := recordtestutil.CreateRecordExec(s.T(),
// 		s.network,
// 		clientCtx, from.String(), digest, digestAlgo, args...)
// 	s.Require().Equal(expectedCode, txResult.Code)

// 	recordID := s.network.GetAttribute(recordtypes.EventTypeCreateRecord, recordtypes.AttributeKeyRecordID, txResult.Events)

// 	// ---------------------------------------------------------------------------

// 	record := &recordtypes.Record{}
// 	expectedContents := []recordtypes.Content{{
// 		Digest:     digest,
// 		DigestAlgo: digestAlgo,
// 		URI:        uri,
// 		Meta:       meta,
// 	}}

// 	recordtestutil.QueryRecordExec(s.T(), s.network, clientCtx, recordID, record)
// 	s.Require().Equal(expectedContents, record.Contents)
// }
