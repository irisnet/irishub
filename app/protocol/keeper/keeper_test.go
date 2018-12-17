package keeper

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestKeeper(t *testing.T) {
	require.Equal(t,false, isValidProtocolVersion(1,1,1))
	require.Equal(t,true, isValidProtocolVersion(1,4,4))
	require.Equal(t,true, isValidProtocolVersion(1,4,5))
	require.Equal(t,true, isValidProtocolVersion(1,1,2))
	require.Equal(t,true, isValidProtocolVersion(2,1,3))
}
