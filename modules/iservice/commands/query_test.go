package commands

import ("testing"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	is "github.com/irisnet/iris-hub/modules/iservice"
	"strings"
	"log"
)

func TestQueryTx(t *testing.T) {
	c := client.GetNode("localhost:46657")
	//data, err := hex.DecodeString(cmn.StripHex("F91E63BDC83858234B4AAC0B2DFE0B74E70F6CB3"))
	//findSender := is.QSender(data)
	//findSender := is.QKeyTag("iriskyc")
	//findSender := is.Qheight(16)
	findSender := is.QServiceName("drgs")

	result,err := c.TxSearch(findSender,false)
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println(len(result))
}

func Test1(t *testing.T){
	tx_tags := strings.Split("", ",")
	for _, tag := range tx_tags {
		kv := strings.Split(tag, "=")
		if len(kv) == 2 {
			log.Printf("query:%s",is.QKVTag(kv[0], kv[1]))
		} else {
			log.Printf("query:%s",is.QKeyTag(kv[0]))
		}
	}
}
