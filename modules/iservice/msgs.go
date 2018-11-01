package iservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/irisnet/irishub/tools/protoidl"
)

const (
	// name to idetify transaction types
	MsgType       = "iservice"
	outputPrivacy = "output_privacy"
	outputCached  = "output_cached"
	description   = "description"
)

var _ sdk.Msg = MsgSvcDef{}

//______________________________________________________________________

// MsgSvcDef - struct for define a service
type MsgSvcDef struct {
	SvcDef
}

func NewMsgSvcDef(name, chainId, description string, tags []string, author sdk.AccAddress, authorDescription, idlContent string, messaging MessagingType) MsgSvcDef {
	return MsgSvcDef{
		SvcDef{
			Name:              name,
			ChainId:           chainId,
			Description:       description,
			Tags:              tags,
			Author:            author,
			AuthorDescription: authorDescription,
			IDLContent:        idlContent,
			Messaging:         messaging,
		},
	}
}


func (msg MsgSvcDef) Route() string { return MsgType }
func (msg MsgSvcDef) Type() string {return "iservice definition"}

func (msg MsgSvcDef) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcDef) ValidateBasic() sdk.Error {
	if len(msg.ChainId) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}
	if len(msg.Name) == 0 {
		return ErrInvalidServiceName(DefaultCodespace)
	}
	if valid, err := validateTags(msg.Tags); !valid {
		return err
	}
	if len(msg.Author) == 0 {
		return ErrInvalidAuthor(DefaultCodespace)
	}
	if !validMessagingType(msg.Messaging) {
		return ErrInvalidMessagingType(DefaultCodespace, msg.Messaging)
	}

	if len(msg.IDLContent) == 0 {
		return ErrInvalidIDL(DefaultCodespace, "content is empty")
	}
	methods, err := protoidl.GetMethods(msg.IDLContent)
	if err != nil {
		return ErrInvalidIDL(DefaultCodespace, err.Error())
	}
	if valid, err := validateMethods(methods); !valid {
		return err
	}

	return nil
}

func (msg MsgSvcDef) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Author}
}

func validateMethods(methods []protoidl.Method) (bool, sdk.Error) {
	for _, method := range methods {
		if len(method.Name) == 0 {
			return false, ErrInvalidMethodName(DefaultCodespace)
		}
		if _, ok := method.Attributes[outputPrivacy]; ok {
			_, err := OutputPrivacyEnumFromString(method.Attributes[outputPrivacy])
			if err != nil {
				return false, ErrInvalidOutputPrivacyEnum(DefaultCodespace, method.Attributes[outputPrivacy])
			}
		}
		if _, ok := method.Attributes[outputCached]; ok {
			_, err := OutputCachedEnumFromString(method.Attributes[outputCached])
			if err != nil {
				return false, ErrInvalidOutputCachedEnum(DefaultCodespace, method.Attributes[outputCached])
			}
		}
	}
	return true, nil
}

func validateTags(tags []string) (bool, sdk.Error) {
	if len(tags) > iserviceParams.MaxTagsNum {
		return false, ErrMoreTags(DefaultCodespace)
	}
	if len(tags) > 0 {
		for i, tag := range tags {
			for _, tag1 := range tags[i+1:] {
				if tag == tag1 {
					return false, ErrDuplicateTags(DefaultCodespace)
				}
			}
		}
	}
	return true, nil
}

func methodToMethodProperty(index int, method protoidl.Method) (methodProperty MethodProperty, err sdk.Error) {
	// set default value
	opp := NoPrivacy
	opc := NoCached

	var err1 error
	if _, ok := method.Attributes[outputPrivacy]; ok {
		opp, err1 = OutputPrivacyEnumFromString(method.Attributes[outputPrivacy])
		if err1 != nil {
			return methodProperty, ErrInvalidOutputPrivacyEnum(DefaultCodespace, method.Attributes[outputPrivacy])
		}
	}
	if _, ok := method.Attributes[outputCached]; ok {
		opc, err1 = OutputCachedEnumFromString(method.Attributes[outputCached])
		if err != nil {
			return methodProperty, ErrInvalidOutputCachedEnum(DefaultCodespace, method.Attributes[outputCached])
		}
	}
	methodProperty = MethodProperty{
		ID:            index,
		Name:          method.Name,
		Description:   method.Attributes[description],
		OutputPrivacy: opp,
		OutputCached:  opc,
	}
	return
}

//______________________________________________________________________

// MsgSvcBinding - struct for bind a service
type MsgSvcBind struct {
	SvcBinding
}

func NewMsgSvcBind(defChainID, defName, bindChainID string, provider sdk.AccAddress, bindingType BindingType, deposit sdk.Coins, prices []sdk.Coin, level Level, expiration int64) MsgSvcBind {
	return MsgSvcBind{
		SvcBinding{
			DefChainID:  defChainID,
			DefName:     defName,
			BindChainID: bindChainID,
			Provider:    provider,
			BindingType: bindingType,
			Deposit:     deposit,
			Expiration:  expiration,
			Prices:      prices,
			Level:       level,
		},
	}
}

func (msg MsgSvcBind) Type() string {
	return MsgType
}

func (msg MsgSvcBind) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcBind) ValidateBasic() sdk.Error {
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}
	if len(msg.DefName) == 0 {
		return ErrInvalidServiceName(DefaultCodespace)
	}
	if !validBindingType(msg.BindingType) {
		return ErrInvalidBindingType(DefaultCodespace, msg.BindingType)
	}
	if len(msg.Provider) == 0 {
		sdk.ErrInvalidAddress(msg.Provider.String())
	}
	if !msg.Deposit.IsValid() {
		return sdk.ErrInvalidCoins(msg.Deposit.String())
	}
	if !msg.Deposit.IsNotNegative() {
		return sdk.ErrInvalidCoins(msg.Deposit.String())
	}
	for _, price := range msg.Prices {
		if !price.IsNotNegative() {
			return sdk.ErrInvalidCoins(price.String())
		}
	}
	if !validLevel(msg.Level) {
		return ErrInvalidLevel(DefaultCodespace, msg.Level)
	}
	return nil
}

func (msg MsgSvcBind) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}

//______________________________________________________________________

// MsgSvcBindingUpdate - struct for update a service binding
type MsgSvcBindingUpdate struct {
	SvcBinding
}

func NewMsgSvcBindingUpdate(defChainID, defName, bindChainID string, provider sdk.AccAddress, bindingType BindingType, deposit sdk.Coins, prices []sdk.Coin, level Level, expiration int64) MsgSvcBindingUpdate {
	return MsgSvcBindingUpdate{
		SvcBinding{
			DefChainID:  defChainID,
			DefName:     defName,
			BindChainID: bindChainID,
			Provider:    provider,
			BindingType: bindingType,
			Deposit:     deposit,
			Expiration:  expiration,
			Prices:      prices,
			Level:       level,
		},
	}
}

func (msg MsgSvcBindingUpdate) Type() string {
	return MsgType
}

func (msg MsgSvcBindingUpdate) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcBindingUpdate) ValidateBasic() sdk.Error {
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}
	if len(msg.DefName) == 0 {
		return ErrInvalidServiceName(DefaultCodespace)
	}
	if len(msg.Provider) == 0 {
		sdk.ErrInvalidAddress(msg.Provider.String())
	}
	if !validBindingType(msg.BindingType) {
		return ErrInvalidBindingType(DefaultCodespace, msg.BindingType)
	}
	if !msg.Deposit.IsValid() {
		return sdk.ErrInvalidCoins(msg.Deposit.String())
	}
	for _, price := range msg.Prices {
		if !price.IsNotNegative() {
			return sdk.ErrInvalidCoins(price.String())
		}
	}
	if !validUpdateLevel(msg.Level) {
		return ErrInvalidLevel(DefaultCodespace, msg.Level)
	}
	return nil
}

func (msg MsgSvcBindingUpdate) GetSigners() []sdk.AccAddress {
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

func NewMsgSvcRefundDeposit(defChainID, defName, bindChainID string, provider sdk.AccAddress) MsgSvcRefundDeposit {
	return MsgSvcRefundDeposit{
		DefChainID:  defChainID,
		DefName:     defName,
		BindChainID: bindChainID,
		Provider:    provider,
	}
}

func (msg MsgSvcRefundDeposit) Type() string {
	return MsgType
}

func (msg MsgSvcRefundDeposit) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSvcRefundDeposit) ValidateBasic() sdk.Error {
	if len(msg.DefChainID) == 0 {
		return ErrInvalidDefChainId(DefaultCodespace)
	}
	if len(msg.BindChainID) == 0 {
		return ErrInvalidChainId(DefaultCodespace)
	}
	if len(msg.DefName) == 0 {
		return ErrInvalidServiceName(DefaultCodespace)
	}
	if len(msg.Provider) == 0 {
		sdk.ErrInvalidAddress(msg.Provider.String())
	}
	return nil
}

func (msg MsgSvcRefundDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Provider}
}
