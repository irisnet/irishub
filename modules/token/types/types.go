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

func LossLessSwap(input sdk.Int, ratio sdk.Dec, inputScale, outputScale uint32) (sdk.Int, sdk.Int) {
	inputDec := sdk.NewDecFromInt(input)
	if inputScale >= outputScale {
		scaleFactor := inputScale - outputScale
		scaleMultipler := sdk.NewDecWithPrec(1, int64(scaleFactor))
		outputDec := scaleMultipler.Clone().Mul(inputDec).Mul(ratio)
		outputInt := outputDec.Clone().TruncateDec()
		if outputDec.Equal(outputInt) {
			return input, outputInt.TruncateInt()
		}

		//If there are decimal places, the decimal places need to be subtracted from input
		outputFrac := outputDec.Clone().Sub(outputInt)
		scaleReverseMultipler := sdkmath.NewIntWithDecimal(1, int(scaleFactor))
		inputFrac := outputFrac.MulInt(scaleReverseMultipler)
		return inputDec.Sub(inputFrac).TruncateInt(), outputInt.TruncateInt()
	}

	// When a large unit wants to convert a small unit, there is no case of discarding decimal places
	scaleFactor := outputScale - inputScale
	scaleMultipler := sdkmath.NewIntWithDecimal(1, int(scaleFactor))
	outputDec := inputDec.Clone().Mul(sdk.NewDecFromInt(scaleMultipler)).Mul(ratio)
	outputInt := outputDec.Clone().TruncateDec()
	if outputDec.Equal(outputInt) {
		return input, outputInt.TruncateInt()
	}

	//If there are decimal places, the decimal places need to be subtracted from input
	outputFrac := outputDec.Clone().Sub(outputInt)
	scaleReverseMultipler := sdk.NewDecWithPrec(1, int64(scaleFactor))
	inputFrac := outputFrac.Mul(scaleReverseMultipler)
	return inputDec.Sub(inputFrac).TruncateInt(), outputInt.TruncateInt()
}
