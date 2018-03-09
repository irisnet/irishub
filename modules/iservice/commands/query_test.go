package commands

import ("testing"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	is "github.com/irisnet/iris-hub/modules/iservice"
)

func TestQueryTx(t *testing.T) {
	c := client.GetNode("localhost:46657")

	//findSender := is.QSender("CEA238B4D6231EBE6305DF65BC94AC3DEDD477D3")
	findSvcName := is.QServiceName("kyc5")

	result,err := c.TxSearch(findSvcName,false)
	if err != nil {
		fmt.Print("error")
	}
	fmt.Print(len(result))
}
