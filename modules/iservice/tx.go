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
	ServiceDefinition
}

// NewTxDefineService - new TxDefineService
func NewTxDefineService( service ServiceDefinition) sdk.Tx {
	return TxDefineService{
		service,
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
	if tx.ChainID == "" {
		return errServiceChainID
	}
	if tx.Messaging != "Unicast" || tx.Messaging != "Multicast" {
		return errServiceMessaging
	}
	if len(tx.Methods) < 0 {
		return errServiceMethods
	}

	var mthId = make(map[int64]bool)
	var mthNm = make(map[string]bool)

	for _,method := range tx.Methods {
		if method.ID < 0 {
			return errServiceMethodID
		}
		if mthId[method.ID]{
            return errServiceMethodIDNotUnique
		}
		mthId[method.ID] = true

		if mthNm[method.Name]{
			return errServiceMethodNameNotUnique
		}
		if method.Name == "" {
			return errServiceMethodName
		}
		mthNm[method.Name] = true

		if method.Description == "" {
			return errServiceMethodDescription
		}
		if method.Input == "" {
			return errServiceMethodInput
		}
		if method.Output == "" {
			return errServiceMethodOutput
		}
		if method.OutputPrivacy != "NoPrivacy" || method.OutputPrivacy != "PubKeyEncryption"{
			return errServiceMethodOutputPrivacy
		}
	}
	return nil
}
