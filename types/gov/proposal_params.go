package gov

const (
	Insert string = "insert"
	Update string = "update"
)

type Param struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Op    string `json:"op"`
}

// Implements Proposal Interface
var _ Proposal = (*ParameterProposal)(nil)

type ParameterProposal struct {
	TextProposal
	Param Param `json:"params"`
}
