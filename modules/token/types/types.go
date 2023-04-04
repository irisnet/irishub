package types

import (
	"encoding/json"
	"strconv"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Bool string

const (
	False Bool = "false"
	True  Bool = "true"
	Nil   Bool = ""
)

func (b Bool) ToBool() bool {
	v := string(b)
	if len(v) == 0 {
		return false
	}
	result, _ := strconv.ParseBool(v)
	return result
}

func (b Bool) String() string {
	return string(b)
}

// Marshal needed for protobuf compatibility
func (b Bool) Marshal() ([]byte, error) {
	return []byte(b), nil
}

// Unmarshal needed for protobuf compatibility
func (b *Bool) Unmarshal(data []byte) error {
	*b = Bool(data[:])
	return nil
}

// Marshals to JSON using string
func (b Bool) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.String())
}

// UnmarshalJSON from using string
func (b *Bool) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	*b = Bool(s)
	return nil
}
func ParseBool(v string) (Bool, error) {
	if len(v) == 0 {
		return Nil, nil
	}
	result, err := strconv.ParseBool(v)
	if err != nil {
		return Nil, err
	}
	if result {
		return True, nil
	}
	return False, nil
}

// LossLessSwap calculates the output amount of a swap, ensuring no loss.
// input: input amount
// ratio: swap rate
// inputScale: the decimal scale of input amount
// outputScale: the decimal scale of output amount
func LossLessSwap(input sdk.Int, ratio sdk.Dec, inputScale, outputScale uint32) (sdk.Int, sdk.Int) {
	inputDec := sdk.NewDecFromInt(input)
	scaleFactor := int64(inputScale) - int64(outputScale)
	var scaleMultipler, scaleReverseMultipler sdk.Dec

	if scaleFactor >= 0 {
		scaleMultipler = sdk.NewDecWithPrec(1, scaleFactor)
		scaleReverseMultipler = sdk.NewDecFromInt(sdkmath.NewIntWithDecimal(1, int(scaleFactor)))
	} else {
		scaleMultipler = sdk.NewDecFromInt(sdkmath.NewIntWithDecimal(1, int(-scaleFactor)))
		scaleReverseMultipler = sdk.NewDecWithPrec(1, -scaleFactor)
	}

	// Calculate output
	outputDec := inputDec.Clone().Mul(scaleMultipler).Mul(ratio)
	outputInt := outputDec.Clone().TruncateDec()

	// Adjust input if there are decimal places in the output
	if !outputDec.Equal(outputInt) {
		outputFrac := outputDec.Clone().Sub(outputInt)
		inputFrac := outputFrac.Mul(scaleReverseMultipler)
		input = inputDec.Sub(inputFrac).TruncateInt()
	}

	return input, outputInt.TruncateInt()
}
