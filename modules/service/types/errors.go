package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// service module sentinel errors
var (
	ErrInvalidServiceName       = sdkerrors.Register(ModuleName, 2, "invalid service name, must contain alphanumeric characters, _ and - only，length greater than 0 and less than or equal to 70")
	ErrInvalidDescription       = sdkerrors.Register(ModuleName, 3, "invalid description")
	ErrInvalidTags              = sdkerrors.Register(ModuleName, 4, "invalid tags")
	ErrInvalidSchemas           = sdkerrors.Register(ModuleName, 5, "invalid schemas")
	ErrUnknownServiceDefinition = sdkerrors.Register(ModuleName, 6, "unknown service definition")
	ErrServiceDefinitionExists  = sdkerrors.Register(ModuleName, 7, "service definition already exists")

	ErrInvalidDeposit            = sdkerrors.Register(ModuleName, 8, "invalid deposit")
	ErrInvalidMinDeposit         = sdkerrors.Register(ModuleName, 9, "invalid minimum deposit")
	ErrInvalidPricing            = sdkerrors.Register(ModuleName, 10, "invalid pricing")
	ErrInvalidQoS                = sdkerrors.Register(ModuleName, 11, "invalid qos")
	ErrInvalidOptions            = sdkerrors.Register(ModuleName, 12, "invalid options")
	ErrServiceBindingExists      = sdkerrors.Register(ModuleName, 13, "service binding already exists")
	ErrUnknownServiceBinding     = sdkerrors.Register(ModuleName, 14, "unknown service binding")
	ErrServiceBindingUnavailable = sdkerrors.Register(ModuleName, 15, "service binding unavailable")
	ErrServiceBindingAvailable   = sdkerrors.Register(ModuleName, 16, "service binding available")
	ErrIncorrectRefundTime       = sdkerrors.Register(ModuleName, 17, "incorrect refund time")

	ErrInvalidServiceFeeCap      = sdkerrors.Register(ModuleName, 18, "invalid service fee cap")
	ErrInvalidProviders          = sdkerrors.Register(ModuleName, 19, "invalid providers")
	ErrInvalidTimeout            = sdkerrors.Register(ModuleName, 20, "invalid timeout")
	ErrInvalidRepeatedFreq       = sdkerrors.Register(ModuleName, 21, "invalid repeated frequency")
	ErrInvalidRepeatedTotal      = sdkerrors.Register(ModuleName, 22, "invalid repeated total count")
	ErrInvalidResponseThreshold  = sdkerrors.Register(ModuleName, 23, "invalid response threshold")
	ErrInvalidResponse           = sdkerrors.Register(ModuleName, 24, "invalid response")
	ErrInvalidRequestID          = sdkerrors.Register(ModuleName, 25, "invalid request ID")
	ErrUnknownRequest            = sdkerrors.Register(ModuleName, 26, "unknown request")
	ErrUnknownResponse           = sdkerrors.Register(ModuleName, 27, "unknown response")
	ErrUnknownRequestContext     = sdkerrors.Register(ModuleName, 28, "unknown request context")
	ErrInvalidRequestContextID   = sdkerrors.Register(ModuleName, 29, "invalid request context ID")
	ErrRequestContextNonRepeated = sdkerrors.Register(ModuleName, 30, "request context non repeated")
	ErrRequestContextNotRunning  = sdkerrors.Register(ModuleName, 31, "request context not running")
	ErrRequestContextNotPaused   = sdkerrors.Register(ModuleName, 32, "request context not paused")
	ErrRequestContextCompleted   = sdkerrors.Register(ModuleName, 33, "request context completed")
	ErrCallbackRegistered        = sdkerrors.Register(ModuleName, 34, "callback registered")
	ErrCallbackNotRegistered     = sdkerrors.Register(ModuleName, 35, "callback not registered")
	ErrNoEarnedFees              = sdkerrors.Register(ModuleName, 36, "no earned fees")

	ErrInvalidRequestInput   = sdkerrors.Register(ModuleName, 37, "invalid request input")
	ErrInvalidResponseOutput = sdkerrors.Register(ModuleName, 38, "invalid response output")
	ErrInvalidResponseResult = sdkerrors.Register(ModuleName, 39, "invalid response result")

	ErrInvalidSchemaName = sdkerrors.Register(ModuleName, 40, "invalid service schema name")
	ErrNotAuthorized     = sdkerrors.Register(ModuleName, 41, "not authorized")

	ErrModuleServiceRegistered   = sdkerrors.Register(ModuleName, 42, "module service registered")
	ErrInvalidModuleService      = sdkerrors.Register(ModuleName, 43, "invalid module service")
	ErrBindModuleService         = sdkerrors.Register(ModuleName, 44, "can not bind module service")
	ErrInvalidRequestInputBody   = sdkerrors.Register(ModuleName, 45, "invalid request input body")
	ErrInvalidResponseOutputBody = sdkerrors.Register(ModuleName, 46, "invalid response output body")
)
