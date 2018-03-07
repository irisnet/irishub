package iservice

import (
	"github.com/cosmos/cosmos-sdk"
)

// Tx
//--------------------------------------------------------------------------------

// register the tx type with its validation logic
// make sure to use the name of the handler as the prefix in the tx type,
// so it gets routed properly
const (
	ByteTxDefineService = 0x59
	TypeTxDefineService = iserviceModuleName + "/defineService"
)

func init() {
	sdk.TxMapper.RegisterImplementation(TxDefineService{}, TypeTxDefineService, ByteTxDefineService)
}

//Verify interface at compile time
var _ sdk.TxInner = &TxDefineService{}

// TxDefineService - struct for service define transactions
type TxDefineService struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NewTxDefineService - new TxDefineService
func NewTxDefineService(name string, description string) sdk.Tx {
	return TxDefineService{
		Name:        name,
		Description: description,
	}.Wrap()
}

// Wrap - Wrap a Tx as a Basecoin Tx
func (tx TxDefineService) Wrap() sdk.Tx { return sdk.Tx{tx} }

// ValidateBasic
func (tx TxDefineService) ValidateBasic() error {
	if tx.Name == "" {
		return errServiceNameEmpty
	}
	if tx.Description == "" {
		return errServiceDescEmpty
	}
	return nil
}
