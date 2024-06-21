package contracts_test

import (
	"testing"

	"mods.irisnet.org/contracts"
)

func TestERC20(t *testing.T) {
	for k, v := range contracts.ERC20TokenContract.ABI.Events {
		t.Logf("event:%s ,id: %s", k, v.ID)
	}
}
