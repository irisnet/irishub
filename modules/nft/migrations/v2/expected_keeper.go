package v2

import sdk "github.com/cosmos/cosmos-sdk/types"

type NFTKeeper interface {
	SaveDenom(ctx sdk.Context, id,
		name,
		schema,
		symbol string,
		creator sdk.AccAddress,
		mintRestricted,
		updateRestricted bool,
		description,
		uri,
		uriHash,
		data string,
	) error

	SaveNFT(ctx sdk.Context, denomID,
		tokenID,
		tokenNm,
		tokenURI,
		tokenUriHash,
		tokenData string,
		receiver sdk.AccAddress,
	) error
}
