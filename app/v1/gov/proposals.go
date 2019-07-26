package gov

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdk "github.com/irisnet/irishub/types"
	"github.com/pkg/errors"
)

//-----------------------------------------------------------
// Proposal interface
type Proposal interface {
	GetProposalID() uint64
	SetProposalID(uint64)

	GetTitle() string
	SetTitle(string)

	GetDescription() string
	SetDescription(string)

	GetProposalType() ProposalKind
	SetProposalType(ProposalKind)

	GetStatus() ProposalStatus
	SetStatus(ProposalStatus)

	GetTallyResult() TallyResult
	SetTallyResult(TallyResult)

	GetSubmitTime() time.Time
	SetSubmitTime(time.Time)

	GetDepositEndTime() time.Time
	SetDepositEndTime(time.Time)

	GetTotalDeposit() sdk.Coins
	SetTotalDeposit(sdk.Coins)

	GetVotingStartTime() time.Time
	SetVotingStartTime(time.Time)

	GetVotingEndTime() time.Time
	SetVotingEndTime(time.Time)

	GetProposalLevel() ProposalLevel
	GetProposer() sdk.AccAddress

	String() string
	Validate(ctx sdk.Context, gk Keeper, verifyPropNum bool) sdk.Error
	Execute(ctx sdk.Context, gk Keeper) sdk.Error
}

//-----------------------------------------------------------
// Basic Proposals
type BasicProposal struct {
	ProposalID   uint64       `json:"proposal_id"`   //  ID of the proposal
	Title        string       `json:"title"`         //  Title of the proposal
	Description  string       `json:"description"`   //  Description of the proposal
	ProposalType ProposalKind `json:"proposal_type"` //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}

	Status      ProposalStatus `json:"proposal_status"` //  Status of the Proposal {Pending, Active, Passed, Rejected}
	TallyResult TallyResult    `json:"tally_result"`    //  Result of Tallys

	SubmitTime     time.Time `json:"submit_time"`      //  Time of the block where TxGovSubmitProposal was included
	DepositEndTime time.Time `json:"deposit_end_time"` // Time that the Proposal would expire if deposit amount isn't met
	TotalDeposit   sdk.Coins `json:"total_deposit"`    //  Current deposit on this proposal. Initial value is set at InitialDeposit

	VotingStartTime time.Time      `json:"voting_start_time"` //  Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	VotingEndTime   time.Time      `json:"voting_end_time"`   // Time that the VotingPeriod for this proposal will end and votes will be tallied
	Proposer        sdk.AccAddress `json:"proposer"`
}

func (bp BasicProposal) String() string {
	return fmt.Sprintf(`Proposal %d:
  Title:              %s
  Type:               %s
  Proposer:           %s
  Status:             %s
  Submit Time:        %s
  Deposit End Time:   %s
  Total Deposit:      %s
  Voting Start Time:  %s
  Voting End Time:    %s
  Description:        %s`,
		bp.ProposalID, bp.Title, bp.ProposalType, bp.Proposer.String(),
		bp.Status, bp.SubmitTime, bp.DepositEndTime,
		bp.TotalDeposit.MainUnitString(), bp.VotingStartTime, bp.VotingEndTime, bp.GetDescription(),
	)
}

func (bp BasicProposal) HumanString(converter sdk.CoinsConverter) string {
	return fmt.Sprintf(`Proposal %d:
  Title:              %s
  Type:               %s
  Status:             %s
  Submit Time:        %s
  Deposit End Time:   %s
  Total Deposit:      %s
  Voting Start Time:  %s
  Voting End Time:    %s
  Description:        %s`,
		bp.ProposalID, bp.Title, bp.ProposalType,
		bp.Status, bp.SubmitTime, bp.DepositEndTime,
		converter.ToMainUnit(bp.TotalDeposit), bp.VotingStartTime, bp.VotingEndTime, bp.GetDescription(),
	)
}

// Proposals is an array of proposal
type Proposals []Proposal

// nolint
func (p Proposals) String() string {
	if len(p) == 0 {
		return "[]"
	}
	out := "ID - (Status) [Type] [TotalDeposit] Title\n"
	for _, prop := range p {
		out += fmt.Sprintf("%d - (%s) [%s] [%s] %s\n",
			prop.GetProposalID(), prop.GetStatus(),
			prop.GetProposalType(), prop.GetTotalDeposit().MainUnitString(), prop.GetTitle())
	}
	return strings.TrimSpace(out)
}

func (p Proposals) HumanString(converter sdk.CoinsConverter) string {
	if len(p) == 0 {
		return "[]"
	}
	out := "ID - (Status) [Type] [TotalDeposit] Title\n"
	for _, prop := range p {
		out += fmt.Sprintf("%d - (%s) [%s] [%s] %s\n",
			prop.GetProposalID(), prop.GetStatus(),
			prop.GetProposalType(), converter.ToMainUnit(prop.GetTotalDeposit()), prop.GetTitle())
	}
	return strings.TrimSpace(out)
}

// Implements Proposal Interface
var _ Proposal = (*BasicProposal)(nil)

// nolint
func (bp BasicProposal) GetProposalID() uint64                      { return bp.ProposalID }
func (bp *BasicProposal) SetProposalID(proposalID uint64)           { bp.ProposalID = proposalID }
func (bp BasicProposal) GetTitle() string                           { return bp.Title }
func (bp *BasicProposal) SetTitle(title string)                     { bp.Title = title }
func (bp BasicProposal) GetDescription() string                     { return bp.Description }
func (bp *BasicProposal) SetDescription(description string)         { bp.Description = description }
func (bp BasicProposal) GetProposalType() ProposalKind              { return bp.ProposalType }
func (bp *BasicProposal) SetProposalType(proposalType ProposalKind) { bp.ProposalType = proposalType }
func (bp BasicProposal) GetStatus() ProposalStatus                  { return bp.Status }
func (bp *BasicProposal) SetStatus(status ProposalStatus)           { bp.Status = status }
func (bp BasicProposal) GetTallyResult() TallyResult                { return bp.TallyResult }
func (bp *BasicProposal) SetTallyResult(tallyResult TallyResult)    { bp.TallyResult = tallyResult }
func (bp BasicProposal) GetSubmitTime() time.Time                   { return bp.SubmitTime }
func (bp *BasicProposal) SetSubmitTime(submitTime time.Time)        { bp.SubmitTime = submitTime }
func (bp BasicProposal) GetDepositEndTime() time.Time               { return bp.DepositEndTime }
func (bp *BasicProposal) SetDepositEndTime(depositEndTime time.Time) {
	bp.DepositEndTime = depositEndTime
}
func (bp BasicProposal) GetTotalDeposit() sdk.Coins              { return bp.TotalDeposit }
func (bp *BasicProposal) SetTotalDeposit(totalDeposit sdk.Coins) { bp.TotalDeposit = totalDeposit }
func (bp BasicProposal) GetVotingStartTime() time.Time           { return bp.VotingStartTime }
func (bp *BasicProposal) SetVotingStartTime(votingStartTime time.Time) {
	bp.VotingStartTime = votingStartTime
}
func (bp BasicProposal) GetVotingEndTime() time.Time { return bp.VotingEndTime }
func (bp *BasicProposal) SetVotingEndTime(votingEndTime time.Time) {
	bp.VotingEndTime = votingEndTime
}
func (bp BasicProposal) GetProtocolDefinition() sdk.ProtocolDefinition {
	return sdk.ProtocolDefinition{}
}
func (bp *BasicProposal) SetProtocolDefinition(sdk.ProtocolDefinition) {}
func (bp BasicProposal) GetTaxUsage() TaxUsage                         { return TaxUsage{} }
func (bp *BasicProposal) SetTaxUsage(taxUsage TaxUsage)                {}
func (bp *BasicProposal) Validate(ctx sdk.Context, k Keeper, verify bool) sdk.Error {
	if !verify {
		return nil
	}
	pLevel := bp.ProposalType.GetProposalLevel()
	if num, ok := k.HasReachedTheMaxProposalNum(ctx, pLevel); ok {
		return ErrMoreThanMaxProposal(k.codespace, num, pLevel.string())
	}
	return nil
}
func (bp *BasicProposal) GetProposalLevel() ProposalLevel {
	return bp.ProposalType.GetProposalLevel()
}

func (bp *BasicProposal) GetProposer() sdk.AccAddress {
	return bp.Proposer
}
func (bp *BasicProposal) Execute(ctx sdk.Context, gk Keeper) sdk.Error {
	return sdk.MarshalResultErr(errors.New("BasicProposal can not execute 'Execute' method"))
}

//-----------------------------------------------------------
// ProposalQueue
type ProposalQueue []uint64

//-----------------------------------------------------------
// ProposalKind

// Type that represents Proposal Type as a byte
type ProposalKind byte

//nolint
const (
	ProposalTypeNil               ProposalKind = 0x00
	ProposalTypeParameter         ProposalKind = 0x01
	ProposalTypeSoftwareUpgrade   ProposalKind = 0x02
	ProposalTypeSystemHalt        ProposalKind = 0x03
	ProposalTypeCommunityTaxUsage ProposalKind = 0x04
	ProposalTypePlainText         ProposalKind = 0x05
	ProposalTypeTokenAddition     ProposalKind = 0x06
)

var pTypeMap = map[string]pTypeInfo{
	"PlainText":         createPlainTextInfo(),
	"Parameter":         createParameterInfo(),
	"SoftwareUpgrade":   createSoftwareUpgradeInfo(),
	"SystemHalt":        createSystemHaltInfo(),
	"CommunityTaxUsage": createCommunityTaxUsageInfo(),
	"TokenAddition":     createTokenAdditionInfo(),
}

// String to proposalType byte.  Returns ff if invalid.
func ProposalTypeFromString(str string) (ProposalKind, error) {
	kind, ok := pTypeMap[str]
	if !ok {
		return ProposalKind(0xff), errors.Errorf("'%s' is not a valid proposal type", str)
	}
	return kind.Type, nil
}

// is defined ProposalType?
func ValidProposalType(pt ProposalKind) bool {
	_, ok := pTypeMap[pt.String()]
	return ok
}

// Marshal needed for protobuf compatibility
func (pt ProposalKind) Marshal() ([]byte, error) {
	return []byte{byte(pt)}, nil
}

// Unmarshal needed for protobuf compatibility
func (pt *ProposalKind) Unmarshal(data []byte) error {
	*pt = ProposalKind(data[0])
	return nil
}

// Marshals to JSON using string
func (pt ProposalKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(pt.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (pt *ProposalKind) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := ProposalTypeFromString(s)
	if err != nil {
		return err
	}
	*pt = bz2
	return nil
}

// Turns VoteOption byte to String
func (pt ProposalKind) String() string {
	for k, v := range pTypeMap {
		if v.Type == pt {
			return k
		}
	}
	return ""
}

func (pt ProposalKind) NewProposal(content Content) (Proposal, sdk.Error) {
	typInfo, ok := pTypeMap[pt.String()]
	if !ok {
		return nil, ErrInvalidProposalType(DefaultCodespace, pt)
	}
	return typInfo.createProposal(content), nil
}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (pt ProposalKind) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(pt.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(pt))))
	}
}

func (pt ProposalKind) GetProposalLevel() ProposalLevel {
	return pTypeMap[pt.String()].Level
}

//-----------------------------------------------------------
// ProposalStatus

// Type that represents Proposal Status as a byte
type ProposalStatus byte

//nolint
const (
	StatusNil           ProposalStatus = 0x00
	StatusDepositPeriod ProposalStatus = 0x01
	StatusVotingPeriod  ProposalStatus = 0x02
	StatusPassed        ProposalStatus = 0x03
	StatusRejected      ProposalStatus = 0x04
)

var pStatusMap = map[string]ProposalStatus{
	"DepositPeriod": StatusDepositPeriod,
	"VotingPeriod":  StatusVotingPeriod,
	"Passed":        StatusPassed,
	"Rejected":      StatusRejected,
}

// ProposalStatusToString turns a string into a ProposalStatus
func ProposalStatusFromString(str string) (ProposalStatus, error) {
	status, ok := pStatusMap[str]
	if !ok {
		return ProposalStatus(0xff), errors.Errorf("'%s' is not a valid proposal status", str)
	}
	return status, nil
}

// is defined ProposalType?
func ValidProposalStatus(status ProposalStatus) bool {
	_, ok := pStatusMap[status.String()]
	return ok
}

// Marshal needed for protobuf compatibility
func (status ProposalStatus) Marshal() ([]byte, error) {
	return []byte{byte(status)}, nil
}

// Unmarshal needed for protobuf compatibility
func (status *ProposalStatus) Unmarshal(data []byte) error {
	*status = ProposalStatus(data[0])
	return nil
}

// Marshals to JSON using string
func (status ProposalStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(status.String())
}

// Unmarshals from JSON assuming Bech32 encoding
func (status *ProposalStatus) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil
	}

	bz2, err := ProposalStatusFromString(s)
	if err != nil {
		return err
	}
	*status = bz2
	return nil
}

// Turns VoteStatus byte to String
func (status ProposalStatus) String() string {
	for k, v := range pStatusMap {
		if v == status {
			return k
		}
	}
	return ""
}

// For Printf / Sprintf, returns bech32 when using %s
// nolint: errcheck
func (status ProposalStatus) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		s.Write([]byte(status.String()))
	default:
		// TODO: Do this conversion more directly
		s.Write([]byte(fmt.Sprintf("%v", byte(status))))
	}
}

//-----------------------------------------------------------
// Tally Results
type TallyResult struct {
	Yes               sdk.Dec `json:"yes"`
	Abstain           sdk.Dec `json:"abstain"`
	No                sdk.Dec `json:"no"`
	NoWithVeto        sdk.Dec `json:"no_with_veto"`
	SystemVotingPower sdk.Dec `json:"system_voting_power"`
}

// checks if two proposals are equal
func EmptyTallyResult() TallyResult {
	return TallyResult{
		Yes:               sdk.ZeroDec(),
		Abstain:           sdk.ZeroDec(),
		No:                sdk.ZeroDec(),
		NoWithVeto:        sdk.ZeroDec(),
		SystemVotingPower: sdk.ZeroDec(),
	}
}

// checks if two proposals are equal
func (tr TallyResult) Equals(resultB TallyResult) bool {
	return tr.Yes.Equal(resultB.Yes) &&
		tr.Abstain.Equal(resultB.Abstain) &&
		tr.No.Equal(resultB.No) &&
		tr.NoWithVeto.Equal(resultB.NoWithVeto) &&
		tr.SystemVotingPower.Equal(resultB.SystemVotingPower)
}

func (tr TallyResult) String() string {
	return fmt.Sprintf(`Tally Result:
  Yes:                %s
  Abstain:            %s
  No:                 %s
  NoWithVeto:         %s
  SystemVotingPower:  %s`, tr.Yes.String(), tr.Abstain.String(), tr.No.String(), tr.NoWithVeto.String(), tr.SystemVotingPower.String())
}
