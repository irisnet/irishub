package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Message types for the service module
const (
	TypeMsgDefineService         = "define_service"          // type for MsgDefineService
	TypeMsgBindService           = "bind_service"            // type for MsgBindService
	TypeMsgUpdateServiceBinding  = "update_service_binding"  // type for MsgUpdateServiceBinding
	TypeMsgSetWithdrawAddress    = "set_withdraw_address"    // type for MsgSetWithdrawAddress
	TypeMsgDisableServiceBinding = "disable_service_binding" // type for MsgDisableServiceBinding
	TypeMsgEnableServiceBinding  = "enable_service_binding"  // type for MsgEnableServiceBinding
	TypeMsgRefundServiceDeposit  = "refund_service_deposit"  // type for MsgRefundServiceDeposit
	TypeMsgCallService           = "call_service"            // type for MsgCallService
	TypeMsgRespondService        = "respond_service"         // type for MsgRespondService
	TypeMsgPauseRequestContext   = "pause_request_context"   // type for MsgPauseRequestContext
	TypeMsgStartRequestContext   = "start_request_context"   // type for MsgStartRequestContext
	TypeMsgKillRequestContext    = "kill_request_context"    // type for MsgKillRequestContext
	TypeMsgUpdateRequestContext  = "update_request_context"  // type for MsgUpdateRequestContext
	TypeMsgWithdrawEarnedFees    = "withdraw_earned_fees"    // type for MsgWithdrawEarnedFees
)

var (
	_ sdk.Msg = &MsgDefineService{}
	_ sdk.Msg = &MsgBindService{}
	_ sdk.Msg = &MsgUpdateServiceBinding{}
	_ sdk.Msg = &MsgSetWithdrawAddress{}
	_ sdk.Msg = &MsgDisableServiceBinding{}
	_ sdk.Msg = &MsgEnableServiceBinding{}
	_ sdk.Msg = &MsgRefundServiceDeposit{}
	_ sdk.Msg = &MsgCallService{}
	_ sdk.Msg = &MsgStartRequestContext{}
	_ sdk.Msg = &MsgPauseRequestContext{}
	_ sdk.Msg = &MsgKillRequestContext{}
	_ sdk.Msg = &MsgUpdateRequestContext{}
	_ sdk.Msg = &MsgRespondService{}
	_ sdk.Msg = &MsgWithdrawEarnedFees{}
)

// ______________________________________________________________________

// NewMsgDefineService creates a new MsgDefineService instance
func NewMsgDefineService(
	name string,
	description string,
	tags []string,
	author string,
	authorDescription,
	schemas string,
) *MsgDefineService {
	return &MsgDefineService{
		Name:              name,
		Description:       description,
		Tags:              tags,
		Author:            author,
		AuthorDescription: authorDescription,
		Schemas:           schemas,
	}
}

// Route implements Msg
func (msg MsgDefineService) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgDefineService) Type() string { return TypeMsgDefineService }

// ValidateBasic implements Msg
func (msg MsgDefineService) ValidateBasic() error {
	if err := ValidateAuthor(msg.Author); err != nil {
		return err
	}

	if err := ValidateServiceName(msg.Name); err != nil {
		return err
	}

	if err := ValidateServiceDescription(msg.Description); err != nil {
		return err
	}

	if err := ValidateAuthorDescription(msg.AuthorDescription); err != nil {
		return err
	}

	if err := ValidateTags(msg.Tags); err != nil {
		return err
	}
	return ValidateServiceSchemas(msg.Schemas)
}

// GetSignBytes implements Msg
func (msg MsgDefineService) GetSignBytes() []byte {
	if len(msg.Tags) == 0 {
		msg.Tags = nil
	}

	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgDefineService) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Author)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ______________________________________________________________________

// NewMsgBindService creates a new MsgBindService instance
func NewMsgBindService(
	serviceName string,
	provider string,
	deposit sdk.Coins,
	pricing string,
	qos uint64,
	options string,
	owner string,
) *MsgBindService {
	return &MsgBindService{
		ServiceName: serviceName,
		Provider:    provider,
		Deposit:     deposit,
		Pricing:     pricing,
		QoS:         qos,
		Options:     options,
		Owner:       owner,
	}
}

// Route implements Msg.
func (msg MsgBindService) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgBindService) Type() string { return TypeMsgBindService }

// GetSignBytes implements Msg.
func (msg MsgBindService) GetSignBytes() []byte {
	if msg.Deposit.Empty() {
		msg.Deposit = nil
	}

	b := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgBindService) ValidateBasic() error {
	if err := ValidateProvider(msg.Provider); err != nil {
		return err
	}

	if err := ValidateOwner(msg.Owner); err != nil {
		return err
	}

	if err := ValidateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if err := ValidateServiceDeposit(msg.Deposit); err != nil {
		return err
	}

	if err := ValidateQoS(msg.QoS); err != nil {
		return err
	}

	if err := ValidateOptions(msg.Options); err != nil {
		return err
	}
	return ValidatePricing(msg.Pricing)
}

// GetSigners implements Msg.
func (msg MsgBindService) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ______________________________________________________________________

// NewMsgUpdateServiceBinding creates a new MsgUpdateServiceBinding instance
func NewMsgUpdateServiceBinding(
	serviceName string,
	provider string,
	deposit sdk.Coins,
	pricing string,
	qos uint64,
	options string,
	owner string,
) *MsgUpdateServiceBinding {
	return &MsgUpdateServiceBinding{
		ServiceName: serviceName,
		Provider:    provider,
		Deposit:     deposit,
		Pricing:     pricing,
		QoS:         qos,
		Options:     options,
		Owner:       owner,
	}
}

// Route implements Msg.
func (msg MsgUpdateServiceBinding) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgUpdateServiceBinding) Type() string { return TypeMsgUpdateServiceBinding }

// GetSignBytes implements Msg.
func (msg MsgUpdateServiceBinding) GetSignBytes() []byte {
	if msg.Deposit.Empty() {
		msg.Deposit = nil
	}

	b := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateServiceBinding) ValidateBasic() error {
	if err := ValidateProvider(msg.Provider); err != nil {
		return err
	}

	if err := ValidateOwner(msg.Owner); err != nil {
		return err
	}

	if err := ValidateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if !msg.Deposit.Empty() {
		if err := ValidateServiceDeposit(msg.Deposit); err != nil {
			return err
		}
	}

	if len(msg.Options) != 0 {
		if err := ValidateOptions(msg.Options); err != nil {
			return err
		}
	}

	if len(msg.Pricing) != 0 {
		return ValidatePricing(msg.Pricing)
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgUpdateServiceBinding) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ______________________________________________________________________

// NewMsgSetWithdrawAddress creates a new MsgSetWithdrawAddress instance
func NewMsgSetWithdrawAddress(owner, withdrawAddr string) *MsgSetWithdrawAddress {
	return &MsgSetWithdrawAddress{
		Owner:           owner,
		WithdrawAddress: withdrawAddr,
	}
}

// Route implements Msg.
func (msg MsgSetWithdrawAddress) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgSetWithdrawAddress) Type() string { return TypeMsgSetWithdrawAddress }

// GetSignBytes implements Msg.
func (msg MsgSetWithdrawAddress) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSetWithdrawAddress) ValidateBasic() error {
	if err := ValidateOwner(msg.Owner); err != nil {
		return err
	}
	return ValidateWithdrawAddress(msg.WithdrawAddress)
}

// Normalize return a string with spaces removed and lowercase
func (msg MsgSetWithdrawAddress) Normalize() MsgSetWithdrawAddress {
	return msg
}

// GetSigners implements Msg.
func (msg MsgSetWithdrawAddress) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ______________________________________________________________________

// NewMsgDisableServiceBinding creates a new MsgDisableServiceBinding instance
func NewMsgDisableServiceBinding(serviceName, provider, owner string) *MsgDisableServiceBinding {
	return &MsgDisableServiceBinding{
		ServiceName: serviceName,
		Provider:    provider,
		Owner:       owner,
	}
}

// Route implements Msg.
func (msg MsgDisableServiceBinding) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgDisableServiceBinding) Type() string { return TypeMsgDisableServiceBinding }

// GetSignBytes implements Msg.
func (msg MsgDisableServiceBinding) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgDisableServiceBinding) ValidateBasic() error {
	if err := ValidateProvider(msg.Provider); err != nil {
		return err
	}

	if err := ValidateOwner(msg.Owner); err != nil {
		return err
	}

	return ValidateServiceName(msg.ServiceName)
}

// GetSigners implements Msg.
func (msg MsgDisableServiceBinding) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ______________________________________________________________________

// NewMsgEnableServiceBinding creates a new MsgEnableServiceBinding instance
func NewMsgEnableServiceBinding(
	serviceName string,
	provider string,
	deposit sdk.Coins,
	owner string,
) *MsgEnableServiceBinding {
	return &MsgEnableServiceBinding{
		ServiceName: serviceName,
		Provider:    provider,
		Deposit:     deposit,
		Owner:       owner,
	}
}

// Route implements Msg.
func (msg MsgEnableServiceBinding) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgEnableServiceBinding) Type() string { return TypeMsgEnableServiceBinding }

// GetSignBytes implements Msg.
func (msg MsgEnableServiceBinding) GetSignBytes() []byte {
	if msg.Deposit.Empty() {
		msg.Deposit = nil
	}

	b := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgEnableServiceBinding) ValidateBasic() error {
	if err := ValidateProvider(msg.Provider); err != nil {
		return err
	}

	if err := ValidateOwner(msg.Owner); err != nil {
		return err
	}

	if err := ValidateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if !msg.Deposit.Empty() {
		if err := ValidateServiceDeposit(msg.Deposit); err != nil {
			return err
		}
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgEnableServiceBinding) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ______________________________________________________________________

// NewMsgRefundServiceDeposit creates a new MsgRefundServiceDeposit instance
func NewMsgRefundServiceDeposit(serviceName, provider, owner string) *MsgRefundServiceDeposit {
	return &MsgRefundServiceDeposit{
		ServiceName: serviceName,
		Provider:    provider,
		Owner:       owner,
	}
}

// Route implements Msg.
func (msg MsgRefundServiceDeposit) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgRefundServiceDeposit) Type() string { return TypeMsgRefundServiceDeposit }

// GetSignBytes implements Msg.
func (msg MsgRefundServiceDeposit) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgRefundServiceDeposit) ValidateBasic() error {
	if err := ValidateProvider(msg.Provider); err != nil {
		return err
	}

	if err := ValidateOwner(msg.Owner); err != nil {
		return err
	}
	return ValidateServiceName(msg.ServiceName)
}

// GetSigners implements Msg.
func (msg MsgRefundServiceDeposit) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ______________________________________________________________________

// NewMsgCallService creates a new MsgCallService instance
func NewMsgCallService(
	serviceName string,
	providers []string,
	consumer string,
	input string,
	serviceFeeCap sdk.Coins,
	timeout int64,
	repeated bool,
	repeatedFrequency uint64,
	repeatedTotal int64,
) *MsgCallService {
	return &MsgCallService{
		ServiceName:       serviceName,
		Providers:         providers,
		Consumer:          consumer,
		Input:             input,
		ServiceFeeCap:     serviceFeeCap,
		Timeout:           timeout,
		Repeated:          repeated,
		RepeatedFrequency: repeatedFrequency,
		RepeatedTotal:     repeatedTotal,
	}
}

// Route implements Msg.
func (msg MsgCallService) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgCallService) Type() string { return TypeMsgCallService }

// GetSignBytes implements Msg.
func (msg MsgCallService) GetSignBytes() []byte {
	if len(msg.Providers) == 0 {
		msg.Providers = nil
	}

	if msg.ServiceFeeCap.Empty() {
		msg.ServiceFeeCap = nil
	}

	b := ModuleCdc.MustMarshalJSON(&msg)

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgCallService) ValidateBasic() error {
	if err := ValidateConsumer(msg.Consumer); err != nil {
		return err
	}

	pds := make([]sdk.AccAddress, len(msg.Providers))
	for i, provider := range msg.Providers {
		pd, err := sdk.AccAddressFromBech32(provider)
		if err != nil {
			return err
		}
		pds[i] = pd
	}

	return ValidateRequest(
		msg.ServiceName,
		msg.ServiceFeeCap,
		pds,
		msg.Input,
		msg.Timeout,
		msg.Repeated,
		msg.RepeatedFrequency,
		msg.RepeatedTotal,
	)
}

// GetSigners implements Msg.
func (msg MsgCallService) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ______________________________________________________________________

// NewMsgRespondService creates a new MsgRespondService instance
func NewMsgRespondService(requestID, provider, result, output string) *MsgRespondService {
	return &MsgRespondService{
		RequestId: requestID,
		Provider:  provider,
		Result:    result,
		Output:    output,
	}
}

// Route implements Msg.
func (msg MsgRespondService) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgRespondService) Type() string { return TypeMsgRespondService }

// GetSignBytes implements Msg.
func (msg MsgRespondService) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgRespondService) ValidateBasic() error {
	if err := ValidateProvider(msg.Provider); err != nil {
		return err
	}

	if err := ValidateRequestID(msg.RequestId); err != nil {
		return err
	}

	if err := ValidateResponseResult(msg.Result); err != nil {
		return err
	}

	result, err := ParseResult(msg.Result)
	if err != nil {
		return err
	}
	return ValidateOutput(result.Code, msg.Output)
}

// GetSigners implements Msg.
func (msg MsgRespondService) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ______________________________________________________________________

// NewMsgPauseRequestContext creates a new MsgPauseRequestContext instance
func NewMsgPauseRequestContext(requestContextID, consumer string) *MsgPauseRequestContext {
	return &MsgPauseRequestContext{
		RequestContextId: requestContextID,
		Consumer:         consumer,
	}
}

// Route implements Msg.
func (msg MsgPauseRequestContext) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgPauseRequestContext) Type() string { return TypeMsgPauseRequestContext }

// GetSignBytes implements Msg.
func (msg MsgPauseRequestContext) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgPauseRequestContext) ValidateBasic() error {
	if err := ValidateConsumer(msg.Consumer); err != nil {
		return err
	}
	return ValidateContextID(msg.RequestContextId)
}

// GetSigners implements Msg.
func (msg MsgPauseRequestContext) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ______________________________________________________________________

// NewMsgStartRequestContext creates a new MsgStartRequestContext instance
func NewMsgStartRequestContext(requestContextID, consumer string) *MsgStartRequestContext {
	return &MsgStartRequestContext{
		RequestContextId: requestContextID,
		Consumer:         consumer,
	}
}

// Route implements Msg.
func (msg MsgStartRequestContext) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgStartRequestContext) Type() string { return TypeMsgStartRequestContext }

// GetSignBytes implements Msg.
func (msg MsgStartRequestContext) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgStartRequestContext) ValidateBasic() error {
	if err := ValidateConsumer(msg.Consumer); err != nil {
		return err
	}
	return ValidateContextID(msg.RequestContextId)
}

// GetSigners implements Msg.
func (msg MsgStartRequestContext) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ______________________________________________________________________

// NewMsgKillRequestContext creates a new MsgKillRequestContext instance
func NewMsgKillRequestContext(requestContextID, consumer string) *MsgKillRequestContext {
	return &MsgKillRequestContext{
		RequestContextId: requestContextID,
		Consumer:         consumer,
	}
}

// Route implements Msg.
func (msg MsgKillRequestContext) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgKillRequestContext) Type() string { return TypeMsgKillRequestContext }

// GetSignBytes implements Msg.
func (msg MsgKillRequestContext) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgKillRequestContext) ValidateBasic() error {
	if err := ValidateConsumer(msg.Consumer); err != nil {
		return err
	}
	return ValidateContextID(msg.RequestContextId)
}

// GetSigners implements Msg.
func (msg MsgKillRequestContext) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ______________________________________________________________________

// NewMsgUpdateRequestContext creates a new MsgUpdateRequestContext instance
func NewMsgUpdateRequestContext(
	requestContextID string,
	providers []string,
	serviceFeeCap sdk.Coins,
	timeout int64,
	repeatedFrequency uint64,
	repeatedTotal int64,
	consumer string,
) *MsgUpdateRequestContext {
	return &MsgUpdateRequestContext{
		RequestContextId:  requestContextID,
		Providers:         providers,
		ServiceFeeCap:     serviceFeeCap,
		Timeout:           timeout,
		RepeatedFrequency: repeatedFrequency,
		RepeatedTotal:     repeatedTotal,
		Consumer:          consumer,
	}
}

// Route implements Msg.
func (msg MsgUpdateRequestContext) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgUpdateRequestContext) Type() string { return TypeMsgUpdateRequestContext }

// GetSignBytes implements Msg.
func (msg MsgUpdateRequestContext) GetSignBytes() []byte {
	if len(msg.Providers) == 0 {
		msg.Providers = nil
	}

	if msg.ServiceFeeCap.Empty() {
		msg.ServiceFeeCap = nil
	}

	b := ModuleCdc.MustMarshalJSON(&msg)

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateRequestContext) ValidateBasic() error {
	if err := ValidateConsumer(msg.Consumer); err != nil {
		return err
	}

	if err := ValidateContextID(msg.RequestContextId); err != nil {
		return err
	}

	pds := make([]sdk.AccAddress, len(msg.Providers))
	for i, provider := range msg.Providers {
		pd, err := sdk.AccAddressFromBech32(provider)
		if err != nil {
			return err
		}
		pds[i] = pd
	}

	return ValidateRequestContextUpdating(
		pds,
		msg.ServiceFeeCap,
		msg.Timeout,
		msg.RepeatedFrequency,
		msg.RepeatedTotal,
	)
}

// GetSigners implements Msg.
func (msg MsgUpdateRequestContext) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Consumer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ______________________________________________________________________

// NewMsgWithdrawEarnedFees creates a new MsgWithdrawEarnedFees instance
func NewMsgWithdrawEarnedFees(owner, provider string) *MsgWithdrawEarnedFees {
	return &MsgWithdrawEarnedFees{
		Owner:    owner,
		Provider: provider,
	}
}

// Route implements Msg.
func (msg MsgWithdrawEarnedFees) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgWithdrawEarnedFees) Type() string { return TypeMsgWithdrawEarnedFees }

// GetSignBytes implements Msg.
func (msg MsgWithdrawEarnedFees) GetSignBytes() []byte {
	b := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgWithdrawEarnedFees) ValidateBasic() error {
	return ValidateOwner(msg.Owner)
}

// GetSigners implements Msg.
func (msg MsgWithdrawEarnedFees) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}
