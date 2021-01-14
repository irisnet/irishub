package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

	MaxNameLength        = 70  // maximum length of the service name
	MaxDescriptionLength = 280 // maximum length of the service and author description
	MaxTagsNum           = 10  // maximum total number of the tags
	MaxTagLength         = 70  // maximum length of the tag

	MaxProvidersNum = 10 // maximum total number of the providers to request
)

// the service name only accepts alphanumeric characters, _ and -, beginning with alpha character
var reServiceName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

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
	if _, err := sdk.AccAddressFromBech32(msg.Author); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid author address (%s)", err)
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
	return ValidateBindingPricing(msg.Pricing)
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
		return ValidateBindingPricing(msg.Pricing)
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
	superMode bool,
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
		SuperMode:         superMode,
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

func ValidateAuthor(author string) error {
	if _, err := sdk.AccAddressFromBech32(author); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid author address (%s)", err)
	}
	return nil
}

// ValidateServiceName validates the service name
func ValidateServiceName(name string) error {
	if !reServiceName.MatchString(name) || len(name) > MaxNameLength {
		return sdkerrors.Wrap(ErrInvalidServiceName, name)
	}
	return nil
}

func ValidateTags(tags []string) error {
	if len(tags) > MaxTagsNum {
		return sdkerrors.Wrap(ErrInvalidTags, fmt.Sprintf("invalid tags size; got: %d, max: %d", len(tags), MaxTagsNum))
	}
	if HasDuplicate(tags) {
		return sdkerrors.Wrap(ErrInvalidTags, "duplicate tag")
	}
	for i, tag := range tags {
		if len(tag) == 0 {
			return sdkerrors.Wrap(ErrInvalidTags, fmt.Sprintf("invalid tag[%d] length: tag must not be empty", i))
		}
		if len(tag) > MaxTagLength {
			return sdkerrors.Wrap(ErrInvalidTags, fmt.Sprintf("invalid tag[%d] length; got: %d, max: %d", i, len(tag), MaxTagLength))
		}
	}
	return nil
}

func ValidateServiceDescription(svcDescription string) error {
	if len(svcDescription) > MaxDescriptionLength {
		return sdkerrors.Wrap(ErrInvalidDescription, fmt.Sprintf("invalid service description length; got: %d, max: %d", len(svcDescription), MaxDescriptionLength))
	}
	return nil
}

func ValidateAuthorDescription(authorDescription string) error {
	if len(authorDescription) > MaxDescriptionLength {
		return sdkerrors.Wrap(ErrInvalidDescription, fmt.Sprintf("invalid author description length; got: %d, max: %d", len(authorDescription), MaxDescriptionLength))
	}
	return nil
}

func ValidateProvider(provider string) error {
	if _, err := sdk.AccAddressFromBech32(provider); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid provider address (%s)", err)
	}
	return nil
}

func ValidateOwner(owner string) error {
	if _, err := sdk.AccAddressFromBech32(owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	return nil
}

func ValidateServiceDeposit(deposit sdk.Coins) error {
	if !deposit.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid deposit")
	}
	if deposit.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "invalid deposit")
	}
	return nil
}

func ValidateQoS(qos uint64) error {
	if qos == 0 {
		return sdkerrors.Wrap(ErrInvalidQoS, "qos must be greater than 0")
	}
	return nil
}

func ValidateOptions(options string) error {
	if !json.Valid([]byte(options)) {
		return sdkerrors.Wrap(ErrInvalidOptions, "options is not valid JSON")
	}
	return nil
}

func ValidateWithdrawAddress(withdrawAddress string) error {
	if _, err := sdk.AccAddressFromBech32(withdrawAddress); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid withdrawal address (%s)", err)
	}
	return nil
}

// ______________________________________________________________________

// ValidateRequest validates the request params
func ValidateRequest(
	serviceName string,
	serviceFeeCap sdk.Coins,
	providers []sdk.AccAddress,
	input string,
	timeout int64,
	repeated bool,
	repeatedFrequency uint64,
	repeatedTotal int64,
) error {
	if err := ValidateServiceName(serviceName); err != nil {
		return err
	}
	if err := ValidateServiceFeeCap(serviceFeeCap); err != nil {
		return err
	}
	if err := ValidateProviders(providers); err != nil {
		return err
	}
	if err := ValidateInput(input); err != nil {
		return err
	}
	if timeout <= 0 {
		return sdkerrors.Wrapf(ErrInvalidTimeout, "timeout [%d] must be greater than 0", timeout)
	}
	if repeated {
		if repeatedFrequency > 0 && repeatedFrequency < uint64(timeout) {
			return sdkerrors.Wrapf(ErrInvalidRepeatedFreq, "repeated frequency [%d] must not be less than timeout [%d]", repeatedFrequency, timeout)
		}
		if repeatedTotal < -1 || repeatedTotal == 0 {
			return sdkerrors.Wrapf(ErrInvalidRepeatedTotal, "repeated total number [%d] must be greater than 0 or equal to -1", repeatedTotal)
		}
	}
	return nil
}

// ValidateRequestContextUpdating validates the request context updating operation
func ValidateRequestContextUpdating(
	providers []sdk.AccAddress,
	serviceFeeCap sdk.Coins,
	timeout int64,
	repeatedFrequency uint64,
	repeatedTotal int64,
) error {
	if err := ValidateProvidersCanEmpty(providers); err != nil {
		return err
	}
	if !serviceFeeCap.Empty() {
		if err := ValidateServiceFeeCap(serviceFeeCap); err != nil {
			return err
		}
	}
	if timeout < 0 {
		return sdkerrors.Wrapf(ErrInvalidTimeout, "timeout must not be less than 0: %d", timeout)
	}
	if timeout != 0 && repeatedFrequency != 0 && repeatedFrequency < uint64(timeout) {
		return sdkerrors.Wrapf(ErrInvalidRepeatedFreq, "frequency [%d] must not be less than timeout [%d]", repeatedFrequency, timeout)
	}
	if repeatedTotal < -1 {
		return sdkerrors.Wrapf(ErrInvalidRepeatedFreq, "repeated total number must not be less than -1: %d", repeatedTotal)
	}
	return nil
}

func ValidateConsumer(consumer string) error {
	if _, err := sdk.AccAddressFromBech32(consumer); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid consumer address (%s)", err)
	}
	return nil
}

func ValidateProviders(providers []sdk.AccAddress) error {
	if len(providers) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "providers missing")
	}
	if len(providers) > MaxProvidersNum {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "total number of the providers must not be greater than %d", MaxProvidersNum)
	}
	if err := checkDuplicateProviders(providers); err != nil {
		return err
	}
	return nil
}

func ValidateProvidersCanEmpty(providers []sdk.AccAddress) error {
	if len(providers) > MaxProvidersNum {
		return sdkerrors.Wrapf(ErrInvalidProviders, "total number of the providers must not be greater than %d", MaxProvidersNum)
	}
	if len(providers) > 0 {
		if err := checkDuplicateProviders(providers); err != nil {
			return err
		}
	}
	return nil
}

func ValidateServiceFeeCap(serviceFeeCap sdk.Coins) error {
	if !serviceFeeCap.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, fmt.Sprintf("invalid service fee cap: %s", serviceFeeCap))
	}
	return nil
}

func ValidateRequestID(reqID string) error {
	if len(reqID) != RequestIDLen {
		return sdkerrors.Wrapf(ErrInvalidRequestID, "length of the request ID must be %d", RequestIDLen)
	}
	if _, err := hex.DecodeString(reqID); err != nil {
		return sdkerrors.Wrap(ErrInvalidRequestID, "request ID must be a hex encoded string")
	}
	return nil
}

func ValidateContextID(contextID string) error {
	if len(contextID) != ContextIDLen {
		return sdkerrors.Wrapf(ErrInvalidRequestContextID, "length of the request context ID must be %d in bytes", ContextIDLen)
	}
	if _, err := hex.DecodeString(contextID); err != nil {
		return sdkerrors.Wrap(ErrInvalidRequestContextID, "request context ID must be a hex encoded string")
	}
	return nil
}

func ValidateInput(input string) error {
	if len(input) == 0 {
		return sdkerrors.Wrap(ErrInvalidRequestInput, "input missing")
	}

	if ValidateRequestInput(input) != nil {
		return sdkerrors.Wrap(ErrInvalidRequestInput, "invalid input")
	}

	return nil
}

func ValidateOutput(code ResultCode, output string) error {
	if code == ResultOK && len(output) == 0 {
		return sdkerrors.Wrapf(ErrInvalidResponse, "output must be specified when the result code is %v", ResultOK)
	}

	if code != ResultOK && len(output) != 0 {
		return sdkerrors.Wrapf(ErrInvalidResponse, "output should not be specified when the result code is not %v", ResultOK)
	}

	if len(output) > 0 && ValidateResponseOutput(output) != nil {
		return sdkerrors.Wrap(ErrInvalidResponse, "invalid output")
	}

	return nil
}

func checkDuplicateProviders(providers []sdk.AccAddress) error {
	providerArr := make([]string, len(providers))

	for i, provider := range providers {
		providerArr[i] = provider.String()
	}

	if HasDuplicate(providerArr) {
		return sdkerrors.Wrap(ErrInvalidProviders, "there exists duplicate providers")
	}

	return nil
}
