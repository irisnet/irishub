package types

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
func ParamKeyTable() KeyTable {
	return NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements paramstypes.ParamSet
func (p *Params) ParamSetPairs() ParamSetPairs {
	return ParamSetPairs{
		NewParamSetPair(
			KeyMaxRequestTimeout,
			&p.MaxRequestTimeout,
			validateMaxRequestTimeout,
		),
		NewParamSetPair(
			KeyMinDepositMultiple,
			&p.MinDepositMultiple,
			validateMinDepositMultiple,
		),
		NewParamSetPair(KeyMinDeposit, &p.MinDeposit, validateMinDeposit),
		NewParamSetPair(KeyServiceFeeTax, &p.ServiceFeeTax, validateServiceFeeTax),
		NewParamSetPair(KeySlashFraction, &p.SlashFraction, validateSlashFraction),
		NewParamSetPair(
			KeyComplaintRetrospect,
			&p.ComplaintRetrospect,
			validateComplaintRetrospect,
		),
		NewParamSetPair(
			KeyArbitrationTimeLimit,
			&p.ArbitrationTimeLimit,
			validateArbitrationTimeLimit,
		),
		NewParamSetPair(KeyTxSizeLimit, &p.TxSizeLimit, validateTxSizeLimit),
		NewParamSetPair(KeyBaseDenom, &p.BaseDenom, validateBaseDenom),
		NewParamSetPair(
			KeyRestrictedServiceFeeDenom,
			&p.RestrictedServiceFeeDenom,
			validateRestrictedServiceFeeDenom,
		),
	}
}
