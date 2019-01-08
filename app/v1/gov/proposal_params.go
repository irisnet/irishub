package gov

const (
	Insert string = "insert"
	Update string = "update"
)

type Param struct {
	Subspace string `json:"subspace"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type Params []Param

// Implements Proposal Interface
var _ Proposal = (*ParameterProposal)(nil)

type ParameterProposal struct {
	TextProposal
	Params Params `json:"params"`
}
