package client

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"fmt"
	"encoding/json"
)

func Test_getValidator(t *testing.T) {
	address, _ := sdk.AccAddressFromBech32("faa1r4dnf8lnakw743dwhd4nnpxatcx5v40n0vntc6")
	//ctx := app.NewContext().Ctx
	//cdc := app.MakeCodec()


	type args struct {
		address sdk.AccAddress
		ctx     context.CoreContext
		cdc     *wire.Codec
	}
	tests := []struct {
		name    string
		args    args
	}{
		{
			name: "test get validator",
			args: args{
				address: address,
				//ctx: ctx,
				//cdc: cdc,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//res, err := getValidator(tt.args.address, tt.args.ctx, tt.args.cdc)
			res := ExRateResponse{
				ExRate: 1.0,
			}
			resStr, err := json.Marshal(res)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(resStr))
		})
	}
}
