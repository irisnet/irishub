package v1

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto/tmhash"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tokentypes "mods.irisnet.org/modules/token/types"
)

var (
	emptyAddr string

	addr1 = sdk.AccAddress(tmhash.SumTruncated([]byte("addr1"))).String()
	addr2 = sdk.AccAddress(tmhash.SumTruncated([]byte("addr2"))).String()
)

// test ValidateBasic for MsgIssueToken
func TestMsgIssueAsset(t *testing.T) {
	addr := sdk.AccAddress(tmhash.SumTruncated([]byte("test"))).String()

	tests := []struct {
		testCase string
		*MsgIssueToken
		expectPass bool
	}{
		{
			"basic good",
			NewMsgIssueToken("stake", "satoshi", "Bitcoin Network", 18, 1, 1, true, addr),
			true,
		},
		{
			"symbol empty",
			NewMsgIssueToken("", "satoshi", "Bitcoin Network", 18, 1, 1, true, addr),
			false,
		},
		{
			"symbol error",
			NewMsgIssueToken("b&stake", "satoshi", "Bitcoin Network", 18, 1, 1, true, addr),
			false,
		},
		{
			"symbol first letter is num",
			NewMsgIssueToken("4stake", "satoshi", "Bitcoin Network", 18, 1, 1, true, addr),
			false,
		},
		{
			"symbol too long",
			NewMsgIssueToken(
				"stake123456789012345678901234567890123456789012345678901234567890",
				"satoshi",
				"Bitcoin Network",
				18,
				1,
				1,
				true,
				addr,
			),
			false,
		},
		{
			"symbol too short",
			NewMsgIssueToken("ht", "satoshi", "Bitcoin Network", 18, 1, 1, true, addr),
			false,
		},
		{"name empty", NewMsgIssueToken("stake", "satoshi", "", 18, 1, 1, true, addr), false},
		{
			"name too long",
			NewMsgIssueToken(
				"stake",
				"satoshi",
				"Bitcoin Network aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
				18,
				1,
				1,
				true,
				addr,
			),
			false,
		},
		{
			"initial supply is zero",
			NewMsgIssueToken("stake", "satoshi", "Bitcoin Network", 18, 0, 1, true, addr),
			true,
		},
		{
			"max supply is zero",
			NewMsgIssueToken("stake", "satoshi", "Bitcoin Network", 18, 1, 0, true, addr),
			true,
		},
		{
			"init supply bigger than max supply",
			NewMsgIssueToken("stake", "satoshi", "Bitcoin Network", 18, 2, 1, true, addr),
			false,
		},
		{
			"decimal error",
			NewMsgIssueToken("stake", "satoshi", "Bitcoin Network", 19, 1, 1, true, addr),
			false,
		},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.MsgIssueToken.ValidateBasic(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.MsgIssueToken.ValidateBasic(), "test: %v", tc.testCase)
		}
	}
}

// test ValidateBasic for MsgIssueToken
func TestMsgEditToken(t *testing.T) {
	owner := sdk.AccAddress(tmhash.SumTruncated([]byte("owner"))).String()
	mintable := tokentypes.False

	tests := []struct {
		testCase string
		*MsgEditToken
		expectPass bool
	}{
		{"native basic good", NewMsgEditToken("BTC Token", "btc", 10000, mintable, owner), true},
		{"wrong symbol", NewMsgEditToken("BTC Token", "BT", 10000, mintable, owner), false},
		{"loss owner", NewMsgEditToken("BTC Token", "btc", 10000, mintable, ""), false},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.MsgEditToken.ValidateBasic(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.MsgEditToken.ValidateBasic(), "test: %v", tc.testCase)
		}
	}
}

func TestMsgEditTokenRoute(t *testing.T) {
	symbol := "btc"
	mintable := tokentypes.False

	// build a MsgEditToken
	msg := MsgEditToken{
		Symbol:    symbol,
		MaxSupply: 10000000,
		Mintable:  mintable,
	}

	require.Equal(t, "token", msg.Route())
}

func TestMsgEditTokenGetSignBytes(t *testing.T) {
	mintable := tokentypes.False

	msg := MsgEditToken{
		Name:      "BTC TOKEN",
		Owner:     sdk.AccAddress(tmhash.SumTruncated([]byte("owner"))).String(),
		Symbol:    "btc",
		MaxSupply: 21000000,
		Mintable:  mintable,
	}

	res := msg.GetSignBytes()

	expected := `{"type":"irismod/token/v1/MsgEditToken","value":{"max_supply":"21000000","mintable":"false","name":"BTC TOKEN","owner":"cosmos1fsgzj6t7udv8zhf6zj32mkqhcjcpv52ygswxa5","symbol":"btc"}}`
	require.Equal(t, expected, string(res))
}

func TestMsgMintTokenValidateBasic(t *testing.T) {
	testData := []struct {
		msg        string
		minUnit    string
		owner      string
		to         string
		amount     uint64
		expectPass bool
	}{
		{"empty minUnit", "", addr1, addr2, 1000, false},
		{"wrong minUnit", "bt", addr1, addr2, 1000, false},
		{"empty owner", "btc", emptyAddr, addr2, 1000, false},
		{"empty to", "btc", addr1, emptyAddr, 1000, true},
		{"not empty to", "btc", addr1, addr2, 1000, true},
		{"invalid amount", "btc", addr1, addr2, 0, false},
		{"basic good", "btc", addr1, addr2, 1000, true},
	}

	for _, td := range testData {
		msg := &MsgMintToken{
			Coin: sdk.Coin{
				Denom:  td.minUnit,
				Amount: sdkmath.NewIntFromUint64(td.amount),
			},
			Receiver: td.to,
			Owner:    td.owner,
		}
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.msg)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.msg)
		}
	}
}

func TestMsgBurnTokenValidateBasic(t *testing.T) {
	testData := []struct {
		msg        string
		minUnit    string
		sender     string
		amount     uint64
		expectPass bool
	}{
		{"basic good", "btc", addr1, 1000, true},
		{"empty minUnit", "", addr1, 1000, false},
		{"wrong minUnit", "bt", addr1, 1000, false},
		{"empty sender", "btc", emptyAddr, 1000, false},
		{"invalid amount", "btc", addr1, 0, false},
	}

	for _, td := range testData {
		msg := MsgBurnToken{
			Coin: sdk.Coin{
				Denom:  td.minUnit,
				Amount: sdkmath.NewIntFromUint64(td.amount),
			},
			Sender: td.sender,
		}
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.msg)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.msg)
		}
	}
}

func TestMsgTransferTokenOwnerValidation(t *testing.T) {
	testData := []struct {
		name       string
		srcOwner   string
		symbol     string
		dstOwner   string
		expectPass bool
	}{
		{"empty srcOwner", emptyAddr, "btc", addr1, false},
		{"empty symbol", addr1, "", addr2, false},
		{"empty dstOwner", addr1, "btc", emptyAddr, false},
		{"invalid symbol", addr1, "btc_min", addr2, false},
		{"basic good", addr1, "btc", addr2, true},
	}

	for _, td := range testData {
		msg := NewMsgTransferTokenOwner(td.srcOwner, td.dstOwner, td.symbol)
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}

func TestMsgDeployERC20(t *testing.T) {
	testData := []struct {
		symbol     string
		name       string
		scale      uint32
		minUnit    string
		authority  string
		expectPass bool
	}{
		{symbol: "btc", name: "BTC TOKEN", scale: 18, minUnit: "staoshi", authority: addr1, expectPass: true},
		{symbol: "BTC", name: "BTC TOKEN", scale: 18, minUnit: "staoshi", authority: addr1, expectPass: false},
		{symbol: "bTC", name: "BTC TOKEN", scale: 18, minUnit: "staoshi", authority: addr1, expectPass: true},
		{symbol: "stake", name: "Stake Token", scale: 18, minUnit: "ibc/3C3D7B3BE4ECC85A0E5B52A3AEC3B7DFC2AA9CA47C37821E57020D6807043BE9", authority: addr1, expectPass: true},
		{symbol: "ibc/3C3D7B3BE4ECC85A0E5B52A3AEC3B7DFC2AA9CA47C37821E57020D6807043BE9", name: "Stake Token", scale: 18, minUnit: "ibc/3C3D7B3BE4ECC85A0E5B52A3AEC3B7DFC2AA9CA47C37821E57020D6807043BE9", authority: addr1, expectPass: true},
	}

	for _, td := range testData {
		msg := MsgDeployERC20{
			Symbol:    td.symbol,
			Name:      td.name,
			Scale:     td.scale,
			MinUnit:   td.minUnit,
			Authority: td.authority,
		}
		if td.expectPass {
			require.Nil(t, msg.ValidateBasic(), "test: %v", td.name)
		} else {
			require.NotNil(t, msg.ValidateBasic(), "test: %v", td.name)
		}
	}
}
