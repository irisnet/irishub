package types

import (
	"errors"
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"math"
	"strconv"
)

var (
	router methodRouter
)

type methodRouter map[string]Aggregate
type Aggregate func(args []Value) Value

func init() {
	router = make(methodRouter)
	_ = RegisterAggregateMethod("max", Max)
	_ = RegisterAggregateMethod("min", Min)
	_ = RegisterAggregateMethod("avg", Avg)
}

func GetAggregateMethod(methodNm string) (Aggregate, sdk.Error) {
	fun, ok := router[methodNm]
	if !ok {
		return nil, ErrNotRegisterMethod(DefaultCodespace, methodNm)
	}
	return fun, nil
}

func RegisterAggregateMethod(methodNm string, fun Aggregate) error {
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
	switch args.(type) {
	case string:
		return strconv.ParseFloat(args.(string), 64)
	case int:
		return float64(args.(int)), nil
	case int16:
		return float64(args.(int16)), nil
	case int32:
		return float64(args.(int32)), nil
	case uint:
		return float64(args.(uint)), nil
	case uint16:
		return float64(args.(uint16)), nil
	case uint32:
		return float64(args.(uint32)), nil
	case uint64:
		return float64(args.(uint64)), nil
	case float32:
		return float64(args.(float32)), nil
	case float64:
		return args.(float64), nil
	}
	return 0.0, errors.New("invalid number")
}
