package record

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/e2e"
	recordcli "mods.irisnet.org/modules/record/client/cli"
	recordtypes "mods.irisnet.org/modules/record/types"
)

// TxTestSuite is a suite of end-to-end tests for the nft module
type TxTestSuite struct {
	e2e.TestSuite
}

// TestTxCmd tests all tx command in the nft module
func (s *TxTestSuite) TestTxCmd() {
	val := s.Network.Validators[0]
	clientCtx := val.ClientCtx

	// ---------------------------------------------------------------------------

	from := val.Address
	digest := "digest"
	digestAlgo := "digest-algo"
	uri := "uri"
	meta := "meta"

	args := []string{
		fmt.Sprintf("--%s=%s", recordcli.FlagURI, uri),
		fmt.Sprintf("--%s=%s", recordcli.FlagMeta, meta),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String()),
	}

	expectedCode := uint32(0)

	txResult := CreateRecordExec(s.T(),
		s.Network,
		clientCtx, from.String(), digest, digestAlgo, args...)
	s.Require().Equal(expectedCode, txResult.Code)

	recordID := s.Network.GetAttribute(recordtypes.EventTypeCreateRecord, recordtypes.AttributeKeyRecordID, txResult.Events)

	// ---------------------------------------------------------------------------

	record := &recordtypes.Record{}
	expectedContents := []recordtypes.Content{{
		Digest:     digest,
		DigestAlgo: digestAlgo,
		URI:        uri,
		Meta:       meta,
	}}

	QueryRecordExec(s.T(), s.Network, clientCtx, recordID, record)
	s.Require().Equal(expectedContents, record.Contents)
}
