package types

import (
	"fmt"
	"regexp"

	sdk "github.com/irisnet/irishub/types"
)

const (
	MsgRoute = "service" // route for service msgs

	TypeMsgDefineService        = "define_service"         // type for MsgDefineService
	TypeMsgBindService          = "bind_service"           // type for MsgBindService
	TypeMsgUpdateServiceBinding = "update_service_binding" // type for MsgUpdateServiceBinding
	TypeMsgSetWithdrawAddress   = "set_withdraw_address"   // type for MsgSetWithdrawAddress
	TypeMsgDisableService       = "disable_service"        // type for MsgDisableService
	TypeMsgEnableService        = "enable_service"         // type for MsgEnableService
	TypeMsgRefundServiceDeposit = "refund_service_deposit" // type for MsgRefundServiceDeposit
	TypeMsgSvcRequest           = "call_service"           // type for MsgSvcRequest
	TypeMsgSvcResponse          = "respond_service"        // type for MsgSvcResponse
	TypeMsgSvcRefundFees        = "refund_service_fees"    // type for MsgSvcRefundFees
	TypeMsgSvcWithdrawFees      = "withdraw_service_fees"  // type for MsgSvcWithdrawFees
	TypeMsgSvcWithdrawTax       = "withdraw_service_tax"   // type for MsgSvcWithdrawTax

	MaxNameLength        = 70  // max length of the service name
	MaxDescriptionLength = 280 // max length of the service and author description
	MaxTagsNum           = 10  // max total number of the tags
	MaxTagLength         = 70  // max length of the tag
)

// the service name only accepts alphanumeric characters, _ and -
var reServiceName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

var (
	_ sdk.Msg = MsgDefineService{}
	_ sdk.Msg = &MsgBindService{}
	_ sdk.Msg = &MsgUpdateServiceBinding{}
	_ sdk.Msg = &MsgSetWithdrawAddress{}
	_ sdk.Msg = &MsgDisableService{}
	_ sdk.Msg = &MsgEnableService{}
	_ sdk.Msg = &MsgRefundServiceDeposit{}
	_ sdk.Msg = &MsgSvcRequest{}
	_ sdk.Msg = &MsgSvcResponse{}
	_ sdk.Msg = &MsgSvcRefundFees{}
	_ sdk.Msg = &MsgSvcWithdrawFees{}
	_ sdk.Msg = &MsgSvcWithdrawTax{}
)

//______________________________________________________________________

// MsgDefineService defines a message to define a service
type MsgDefineService struct {
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Tags              []string       `json:"tags"`
	Author            sdk.AccAddress `json:"author"`
	AuthorDescription string         `json:"author_description"`
	Schemas           string         `json:"schemas"`
}

// NewMsgDefineService creates a new MsgDefineService instance
func NewMsgDefineService(name, description string, tags []string, author sdk.AccAddress, authorDescription, schemas string) MsgDefineService {
	return MsgDefineService{
		Name:              name,
		Description:       description,
		Tags:              tags,
		Author:            author,
		AuthorDescription: authorDescription,
		Schemas:           schemas,
	}
}

// Route implements Msg
func (msg MsgDefineService) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgDefineService) Type() string { return TypeMsgDefineService }

// GetSignBytes implements Msg
func (msg MsgDefineService) GetSignBytes() []byte {
	if len(msg.Tags) == 0 {
		msg.Tags = nil
	}

	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg
func (msg MsgDefineService) ValidateBasic() sdk.Error {
	if len(msg.Author) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "author missing")
	}

	if !validServiceName(msg.Name) {
		return ErrInvalidServiceName(DefaultCodespace, msg.Name)
	}

	if err := ensureServiceDefLength(msg); err != nil {
		return err
	}

	if len(msg.Schemas) == 0 {
		return ErrInvalidSchemas(DefaultCodespace, "schemas missing")
	}

	return ValidateServiceSchemas(msg.Schemas)
}

// GetSigners implements Msg
func (msg MsgDefineService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Author}
}

//______________________________________________________________________

// MsgBindService defines a message to bind a service
type MsgBindService struct {
	ServiceName     string         `json:"service_name"`
	Provider        sdk.AccAddress `json:"provider"`
	Deposit         sdk.Coins      `json:"deposit"`
	Pricing         string         `json:"pricing"`
	WithdrawAddress sdk.AccAddress `json:"withdraw_address"`
}

// NewMsgBindService creates a new MsgBindService instance
func NewMsgBindService(serviceName string, provider sdk.AccAddress, deposit sdk.Coins, pricing string, withdrawAddr sdk.AccAddress) MsgBindService {
	return MsgBindService{
		ServiceName:     serviceName,
		Provider:        provider,
		Deposit:         deposit,
		Pricing:         pricing,
		WithdrawAddress: withdrawAddr,
	}
}

// Route implements Msg.
func (msg MsgBindService) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgBindService) Type() string { return TypeMsgBindService }

// GetSignBytes implements Msg.
func (msg MsgBindService) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgBindService) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	if err := ValidateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if !validServiceCoins(msg.Deposit) {
		return ErrInvalidDeposit(DefaultCodespace, fmt.Sprintf("invalid deposit: %s", msg.Deposit))
	}

	if len(msg.Pricing) == 0 {
		return ErrInvalidPricing(DefaultCodespace, "pricing missing")
	}

	return validatePricing(msg.Pricing)
}

// GetSigners implements Msg.
func (msg MsgBindService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgUpdateServiceBinding defines a message to update a service binding
type MsgUpdateServiceBinding struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
	Deposit     sdk.Coins      `json:"deposit"`
	Pricing     string         `json:"pricing"`
}

// NewMsgUpdateServiceBinding creates a new MsgUpdateServiceBinding instance
func NewMsgUpdateServiceBinding(serviceName string, provider sdk.AccAddress, deposit sdk.Coins, pricing string) MsgUpdateServiceBinding {
	return MsgUpdateServiceBinding{
		ServiceName: serviceName,
		Provider:    provider,
		Deposit:     deposit,
		Pricing:     pricing,
	}
}

// Route implements Msg.
func (msg MsgUpdateServiceBinding) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgUpdateServiceBinding) Type() string { return TypeMsgUpdateServiceBinding }

// GetSignBytes implements Msg.
func (msg MsgUpdateServiceBinding) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgUpdateServiceBinding) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	if err := ValidateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if !msg.Deposit.Empty() && !validServiceCoins(msg.Deposit) {
		return ErrInvalidDeposit(DefaultCodespace, fmt.Sprintf("invalid deposit: %s", msg.Deposit))
	}

	if len(msg.Pricing) != 0 {
		return validatePricing(msg.Pricing)
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgUpdateServiceBinding) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSetWithdrawAddress defines a message to set the withdrawal address for a service binding
type MsgSetWithdrawAddress struct {
	ServiceName     string         `json:"service_name"`
	Provider        sdk.AccAddress `json:"provider"`
	WithdrawAddress sdk.AccAddress `json:"withdraw_address"`
}

// NewMsgSetWithdrawAddress creates a new MsgSetWithdrawAddress instance
func NewMsgSetWithdrawAddress(serviceName string, provider sdk.AccAddress, withdrawAddr sdk.AccAddress) MsgSetWithdrawAddress {
	return MsgSetWithdrawAddress{
		ServiceName:     serviceName,
		Provider:        provider,
		WithdrawAddress: withdrawAddr,
	}
}

// Route implements Msg.
func (msg MsgSetWithdrawAddress) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgSetWithdrawAddress) Type() string { return TypeMsgSetWithdrawAddress }

// GetSignBytes implements Msg.
func (msg MsgSetWithdrawAddress) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSetWithdrawAddress) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	if err := ValidateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if len(msg.WithdrawAddress) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "withdrawal address missing")
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgSetWithdrawAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgDisableService defines a message to disable a service binding
type MsgDisableService struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
}

// NewMsgDisableService creates a new MsgDisableService instance
func NewMsgDisableService(serviceName string, provider sdk.AccAddress) MsgDisableService {
	return MsgDisableService{
		ServiceName: serviceName,
		Provider:    provider,
	}
}

// Route implements Msg.
func (msg MsgDisableService) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgDisableService) Type() string { return TypeMsgDisableService }

// GetSignBytes implements Msg.
func (msg MsgDisableService) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgDisableService) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	return ValidateServiceName(msg.ServiceName)
}

// GetSigners implements Msg.
func (msg MsgDisableService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgEnableService defines a message to enable a service binding
type MsgEnableService struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
	Deposit     sdk.Coins      `json:"deposit"`
}

// NewMsgEnableService creates a new MsgEnableService instance
func NewMsgEnableService(serviceName string, provider sdk.AccAddress, deposit sdk.Coins) MsgEnableService {
	return MsgEnableService{
		ServiceName: serviceName,
		Provider:    provider,
		Deposit:     deposit,
	}
}

// Route implements Msg.
func (msg MsgEnableService) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgEnableService) Type() string { return TypeMsgEnableService }

// GetSignBytes implements Msg.
func (msg MsgEnableService) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgEnableService) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	if err := ValidateServiceName(msg.ServiceName); err != nil {
		return err
	}

	if !msg.Deposit.Empty() && !validServiceCoins(msg.Deposit) {
		return ErrInvalidDeposit(DefaultCodespace, fmt.Sprintf("invalid deposit: %s", msg.Deposit))
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgEnableService) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgRefundServiceDeposit defines a message to refund deposit from a service binding
type MsgRefundServiceDeposit struct {
	ServiceName string         `json:"service_name"`
	Provider    sdk.AccAddress `json:"provider"`
}

// NewMsgRefundServiceDeposit creates a new MsgRefundServiceDeposit instance
func NewMsgRefundServiceDeposit(serviceName string, provider sdk.AccAddress) MsgRefundServiceDeposit {
	return MsgRefundServiceDeposit{
		ServiceName: serviceName,
		Provider:    provider,
	}
}

// Route implements Msg.
func (msg MsgRefundServiceDeposit) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgRefundServiceDeposit) Type() string { return TypeMsgRefundServiceDeposit }

// GetSignBytes implements Msg.
func (msg MsgRefundServiceDeposit) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgRefundServiceDeposit) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	return ValidateServiceName(msg.ServiceName)
}

// GetSigners implements Msg.
func (msg MsgRefundServiceDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcRequest - struct for call a service
type MsgSvcRequest struct {
	DefChainID  string         `json:"def_chain_id"`
	DefName     string         `json:"def_name"`
	BindChainID string         `json:"bind_chain_id"`
	ReqChainID  string         `json:"req_chain_id"`
	MethodID    int16          `json:"method_id"`
	Provider    sdk.AccAddress `json:"provider"`
	Consumer    sdk.AccAddress `json:"consumer"`
	Input       []byte         `json:"input"`
	ServiceFee  sdk.Coins      `json:"service_fee"`
	Profiling   bool           `json:"profiling"`
}

// NewMsgSvcRequest constructs a MsgSvcRequest
func NewMsgSvcRequest(defChainID, defName, bindChainID, reqChainID string, consumer, provider sdk.AccAddress, methodID int16, input []byte, serviceFee sdk.Coins, profiling bool) MsgSvcRequest {
	return MsgSvcRequest{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		ReqChainID:  reqChainID,
		Consumer:    consumer,
		Provider:    provider,
		MethodID:    methodID,
		Input:       input,
		ServiceFee:  serviceFee,
		Profiling:   profiling,
	}
}

// Route implements Msg.
func (msg MsgSvcRequest) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgSvcRequest) Type() string { return TypeMsgSvcRequest }

// GetSignBytes implements Msg.
func (msg MsgSvcRequest) GetSignBytes() []byte {
	if len(msg.Input) == 0 {
		msg.Input = nil
	}
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSvcRequest) ValidateBasic() sdk.Error {
	if err := ValidateServiceName(msg.DefName); err != nil {
		return err
	}

	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}
	if len(msg.Consumer) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "consumer missing")
	}
	if !msg.ServiceFee.IsValidIrisAtto() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid service fee [%s]", msg.ServiceFee))
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgSvcRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgSvcResponse - struct for respond a service call
type MsgSvcResponse struct {
	ReqChainID string         `json:"req_chain_id"`
	RequestID  string         `json:"request_id"`
	Provider   sdk.AccAddress `json:"provider"`
	Output     []byte         `json:"output"`
	ErrorMsg   []byte         `json:"error_msg"`
}

// NewMsgSvcResponse constructs a MsgSvcResponse
func NewMsgSvcResponse(reqChainID string, requestID string, provider sdk.AccAddress, output, errorMsg []byte) MsgSvcResponse {
	return MsgSvcResponse{
		ReqChainID: reqChainID,
		RequestID:  requestID,
		Provider:   provider,
		Output:     output,
		ErrorMsg:   errorMsg,
	}
}

// Route implements Msg.
func (msg MsgSvcResponse) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgSvcResponse) Type() string { return TypeMsgSvcResponse }

// GetSignBytes implements Msg.
func (msg MsgSvcResponse) GetSignBytes() []byte {
	if len(msg.Output) == 0 {
		msg.Output = nil
	}
	if len(msg.ErrorMsg) == 0 {
		msg.ErrorMsg = nil
	}
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSvcResponse) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}

	_, _, _, err := ConvertRequestID(msg.RequestID)
	if err != nil {
		return ErrInvalidReqId(DefaultCodespace, msg.RequestID)
	}

	return nil
}

// GetSigners implements Msg.
func (msg MsgSvcResponse) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcRefundFees - struct for refund fees
type MsgSvcRefundFees struct {
	Consumer sdk.AccAddress `json:"consumer"`
}

// NewMsgSvcRefundFees constructs a MsgSvcRefundFees
func NewMsgSvcRefundFees(consumer sdk.AccAddress) MsgSvcRefundFees {
	return MsgSvcRefundFees{Consumer: consumer}
}

// Route implements Msg.
func (msg MsgSvcRefundFees) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgSvcRefundFees) Type() string { return TypeMsgSvcRefundFees }

// GetSignBytes implements Msg.
func (msg MsgSvcRefundFees) GetSignBytes() []byte {
	b := msgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSvcRefundFees) ValidateBasic() sdk.Error {
	if len(msg.Consumer) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "consumer missing")
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgSvcRefundFees) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Consumer}
}

//______________________________________________________________________

// MsgSvcWithdrawFees - struct for withdraw fees
type MsgSvcWithdrawFees struct {
	Provider sdk.AccAddress `json:"provider"`
}

// NewMsgSvcWithdrawFees constructs a MsgSvcWithdrawFees
func NewMsgSvcWithdrawFees(provider sdk.AccAddress) MsgSvcWithdrawFees {
	return MsgSvcWithdrawFees{Provider: provider}
}

// Route implements Msg.
func (msg MsgSvcWithdrawFees) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgSvcWithdrawFees) Type() string { return TypeMsgSvcWithdrawFees }

// GetSignBytes implements Msg.
func (msg MsgSvcWithdrawFees) GetSignBytes() []byte {
	b := msgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSvcWithdrawFees) ValidateBasic() sdk.Error {
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgSvcWithdrawFees) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcWithdrawTax - struct for withdraw tax
type MsgSvcWithdrawTax struct {
	Trustee     sdk.AccAddress `json:"trustee"`
	DestAddress sdk.AccAddress `json:"dest_address"`
	Amount      sdk.Coins      `json:"amount"`
}

// NewMsgSvcWithdrawTax constructs a MsgSvcWithdrawTax
func NewMsgSvcWithdrawTax(trustee, destAddress sdk.AccAddress, amount sdk.Coins) MsgSvcWithdrawTax {
	return MsgSvcWithdrawTax{
		Trustee:     trustee,
		DestAddress: destAddress,
		Amount:      amount,
	}
}

// Route implements Msg.
func (msg MsgSvcWithdrawTax) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgSvcWithdrawTax) Type() string { return TypeMsgSvcWithdrawTax }

// GetSignBytes implements Msg.
func (msg MsgSvcWithdrawTax) GetSignBytes() []byte {
	b := msgCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSvcWithdrawTax) ValidateBasic() sdk.Error {
	if len(msg.Trustee) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "trustee missing")
	}
	if len(msg.DestAddress) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "destination address missing")
	}
	if !msg.Amount.IsValidIrisAtto() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid withdraw amount [%s]", msg.Amount))
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgSvcWithdrawTax) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Trustee}
}

//______________________________________________________________________

// ValidateServiceName validates the service name
func ValidateServiceName(name string) sdk.Error {
	if !validServiceName(name) {
		return ErrInvalidServiceName(DefaultCodespace, name)
	}

	return ensureServiceNameLength(name)
}

func validServiceName(name string) bool {
	return reServiceName.MatchString(name)
}

func ensureServiceNameLength(name string) sdk.Error {
	if len(name) > MaxNameLength {
		return ErrInvalidServiceName(DefaultCodespace, name)
	}

	return nil
}

func ensureServiceDefLength(msg MsgDefineService) sdk.Error {
	if err := ensureServiceNameLength(msg.Name); err != nil {
		return err
	}

	if len(msg.Description) > MaxDescriptionLength {
		return ErrInvalidLength(DefaultCodespace, fmt.Sprintf("invalid description length; got: %d, max: %d", len(msg.Description), MaxDescriptionLength))
	}

	if len(msg.Tags) > MaxTagsNum {
		return ErrInvalidLength(DefaultCodespace, fmt.Sprintf("invalid tags size; got: %d, max: %d", len(msg.Tags), MaxTagsNum))
	}

	for i, tag := range msg.Tags {
		if len(tag) > MaxTagLength {
			return ErrInvalidLength(DefaultCodespace, fmt.Sprintf("invalid tag[%d] length; got: %d, max: %d", i, len(tag), MaxTagLength))
		}
	}

	if len(msg.AuthorDescription) > MaxDescriptionLength {
		return ErrInvalidLength(DefaultCodespace, fmt.Sprintf("invalid author description length; got: %d, max: %d", len(msg.AuthorDescription), MaxDescriptionLength))
	}

	return nil
}

func validatePricing(pricing string) sdk.Error {
	if err := ValidateBindingPricing(pricing); err != nil {
		return err
	}

	p, err := ParsePricing(pricing)
	if err != nil {
		return err
	}

	if !validServiceCoins(p.Price) {
		return ErrInvalidPricing(DefaultCodespace, fmt.Sprintf("invalid pricing coins: %s", p.Price))
	}

	return nil
}

func validServiceCoins(coins sdk.Coins) bool {
	return coins.IsValidIrisAtto()
}
