package types

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestMax(t *testing.T) {
	var data = []ArgsType{
		StringArgsType("1"),
		StringArgsType("2"),
		StringArgsType("3"),
		StringArgsType("2"),
		NumArgsType(7.0),
	}
	max, _ := GetAggregateFunc("max")
	result := max(data)
	require.Equal(t, "7.00000000", result)
}

func TestMin(t *testing.T) {
	var data = []ArgsType{
		StringArgsType("1"),
		StringArgsType("2"),
		StringArgsType("3"),
		StringArgsType("2"),
		StringArgsType("-1"),
		NumArgsType(7.0),
	}
	min, _ := GetAggregateFunc("min")
	result := min(data)
	require.Equal(t, "-1.00000000", result)
}

func TestAvg(t *testing.T) {
	var data = []ArgsType{
		StringArgsType("1"),
		StringArgsType("2"),
		StringArgsType("3"),
		StringArgsType("4"),
		StringArgsType("5"),
		StringArgsType("6"),
		StringArgsType("7"),
		StringArgsType("8"),
		StringArgsType("9"),
		StringArgsType("10"),
	}
	avg, _ := GetAggregateFunc("avg")
	result := avg(data)
	require.Equal(t, "5.50000000", result)
}

func StringArgsType(v string) ArgsType {
	return ArgsType{
		Type: gjson.String,
		Str:  v,
	}
}

func NumArgsType(f float64) ArgsType {
	return ArgsType{
		Type: gjson.Number,
		Num:  f,
	}
}
