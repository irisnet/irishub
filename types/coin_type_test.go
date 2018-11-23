package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConvert(t *testing.T) {
	irisToken := NewDefaultCoinType("iris")

	result, err := irisToken.Convert("1500000000000000001iris-atto", "iris-nano")
	require.Nil(t, err)
	require.Equal(t,"1500000000.000000001iris-nano",result)
	t.Log(result)

	result, err = irisToken.Convert("15iris", "iris-atto")
	require.Nil(t, err)
	require.Equal(t,"15000000000000000000iris-atto",result)
	t.Log(result)

	result, err = irisToken.Convert("1.5iris", "iris-nano")
	require.Nil(t, err)
	require.Equal(t,"1500000000iris-nano",result)
	t.Log(result)

	result, err = irisToken.Convert("1500000000000000001iris-atto", "iris-nano")
	require.Nil(t, err)
	require.Equal(t,"1500000000.000000001iris-nano",result)
	t.Log(result)

	result, err = irisToken.Convert("1500000001.123iris-nano", "iris")
	require.Nil(t, err)
	require.Equal(t,"1.500000001123iris",result)
	t.Log(result)

}
