package types

import (
	"fmt"
	"testing"
)

func Test_htlc(t *testing.T) {
	htlc := HTLC{
		Sender: senderAddr,
	}
	ss, _ := msgCdc.MarshalJSON(htlc)
	fmt.Println(string(ss))
}
