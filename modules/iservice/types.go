package iservice

// Params defines the high level settings for staking
type Params struct {
	GasDefineService int64 `json:"gas_define_service"`
}

func defaultParams() Params {
	return Params{
		GasDefineService: 20,
	}
}

//_________________________________________________________________________

type ServiceDefinition struct {
	Name        string
	Description string
}
