package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// service module sentinel errors
var (
	ErrUnknownSvcDef         = sdkerrors.Register(ModuleName, 1, "unknown service definition")
	ErrUnknownSvcBinding     = sdkerrors.Register(ModuleName, 2, "unknown service binding")
	ErrInvalidServiceName    = sdkerrors.Register(ModuleName, 3, "invalid service name, must contain alphanumeric characters, _ and - only，length greater than 0 and less than or equal to 128")
	ErrInvalidSchemas        = sdkerrors.Register(ModuleName, 4, "invalid schemas")
	ErrSvcBindingExists      = sdkerrors.Register(ModuleName, 5, "service binding already exists")
	ErrLtMinProviderDeposit  = sdkerrors.Register(ModuleName, 6, "deposit amount must be equal or greater than min provider deposit")
	ErrUnknownMethod         = sdkerrors.Register(ModuleName, 7, "unknown service method")
	ErrUnavailable           = sdkerrors.Register(ModuleName, 8, "service binding is unavailable")
	ErrAvailable             = sdkerrors.Register(ModuleName, 9, "service binding is available")
	ErrRefundDeposit         = sdkerrors.Register(ModuleName, 10, "can't refund deposit")
	ErrIntOverflow           = sdkerrors.Register(ModuleName, 11, "Int overflow")
	ErrUnknownProfiler       = sdkerrors.Register(ModuleName, 12, "unknown profiler")
	ErrUnknownTrustee        = sdkerrors.Register(ModuleName, 13, "unknown trustee")
	ErrLtServiceFee          = sdkerrors.Register(ModuleName, 14, "service fee amount must be equal or greater than price")
	ErrUnknownActiveRequest  = sdkerrors.Register(ModuleName, 15, "unknown active request")
	ErrNotMatchingProvider   = sdkerrors.Register(ModuleName, 16, "not a matching provider")
	ErrNotMatchingReqChainID = sdkerrors.Register(ModuleName, 17, "not a matching request chain-id")
	ErrUnknownReturnFee      = sdkerrors.Register(ModuleName, 18, "no service refund fees")
	ErrUnknownWithdrawFee    = sdkerrors.Register(ModuleName, 19, "no service withdraw fees")
	ErrUnknownResponse       = sdkerrors.Register(ModuleName, 20, "unknown response")
	ErrInvalidRequestInput   = sdkerrors.Register(ModuleName, 21, "invalid request input")
	ErrInvalidResponseOutput = sdkerrors.Register(ModuleName, 22, "invalid response output")
	ErrInvalidResponseErr    = sdkerrors.Register(ModuleName, 23, "invalid response err")
)
