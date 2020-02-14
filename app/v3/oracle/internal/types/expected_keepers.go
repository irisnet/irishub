package types

import sdk "github.com/irisnet/irishub/types"

//expected Service keeper
type ServiceKeeper interface {
	RegisterResponseHandler(moduleName string,
		callback func(ctx sdk.Context, requestContextID []byte, responseOutput []string)) error

	GetRequestContext(ctx sdk.Context, requestContextID []byte) (RequestContext, bool)

	CreateRequestContext(ctx sdk.Context,
		serviceName string,
		providers []sdk.AccAddress,
		consumer sdk.AccAddress,
		input string,
		serviceFeeCap sdk.Coins,
		timeout int64,
		repeated bool,
		repeatedFrequency uint64,
		repeatedTotal int64,
		state RequestContextState,
		respThreshold uint16,
		respHandler string) ([]byte, sdk.Error)

	UpdateRequestContext(ctx sdk.Context, requestContextID []byte) error

	StartRequestContext(ctx sdk.Context, requestContextID []byte) error

	PauseRequestContext(ctx sdk.Context, requestContextID []byte) error

	KillRequestContext(ctx sdk.Context, requestContextID []byte) error
}

type RequestContext = MockRequestContext

const (
	Running RequestContextState = "running"
	Pause   RequestContextState = "pause"
)

type MockRequestContext struct {
	ServiceName       string              `json:"service_name"`
	Providers         []sdk.AccAddress    `json:"providers"`
	Consumer          sdk.AccAddress      `json:"consumer"`
	Input             string              `json:"input"`
	ServiceFeeCap     sdk.Coins           `json:"service_fee_cap"`
	Timeout           int64               `json:"timeout"`
	Repeated          bool                `json:"repeated"`
	RepeatedFrequency uint64              `json:"repeated_frequency"`
	RepeatedTotal     int64               `json:"repeated_total"`
	BatchCounter      uint64              `json:"batch_counter"`
	State             RequestContextState `json:"state"`
	ResponseThreshold uint16              `json:"response_threshold"`
	ResponseHandler   string              `json:"response_handler"`
}

type RequestContextState string
