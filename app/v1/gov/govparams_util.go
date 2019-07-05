package gov

//-----------------------------------------------------------
// ProposalLevel

// Type that represents Proposal Level as a byte
type ProposalLevel byte

//nolint
const (
	ProposalLevelNil       ProposalLevel = 0x00
	ProposalLevelCritical  ProposalLevel = 0x01
	ProposalLevelImportant ProposalLevel = 0x02
	ProposalLevelNormal    ProposalLevel = 0x03
)

func (p ProposalLevel) string() string {
	switch p {
	case ProposalLevelCritical:
		return "critical"
	case ProposalLevelImportant:
		return "important"
	case ProposalLevelNormal:
		return "normal"
	default:
		return " "
	}
}
