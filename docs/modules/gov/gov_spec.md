# Governance Module

本模块的主要作用是增加区块链的链上治理能力，通过该某块可以实现链上配置数据的动态修改，软件升级等功能，主要分为以下几个步骤：

* 提交提议：当某个需求需要大多数人达成一致时，可以由提议者提交提议，将需求以交易的形式广播到区块链上。当然这个提议不是可以随便提交的，除了需要扣除部分手续费之外，还需要抵押你的一部分代币，为了防止提交恶意请求。但是这部分抵押的代币在提议通过之后，会等额的返还到你的账户。如果抵押的金额没有达到系统预设的最小抵押金额，这个提议将进入抵押代币阶段；如果大于了最小抵押金额，该提议将自动进入投票阶段。该阶段的相关数据模型如下：
    
```golang
//提议交易数据结构
type MsgSubmitProposal struct {
	Title          string         //  Title of the proposal
	Description    string         //  Description of the proposal
	ProposalType   ProposalKind   //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       sdk.AccAddress //  Address of the proposer
	InitialDeposit sdk.Coins      //  Initial deposit paid by sender. Must be strictly positive.
	Params         Params
}

```
  
* 抵押代币：在提交提议的时候，需要抵押部分代币，这个最小的抵押金额可能会很大，造成单方的代币不够，所以，这个时候你可以邀请别人来帮你抵押一部分代币，当然在提议通过之后，双方抵押的代币也都会等额的返还。
抵押过程中，如果当前的抵押金额大于最小抵押金额，该提议进入投票阶段。当提议进入投票阶段之后，这个提议就不能再被抵押金额。在抵押阶段也有时间限制。当到达规定时间，抵押的金额还没能够大于等于最小抵押金额，系统会丢弃该提议，并且会被系统认为是垃圾提议，不会退还任何抵押者的代币。该阶段的相关数据模型如下：

```golang
//抵押交易数据结构
type MsgDeposit struct {
	ProposalID int64          `json:"proposalID"` // ID of the proposal
	Depositer  sdk.AccAddress `json:"depositer"`  // Address of the depositer
	Amount     sdk.Coins      `json:"amount"`     // Coins to add to the proposal's deposit
}
//抵押期的系统预设参数
type DepositProcedure struct {
	MinDeposit       sdk.Coins `json:"min_deposit"`        //  Minimum deposit for a proposal to enter voting period.
	MaxDepositPeriod int64     `json:"max_deposit_period"` //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months
}

```


* 投票：当提议进入投票阶段之后，验证人需要在规定时间(系统预设)内，对提议投票。投票有以下几种可选项：Yes(同意)，Abstain(弃权)，No(反对)，NoWithVeto(强烈反对)。当然你也可以不参与投票，但是你和你的委托人要承担失去部分绑定金额(比例由系统预设)的风险。当达到规定时间，系统开始自动统计投票信息，根据投票结果，会出现以下行为：
投票通过：退回所有抵押者的全部金额，如果存在未投票的验证人，会受到相应的惩罚(具体由slash模块处理)。投票未通过：不退还抵押者的代币，同时对未参与投票的验证人进行slash机制。该阶段的相关数据模型如下：

```golang

//投票交易数据结构
type MsgVote struct {
	ProposalID int64          //  proposalID of the proposal
	Voter      sdk.AccAddress //  address of the voter
	Option     VoteOption     //  option from OptionSet chosen by the voter
}
//投票期的系统预设参数
// Procedure around Voting in governance
type VotingProcedure struct {
	VotingPeriod int64 `json:"voting_period"` //  Length of the voting period.
}

// Procedure around Tallying votes in governance
type TallyingProcedure struct {
	Threshold         sdk.Rat `json:"threshold"`          //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	Veto              sdk.Rat `json:"veto"`               //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	GovernancePenalty sdk.Rat `json:"governance_penalty"` //  Penalty if validator does not vote
}

```

* 提议执行：该阶段是新增的功能，是投票通过的后续动作，根据用户提出的提议的不同类型，会触发不同的后续操作，目前有以下几种提议：
    * 文本提议：一种只读的提议，提议的执行不会对系统产生任何实际的影响
    * 参数变更提议：提议执行的时候，根据提议者参入的修改内容，可以修改系统里预设的一些配置参数。例如提议的最小抵押金额，由原来的1000iris，修改为2000iris
    * 软件升级提议：当提议被通过后，软件升级提议会变更当前软件运行的版本，触发软件升级模块执行升级动作
  提议被设计为借口类型，具体实现数据模型如下：
  
```golang

// 提议接口定义
type Proposal interface {
	GetProposalID() int64
	SetProposalID(int64)

	GetTitle() string
	SetTitle(string)

	GetDescription() string
	SetDescription(string)

	GetProposalType() ProposalKind
	SetProposalType(ProposalKind)

	GetStatus() ProposalStatus
	SetStatus(ProposalStatus)

	GetSubmitBlock() int64
	SetSubmitBlock(int64)

	GetTotalDeposit() sdk.Coins
	SetTotalDeposit(sdk.Coins)

	GetVotingStartBlock() int64
	SetVotingStartBlock(int64)

	Execute(ctx sdk.Context, k Keeper) error
}

//文本提议
type TextProposal struct {
	ProposalID   int64        `json:"proposal_id"`   //  ID of the proposal
	Title        string       `json:"title"`         //  Title of the proposal
	Description  string       `json:"description"`   //  Description of the proposal
	ProposalType ProposalKind `json:"proposal_type"` //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}

	Status ProposalStatus `json:"proposal_status"` //  Status of the Proposal {Pending, Active, Passed, Rejected}

	SubmitBlock  int64     `json:"submit_block"`  //  Height of the block where TxGovSubmitProposal was included
	TotalDeposit sdk.Coins `json:"total_deposit"` //  Current deposit on this proposal. Initial value is set at InitialDeposit

	VotingStartBlock int64 `json:"voting_start_block"` //  Height of the block where MinDeposit was reached. -1 if MinDeposit is not reached
}

//参数变更提议
type ParameterProposal struct {
	TextProposal
	Params Params `json:"params"`
}

type Op string

const (
	Add    Op = "add"
	Update Op = "update"
)

type Param struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Op    Op     `json:"op"`
}

type Params []Param

//软件升级提议
type SoftwareUpgradeProposal struct {
	TextProposal
}

```