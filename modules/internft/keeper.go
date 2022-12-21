package internft

import (
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"

	nfttransfer "github.com/bianjieai/nft-transfer/types"

	irismodnftkeeper "github.com/irisnet/irismod/modules/nft/keeper"
	"github.com/irisnet/irismod/modules/nft/types"
)

// NewInterNftKeeper creates a new ics721 Keeper instance
func NewInterNftKeeper(cdc codec.Codec,
	k irismodnftkeeper.Keeper,
	ak AccountKeeper,
) InterNftKeeper {
	return InterNftKeeper{
		nk:  k.NFTkeeper(),
		cdc: cdc,
		ak:  ak,
		cmr: types.NewClassMetadataResolver(cdc, ak.GetModuleAddress),
		tmr: types.NewTokenMetadataResolver(cdc),
	}
}

// CreateOrUpdateClass implement the method of ICS721Keeper.CreateOrUpdateClass
func (ik InterNftKeeper) CreateOrUpdateClass(ctx sdk.Context,
	classID,
	classURI,
	classData string,
) error {
	var (
		metadata *codectypes.Any
		err      error
	)
	if len(classData) != 0 {
		metadata, err = ik.cmr.Decode(classData)
		if err != nil {
			return err
		}
	} else {
		var denomMetadata = &types.DenomMetadata{
			Creator:          ik.ak.GetModuleAddress(types.ModuleName).String(),
			MintRestricted:   true,
			UpdateRestricted: true,
		}

		metadata, err = codectypes.NewAnyWithValue(denomMetadata)
		if err != nil {
			return err
		}
	}

	if ik.nk.HasClass(ctx, classID) {
		return ik.nk.UpdateClass(ctx, nft.Class{
			Id:   classID,
			Uri:  classURI,
			Data: metadata,
		})
	}

	return ik.nk.SaveClass(ctx, nft.Class{
		Id:   classID,
		Uri:  classURI,
		Data: metadata,
	})
}

func (ik InterNftKeeper) Mint(ctx sdk.Context,
	classID,
	tokenID,
	tokenURI,
	tokenData string,
	receiver sdk.AccAddress,
) error {
	data, err := ik.tmr.Decode(tokenData)
	if err != nil {
		return err
	}
	return ik.nk.Mint(ctx, nft.NFT{
		ClassId: classID,
		Id:      tokenID,
		Uri:     tokenURI,
		Data:    data,
	}, receiver)
}

func (ik InterNftKeeper) Transfer(
	ctx sdk.Context,
	classID,
	tokenID,
	tokenData string,
	receiver sdk.AccAddress,
) error {
	if err := ik.nk.Transfer(ctx, classID, tokenID, receiver); err != nil {
		return err
	}
	if len(tokenData) == 0 {
		return nil
	}
	data, err := ik.tmr.Decode(tokenData)
	if err != nil {
		return err
	}
	nft, _ := ik.nk.GetNFT(ctx, classID, tokenID)
	nft.Data = data
	return ik.nk.Update(ctx, nft)
}

func (ik InterNftKeeper) GetClass(ctx sdk.Context, classID string) (nfttransfer.Class, bool) {
	class, exist := ik.nk.GetClass(ctx, classID)
	if !exist {
		return nil, false
	}

	classDataEncoded, err := ik.cmr.Encode(class.Data)
	if err != nil {
		ik.Logger(ctx).Error("encode class data failed")
		return nil, false
	}
	return InterClass{
		ID:   classID,
		URI:  class.Uri,
		Data: classDataEncoded,
	}, true
}

func (ik InterNftKeeper) GetNFT(ctx sdk.Context, classID, tokenID string) (nfttransfer.NFT, bool) {
	nft, has := ik.nk.GetNFT(ctx, classID, tokenID)
	if !has {
		return nil, false
	}
	nftDataEncoded, err := ik.tmr.Encode(nft.Data)
	if err != nil {
		ik.Logger(ctx).Error("encode nft data failed")
		return nil, false
	}
	return InterToken{
		ClassID: classID,
		ID:      tokenID,
		URI:     nft.Uri,
		Data:    nftDataEncoded,
	}, true
}

func (ik InterNftKeeper) Burn(ctx sdk.Context, classID string, tokenID string) error {
	return ik.nk.Burn(ctx, classID, tokenID)
}

func (ik InterNftKeeper) GetOwner(ctx sdk.Context, classID string, tokenID string) sdk.AccAddress {
	return ik.nk.GetOwner(ctx, classID, tokenID)
}

func (ik InterNftKeeper) HasClass(ctx sdk.Context, classID string) bool {
	return ik.nk.HasClass(ctx, classID)
}

// Logger returns a module-specific logger.
func (ik InterNftKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "ics721/NFTKeeper")
}
