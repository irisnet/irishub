package types

import (
	"encoding/json"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

func TestTokenBuilder_Build(t *testing.T) {
	nftMetadata := NFTMetadata{
		Name: "My Cat",
		Data: "{\"key1\":\"value1\",\"key2\":\"value2\"}",
	}

	bz, err := json.Marshal(nftMetadata)
	require.NoError(t, err, " nftMetadata json.Marshal failed")
	t.Logf("%s", bz)

	any, err := codectypes.NewAnyWithValue(&nftMetadata)
	require.NoError(t, err, " nftMetadata codectypes.NewAnyWithValue failed")

	token := nft.NFT{
		ClassId: "kitty",
		Id:      "cat",
		Uri:     "uri",
		UriHash: "uri_hash",
		Data:    any,
	}

	cdc := GetEncoding()
	bz, err = cdc.MarshalJSON(&token)
	require.NoError(t, err, " token MarshalJSON failed")
	t.Logf("%s", bz)

	builder := NewTokenBuilder(cdc)
	result, err := builder.BuildMetadata(token)
	require.NoError(t, err, " token builder.BuildMetadata failed")
	t.Log(result)

	expToken, err := builder.Build(token.ClassId, token.Id, token.Uri, result)
	require.NoError(t, err, " token builder.Build failed")

	exp, err := cdc.MarshalInterfaceJSON(&token)
	require.NoError(t, err, " token cdc.MarshalInterfaceJSON failed")
	t.Logf("%s", exp)

	act, err := cdc.MarshalInterfaceJSON(&expToken)
	require.NoError(t, err, " token cdc.MarshalInterfaceJSON failed")
	t.Logf("%s", act)

	require.Equal(t, act, exp, "not equal")
}

func GetEncoding() codec.Codec {
	interfaceRegistry := types.NewInterfaceRegistry()
	interfaceRegistry.RegisterImplementations(
		(*proto.Message)(nil),
		&nft.Class{},
		&nft.NFT{},
		&DenomMetadata{},
		&NFTMetadata{},
	)
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	return marshaler
}

func TestClassBuilder_BuildMetadata(t *testing.T) {
	creator, err := sdk.AccAddressFromHexUnsafe(crypto.AddressHash([]byte("test_consumer")).String())
	require.NoError(t, err, "AccAddressFromHexUnsafe failed")

	cdc := GetEncoding()
	getModuleAddress := func(_ string) sdk.AccAddress {
		return creator
	}
	class := nft.Class{
		Name:        "kitty",
		Symbol:      "symbol",
		Description: "digital cat",
		Uri:         "uri",
		UriHash:     "uri_hash",
	}
	denomMetadata := DenomMetadata{
		Creator:          creator.String(),
		Schema:           "{}",
		MintRestricted:   true,
		UpdateRestricted: true,
	}

	type args struct {
		classID   string
		classData string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty classData",
			args: args{
				classID:   "cat",
				classData: "",
			},
			want:    `{"irismod:creator":{"value":"f8a9eee6bce5bc043e5feec2baef355f87dbfcdf"},"irismod:description":{"value":"digital cat"},"irismod:mint_restricted":{"value":true},"irismod:name":{"value":"kitty"},"irismod:schema":{"value":"{}"},"irismod:symbol":{"value":"symbol"},"irismod:update_restricted":{"value":true},"irismod:uri_hash":{"value":"uri_hash"}}`,
			wantErr: false,
		},
		{
			name: "classData is invalid json string",
			args: args{
				classID:   "cat",
				classData: "hhaahha",
			},
			want:    `{"irismod:creator":{"value":"f8a9eee6bce5bc043e5feec2baef355f87dbfcdf"},"irismod:description":{"value":"digital cat"},"irismod:mint_restricted":{"value":true},"irismod:name":{"value":"kitty"},"irismod:schema":{"value":"{}"},"irismod:symbol":{"value":"symbol"},"irismod:update_restricted":{"value":true},"irismod:uri_hash":{"value":"uri_hash"}}`,
			wantErr: false,
		},
		{
			name: "classData is valid json string",
			args: args{
				classID:   "cat",
				classData: "{\"key1\":\"value1\",\"key2\":\"value2\"}",
			},
			want:    `{"irismod:creator":{"value":"f8a9eee6bce5bc043e5feec2baef355f87dbfcdf"},"irismod:description":{"value":"digital cat"},"irismod:mint_restricted":{"value":true},"irismod:name":{"value":"kitty"},"irismod:schema":{"value":"{}"},"irismod:symbol":{"value":"symbol"},"irismod:update_restricted":{"value":true},"irismod:uri_hash":{"value":"uri_hash"},"key1":"value1","key2":"value2"}`,
			wantErr: false,
		},
		{
			name: "class is IBC assets and classData is invalid json string",
			args: args{
				classID:   "ibc/943B966B2B8A53C50A198EDAB7C9A41FCEAF24400A94167846679769D8BF8311",
				classData: "hahhahha",
			},
			want:    `hahhahha`,
			wantErr: false,
		},
		{
			name: "class is IBC assets and classData is valid json string",
			args: args{
				classID:   "ibc/943B966B2B8A53C50A198EDAB7C9A41FCEAF24400A94167846679769D8BF8311",
				classData: "{\"key1\":\"value1\",\"key2\":\"value2\"}",
			},
			want:    `{"irismod:creator":{"value":"f8a9eee6bce5bc043e5feec2baef355f87dbfcdf"},"irismod:description":{"value":"digital cat"},"irismod:mint_restricted":{"value":true},"irismod:name":{"value":"kitty"},"irismod:schema":{"value":"{}"},"irismod:symbol":{"value":"symbol"},"irismod:update_restricted":{"value":true},"irismod:uri_hash":{"value":"uri_hash"},"key1":"value1","key2":"value2"}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := ClassBuilder{
				cdc:              cdc,
				getModuleAddress: getModuleAddress,
			}

			denomMetadata.Data = tt.args.classData
			any, err := codectypes.NewAnyWithValue(&denomMetadata)
			require.NoError(t, err, " denomMetadata codectypes.NewAnyWithValue failed")

			class.Id = tt.args.classID
			class.Data = any

			got, err := cb.BuildMetadata(class)
			want := Base64.EncodeToString([]byte(tt.want))
			if (err != nil) != tt.wantErr {
				t.Errorf("ClassBuilder.BuildMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != want {
				t.Errorf("ClassBuilder.BuildMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClassBuilder_Build(t *testing.T) {
	creator, err := sdk.AccAddressFromHexUnsafe(crypto.AddressHash([]byte("test_consumer")).String())
	require.NoError(t, err, "AccAddressFromHexUnsafe failed")

	cdc := GetEncoding()
	getModuleAddress := func(_ string) sdk.AccAddress {
		return creator
	}
	classID := "cat"
	classURI := "uri"

	type args struct {
		classData string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty classData",
			args: args{
				classData: `{"irismod:creator":{"value":"f8a9eee6bce5bc043e5feec2baef355f87dbfcdf"},"irismod:description":{"value":"digital cat"},"irismod:mint_restricted":{"value":true},"irismod:name":{"value":"kitty"},"irismod:schema":{"value":"{}"},"irismod:symbol":{"value":"symbol"},"irismod:update_restricted":{"value":true},"irismod:uri_hash":{"value":"uri_hash"}}`,
			},
			want:    `{"@type":"/cosmos.nft.v1beta1.Class","id":"cat","name":"kitty","symbol":"symbol","description":"digital cat","uri":"uri","uri_hash":"uri_hash","data":{"@type":"/irismod.nft.DenomMetadata","creator":"cosmos1lz57ae4uuk7qg0jlampt4me4t7rahlxl5pnn3y","schema":"{}","mint_restricted":true,"update_restricted":true,"data":""}}`,
			wantErr: false,
		},
		{
			name: "classData is invalid json string",
			args: args{
				classData: `this is empty class data`,
			},
			want:    `{"@type":"/cosmos.nft.v1beta1.Class","id":"cat","name":"","symbol":"","description":"","uri":"uri","uri_hash":"","data":{"@type":"/irismod.nft.DenomMetadata","creator":"cosmos1lz57ae4uuk7qg0jlampt4me4t7rahlxl5pnn3y","schema":"","mint_restricted":true,"update_restricted":true,"data":"this is empty class data"}}`,
			wantErr: false,
		},
		{
			name: "classData is valid json string",
			args: args{
				classData: `{"irismod:creator":{"value":"f8a9eee6bce5bc043e5feec2baef355f87dbfcdf"},"irismod:description":{"value":"digital cat"},"irismod:mint_restricted":{"value":true},"irismod:name":{"value":"kitty"},"irismod:schema":{"value":"{}"},"irismod:symbol":{"value":"symbol"},"irismod:update_restricted":{"value":true},"irismod:uri_hash":{"value":"uri_hash"},"key1":"value1","key2":"value2"}`,
			},
			want:    `{"@type":"/cosmos.nft.v1beta1.Class","id":"cat","name":"kitty","symbol":"symbol","description":"digital cat","uri":"uri","uri_hash":"uri_hash","data":{"@type":"/irismod.nft.DenomMetadata","creator":"cosmos1lz57ae4uuk7qg0jlampt4me4t7rahlxl5pnn3y","schema":"{}","mint_restricted":true,"update_restricted":true,"data":"{\"key1\":\"value1\",\"key2\":\"value2\"}"}}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cb := ClassBuilder{
				cdc:              cdc,
				getModuleAddress: getModuleAddress,
			}
			classDataRaw := Base64.EncodeToString([]byte(tt.args.classData))
			result, err := cb.Build(classID, classURI, classDataRaw)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClassBuilder.BuildMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got, err := cdc.MarshalInterfaceJSON(&result)
			require.NoError(t, err, " class cdc.MarshalInterfaceJSON failed")

			if string(got) != tt.want {
				t.Errorf("ClassBuilder.BuildMetadata() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
