package types

import (
	"github.com/irisnet/irismod/types/exported"
)

// Keys for parameter access
// nolint
var (
	KeyMaxRequestTimeout         = []byte("MaxRequestTimeout")
	KeyMinDepositMultiple        = []byte("MinDepositMultiple")
	KeyMinDeposit                = []byte("MinDeposit")
	KeyServiceFeeTax             = []byte("ServiceFeeTax")
	KeySlashFraction             = []byte("SlashFraction")
	KeyComplaintRetrospect       = []byte("ComplaintRetrospect")
	KeyArbitrationTimeLimit      = []byte("ArbitrationTimeLimit")
	KeyTxSizeLimit               = []byte("TxSizeLimit")
	KeyBaseDenom                 = []byte("BaseDenom")
	KeyRestrictedServiceFeeDenom = []byte("RestrictedServiceFeeDenom")
)

// ParamKeyTable for service module
func ParamKeyTable() exported.KeyTable {
	return exported.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements paramstypes.ParamSet
func (p *Params) ParamSetPairs() exported.ParamSetPairs {
	return exported.ParamSetPairs{
		exported.NewParamSetPair(
			KeyMaxRequestTimeout,
			&p.MaxRequestTimeout,
			validateMaxRequestTimeout,
		),
		exported.NewParamSetPair(
			KeyMinDepositMultiple,
			&p.MinDepositMultiple,
			validateMinDepositMultiple,
		),
		exported.NewParamSetPair(KeyMinDeposit, &p.MinDeposit, validateMinDeposit),
		exported.NewParamSetPair(KeyServiceFeeTax, &p.ServiceFeeTax, validateServiceFeeTax),
		exported.NewParamSetPair(KeySlashFraction, &p.SlashFraction, validateSlashFraction),
		exported.NewParamSetPair(
			KeyComplaintRetrospect,
			&p.ComplaintRetrospect,
			validateComplaintRetrospect,
		),
		exported.NewParamSetPair(
			KeyArbitrationTimeLimit,
			&p.ArbitrationTimeLimit,
			validateArbitrationTimeLimit,
		),
		exported.NewParamSetPair(KeyTxSizeLimit, &p.TxSizeLimit, validateTxSizeLimit),
		exported.NewParamSetPair(KeyBaseDenom, &p.BaseDenom, validateBaseDenom),
		exported.NewParamSetPair(
			KeyRestrictedServiceFeeDenom,
			&p.RestrictedServiceFeeDenom,
			validateRestrictedServiceFeeDenom,
		),
	}
}
