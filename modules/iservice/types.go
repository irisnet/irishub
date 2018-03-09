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
//服务定义
type ServiceDefinition struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Tags        string          `json:"tags"` // ('k1'='v1','k2'='v2')
	Creator     string          `json:"creator"`
	ChainID     string          `json:"chainID"`
	Messaging   string          `json:"messaging"`
	Methods     []ServiceMethod `json:"methods"`
}

type ServiceMethod struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Input         string `json:"input"`
	Output        string `json:"output"`
	Error         string `json:"error"`
	OutputPrivacy string `json:"outputPrivacy"`
}

//服务绑定
type ServiceBind struct {
	DefinitionHash    []byte `json:"definitionHash"`
	ChainID           string `json:"chainID"`
	ProviderAddress   []byte `json:"providerAddress"`
	BindingType       string `json:"bindingType"`
	ProviderDeposit   int64  `json:"providerDeposit"`
	ServicePricing    string `json:"servicePricing"`
	ServiceLevel      string `json:"serviceLevel"`
	BindingExpiration int64  `json:"bindingExpiration"`
	IsValid           bool   `json:"isValid"`
}

//服务调用
type ServiceRequest struct {
	ChainID         string `json:"chainID"`
	ConsumerAddress []byte `json:"consumerAddress"`
	DefinitionHash  []byte `json:"definitionHash"`
	MethodID        int    `json:"methodID"`
	InputValue      string `json:"inputValue"`
	BindingHash     []byte `json:"bindingHash"`
	MaxServiceFee   int64  `json:"maxServiceFee"`
	Timeout         int64  `json:"timeout"`
}

