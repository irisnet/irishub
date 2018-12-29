package slashing

import (
	"time"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
	sdk "github.com/irisnet/irishub/types"
	"fmt"
	"strconv"
)

var _ params.ParamSet = (*Params)(nil)

// Default parameter namespace
const (
	DefaultParamspace = "slashing"
)

// Parameter store key
var (
	KeyMaxEvidenceAge          = []byte("MaxEvidenceAge")
	KeySignedBlocksWindow      = []byte("SignedBlocksWindow")
	KeyMinSignedPerWindow      = []byte("MinSignedPerWindow")
	KeyDoubleSignJailDuration  = []byte("DoubleSignJailDuration")
	KeyDowntimeJailDuration    = []byte("DowntimeJailDuration")
	KeySlashFractionDoubleSign = []byte("SlashFractionDoubleSign")
	KeySlashFractionDowntime   = []byte("SlashFractionDowntime")
)

// ParamTypeTable for slashing module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for slashing at genesis
type Params struct {
	MaxEvidenceAge          time.Duration `json:"max-evidence-age"`
	SignedBlocksWindow      int64         `json:"signed-blocks-window"`
	MinSignedPerWindow      sdk.Dec       `json:"min-signed-per-window"`
	DoubleSignJailDuration  time.Duration `json:"double-sign-unbond-duration"`
	DowntimeJailDuration    time.Duration `json:"downtime-unbond-duration"`
	SlashFractionDoubleSign sdk.Dec       `json:"slash-fraction-double-sign"`
	SlashFractionDowntime   sdk.Dec       `json:"slash-fraction-downtime"`
}

// Implements params.ParamStruct
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{KeyMaxEvidenceAge, &p.MaxEvidenceAge},
		{KeySignedBlocksWindow, &p.SignedBlocksWindow},
		{KeyMinSignedPerWindow, &p.MinSignedPerWindow},
		{KeyDoubleSignJailDuration, &p.DoubleSignJailDuration},
		{KeyDowntimeJailDuration, &p.DowntimeJailDuration},
		{KeySlashFractionDoubleSign, &p.SlashFractionDoubleSign},
		{KeySlashFractionDowntime, &p.SlashFractionDowntime},
	}
}

func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(KeyMaxEvidenceAge):
		maxEvidenceAge, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateMaxEvidenceAge(p.MaxEvidenceAge); err != nil {
			return nil, err
		}
		return maxEvidenceAge, nil
	case string(KeySignedBlocksWindow):
		signedBlocksWindow, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateSignedBlocksWindow(signedBlocksWindow); err != nil {
			return nil, err
		}
		return signedBlocksWindow, nil
	case string(KeyMinSignedPerWindow):
		minSignedPerWindow, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateDecimalRange(minSignedPerWindow); err != nil {
			return nil, err
		}
		return minSignedPerWindow, nil
	case string(KeyDoubleSignJailDuration):
		doubleSignJailDuration, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateJailDuration(doubleSignJailDuration); err != nil {
			return nil, err
		}
		return doubleSignJailDuration, nil
	case string(KeyDowntimeJailDuration):
		downtimeJailDuration, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateJailDuration(downtimeJailDuration); err != nil {
			return nil, err
		}
		return downtimeJailDuration, nil
	case string(KeySlashFractionDoubleSign):
		slashFractionDoubleSign, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateDecimalRange(slashFractionDoubleSign); err != nil {
			return nil, err
		}
		return slashFractionDoubleSign, nil
	case string(KeySlashFractionDowntime):
		slashFractionDowntime, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateDecimalRange(slashFractionDowntime); err != nil {
			return nil, err
		}
		return slashFractionDowntime, nil
	default:
		return nil, sdk.NewError(params.DefaultCodespace, params.CodeInvalidKey, fmt.Sprintf("%s is not found", key))
	}
}

func (p *Params) GetParamSpace() string {
	return DefaultParamspace
}

func (p *Params) StringFromBytes(cdc *codec.Codec, key string, bytes []byte) (string, error) {
	switch key {
	case string(KeyMaxEvidenceAge):
		err := cdc.UnmarshalJSON(bytes, &p.MaxEvidenceAge)
		return p.MaxEvidenceAge.String(), err
	case string(KeySignedBlocksWindow):
		err := cdc.UnmarshalJSON(bytes, &p.SignedBlocksWindow)
		return strconv.FormatInt(p.SignedBlocksWindow, 10), err
	case string(KeyMinSignedPerWindow):
		err := cdc.UnmarshalJSON(bytes, &p.MinSignedPerWindow)
		return p.MinSignedPerWindow.String(), err
	case string(KeyDoubleSignJailDuration):
		err := cdc.UnmarshalJSON(bytes, &p.DoubleSignJailDuration)
		return p.DoubleSignJailDuration.String(), err
	case string(KeyDowntimeJailDuration):
		err := cdc.UnmarshalJSON(bytes, &p.DowntimeJailDuration)
		return p.DowntimeJailDuration.String(), err
	case string(KeySlashFractionDoubleSign):
		err := cdc.UnmarshalJSON(bytes, &p.SlashFractionDoubleSign)
		return p.SlashFractionDoubleSign.String(), err
	case string(KeySlashFractionDowntime):
		err := cdc.UnmarshalJSON(bytes, &p.SlashFractionDowntime)
		return p.SlashFractionDowntime.String(), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

// Default parameters used by Iris Hub
func DefaultParams() Params {
	return Params{
		MaxEvidenceAge: sdk.Day,
		DoubleSignJailDuration: sdk.Week,
		SignedBlocksWindow: 20000,
		DowntimeJailDuration: sdk.Week,
		MinSignedPerWindow: sdk.NewDecWithPrec(5, 1),
		SlashFractionDoubleSign: sdk.NewDec(1).Quo(sdk.NewDec(20)),
		SlashFractionDowntime: sdk.NewDec(1).Quo(sdk.NewDec(100)),
	}
}

func DefaultParamsForTestnet() Params {
	return Params{
		MaxEvidenceAge: 60 * 2 * time.Second,
		DoubleSignJailDuration: 60 * 5 * time.Second,
		SignedBlocksWindow: 100,
		DowntimeJailDuration: 60 * 10 * time.Second,
		MinSignedPerWindow: sdk.NewDecWithPrec(5, 1),
		SlashFractionDoubleSign: sdk.NewDec(1).Quo(sdk.NewDec(20)),
		SlashFractionDowntime: sdk.NewDec(1).Quo(sdk.NewDec(100)),
	}
}

func validateParams(p Params) sdk.Error {
	if err := validateMaxEvidenceAge(p.MaxEvidenceAge); err != nil {
		return err
	}
	if err := validateJailDuration(p.DoubleSignJailDuration); err != nil {
		return err
	}
	if err := validateJailDuration(p.DowntimeJailDuration); err != nil {
		return err
	}
	if err := validateSignedBlocksWindow(p.SignedBlocksWindow); err != nil {
		return err
	}
	if err := validateDecimalRange(p.MinSignedPerWindow); err != nil {
		return err
	}
	if err := validateDecimalRange(p.SlashFractionDoubleSign); err != nil {
		return err
	}
	if err := validateDecimalRange(p.SlashFractionDowntime); err != nil {
		return err
	}

	return nil
}

func validateMaxEvidenceAge(p time.Duration) sdk.Error {
	if p < 2 * time.Minute || p > sdk.Week {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash MaxEvidenceAge [%s] should be between [10min, 1week] ", p.String()))
	}
	return nil
}

func validateJailDuration(p time.Duration) sdk.Error {
	if p <= 0 || p >= 4 * sdk.Week {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash DoubleSignJailDuration and DowntimeJailDuration [%s] should be between (0, 4week) ", p.String()))
	}
	return nil
}

func validateSignedBlocksWindow(p int64) sdk.Error {
	if p < 100 || p > 140000 {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash SignedBlocksWindow [%d] should be between [100, 140000] ", p))
	}
	return nil
}

func validateDecimalRange(p sdk.Dec) sdk.Error {
	if p.LTE(sdk.NewDec(0)) || p.GTE(sdk.NewDec(1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash MinSignedPerWindow/SlashFractionDoubleSign/SlashFractionDowntime [%s] should be between (0, 1) ", p))
	}
	return nil
}

//______________________________________________________________________

// get inflation params from the global param store
func (k Keeper) GetParamSet(ctx sdk.Context) Params {
	var params Params
	k.paramspace.GetParamSet(ctx, &params)
	return params
}

// set inflation params from the global param store
func (k Keeper) SetParamSet(ctx sdk.Context, params Params) {
	k.paramspace.SetParamSet(ctx, &params)
}

// MaxEvidenceAge - Max age for evidence - 21 days (3 weeks)
// MaxEvidenceAge = 60 * 60 * 24 * 7 * 3
func (k Keeper) MaxEvidenceAge(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, KeyMaxEvidenceAge, &res)
	return
}

// SignedBlocksWindow - sliding window for downtime slashing
func (k Keeper) SignedBlocksWindow(ctx sdk.Context) (res int64) {
	k.paramspace.Get(ctx, KeySignedBlocksWindow, &res)
	return
}

// Downtime slashing thershold - default 50% of the SignedBlocksWindow
func (k Keeper) MinSignedPerWindow(ctx sdk.Context) int64 {
	var minSignedPerWindow sdk.Dec
	k.paramspace.Get(ctx, KeyMinSignedPerWindow, &minSignedPerWindow)
	signedBlocksWindow := k.SignedBlocksWindow(ctx)
	return sdk.NewDec(signedBlocksWindow).Mul(minSignedPerWindow).RoundInt64()
}

// Double-sign unbond duration
func (k Keeper) DoubleSignUnbondDuration(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, KeyDoubleSignJailDuration, &res)
	return
}

// Downtime unbond duration
func (k Keeper) DowntimeUnbondDuration(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, KeyDowntimeJailDuration, &res)
	return
}

// SlashFractionDoubleSign - currently default 5%
func (k Keeper) SlashFractionDoubleSign(ctx sdk.Context) (res sdk.Dec) {
	k.paramspace.Get(ctx, KeySlashFractionDoubleSign, &res)
	return
}

// SlashFractionDowntime - currently default 1%
func (k Keeper) SlashFractionDowntime(ctx sdk.Context) (res sdk.Dec) {
	k.paramspace.Get(ctx, KeySlashFractionDowntime, &res)
	return
}
