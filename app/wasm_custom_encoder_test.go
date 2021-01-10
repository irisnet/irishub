package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"

	nfttypes "github.com/irisnet/irismod/modules/nft/types"
	tokentypes "github.com/irisnet/irismod/modules/token/types"
)

func TestWasmCustomEncoder_Encode(t *testing.T) {
	encodingConfig := simapp.MakeTestEncodingConfig()
	registry := encodingConfig.InterfaceRegistry

	encoder := WasmCustomEncoder{registry}

	tokentypes.RegisterInterfaces(registry)
	nfttypes.RegisterInterfaces(registry)

	msgMintToken := &tokentypes.MsgMintToken{
		Symbol: "uiris",
		Amount: 1000,
		To:     "iaa1xl9q8mpgm5mr3f0vlha3mm94qj66w83608m8e6",
		Owner:  "iaa1xl9q8mpgm5mr3f0vlha3mm94qj66w83608m8e6",
	}

	jsonMsgMintToken, err := json.Marshal(msgMintToken)
	require.NoError(t, err)

	msg, err := json.Marshal(MsgWasmCustom{
		Router: fmt.Sprintf("/%s", proto.MessageName(msgMintToken)),
		Data:   base64.StdEncoding.EncodeToString(jsonMsgMintToken),
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		data    json.RawMessage
		want    []sdk.Msg
		wantErr bool
	}{
		{name: "Mint Token Msg", want: []sdk.Msg{msgMintToken}, wantErr: false, data: msg},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encoder.Encode(nil, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("WasmCustomEncoder.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WasmCustomEncoder.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
