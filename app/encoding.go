package app

import (
	enccodec "github.com/evmos/ethermint/encoding/codec"

	"github.com/irisnet/irishub/v3/app/params"
)

// RegisterEncodingConfig registers concrete types on codec
func RegisterEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	enccodec.RegisterLegacyAminoCodec(encodingConfig.Amino)
	enccodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
