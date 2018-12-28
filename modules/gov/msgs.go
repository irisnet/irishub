package gov

import (
	"fmt"

	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/irishub/modules/params"
	govtypes "github.com/irisnet/irishub/types/gov"
)

// name to idetify transaction types
const MsgRoute = "gov"

var _, _, _, _ sdk.Msg = MsgSubmitProposal{}, MsgSubmitTxTaxUsageProposal{}, MsgDeposit{}, MsgVote{}

//-----------------------------------------------------------
// MsgSubmitProposal
type MsgSubmitProposal struct {
	Title          string         `json:"title"`           //  Title of the proposal
	Description    string         `json:"description"`     //  Description of the proposal
	ProposalType   govtypes.ProposalKind   `json:"proposal_type"`   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       sdk.AccAddress `json:"proposer"`        //  Address of the proposer
	InitialDeposit sdk.Coins      `json:"initial_deposit"` //  Initial deposit paid by sender. Must be strictly positive.
	////////////////////  iris begin  ///////////////////////////
	Param govtypes.Param
	////////////////////  iris end  /////////////////////////////
}

func NewMsgSubmitProposal(title string, description string, proposalType govtypes.ProposalKind, proposer sdk.AccAddress, initialDeposit sdk.Coins, param govtypes.Param) MsgSubmitProposal {
	return MsgSubmitProposal{
		Title:          title,
		Description:    description,
		ProposalType:   proposalType,
		Proposer:       proposer,
		InitialDeposit: initialDeposit,
		////////////////////  iris begin  ///////////////////////////
		Param: param,
		////////////////////  iris end  /////////////////////////////
	}
}

//nolint
func (msg MsgSubmitProposal) Route() string { return MsgRoute }
func (msg MsgSubmitProposal) Type() string  { return "submit_proposal" }

// Implements Msg.
func (msg MsgSubmitProposal) ValidateBasic() sdk.Error {
	if len(msg.Title) == 0 {
		return govtypes.ErrInvalidTitle(govtypes.DefaultCodespace, msg.Title) // TODO: Proper Error
	}
	if len(msg.Description) == 0 {
		return govtypes.ErrInvalidDescription(govtypes.DefaultCodespace, msg.Description) // TODO: Proper Error
	}
	if !govtypes.ValidProposalType(msg.ProposalType) {
		return govtypes.ErrInvalidProposalType(govtypes.DefaultCodespace, msg.ProposalType)
	}
	if len(msg.Proposer) == 0 {
		return sdk.ErrInvalidAddress(msg.Proposer.String())
	}
	if !msg.InitialDeposit.IsValid() {
		return sdk.ErrInvalidCoins(msg.InitialDeposit.String())
	}
	if !msg.InitialDeposit.IsNotNegative() {
		return sdk.ErrInvalidCoins(msg.InitialDeposit.String())
	}
	////////////////////  iris begin  ///////////////////////////
	if msg.ProposalType == govtypes.ProposalTypeParameterChange {

		if msg.Param.Op != govtypes.Update && msg.Param.Op != govtypes.Insert {
			return govtypes.ErrInvalidParamOp(govtypes.DefaultCodespace, msg.Param.Op)
		}

		if p, ok := params.ParamMapping[msg.Param.Key]; ok {
			return p.Valid(msg.Param.Value)
		} else {
			return govtypes.ErrInvalidParam(govtypes.DefaultCodespace)
		}
	}
	////////////////////  iris end  /////////////////////////////
	return nil
}

func (msg MsgSubmitProposal) String() string {
	return fmt.Sprintf("MsgSubmitProposal{%s, %s, %s, %v}", msg.Title, msg.Description, msg.ProposalType, msg.InitialDeposit)
}

// Implements Msg.
func (msg MsgSubmitProposal) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg MsgSubmitProposal) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proposer}
}

type MsgSubmitSoftwareUpgradeProposal struct {
	MsgSubmitProposal
	Version      uint64        `json:"version"`
	Software     string        `json:"software"`
	SwitchHeight uint64        `json:"switch_height"`
}

func NewMsgSubmitSoftwareUpgradeProposal(msgSubmitProposal MsgSubmitProposal, version uint64, software string, switchHeight uint64) MsgSubmitSoftwareUpgradeProposal {
	return MsgSubmitSoftwareUpgradeProposal {
		MsgSubmitProposal: msgSubmitProposal,
		Version:             version,
		Software:       software,
		SwitchHeight:           switchHeight,
	}
}

func (msg MsgSubmitSoftwareUpgradeProposal) ValidateBasic() sdk.Error {
	err := msg.MsgSubmitProposal.ValidateBasic()
	if err != nil {
		return err
	}
	return nil
}

func (msg MsgSubmitSoftwareUpgradeProposal) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

type MsgSubmitTxTaxUsageProposal struct {
	MsgSubmitProposal
	Usage       govtypes.UsageType      `json:"usage"`
	DestAddress sdk.AccAddress `json:"dest_address"`
	Percent     sdk.Dec        `json:"percent"`
}

func NewMsgSubmitTaxUsageProposal(msgSubmitProposal MsgSubmitProposal, usage govtypes.UsageType, destAddress sdk.AccAddress, percent sdk.Dec) MsgSubmitTxTaxUsageProposal {
	return MsgSubmitTxTaxUsageProposal{
		MsgSubmitProposal: msgSubmitProposal,
		Usage:             usage,
		DestAddress:       destAddress,
		Percent:           percent,
	}
}

func (msg MsgSubmitTxTaxUsageProposal) ValidateBasic() sdk.Error {
	err := msg.MsgSubmitProposal.ValidateBasic()
	if err != nil {
		return err
	}
	if !govtypes.ValidUsageType(msg.Usage) {
		return govtypes.ErrInvalidUsageType(govtypes.DefaultCodespace, msg.Usage)
	}
	if msg.Usage != govtypes.UsageTypeBurn && len(msg.DestAddress) == 0 {
		return sdk.ErrInvalidAddress(msg.DestAddress.String())
	}
	if msg.Percent.LTE(sdk.NewDec(0)) || msg.Percent.GT(sdk.NewDec(1)) {
		return govtypes.ErrInvalidPercent(govtypes.DefaultCodespace, msg.Percent)
	}
	return nil
}

func (msg MsgSubmitTxTaxUsageProposal) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

//-----------------------------------------------------------
// MsgDeposit
type MsgDeposit struct {
	ProposalID uint64         `json:"proposal_id"` // ID of the proposal
	Depositor  sdk.AccAddress `json:"depositor"`   // Address of the depositor
	Amount     sdk.Coins      `json:"amount"`      // Coins to add to the proposal's deposit
}

func NewMsgDeposit(depositor sdk.AccAddress, proposalID uint64, amount sdk.Coins) MsgDeposit {
	return MsgDeposit{
		ProposalID: proposalID,
		Depositor:  depositor,
		Amount:     amount,
	}
}

// Implements Msg.
// nolint
func (msg MsgDeposit) Route() string { return MsgRoute }
func (msg MsgDeposit) Type() string  { return "deposit" }

// Implements Msg.
func (msg MsgDeposit) ValidateBasic() sdk.Error {
	if len(msg.Depositor) == 0 {
		return sdk.ErrInvalidAddress(msg.Depositor.String())
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	if !msg.Amount.IsNotNegative() {
		return sdk.ErrInvalidCoins(msg.Amount.String())
	}
	if msg.ProposalID < 0 {
		return govtypes.ErrUnknownProposal(govtypes.DefaultCodespace, msg.ProposalID)
	}
	return nil
}

func (msg MsgDeposit) String() string {
	return fmt.Sprintf("MsgDeposit{%s=>%v: %v}", msg.Depositor, msg.ProposalID, msg.Amount)
}

// Implements Msg.
func (msg MsgDeposit) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg MsgDeposit) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgDeposit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Depositor}
}

//-----------------------------------------------------------
// MsgVote
type MsgVote struct {
	ProposalID uint64         `json:"proposal_id"` // ID of the proposal
	Voter      sdk.AccAddress `json:"voter"`       //  address of the voter
	Option     govtypes.VoteOption     `json:"option"`      //  option from OptionSet chosen by the voter
}

func NewMsgVote(voter sdk.AccAddress, proposalID uint64, option govtypes.VoteOption) MsgVote {
	return MsgVote{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     option,
	}
}

// Implements Msg.
// nolint
func (msg MsgVote) Route() string { return MsgRoute }
func (msg MsgVote) Type() string  { return "vote" }

// Implements Msg.
func (msg MsgVote) ValidateBasic() sdk.Error {
	if len(msg.Voter.Bytes()) == 0 {
		return sdk.ErrInvalidAddress(msg.Voter.String())
	}
	if msg.ProposalID < 0 {
		return govtypes.ErrUnknownProposal(govtypes.DefaultCodespace, msg.ProposalID)
	}
	if !govtypes.ValidVoteOption(msg.Option) {
		return govtypes.ErrInvalidVote(govtypes.DefaultCodespace, msg.Option)
	}
	return nil
}

func (msg MsgVote) String() string {
	return fmt.Sprintf("MsgVote{%v - %s}", msg.ProposalID, msg.Option)
}

// Implements Msg.
func (msg MsgVote) Get(key interface{}) (value interface{}) {
	return nil
}

// Implements Msg.
func (msg MsgVote) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// Implements Msg.
func (msg MsgVote) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Voter}
}
