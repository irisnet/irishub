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

	GetProtocolDefinition() sdk.ProtocolDefinition
	SetProtocolDefinition(sdk.ProtocolDefinition)

	GetTaxUsage() TaxUsage
	SetTaxUsage(TaxUsage)

	String() string
}

// checks if two proposals are equal
func ProposalEqual(proposalA Proposal, proposalB Proposal) bool {
	if proposalA.GetProposalID() == proposalB.GetProposalID() &&
		proposalA.GetTitle() == proposalB.GetTitle() &&
		proposalA.GetDescription() == proposalB.GetDescription() &&
		proposalA.GetProposalType() == proposalB.GetProposalType() &&
		proposalA.GetStatus() == proposalB.GetStatus() &&
		proposalA.GetTallyResult().Equals(proposalB.GetTallyResult()) &&
		proposalA.GetSubmitTime().Equal(proposalB.GetSubmitTime()) &&
		proposalA.GetDepositEndTime().Equal(proposalB.GetDepositEndTime()) &&
		proposalA.GetTotalDeposit().IsEqual(proposalB.GetTotalDeposit()) &&
		proposalA.GetVotingStartTime().Equal(proposalB.GetVotingStartTime()) &&
		proposalA.GetVotingEndTime().Equal(proposalB.GetVotingEndTime()) {
		return true
	}
	return false
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

	VotingStartTime time.Time `json:"voting_start_time"` //  Time of the block where MinDeposit was reached. -1 if MinDeposit is not reached
	VotingEndTime   time.Time `json:"voting_end_time"`   // Time that the VotingPeriod for this proposal will end and votes will be tallied
}

func (bp BasicProposal) String() string {
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
		bp.TotalDeposit.MainUnitString(), bp.VotingStartTime, bp.VotingEndTime, bp.GetDescription(),
	)
}

// Proposals is an array of proposal
type Proposals []Proposal

// nolint
func (p Proposals) String() string {
	out := "ID - (Status) [Type] [TotalDeposit] Title\n"
	for _, prop := range p {
		out += fmt.Sprintf("%d - (%s) [%s] [%s] %s\n",
			prop.GetProposalID(), prop.GetStatus(),
			prop.GetProposalType(), prop.GetTotalDeposit().MainUnitString(), prop.GetTitle())
	}
	return strings.TrimSpace(out)
}

// Implements Proposal Interface
var _ Proposal = (*BasicProposal)(nil)

// nolint
func (tp BasicProposal) GetProposalID() uint64                      { return tp.ProposalID }
func (tp *BasicProposal) SetProposalID(proposalID uint64)           { tp.ProposalID = proposalID }
func (tp BasicProposal) GetTitle() string                           { return tp.Title }
func (tp *BasicProposal) SetTitle(title string)                     { tp.Title = title }
func (tp BasicProposal) GetDescription() string                     { return tp.Description }
func (tp *BasicProposal) SetDescription(description string)         { tp.Description = description }
func (tp BasicProposal) GetProposalType() ProposalKind              { return tp.ProposalType }
func (tp *BasicProposal) SetProposalType(proposalType ProposalKind) { tp.ProposalType = proposalType }
func (tp BasicProposal) GetStatus() ProposalStatus                  { return tp.Status }
func (tp *BasicProposal) SetStatus(status ProposalStatus)           { tp.Status = status }
func (tp BasicProposal) GetTallyResult() TallyResult                { return tp.TallyResult }
func (tp *BasicProposal) SetTallyResult(tallyResult TallyResult)    { tp.TallyResult = tallyResult }
func (tp BasicProposal) GetSubmitTime() time.Time                   { return tp.SubmitTime }
func (tp *BasicProposal) SetSubmitTime(submitTime time.Time)        { tp.SubmitTime = submitTime }
func (tp BasicProposal) GetDepositEndTime() time.Time               { return tp.DepositEndTime }
func (tp *BasicProposal) SetDepositEndTime(depositEndTime time.Time) {
	tp.DepositEndTime = depositEndTime
}
func (tp BasicProposal) GetTotalDeposit() sdk.Coins              { return tp.TotalDeposit }
func (tp *BasicProposal) SetTotalDeposit(totalDeposit sdk.Coins) { tp.TotalDeposit = totalDeposit }
func (tp BasicProposal) GetVotingStartTime() time.Time           { return tp.VotingStartTime }
func (tp *BasicProposal) SetVotingStartTime(votingStartTime time.Time) {
	tp.VotingStartTime = votingStartTime
}
func (tp BasicProposal) GetVotingEndTime() time.Time { return tp.VotingEndTime }
func (tp *BasicProposal) SetVotingEndTime(votingEndTime time.Time) {
	tp.VotingEndTime = votingEndTime
}
func (tp BasicProposal) GetProtocolDefinition() sdk.ProtocolDefinition {
	return sdk.ProtocolDefinition{}
}
func (tp *BasicProposal) SetProtocolDefinition(sdk.ProtocolDefinition) {}
func (tp BasicProposal) GetTaxUsage() TaxUsage                         { return TaxUsage{} }
func (tp *BasicProposal) SetTaxUsage(taxUsage TaxUsage)                {}

//-----------------------------------------------------------
// ProposalQueue
type ProposalQueue []uint64

//-----------------------------------------------------------
// ProposalKind

// Type that represents Proposal Type as a byte
type ProposalKind byte

//nolint
const (
	ProposalTypeNil             ProposalKind = 0x00
	ProposalTypeParameterChange ProposalKind = 0x01
	ProposalTypeSoftwareUpgrade ProposalKind = 0x02
	ProposalTypeSystemHalt      ProposalKind = 0x03
	ProposalTypeTxTaxUsage      ProposalKind = 0x04
)

// String to proposalType byte.  Returns ff if invalid.
func ProposalTypeFromString(str string) (ProposalKind, error) {
	switch str {
	case "ParameterChange":
		return ProposalTypeParameterChange, nil
	case "SoftwareUpgrade":
		return ProposalTypeSoftwareUpgrade, nil
	case "SystemHalt":
		return ProposalTypeSystemHalt, nil
	case "TxTaxUsage":
		return ProposalTypeTxTaxUsage, nil
	default:
		return ProposalKind(0xff), errors.Errorf("'%s' is not a valid proposal type", str)
	}
}

// is defined ProposalType?
func ValidProposalType(pt ProposalKind) bool {
	if pt == ProposalTypeParameterChange ||
		pt == ProposalTypeSoftwareUpgrade ||
		pt == ProposalTypeSystemHalt ||
		pt == ProposalTypeTxTaxUsage {
		return true
	}
	return false
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
	switch pt {
	case ProposalTypeParameterChange:
		return "ParameterChange"
	case ProposalTypeSoftwareUpgrade:
		return "SoftwareUpgrade"
	case ProposalTypeSystemHalt:
		return "SystemHalt"
	case ProposalTypeTxTaxUsage:
		return "TxTaxUsage"
	default:
		return ""
	}
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

// ProposalStatusToString turns a string into a ProposalStatus
func ProposalStatusFromString(str string) (ProposalStatus, error) {
	switch str {
	case "DepositPeriod":
		return StatusDepositPeriod, nil
	case "VotingPeriod":
		return StatusVotingPeriod, nil
	case "Passed":
		return StatusPassed, nil
	case "Rejected":
		return StatusRejected, nil
	case "":
		return StatusNil, nil
	default:
		return ProposalStatus(0xff), errors.Errorf("'%s' is not a valid proposal status", str)
	}
}

// is defined ProposalType?
func ValidProposalStatus(status ProposalStatus) bool {
	if status == StatusDepositPeriod ||
		status == StatusVotingPeriod ||
		status == StatusPassed ||
		status == StatusRejected {
		return true
	}
	return false
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
	switch status {
	case StatusDepositPeriod:
		return "DepositPeriod"
	case StatusVotingPeriod:
		return "VotingPeriod"
	case StatusPassed:
		return "Passed"
	case StatusRejected:
		return "Rejected"
	default:
		return ""
	}
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
	Yes        sdk.Dec `json:"yes"`
	Abstain    sdk.Dec `json:"abstain"`
	No         sdk.Dec `json:"no"`
	NoWithVeto sdk.Dec `json:"no_with_veto"`
}

// checks if two proposals are equal
func EmptyTallyResult() TallyResult {
	return TallyResult{
		Yes:        sdk.ZeroDec(),
		Abstain:    sdk.ZeroDec(),
		No:         sdk.ZeroDec(),
		NoWithVeto: sdk.ZeroDec(),
	}
}

// checks if two proposals are equal
func (resultA TallyResult) Equals(resultB TallyResult) bool {
	return resultA.Yes.Equal(resultB.Yes) &&
		resultA.Abstain.Equal(resultB.Abstain) &&
		resultA.No.Equal(resultB.No) &&
		resultA.NoWithVeto.Equal(resultB.NoWithVeto)
}

func (tr TallyResult) String() string {
	return fmt.Sprintf(`Tally Result:
  Yes:        %s
  Abstain:    %s
  No:         %s
  NoWithVeto: %s`, tr.Yes, tr.Abstain, tr.No, tr.NoWithVeto)
}
