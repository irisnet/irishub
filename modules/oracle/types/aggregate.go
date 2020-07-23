package types

import (
	"fmt"
	"math"
	"strconv"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tidwall/gjson"
)

var (
	router funcRouter
)

type ArgsType = gjson.Result
type funcRouter map[string]Aggregate
type Aggregate func(args []ArgsType) string

func init() {
	router = make(funcRouter)
	_ = RegisterAggregateFunc("max", Max)
	_ = RegisterAggregateFunc("min", Min)
	_ = RegisterAggregateFunc("avg", Avg)
}

func GetAggregateFunc(methodNm string) (Aggregate, error) {
	fun, ok := router[methodNm]
	if !ok {
		return nil, sdkerrors.Wrapf(ErrNotRegisterFunc, methodNm)
	}
	return fun, nil
}

func RegisterAggregateFunc(methodNm string, fun Aggregate) error {
	_, ok := router[methodNm]
	if ok {
		return fmt.Errorf("%s has existed", methodNm)
	}
	router[methodNm] = fun
	return nil
}

func Max(data []ArgsType) string {
	var maxNumber = math.SmallestNonzeroFloat64
	for _, d := range data {
		f := d.Float()
		if maxNumber < f {
			maxNumber = f
		}
	}
	return strconv.FormatFloat(maxNumber, 'f', 8, 64)
}

func Min(data []ArgsType) string {
	var minNum = math.MaxFloat64
	for _, d := range data {
		f := d.Float()
		if minNum > d.Float() {
			minNum = f
		}
	}
	return strconv.FormatFloat(minNum, 'f', 8, 64)
}

func Avg(data []ArgsType) string {
	var total = 0.0
	for _, d := range data {
		f := d.Float()
		total += f
	}
	return strconv.FormatFloat(total/float64(len(data)), 'f', 8, 64)
}
