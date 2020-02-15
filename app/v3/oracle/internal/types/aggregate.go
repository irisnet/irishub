package types

import (
	"errors"
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"math"
	"strconv"
)

var (
	router funcRouter
)

type funcRouter map[string]Aggregate
type Aggregate func(args []Value) Value

func init() {
	router = make(funcRouter)
	_ = RegisterAggregateFunc("max", Max)
	_ = RegisterAggregateFunc("min", Min)
	_ = RegisterAggregateFunc("avg", Avg)
}

func GetAggregateFunc(methodNm string) (Aggregate, sdk.Error) {
	fun, ok := router[methodNm]
	if !ok {
		return nil, ErrNotRegisterMethod(DefaultCodespace, methodNm)
	}
	return fun, nil
}

func RegisterAggregateFunc(methodNm string, fun Aggregate) error {
	_, ok := router[methodNm]
	if ok {
		return errors.New(fmt.Sprintf("%s has existed", methodNm))
	}
	router[methodNm] = fun
	return nil
}

func Max(data []Value) Value {
	var maxNumber = math.SmallestNonzeroFloat64
	for _, d := range data {
		f, err := ConvertToFloat64(d)
		if err != nil {
			continue
		}
		if maxNumber < f {
			maxNumber = f
		}
	}
	return maxNumber
}

func Min(data []Value) Value {
	var maxNumber = math.MaxFloat64
	for _, d := range data {
		f, err := ConvertToFloat64(d)
		if err != nil {
			continue
		}
		if maxNumber > f {
			maxNumber = f
		}
	}
	return maxNumber
}

func Avg(data []Value) Value {
	var total = 0.0
	for _, d := range data {
		f, err := ConvertToFloat64(d)
		if err != nil {
			continue
		}
		total += f
	}
	return total / float64(len(data))
}

func ConvertToFloat64(args Value) (float64, error) {
	switch args := args.(type) {
	case string:
		return strconv.ParseFloat(args, 64)
	case int:
		return float64(args), nil
	case int16:
		return float64(args), nil
	case int32:
		return float64(args), nil
	case uint:
		return float64(args), nil
	case uint16:
		return float64(args), nil
	case uint32:
		return float64(args), nil
	case uint64:
		return float64(args), nil
	case float32:
		return float64(args), nil
	case float64:
		return args, nil
	}
	return 0.0, errors.New("invalid number")
}
