package types

import (
	errorsmod "cosmossdk.io/errors"
)

// service module sentinel errors
var (
	ErrInvalidServiceName       = errorsmod.Register(ModuleName, 2, "invalid service name, must contain alphanumeric characters, _ and - only，length greater than 0 and less than or equal to 70")
	ErrInvalidDescription       = errorsmod.Register(ModuleName, 3, "invalid description")
	ErrInvalidTags              = errorsmod.Register(ModuleName, 4, "invalid tags")
	ErrInvalidSchemas           = errorsmod.Register(ModuleName, 5, "invalid schemas")
	ErrUnknownServiceDefinition = errorsmod.Register(ModuleName, 6, "unknown service definition")
	ErrServiceDefinitionExists  = errorsmod.Register(ModuleName, 7, "service definition already exists")

	ErrInvalidDeposit            = errorsmod.Register(ModuleName, 8, "invalid deposit")
	ErrInvalidMinDeposit         = errorsmod.Register(ModuleName, 9, "invalid minimum deposit")
	ErrInvalidPricing            = errorsmod.Register(ModuleName, 10, "invalid pricing")
	ErrInvalidQoS                = errorsmod.Register(ModuleName, 11, "invalid QoS")
	ErrInvalidOptions            = errorsmod.Register(ModuleName, 12, "invalid options")
	ErrServiceBindingExists      = errorsmod.Register(ModuleName, 13, "service binding already exists")
	ErrUnknownServiceBinding     = errorsmod.Register(ModuleName, 14, "unknown service binding")
	ErrServiceBindingUnavailable = errorsmod.Register(ModuleName, 15, "service binding unavailable")
	ErrServiceBindingAvailable   = errorsmod.Register(ModuleName, 16, "service binding available")
	ErrIncorrectRefundTime       = errorsmod.Register(ModuleName, 17, "incorrect refund time")

	ErrInvalidServiceFeeCap      = errorsmod.Register(ModuleName, 18, "invalid service fee cap")
	ErrInvalidProviders          = errorsmod.Register(ModuleName, 19, "invalid providers")
	ErrInvalidTimeout            = errorsmod.Register(ModuleName, 20, "invalid timeout")
	ErrInvalidRepeatedFreq       = errorsmod.Register(ModuleName, 21, "invalid repeated frequency")
	ErrInvalidRepeatedTotal      = errorsmod.Register(ModuleName, 22, "invalid repeated total count")
	ErrInvalidResponseThreshold  = errorsmod.Register(ModuleName, 23, "invalid response threshold")
	ErrInvalidResponse           = errorsmod.Register(ModuleName, 24, "invalid response")
	ErrInvalidRequestID          = errorsmod.Register(ModuleName, 25, "invalid request ID")
	ErrUnknownRequest            = errorsmod.Register(ModuleName, 26, "unknown request")
	ErrUnknownResponse           = errorsmod.Register(ModuleName, 27, "unknown response")
	ErrUnknownRequestContext     = errorsmod.Register(ModuleName, 28, "unknown request context")
	ErrInvalidRequestContextID   = errorsmod.Register(ModuleName, 29, "invalid request context ID")
	ErrRequestContextNonRepeated = errorsmod.Register(ModuleName, 30, "request context non repeated")
	ErrRequestContextNotRunning  = errorsmod.Register(ModuleName, 31, "request context not running")
	ErrRequestContextNotPaused   = errorsmod.Register(ModuleName, 32, "request context not paused")
	ErrRequestContextCompleted   = errorsmod.Register(ModuleName, 33, "request context completed")
	ErrCallbackRegistered        = errorsmod.Register(ModuleName, 34, "callback registered")
	ErrCallbackNotRegistered     = errorsmod.Register(ModuleName, 35, "callback not registered")
	ErrNoEarnedFees              = errorsmod.Register(ModuleName, 36, "no earned fees")

	ErrInvalidRequestInput   = errorsmod.Register(ModuleName, 37, "invalid request input")
	ErrInvalidResponseOutput = errorsmod.Register(ModuleName, 38, "invalid response output")
	ErrInvalidResponseResult = errorsmod.Register(ModuleName, 39, "invalid response result")

	ErrInvalidSchemaName = errorsmod.Register(ModuleName, 40, "invalid service schema name")
	ErrNotAuthorized     = errorsmod.Register(ModuleName, 41, "not authorized")

	ErrModuleServiceRegistered   = errorsmod.Register(ModuleName, 42, "module service registered")
	ErrInvalidModuleService      = errorsmod.Register(ModuleName, 43, "invalid module service")
	ErrBindModuleService         = errorsmod.Register(ModuleName, 44, "can not bind module service")
	ErrInvalidRequestInputBody   = errorsmod.Register(ModuleName, 45, "invalid request input body")
	ErrInvalidResponseOutputBody = errorsmod.Register(ModuleName, 46, "invalid response output body")
)
