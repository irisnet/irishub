package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

func TestValidation(t *testing.T) {
	conf := sdk.GetConfig()
	conf.SetBech32PrefixForAccount("iaa", "iva")
	escrowCollector := sdk.AccAddress(crypto.AddressHash([]byte(EscrowCollector)))
	fmt.Println(escrowCollector.String())
}
