package v1beta1

import (
	"math/big"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/gogoproto/proto"
	"gopkg.in/yaml.v2"

	tokentypes "mods.irisnet.org/modules/token/types"
)

var (
	_      proto.Message = &Token{}
	tenInt               = big.NewInt(10)
)

// TokenI defines an interface for Token
type TokenI interface {
	GetSymbol() string
	GetName() string
	GetScale() uint32
	GetMinUnit() string
	GetInitialSupply() uint64
	GetMaxSupply() uint64
	GetMintable() bool
	GetOwner() sdk.AccAddress

	ToMainCoin(coin sdk.Coin) (sdk.DecCoin, error)
	ToMinCoin(coin sdk.DecCoin) (sdk.Coin, error)
}

// NewToken constructs a new Token instance
func NewToken(
	symbol string,
	name string,
	minUnit string,
	scale uint32,
	initialSupply,
	maxSupply uint64,
	mintable bool,
	owner sdk.AccAddress,
) Token {
	if maxSupply == 0 {
		if mintable {
			maxSupply = tokentypes.MaximumMaxSupply
		} else {
			maxSupply = initialSupply
		}
	}

	return Token{
		Symbol:        symbol,
		Name:          name,
		MinUnit:       minUnit,
		Scale:         scale,
		InitialSupply: initialSupply,
		MaxSupply:     maxSupply,
		Mintable:      mintable,
		Owner:         owner.String(),
	}
}

// GetSymbol implements exported.TokenI
func (t Token) GetSymbol() string {
	return t.Symbol
}

// GetName implements exported.TokenI
func (t Token) GetName() string {
	return t.Name
}

// GetScale implements exported.TokenI
func (t Token) GetScale() uint32 {
	return t.Scale
}

// GetMinUnit implements exported.TokenI
func (t Token) GetMinUnit() string {
	return t.MinUnit
}

// GetInitialSupply implements exported.TokenI
func (t Token) GetInitialSupply() uint64 {
	return t.InitialSupply
}

// GetMaxSupply implements exported.TokenI
func (t Token) GetMaxSupply() uint64 {
	return t.MaxSupply
}

// GetMintable implements exported.TokenI
func (t Token) GetMintable() bool {
	return t.Mintable
}

// GetOwner implements exported.TokenI
func (t Token) GetOwner() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(t.Owner)
	return owner
}

func (t Token) String() string {
	bz, _ := yaml.Marshal(t)
	return string(bz)
}

// ToMainCoin returns the main denom coin from args
func (t Token) ToMainCoin(coin sdk.Coin) (sdk.DecCoin, error) {
	if t.Symbol != coin.Denom && t.MinUnit != coin.Denom {
		return sdk.NewDecCoinFromDec(
				coin.Denom,
				sdk.ZeroDec(),
			), errorsmod.Wrapf(
				tokentypes.ErrTokenNotExists,
				"token not match",
			)
	}

	if t.Symbol == coin.Denom {
		return sdk.NewDecCoin(coin.Denom, coin.Amount), nil
	}

	precision := new(big.Int).Exp(tenInt, big.NewInt(int64(t.Scale)), nil)
	// dest amount = src amount / 10^(scale)
	amount := sdk.NewDecFromInt(coin.Amount).Quo(sdk.NewDecFromBigInt(precision))
	return sdk.NewDecCoinFromDec(t.Symbol, amount), nil
}

// ToMinCoin returns the min denom coin from args
func (t Token) ToMinCoin(coin sdk.DecCoin) (newCoin sdk.Coin, err error) {
	if t.Symbol != coin.Denom && t.MinUnit != coin.Denom {
		return sdk.NewCoin(
				coin.Denom,
				sdk.ZeroInt(),
			), errorsmod.Wrapf(
				tokentypes.ErrTokenNotExists,
				"token not match",
			)
	}

	if t.MinUnit == coin.Denom {
		return sdk.NewCoin(coin.Denom, coin.Amount.TruncateInt()), nil
	}

	precision := new(big.Int).Exp(tenInt, big.NewInt(int64(t.Scale)), nil)
	// dest amount = src amount * 10^(dest scale)
	amount := coin.Amount.Mul(sdk.NewDecFromBigInt(precision))
	return sdk.NewCoin(t.MinUnit, amount.TruncateInt()), nil
}

// Validate checks if the given token is valid
func (t Token) Validate() error {
	if len(t.Owner) > 0 {
		if _, err := sdk.AccAddressFromBech32(t.Owner); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
		}
	}
	if err := tokentypes.ValidateName(t.Name); err != nil {
		return err
	}
	if err := tokentypes.ValidateSymbol(t.Symbol); err != nil {
		return err
	}
	if err := tokentypes.ValidateMinUnit(t.MinUnit); err != nil {
		return err
	}
	if err := tokentypes.ValidateInitialSupply(t.InitialSupply); err != nil {
		return err
	}
	if t.MaxSupply < t.InitialSupply {
		return errorsmod.Wrapf(
			tokentypes.ErrInvalidMaxSupply,
			"invalid token max supply %d, only accepts value [%d, %d]",
			t.MaxSupply,
			t.InitialSupply,
			uint64(tokentypes.MaximumMaxSupply),
		)
	}
	return tokentypes.ValidateScale(t.Scale)
}
