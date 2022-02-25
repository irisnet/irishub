package tibc_test

import (
	"testing"

	migratetibc "github.com/irisnet/irishub/migrate/tibc"
	"github.com/irisnet/irishub/simapp"
	"github.com/stretchr/testify/require"
)

func TestLoadClient(t *testing.T) {
	app := simapp.Setup(false)
	clients := migratetibc.LoadClient(app.AppCodec(), "v1.3")
	require.Equal(t, 1, len(clients))
}
