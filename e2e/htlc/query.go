package htlc

import (
	"fmt"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"mods.irisnet.org/e2e"
	htlccli "mods.irisnet.org/modules/htlc/client/cli"
)

// QueryTestSuite is a suite of end-to-end tests for the htlc module
type QueryTestSuite struct {
	e2e.TestSuite
}

// TestQueryCmd tests all query command in the htlc module
func (s *QueryTestSuite) TestQueryCmd() {
	val := s.Network.Validators[0]

	//------test GetCmdCreateHTLC()-------------
	// baseURL := val.APIAddress
	from := val.Address
	to := sdk.AccAddress(crypto.AddressHash([]byte("dgsbl")))
	amount := "1000" + sdk.DefaultBondDenom
	receiverOnOtherChain := "0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826"
	hashLock := "e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561"
	timeLock := uint64(50)
	timestamp := uint64(1580000000)
	// stateOpen := "HTLC_STATE_OPEN"

	args := []string{
		fmt.Sprintf("--%s=%s", htlccli.FlagTo, to),
		fmt.Sprintf("--%s=%s", htlccli.FlagAmount, amount),
		fmt.Sprintf("--%s=%s", htlccli.FlagReceiverOnOtherChain, receiverOnOtherChain),
		fmt.Sprintf("--%s=%s", htlccli.FlagHashLock, hashLock),
		fmt.Sprintf("--%s=%d", htlccli.FlagTimeLock, timeLock),
		fmt.Sprintf("--%s=%d", htlccli.FlagTimestamp, timestamp),

		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastSync),
		fmt.Sprintf(
			"--%s=%s",
			flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.Network.BondDenom, sdk.NewInt(10))).String(),
		),
	}

	_ = CreateHTLCExec(
		s.T(),
		s.Network,
		val.ClientCtx,
		from.String(),
		args...,
	)
}
