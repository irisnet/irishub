package ante

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256r1"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/evmos/ethermint/crypto/ethsecp256k1"
)

const (
	Secp256k1VerifyCost uint64 = 21000
)

// DefaultSigVerificationGasConsumer is the default implementation of SignatureVerificationGasConsumer. It consumes gas
// for signature verification based upon the public key type. The cost is fetched from the given params and is matched
// by the concrete type.
func DefaultSigVerificationGasConsumer(
	meter sdk.GasMeter, sig signing.SignatureV2, params types.Params,
) error {
	pubkey := sig.PubKey
	switch pubkey := pubkey.(type) {
	case *ed25519.PubKey:
		meter.ConsumeGas(params.SigVerifyCostED25519, "ante verify: ed25519")
		return sdkerrors.Wrap(sdkerrors.ErrInvalidPubKey, "ED25519 public keys are unsupported")

	case *secp256k1.PubKey:
		meter.ConsumeGas(params.SigVerifyCostSecp256k1, "ante verify: secp256k1")
		return nil

	case *ethsecp256k1.PubKey:
		// Ethereum keys
		meter.ConsumeGas(Secp256k1VerifyCost, "ante verify: eth_secp256k1")
		return nil

	case *secp256r1.PubKey:
		meter.ConsumeGas(params.SigVerifyCostSecp256r1(), "ante verify: secp256r1")
		return nil

	case multisig.PubKey:
		multisignature, ok := sig.Data.(*signing.MultiSignatureData)
		if !ok {
			return fmt.Errorf("expected %T, got, %T", &signing.MultiSignatureData{}, sig.Data)
		}
		err := ConsumeMultisignatureVerificationGas(meter, multisignature, pubkey, params, sig.Sequence)
		if err != nil {
			return err
		}
		return nil

	default:
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidPubKey, "unrecognized public key type: %T", pubkey)
	}
}

// ConsumeMultisignatureVerificationGas consumes gas from a GasMeter for verifying a multisig pubkey signature
func ConsumeMultisignatureVerificationGas(
	meter sdk.GasMeter, sig *signing.MultiSignatureData, pubkey multisig.PubKey,
	params types.Params, accSeq uint64,
) error {
	size := sig.BitArray.Count()
	sigIndex := 0

	for i := 0; i < size; i++ {
		if !sig.BitArray.GetIndex(i) {
			continue
		}
		sigV2 := signing.SignatureV2{
			PubKey:   pubkey.GetPubKeys()[i],
			Data:     sig.Signatures[sigIndex],
			Sequence: accSeq,
		}
		err := DefaultSigVerificationGasConsumer(meter, sigV2, params)
		if err != nil {
			return err
		}
		sigIndex++
	}

	return nil
}
