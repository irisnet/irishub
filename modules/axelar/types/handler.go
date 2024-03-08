package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GeneralMessageHandler handles messages from axelar
type GeneralMessageHandler interface {
	HandleGeneralMessage(ctx sdk.Context, srcChain, srcAddress string, destAddress string, payload []byte) error
	HandleGeneralMessageWithToken(ctx sdk.Context, srcChain, srcAddress string, destAddress string, payload []byte, coin sdk.Coin) error
}

// EmptyHandler implements the GeneralMessageHandler
type EmptyHandler struct{}

// HandleGeneralMessage implements the GeneralMessageHandler
func (h EmptyHandler) HandleGeneralMessage(ctx sdk.Context, srcChain, srcAddress string, destAddress string, payload []byte) error {
	ctx.Logger().Info("HandleGeneralMessage called",
		"srcChain", srcChain,
		"srcAddress", srcAddress,
		"destAddress", destAddress,
		"payload", payload,
		"module", "axelar",
	)
	return nil
}

// HandleGeneralMessageWithToken implements the GeneralMessageHandler
func (h EmptyHandler) HandleGeneralMessageWithToken(ctx sdk.Context, srcChain, srcAddress string, destAddress string, payload []byte, coin sdk.Coin) error {
	ctx.Logger().Info("HandleGeneralMessageWithToken called",
		"srcChain", srcChain,
		"srcAddress", srcAddress,
		"destAddress", destAddress,
		"payload", payload,
		"coin", coin,
		"module", "axelar",
	)
	return nil
}
