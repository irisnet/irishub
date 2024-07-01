package record

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	"mods.irisnet.org/e2e"
	recordcli "mods.irisnet.org/modules/record/client/cli"
	recordtypes "mods.irisnet.org/modules/record/types"
)

// QueryTestSuite is a suite of end-to-end tests for the nft module
type QueryTestSuite struct {
	e2e.TestSuite
}

// TestQueryCmd tests all query command in the nft module
func (s *QueryTestSuite) TestQueryCmd() {
	val := s.Network.Validators[0]
	clientCtx := val.ClientCtx

	// ---------------------------------------------------------------------------

	from := val.Address
	digest := "digest"
	digestAlgo := "digest-algo"
	uri := "https://example.abc"
	meta := "meta data"

	args := []string{
		fmt.Sprintf("--%s=%s", recordcli.FlagURI, uri),
		fmt.Sprintf("--%s=%s", recordcli.FlagMeta, meta),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	expectedCode := uint32(0)

	txResult := CreateRecordExec(s.T(),
		s.Network,
		clientCtx, from.String(), digest, digestAlgo, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	recordID := s.Network.GetAttribute(
		recordtypes.EventTypeCreateRecord,
		recordtypes.AttributeKeyRecordID,
		txResult.Events,
	)
	// ---------------------------------------------------------------------------

	baseURL := val.APIAddress
	url := fmt.Sprintf("%s/irismod/record/records/%s", baseURL, recordID)

	respType := proto.Message(&recordtypes.QueryRecordResponse{})
	expectedContents := []recordtypes.Content{{
		Digest:     digest,
		DigestAlgo: digestAlgo,
		URI:        uri,
		Meta:       meta,
	}}

	resp, err := testutil.GetRequest(url)
	s.Require().NoError(err)
	s.Require().NoError(clientCtx.Codec.UnmarshalJSON(resp, respType))
	record := respType.(*recordtypes.QueryRecordResponse).Record
	s.Require().Equal(expectedContents, record.Contents)
}
