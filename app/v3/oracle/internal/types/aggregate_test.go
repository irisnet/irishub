package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMax(t *testing.T) {
	var data = []Value{1, "2", "3", "2", 7.0}
	max, _ := GetAggregateFunc("max")
	result := max(data)
	require.Equal(t, 7.0, result)
}

func TestMin(t *testing.T) {
	var data = []Value{1, "2", "3", "2", -1, 7.0}
	min, _ := GetAggregateFunc("min")
	result := min(data)
	require.Equal(t, -1.0, result)
}

func TestAvg(t *testing.T) {
	var data = []Value{1, "2", "3", "4", 5, 6, 7, 8, 9, 10}
	avg, _ := GetAggregateFunc("avg")
	result := avg(data)
	require.Equal(t, 5.5, result)
}
