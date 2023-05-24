package types

import (
	"fmt"
	"testing"

	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidation(t *testing.T) {
	conf := sdk.GetConfig()
	conf.SetBech32PrefixForAccount("iaa", "iva")
	escrowCollector := sdk.AccAddress(crypto.AddressHash([]byte(EscrowCollector)))
	fmt.Println(escrowCollector.String())
}
