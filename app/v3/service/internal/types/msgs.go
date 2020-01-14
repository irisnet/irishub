package types

import (
	"fmt"
	"regexp"

	sdk "github.com/irisnet/irishub/types"
)

const (
	MsgRoute = "service" // route for service msgs

	TypeMsgDefineService    = "define_service"         // type for MsgDefineService
	TypeMsgSvcBind          = "bind_service"           // type for MsgSvcBind
	TypeMsgSvcBindingUpdate = "update_service_binding" // type for MsgSvcBindingUpdate
	TypeMsgSvcDisable       = "disable_service"        // type for MsgSvcDisable
	TypeMsgSvcEnable        = "enable_service"         // type for MsgSvcEnable
	TypeMsgSvcRefundDeposit = "refund_service_deposit" // type for MsgSvcRefundDeposit
	TypeMsgSvcRequest       = "call_service"           // type for MsgSvcRequest
	TypeMsgSvcResponse      = "respond_service"        // type for MsgSvcResponse
	TypeMsgSvcRefundFees    = "refund_service_fees"    // type for MsgSvcRefundFees
	TypeMsgSvcWithdrawFees  = "withdraw_service_fees"  // type for MsgSvcWithdrawFees
	TypeMsgSvcWithdrawTax   = "withdraw_service_tax"   // type for MsgSvcWithdrawTax

	MaxNameLength        = 70  // max length of the service name
	MaxDescriptionLength = 280 // max length of the service and author description
	MaxTagsNum           = 10  // max total number of the tags
	MaxTagLength         = 70  // max length of the tag
)

// the service name only accepts alphanumeric characters, _ and -
var reServiceName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

var (
	_ sdk.Msg = MsgDefineService{}
	_ sdk.Msg = MsgSvcBind{}
	_ sdk.Msg = MsgSvcBindingUpdate{}
	_ sdk.Msg = MsgSvcDisable{}
	_ sdk.Msg = MsgSvcEnable{}
	_ sdk.Msg = MsgSvcRefundDeposit{}
	_ sdk.Msg = MsgSvcRequest{}
	_ sdk.Msg = MsgSvcResponse{}
	_ sdk.Msg = MsgSvcRefundFees{}
	_ sdk.Msg = MsgSvcWithdrawFees{}
	_ sdk.Msg = MsgSvcWithdrawTax{}
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

// MsgSvcBinding - struct for bind a service
type MsgSvcBind struct {
	DefName     string         `json:"def_name"`
	DefChainID  string         `json:"def_chain_id"`
	BindChainID string         `json:"bind_chain_id"`
	Provider    sdk.AccAddress `json:"provider"`
	BindingType BindingType    `json:"binding_type"`
	Deposit     sdk.Coins      `json:"deposit"`
	Prices      []sdk.Coin     `json:"price"`
	Level       Level          `json:"level"`
}

// NewMsgSvcBind constructs a MsgSvcBind
func NewMsgSvcBind(defChainID, defName, bindChainID string, provider sdk.AccAddress, bindingType BindingType, deposit sdk.Coins, prices []sdk.Coin, level Level) MsgSvcBind {
	return MsgSvcBind{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		Provider:    provider,
		BindingType: bindingType,
		Deposit:     deposit,
		Prices:      prices,
		Level:       level,
	}
}

// Route implements Msg.
func (msg MsgSvcBind) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgSvcBind) Type() string { return TypeMsgSvcBind }

// GetSignBytes implements Msg.
func (msg MsgSvcBind) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSvcBind) ValidateBasic() sdk.Error {
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}

	if err := ValidateServiceName(msg.DefName); err != nil {
		return err
	}

	if !validBindingType(msg.BindingType) {
		return ErrInvalidBindingType(DefaultCodespace, msg.BindingType)
	}
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}
	if !msg.Deposit.IsValidIrisAtto() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid deposit [%s]", msg.Deposit))
	}
	for _, price := range msg.Prices {
		if !price.IsValidIrisAtto() {
			return sdk.ErrInvalidCoins(fmt.Sprintf("invalid price [%s]", price))
		}
	}
	if !validLevel(msg.Level) {
		return ErrInvalidLevel(DefaultCodespace, msg.Level)
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgSvcBind) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcBindingUpdate - struct for update a service binding
type MsgSvcBindingUpdate struct {
	DefName     string         `json:"def_name"`
	DefChainID  string         `json:"def_chain_id"`
	BindChainID string         `json:"bind_chain_id"`
	Provider    sdk.AccAddress `json:"provider"`
	BindingType BindingType    `json:"binding_type"`
	Deposit     sdk.Coins      `json:"deposit" `
	Prices      []sdk.Coin     `json:"price" `
	Level       Level          `json:"level"`
}

// NewMsgSvcBindingUpdate constructs a MsgSvcBindingUpdate
func NewMsgSvcBindingUpdate(defChainID, defName, bindChainID string, provider sdk.AccAddress, bindingType BindingType, deposit sdk.Coins, prices []sdk.Coin, level Level) MsgSvcBindingUpdate {
	return MsgSvcBindingUpdate{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		Provider:    provider,
		BindingType: bindingType,
		Deposit:     deposit,
		Prices:      prices,
		Level:       level,
	}
}

// Route implements Msg.
func (msg MsgSvcBindingUpdate) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgSvcBindingUpdate) Type() string { return TypeMsgSvcBindingUpdate }

// GetSignBytes implements Msg.
func (msg MsgSvcBindingUpdate) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSvcBindingUpdate) ValidateBasic() sdk.Error {
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}

	if err := ValidateServiceName(msg.DefName); err != nil {
		return err
	}

	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}
	if msg.BindingType != 0x00 && !validBindingType(msg.BindingType) {
		return ErrInvalidBindingType(DefaultCodespace, msg.BindingType)
	}
	if !msg.Deposit.IsValidIrisAtto() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid deposit [%s]", msg.Deposit))
	}
	for _, price := range msg.Prices {
		if !price.IsValidIrisAtto() {
			return sdk.ErrInvalidCoins(fmt.Sprintf("invalid price [%s]", price))
		}
	}
	if !validUpdateLevel(msg.Level) {
		return ErrInvalidLevel(DefaultCodespace, msg.Level)
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgSvcBindingUpdate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcDisable - struct for disable a service binding
type MsgSvcDisable struct {
	DefName     string         `json:"def_name"`
	DefChainID  string         `json:"def_chain_id"`
	BindChainID string         `json:"bind_chain_id"`
	Provider    sdk.AccAddress `json:"provider"`
}

// NewMsgSvcDisable constructs a MsgSvcDisable
func NewMsgSvcDisable(defChainID, defName, bindChainID string, provider sdk.AccAddress) MsgSvcDisable {
	return MsgSvcDisable{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		Provider:    provider,
	}
}

// Route implements Msg.
func (msg MsgSvcDisable) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgSvcDisable) Type() string { return TypeMsgSvcDisable }

// GetSignBytes implements Msg.
func (msg MsgSvcDisable) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSvcDisable) ValidateBasic() sdk.Error {
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}

	if err := ValidateServiceName(msg.DefName); err != nil {
		return err
	}

	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgSvcDisable) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcEnable - struct for enable a service binding
type MsgSvcEnable struct {
	DefName     string         `json:"def_name"`
	DefChainID  string         `json:"def_chain_id"`
	BindChainID string         `json:"bind_chain_id"`
	Provider    sdk.AccAddress `json:"provider"`
	Deposit     sdk.Coins      `json:"deposit"`
}

// NewMsgSvcEnable constructs a MsgSvcEnable
func NewMsgSvcEnable(defChainID, defName, bindChainID string, provider sdk.AccAddress, deposit sdk.Coins) MsgSvcEnable {
	return MsgSvcEnable{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		Provider:    provider,
		Deposit:     deposit,
	}
}

// Route implements Msg.
func (msg MsgSvcEnable) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgSvcEnable) Type() string { return TypeMsgSvcEnable }

// GetSignBytes implements Msg.
func (msg MsgSvcEnable) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSvcEnable) ValidateBasic() sdk.Error {
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}

	if err := ValidateServiceName(msg.DefName); err != nil {
		return err
	}

	if !msg.Deposit.IsValidIrisAtto() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("invalid deposit [%s]", msg.Deposit))
	}
	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgSvcEnable) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcRefundDeposit - struct for refund deposit from a service binding
type MsgSvcRefundDeposit struct {
	DefName     string         `json:"def_name"`
	DefChainID  string         `json:"def_chain_id"`
	BindChainID string         `json:"bind_chain_id"`
	Provider    sdk.AccAddress `json:"provider"`
}

// NewMsgSvcRefundDeposit constructs a MsgSvcRefundDeposit
func NewMsgSvcRefundDeposit(defChainID, defName, bindChainID string, provider sdk.AccAddress) MsgSvcRefundDeposit {
	return MsgSvcRefundDeposit{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		Provider:    provider,
	}
}

// Route implements Msg.
func (msg MsgSvcRefundDeposit) Route() string { return MsgRoute }

// Type implements Msg.
func (msg MsgSvcRefundDeposit) Type() string { return TypeMsgSvcRefundDeposit }

// GetSignBytes implements Msg.
func (msg MsgSvcRefundDeposit) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// ValidateBasic implements Msg.
func (msg MsgSvcRefundDeposit) ValidateBasic() sdk.Error {
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}

	if err := ValidateServiceName(msg.DefName); err != nil {
		return err
	}

	if len(msg.Provider) == 0 {
		return ErrInvalidAddress(DefaultCodespace, "provider missing")
	}
	return nil
}

// GetSigners implements Msg.
func (msg MsgSvcRefundDeposit) GetSigners() []sdk.AccAddress {
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
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidBindChainId(DefaultCodespace)
	}
	if len(msg.ReqChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}

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
	if len(msg.ReqChainID) == 0 {
		return ErrInvalidReqChainId(DefaultCodespace)
	}
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

	if err := ensureServiceNameLength(name); err != nil {
		return err
	}

	return nil
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
