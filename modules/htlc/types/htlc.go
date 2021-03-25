package types

import (
	fmt "fmt"
	time "time"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHTLC constructs a new HTLC instance
func NewHTLC(
	id tmbytes.HexBytes,
	sender sdk.AccAddress,
	to sdk.AccAddress,
	receiverOnOtherChain string,
	senderOnOtherChain string,
	amount sdk.Coins,
	hashLock tmbytes.HexBytes,
	secret tmbytes.HexBytes,
	timestamp uint64,
	expirationHeight uint64,
	state HTLCState,
	closedBlock uint64,
	transfer bool,
	direction SwapDirection,
) HTLC {
	return HTLC{
		Id:                   id.String(),
		Sender:               sender.String(),
		To:                   to.String(),
		ReceiverOnOtherChain: receiverOnOtherChain,
		SenderOnOtherChain:   senderOnOtherChain,
		Amount:               amount,
		HashLock:             hashLock.String(),
		Secret:               secret.String(),
		Timestamp:            timestamp,
		ExpirationHeight:     expirationHeight,
		State:                state,
		ClosedBlock:          closedBlock,
		Transfer:             transfer,
		Direction:            direction,
	}
}

// Validate validates the HTLC
func (h HTLC) Validate() error {
	if err := ValidateID(h.Id); err != nil {
		return err
	}
	if err := ValidateHashLock(h.HashLock); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(h.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(h.To); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid recipient address (%s)", err)
	}
	if err := ValidateReceiverOnOtherChain(h.ReceiverOnOtherChain); err != nil {
		return err
	}
	if err := ValidateSenderOnOtherChain(h.SenderOnOtherChain); err != nil {
		return err
	}
	if h.ExpirationHeight == 0 {
		return sdkerrors.Wrapf(ErrInvalidExpirationHeight, "expire height cannot be 0")
	}
	if h.Timestamp == 0 {
		return sdkerrors.Wrapf(ErrInvalidTimestamp, "timestamp cannot be 0")
	}
	if err := ValidateAmount(h.Transfer, h.Amount); err != nil {
		return err
	}
	if h.State > Refunded {
		return sdkerrors.Wrapf(ErrInvalidState, "invalid htlc status")
	}
	if h.State == Completed && h.ClosedBlock == 0 {
		return sdkerrors.Wrapf(ErrInvalidClosedBlock, "closed block cannot be 0")
	}
	if !h.Transfer && h.Direction != 0 {
		return sdkerrors.Wrapf(ErrInvalidDirection, "invalid htlc direction")
	}
	if h.Transfer && (h.Direction < Incoming || h.Direction > Outgoing) {
		return sdkerrors.Wrapf(ErrInvalidDirection, "invalid htlt direction")
	}
	if h.State != Completed && len(h.Secret) > 0 {
		return sdkerrors.Wrapf(ErrInvalidSecret, "secret must be empty when the HTLC has not be claimed")
	}
	if h.State == Completed && len(h.Secret) != SecretLength {
		return sdkerrors.Wrapf(ErrInvalidSecret, "length of the secret must be %d", SecretLength)
	}
	return nil
}

// NewAssetSupply constructs a new AssetSupply instance
func NewAssetSupply(
	incomingSupply sdk.Coin,
	outgoingSupply sdk.Coin,
	currentSupply sdk.Coin,
	timeLimitedCurrentSupply sdk.Coin,
	timeElapsed time.Duration,
) AssetSupply {
	return AssetSupply{
		IncomingSupply:           incomingSupply,
		OutgoingSupply:           outgoingSupply,
		CurrentSupply:            currentSupply,
		TimeLimitedCurrentSupply: timeLimitedCurrentSupply,
		TimeElapsed:              timeElapsed,
	}
}

// DefaultAssetSupplies gets the raw asset supplies for testing
func DefaultAssetSupplies() []AssetSupply {
	return []AssetSupply{}
}

// Validate performs a basic validation of an asset supply fields.
func (a AssetSupply) Validate() error {
	if !a.IncomingSupply.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "incoming supply %s", a.IncomingSupply)
	}
	if !a.OutgoingSupply.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "outgoing supply %s", a.OutgoingSupply)
	}
	if !a.CurrentSupply.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "current supply %s", a.CurrentSupply)
	}
	if !a.TimeLimitedCurrentSupply.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "time-limited current supply %s", a.CurrentSupply)
	}
	denom := a.CurrentSupply.Denom
	if (a.IncomingSupply.Denom != denom) || (a.OutgoingSupply.Denom != denom) || (a.TimeLimitedCurrentSupply.Denom != denom) {
		return fmt.Errorf("asset supply denoms do not match %s %s %s %s", a.CurrentSupply.Denom, a.IncomingSupply.Denom, a.OutgoingSupply.Denom, a.TimeLimitedCurrentSupply.Denom)
	}
	return nil
}

// GetHashLock calculates the hash lock from the given secret and timestamp
func GetHashLock(secret tmbytes.HexBytes, timestamp uint64) []byte {
	if timestamp > 0 {
		return tmhash.Sum(append(secret, sdk.Uint64ToBigEndian(timestamp)...))
	}
	return tmhash.Sum(secret)
}

func GetID(
	sender sdk.AccAddress,
	to sdk.AccAddress,
	amount sdk.Coins,
	hashLock tmbytes.HexBytes,
) tmbytes.HexBytes {
	return tmhash.Sum(
		append(append(append(hashLock, sender...), to...), []byte(amount.Sort().String())...),
	)
}
