package cli

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"strings"
	"testing"
)

func TestTxSearch(t *testing.T) {
	txStr := "nALZHnawCn4qlySsCjsKFDvFEHUEs4I/M1dfXMv4UZ0X4TtOEiMKCWlyaXMtYXR0bxIWMjA0MDkyODExMDAwMDAwMDAwMDAwMBI7ChR2IdpSv6fMUpP8LXkvm5zUrM7hgRIjCglpcmlzLWF0dG8SFjIwNDA5MjgxMTAwMDAwMDAwMDAwMDASJAofCglpcmlzLWF0dG8SEjEwMDAwMDAwMDAwMDAwMDAwMBCQThpwCibrWumHIQIANLPdNHAbKaKDVPNPhHSiH4jYKLHBio1R9VmYmWq0oBJAh2F4M4o512LL1ktxrYT/3bzgTQpJDaMtHojOsNTw3ygFo4R3joSH7NxGXxXD+m1brNZnqk9DYFvYEFsG5hPYphjAPiCDDw=="
	txBz, _ := base64.StdEncoding.DecodeString(txStr)

	txHashBz := sha256.Sum256(txBz)

	txHash := strings.ToUpper(hex.EncodeToString(txHashBz[:]))

	fmt.Println(txHash)

	rpc := rpcclient.NewHTTP("tcp://irisnet-rpc.rainbow.one:26657", "/websocket")

	rpc.Start()
	defer rpc.Stop()

	tx, err := rpc.Tx(txHashBz[:], false)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(fmt.Sprintf("%v", tx))
}
