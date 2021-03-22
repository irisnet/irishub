package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

var (
	sender = sdk.AccAddress(tmhash.SumTruncated([]byte("sender"))).String()
)

func TestMsgSwapOrder_ValidateBasic(t *testing.T) {
	type fields struct {
		Input      Input
		Output     Output
		Deadline   int64
		IsBuyOrder bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "right test case", wantErr: false, fields: fields{IsBuyOrder: true, Deadline: 10, Input: Input{Address: sender, Coin: buildCoin("stake", 1000)}, Output: Output{Address: sender, Coin: buildCoin("iris", 1000)}}},
		{name: "invalid input sender", wantErr: true, fields: fields{IsBuyOrder: true, Deadline: 10, Input: Input{Address: "", Coin: buildCoin("stake", 1000)}, Output: Output{Address: sender, Coin: buildCoin("iris", 1000)}}},
		{name: "invalid input coin  denom", wantErr: true, fields: fields{IsBuyOrder: true, Deadline: 10, Input: Input{Address: sender, Coin: buildCoin("131stake", 1000)}, Output: Output{Address: sender, Coin: buildCoin("iris", 1000)}}},
		{name: "invalid input coin amount", wantErr: true, fields: fields{IsBuyOrder: true, Deadline: 10, Input: Input{Address: sender, Coin: buildCoin("stake", -1000)}, Output: Output{Address: sender, Coin: buildCoin("iris", 1000)}}},
		{name: "invalid output sender", wantErr: true, fields: fields{IsBuyOrder: true, Deadline: 10, Input: Input{Address: sender, Coin: buildCoin("stake", 1000)}, Output: Output{Address: "", Coin: buildCoin("iris", 1000)}}},
		{name: "invalid output coin denom", wantErr: true, fields: fields{IsBuyOrder: true, Deadline: 10, Input: Input{Address: sender, Coin: buildCoin("stake", 1000)}, Output: Output{Address: sender, Coin: buildCoin("131iris", 1000)}}},
		{name: "invalid output coin amount", wantErr: true, fields: fields{IsBuyOrder: true, Deadline: 10, Input: Input{Address: sender, Coin: buildCoin("stake", 1000)}, Output: Output{Address: sender, Coin: buildCoin("iris", -1000)}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgSwapOrder{
				Input:      tt.fields.Input,
				Output:     tt.fields.Output,
				Deadline:   tt.fields.Deadline,
				IsBuyOrder: tt.fields.IsBuyOrder,
			}
			if err := msg.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("MsgSwapOrder.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMsgAddLiquidity_ValidateBasic(t *testing.T) {
	type fields struct {
		MaxToken         sdk.Coin
		ExactStandardAmt sdk.Int
		MinLiquidity     sdk.Int
		Deadline         int64
		Sender           string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "invalid MaxToken denom",
			wantErr: true,
			fields: fields{
				MaxToken:         buildCoin("131stake", 1000),
				ExactStandardAmt: sdk.NewInt(100),
				MinLiquidity:     sdk.NewInt(100),
				Deadline:         1611213344,
				Sender:           sender,
			},
		},
		{
			name:    "invalid MaxToken amount",
			wantErr: true,
			fields: fields{
				MaxToken:         buildCoin("stake", -1000),
				ExactStandardAmt: sdk.NewInt(100),
				MinLiquidity:     sdk.NewInt(100),
				Deadline:         1611213344,
				Sender:           sender,
			},
		},
		{
			name:    "invalid ExactStandardAmt",
			wantErr: true,
			fields: fields{
				MaxToken:         buildCoin("stake", 1000),
				ExactStandardAmt: sdk.NewInt(-100),
				MinLiquidity:     sdk.NewInt(100),
				Deadline:         1611213344,
				Sender:           sender,
			},
		},
		{
			name:    "invalid MinLiquidity",
			wantErr: true,
			fields: fields{
				MaxToken:         buildCoin("stake", 1000),
				ExactStandardAmt: sdk.NewInt(100),
				MinLiquidity:     sdk.NewInt(-100),
				Deadline:         1611213344,
				Sender:           sender,
			},
		},
		{
			name:    "invalid Deadline",
			wantErr: true,
			fields: fields{
				MaxToken:         buildCoin("stake", 1000),
				ExactStandardAmt: sdk.NewInt(100),
				MinLiquidity:     sdk.NewInt(100),
				Deadline:         0,
				Sender:           sender,
			},
		},
		{
			name:    "invalid Sender",
			wantErr: true,
			fields: fields{
				MaxToken:         buildCoin("stake", 1000),
				ExactStandardAmt: sdk.NewInt(100),
				MinLiquidity:     sdk.NewInt(100),
				Deadline:         0,
				Sender:           "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgAddLiquidity{
				MaxToken:         tt.fields.MaxToken,
				ExactStandardAmt: tt.fields.ExactStandardAmt,
				MinLiquidity:     tt.fields.MinLiquidity,
				Deadline:         tt.fields.Deadline,
				Sender:           tt.fields.Sender,
			}
			if err := msg.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("MsgAddLiquidity.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMsgRemoveLiquidity_ValidateBasic(t *testing.T) {
	type fields struct {
		WithdrawLiquidity sdk.Coin
		MinToken          sdk.Int
		MinStandardAmt    sdk.Int
		Deadline          int64
		Sender            string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "right test case",
			wantErr: true,
			fields: fields{
				WithdrawLiquidity: buildCoin("stake", 1000),
				MinToken:          sdk.NewInt(100),
				MinStandardAmt:    sdk.NewInt(100),
				Deadline:          1611213344,
				Sender:            sender,
			},
		},
		{
			name:    "invalid WithdrawLiquidity denom",
			wantErr: true,
			fields: fields{
				WithdrawLiquidity: buildCoin("131stake", 1000),
				MinToken:          sdk.NewInt(100),
				MinStandardAmt:    sdk.NewInt(100),
				Deadline:          1611213344,
				Sender:            sender,
			},
		},
		{
			name:    "invalid MinToken",
			wantErr: true,
			fields: fields{
				WithdrawLiquidity: buildCoin("stake", -1000),
				MinToken:          sdk.NewInt(-100),
				MinStandardAmt:    sdk.NewInt(100),
				Deadline:          1611213344,
				Sender:            sender,
			},
		},
		{
			name:    "invalid MinStandardAmt",
			wantErr: true,
			fields: fields{
				WithdrawLiquidity: buildCoin("stake", 1000),
				MinToken:          sdk.NewInt(100),
				MinStandardAmt:    sdk.NewInt(-100),
				Deadline:          1611213344,
				Sender:            sender,
			},
		},
		{
			name:    "invalid Deadline",
			wantErr: true,
			fields: fields{
				WithdrawLiquidity: buildCoin("stake", 1000),
				MinToken:          sdk.NewInt(100),
				MinStandardAmt:    sdk.NewInt(100),
				Deadline:          0,
				Sender:            sender,
			},
		},
		{
			name:    "invalid Sender",
			wantErr: true,
			fields: fields{
				WithdrawLiquidity: buildCoin("stake", 1000),
				MinToken:          sdk.NewInt(100),
				MinStandardAmt:    sdk.NewInt(100),
				Deadline:          0,
				Sender:            "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := MsgRemoveLiquidity{
				WithdrawLiquidity: tt.fields.WithdrawLiquidity,
				MinToken:          tt.fields.MinToken,
				MinStandardAmt:    tt.fields.MinStandardAmt,
				Deadline:          tt.fields.Deadline,
				Sender:            tt.fields.Sender,
			}
			if err := msg.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("MsgRemoveLiquidity.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func buildCoin(denom string, amt int64) sdk.Coin {
	return sdk.Coin{
		Denom:  denom,
		Amount: sdk.NewInt(amt),
	}
}
