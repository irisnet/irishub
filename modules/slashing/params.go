package slashing

import (
	"fmt"
	"strconv"
	"time"

	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/params"
	sdk "github.com/irisnet/irishub/types"
)

var _ params.ParamSet = (*Params)(nil)

// Default parameter namespace
const (
	DefaultParamspace = "slashing"
	BlocksPerMinute   = 12   // 5 seconds a block
	BlocksPerDay      = BlocksPerMinute * 60 * 24   // 17280
)

// Parameter store key
var (
	KeyMaxEvidenceAge          = []byte("MaxEvidenceAge")
	KeySignedBlocksWindow      = []byte("SignedBlocksWindow")
	KeyMinSignedPerWindow      = []byte("MinSignedPerWindow")
	KeyDoubleSignJailDuration  = []byte("DoubleSignJailDuration")
	KeyDowntimeJailDuration    = []byte("DowntimeJailDuration")
	KeyCensorshipJailDuration  = []byte("CensorshipJailDuration")
	KeySlashFractionDoubleSign = []byte("SlashFractionDoubleSign")
	KeySlashFractionDowntime   = []byte("SlashFractionDowntime")
	KeySlashFractionCensorship = []byte("SlashFractionCensorship")
)

// ParamTypeTable for slashing module
func ParamTypeTable() params.TypeTable {
	return params.NewTypeTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for slashing at genesis
type Params struct {
	MaxEvidenceAge          int64         `json:"max_evidence_age"`
	SignedBlocksWindow      int64         `json:"signed_blocks_window"`
	MinSignedPerWindow      sdk.Dec       `json:"min_signed_per_window"`
	DoubleSignJailDuration  time.Duration `json:"double_sign_jail_duration"`
	DowntimeJailDuration    time.Duration `json:"downtime_jail_duration"`
	CensorshipJailDuration  time.Duration `json:"censorship_jail_duration"`
	SlashFractionDoubleSign sdk.Dec       `json:"slash_fraction_double_sign"`
	SlashFractionDowntime   sdk.Dec       `json:"slash_fraction_downtime"`
	SlashFractionCensorship sdk.Dec       `json:"slash_fraction_censorship"`
}

// Implements params.ParamStruct
func (p *Params) KeyValuePairs() params.KeyValuePairs {
	return params.KeyValuePairs{
		{KeyMaxEvidenceAge, &p.MaxEvidenceAge},
		{KeySignedBlocksWindow, &p.SignedBlocksWindow},
		{KeyMinSignedPerWindow, &p.MinSignedPerWindow},
		{KeyDoubleSignJailDuration, &p.DoubleSignJailDuration},
		{KeyDowntimeJailDuration, &p.DowntimeJailDuration},
		{KeyCensorshipJailDuration, &p.CensorshipJailDuration},
		{KeySlashFractionDoubleSign, &p.SlashFractionDoubleSign},
		{KeySlashFractionDowntime, &p.SlashFractionDowntime},
		{KeySlashFractionCensorship, &p.SlashFractionCensorship},
	}
}

func (p *Params) Validate(key string, value string) (interface{}, sdk.Error) {
	switch key {
	case string(KeyMaxEvidenceAge):
		maxEvidenceAge, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateMaxEvidenceAge(maxEvidenceAge); err != nil {
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
		if err := validateMinSignedPerWindow(minSignedPerWindow); err != nil {
			return nil, err
		}
		return minSignedPerWindow, nil
	case string(KeyDoubleSignJailDuration):
		doubleSignJailDuration, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateDoubleSignJailDuration(doubleSignJailDuration); err != nil {
			return nil, err
		}
		return doubleSignJailDuration, nil
	case string(KeyDowntimeJailDuration):
		downtimeJailDuration, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateDowntimeJailDuration(downtimeJailDuration); err != nil {
			return nil, err
		}
		return downtimeJailDuration, nil
	case string(KeyCensorshipJailDuration):
		censorshipJailDuration, err := time.ParseDuration(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateCensorshipJailDuration(censorshipJailDuration); err != nil {
			return nil, err
		}
		return censorshipJailDuration, nil
	case string(KeySlashFractionDoubleSign):
		slashFractionDoubleSign, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateSlashFractionDoubleSign(slashFractionDoubleSign); err != nil {
			return nil, err
		}
		return slashFractionDoubleSign, nil
	case string(KeySlashFractionDowntime):
		slashFractionDowntime, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateSlashFractionDowntime(slashFractionDowntime); err != nil {
			return nil, err
		}
		return slashFractionDowntime, nil
	case string(KeySlashFractionCensorship):
		slashFractionCensorship, err := sdk.NewDecFromStr(value)
		if err != nil {
			return nil, params.ErrInvalidString(value)
		}
		if err := validateSlashFractionCensorship(slashFractionCensorship); err != nil {
			return nil, err
		}
		return slashFractionCensorship, nil
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
		return strconv.FormatInt(p.MaxEvidenceAge, 10), err
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
	case string(KeyCensorshipJailDuration):
		err := cdc.UnmarshalJSON(bytes, &p.CensorshipJailDuration)
		return p.CensorshipJailDuration.String(), err
	case string(KeySlashFractionDoubleSign):
		err := cdc.UnmarshalJSON(bytes, &p.SlashFractionDoubleSign)
		return p.SlashFractionDoubleSign.String(), err
	case string(KeySlashFractionDowntime):
		err := cdc.UnmarshalJSON(bytes, &p.SlashFractionDowntime)
		return p.SlashFractionDowntime.String(), err
	case string(KeySlashFractionCensorship):
		err := cdc.UnmarshalJSON(bytes, &p.SlashFractionCensorship)
		return p.SlashFractionCensorship.String(), err
	default:
		return "", fmt.Errorf("%s is not existed", key)
	}
}

// Default parameters used by Iris Hub
func DefaultParams() Params {
	return Params{
		MaxEvidenceAge:          3 * BlocksPerDay,
		SignedBlocksWindow:      2 * BlocksPerDay,
		MinSignedPerWindow:      sdk.NewDecWithPrec(5, 1),
		DoubleSignJailDuration:  2 * sdk.Day,
		DowntimeJailDuration:    1 * sdk.Day,
		CensorshipJailDuration:  2 * sdk.Day,
		SlashFractionDoubleSign: sdk.NewDecWithPrec(1, 2),
		SlashFractionDowntime:   sdk.ZeroDec(),
		SlashFractionCensorship: sdk.ZeroDec(),
	}
}

func DefaultParamsForTestnet() Params {
	return Params{
		MaxEvidenceAge:          3 * BlocksPerMinute,
		SignedBlocksWindow:      200,
		MinSignedPerWindow:      sdk.NewDecWithPrec(5, 1),
		DoubleSignJailDuration:  20 * time.Minute,
		DowntimeJailDuration:    10 * time.Minute,
		CensorshipJailDuration:  20 * time.Minute,
		SlashFractionDoubleSign: sdk.NewDecWithPrec(5, 2),
		SlashFractionDowntime:   sdk.ZeroDec(),
		SlashFractionCensorship: sdk.ZeroDec(),
	}
}

func validateParams(p Params) sdk.Error {
	if sdk.NetworkType != sdk.Mainnet {
		return nil
	}

	if err := validateMaxEvidenceAge(p.MaxEvidenceAge); err != nil {
		return err
	}
	if err := validateSignedBlocksWindow(p.SignedBlocksWindow); err != nil {
		return err
	}
	if err := validateMinSignedPerWindow(p.MinSignedPerWindow); err != nil {
		return err
	}
	if err := validateDoubleSignJailDuration(p.DoubleSignJailDuration); err != nil {
		return err
	}
	if err := validateDowntimeJailDuration(p.DowntimeJailDuration); err != nil {
		return err
	}
	if err := validateCensorshipJailDuration(p.CensorshipJailDuration); err != nil {
		return err
	}
	if err := validateSlashFractionDoubleSign(p.SlashFractionDoubleSign); err != nil {
		return err
	}
	if err := validateSlashFractionDowntime(p.SlashFractionDowntime); err != nil {
		return err
	}
	if err := validateSlashFractionCensorship(p.SlashFractionCensorship); err != nil {
		return err
	}
	return nil
}

func validateMaxEvidenceAge(p int64) sdk.Error {
	if sdk.NetworkType == sdk.Mainnet {
		if p < 2*BlocksPerDay {
			return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash MaxEvidenceAge [%d] should be between [2days,) ", p))
		}
	} else if p < 2*BlocksPerMinute {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash MaxEvidenceAge [%d] should be between [2minutes,) ", p))
	}
	return nil
}

func validateSignedBlocksWindow(p int64) sdk.Error {
	if p < 100 || p > 140000 {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash SignedBlocksWindow [%d] should be between [100, 140000] ", p))
	}
	return nil
}

func validateMinSignedPerWindow(p sdk.Dec) sdk.Error {
	if p.LT(sdk.NewDecWithPrec(5, 1)) || p.GT(sdk.NewDecWithPrec(9, 1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash MinSignedPerWindow [%s] should be between [0.5, 0.9] ", p.String()))
	}
	return nil
}

func validateDoubleSignJailDuration(p time.Duration) sdk.Error {
	if p <= 0 || p >= 2*sdk.Week {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash DoubleSignJailDuration [%s] should be between (0, 2weeks) ", p.String()))
	}
	return nil
}

func validateDowntimeJailDuration(p time.Duration) sdk.Error {
	if p <= 0 || p >= 1*sdk.Week {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash DowntimeJailDuration [%s] should be between (0, 1week) ", p.String()))
	}
	return nil
}

func validateCensorshipJailDuration(p time.Duration) sdk.Error {
	if p <= 0 || p >= 2*sdk.Week {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash CensorshipJailDuration [%s] should be between (0, 2weeks) ", p.String()))
	}
	return nil
}

func validateSlashFractionDoubleSign(p sdk.Dec) sdk.Error {
	if p.LT(sdk.ZeroDec()) || p.GT(sdk.NewDecWithPrec(1, 1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash SlashFractionDoubleSign [%s] should be between [0, 0.1] ", p.String()))
	}
	return nil
}

func validateSlashFractionDowntime(p sdk.Dec) sdk.Error {
	if p.LT(sdk.ZeroDec()) || p.GT(sdk.NewDecWithPrec(1, 1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash SlashFractionDowntime [%s] should be between [0, 0.1] ", p.String()))
	}
	return nil
}

func validateSlashFractionCensorship(p sdk.Dec) sdk.Error {
	if p.LT(sdk.ZeroDec()) || p.GT(sdk.NewDecWithPrec(1, 1)) {
		return sdk.NewError(params.DefaultCodespace, params.CodeInvalidSlashParams, fmt.Sprintf("Slash SlashFractionCensorship [%s] should be between [0, 0.1] ", p.String()))
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

// MaxEvidenceAge - Max age for evidence
func (k Keeper) MaxEvidenceAge(ctx sdk.Context) (res int64) {
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

// Double-sign jail duration
func (k Keeper) DoubleSignJailDuration(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, KeyDoubleSignJailDuration, &res)
	return
}

// Downtime jail duration
func (k Keeper) DowntimeJailDuration(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, KeyDowntimeJailDuration, &res)
	return
}

// Censorship jail duration
func (k Keeper) CensorshipJailDuration(ctx sdk.Context) (res time.Duration) {
	k.paramspace.Get(ctx, KeyCensorshipJailDuration, &res)
	return
}

// Slash fraction for DoubleSign
func (k Keeper) SlashFractionDoubleSign(ctx sdk.Context) (res sdk.Dec) {
	k.paramspace.Get(ctx, KeySlashFractionDoubleSign, &res)
	return
}

// Slash fraction for Downtime
func (k Keeper) SlashFractionDowntime(ctx sdk.Context) (res sdk.Dec) {
	k.paramspace.Get(ctx, KeySlashFractionDowntime, &res)
	return
}

// Slash fraction for Censorship
func (k Keeper) SlashFractionCensorship(ctx sdk.Context) (res sdk.Dec) {
	k.paramspace.Get(ctx, KeySlashFractionCensorship, &res)
	return
}
