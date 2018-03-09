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

type ServiceDefinition struct{
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        string `json:"tags"` // ('k1'='v1','k2'='v2')
	Creator     string `json:"creator"`
	ChainID     string `json:"chainID"`
	Messaging   string `json:"messaging"`
	Methods     []ServiceMethod `json:"methods"`
}

type ServiceMethod struct {
	ID            int64 `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Input         string `json:"input"`
	Output        string `json:"output"`
	Error         string `json:"error"`
	OutputPrivacy string `json:"outputPrivacy"`
}
