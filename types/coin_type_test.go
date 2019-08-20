package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConvert(t *testing.T) {
	result, err := IrisCoinType.Convert("1500000000000000001iris-atto", "iris-nano")
	require.Nil(t, err)
	require.Equal(t, "1500000000.000000001iris-nano", result)
	t.Log(result)

	result, err = IrisCoinType.Convert("15iris", "iris-atto")
	require.Nil(t, err)
	require.Equal(t, "15000000000000000000iris-atto", result)
	t.Log(result)

	result, err = IrisCoinType.Convert("1.5iris", "iris-nano")
	require.Nil(t, err)
	require.Equal(t, "1500000000iris-nano", result)
	t.Log(result)

	result, err = IrisCoinType.Convert("1500000000000000001iris-atto", "iris-nano")
	require.Nil(t, err)
	require.Equal(t, "1500000000.000000001iris-nano", result)
	t.Log(result)

	result, err = IrisCoinType.Convert("1500000001.123iris-nano", "iris")
	require.Nil(t, err)
	require.Equal(t, "1.500000001123iris", result)
	t.Log(result)

}

func TestGetCoin(t *testing.T) {
	testData := []struct {
		name, coinStr, expectAmount, expectDenom string
		expectPass                   bool
	}{
		{"standard", "1000iris", "1000", "iris", true},
		{"with -", "1000iris-atto", "1000", "iris-atto", true},
		{"with gateway", "1000gdex.btc", "1000", "gdex.btc", true},
		{"with x.", "1000x.btc", "1000", "x.btc", true},
		{"with decimal", "1000.001gdex.btc", "1000.001", "gdex.btc", true},
		{"with decimal and numeric", "1000.001gdex1.btc1d", "1000.001", "gdex1.btc1d", true},
	}

	for _, td := range testData {
		denom, amt, err := ParseCoinParts(td.coinStr)
		if td.expectPass {
			require.Equal(t, td.expectAmount, amt, "test: %v", td.name)
			require.Equal(t, td.expectDenom, denom, "test: %v", td.name)
			require.Nil(t, err, "test: %v", td.name)
		} else {
			require.NotNil(t, err, "test: %v", td.name)
		}
	}
}
