# Governance In IRISHub

## What is Blockchain Governance?

Generally, governance is about decisions that ultimately affect stakeholders of blockchain network. It’s about the processes 
that community members in governance use to make decisions. 
It’s also about how different groups coordinate around decisions and decision-making processes. 
It includes the establishment, maintenance, and revocation of the legitimacy of decisions, decision making processes, norms,
and other mechanisms for coordination.

### Stakeholder of Governance

* Core Developers:a blockchain are software developers who work on the software that implement that protocol. 
Developers have processes that are supposed to assure the quality of the software they release, 

* Validators/Delegators: Validators and delegators are responsible for securing the network. They need to actively involved in governance
proposals or their credit will be damaged. 

* Community Members: They could come up with interesting for improvement of network. block explorers and other low level service providers, exchanges, speculators, application developers, users, journalists and passive observers



## Governance Process
![workflow](../pics/flow.jpg)

The governance process is divided in a few steps that are outlined below:

* Proposal submission: Proposal is submitted to the blockchain with a deposit.
* Vote: Once deposit reaches a certain value (MinDeposit), proposal is confirmed and vote opens. Bonded Atom holders can then send TxGovVote transactions to vote on the proposal.
* If the proposal is passed, then there sill be some changes. 


For `System Parameter Change` proposals, the following parameters are eligible for change:

* Mininimum Depost: `10IRIS`
* Deposit Period: 1440 blocks
* Penalty for non-voting validtors: 1%
* Pass Threshold: 50%
* Voting Period: 20000 blocks


### Types of Governance Proposal

* Text
* System Parameter Change
* Protocol Upgrade 


## How to submit a proposal?

Anyone could submit a governance proposal, but you need to make the deposit for this proposal more than the minimium requirement.

The following command is for submitting a `Text` proposal:

```
iriscli gov submit-proposal --title="Text" --description="name of the proposal" --type="Text" --deposit="100iris" --proposer=<account>  --from=<name>  --chain-id=fuxi-4000 --fee=0.05iris --gas=20000 --node=http://localhost:26657
```

The `<account>` for `proposer` field should start with `faa` which corresponds to `<name>`.


## How to add deposit to a proposal?

To add deposit to some proposal, you could execute this command to add `100IRIS` to the proposal's deposit:

```
iriscli gov deposit --proposalID=1 --depositer=<account> --deposit=1000000000000000000iris   --from=<name>  --chain-id=fuxi-4000  --fee=0.05iris --gas=20000  --node=http://localhost:36657 
```
## How to query proposals?

Run the following command will return all the existing proposals:
```$xslt
iriscli gov query-proposals
```
Example output:
```$xslt
  1 - 78547
  2 - 96866
  3 - 46727
  4 - 92454
  5 - 57682
```

You could also use IRISplorer to see all proposal. 

## How to vote a proposal?

In the current version of governance module, you have the following choices for each proposal:
* Yes
* No
* NoWithVeto
* Abstien

You could put one of the choices in the `--option` field. 

To vote for a proposal, you need to get the correct `<proposal_id>`.You could execute the following command to vote on proposal with ID = 1:


## How to get more information about a proposal?

You could use the following command to get the first proposal:  
```
iriscli gov query-proposal --proposal-id=6 --chain-id=fuxi-3001 --node=http://localhost:26657
```

Example output is the following:
```$xslt
{
  "proposal_id": "1",
  "title": "text_proposal",
  "description": "test_description",
  "proposal_type": "Text",
  "proposal_status": "VotingPeriod",
  "tally_result": {
    "yes": "0",
    "abstain": "0",
    "no": "0",
    "no_with_veto": "0"
  },
  "submit_block": "200981",
  "total_deposit": [
    "20iris"
  ],
  "voting_start_block": "200981",
  "param": {
    "key": "",
    "value": "",
    "op": ""
  }
}
```
## Proposal Examples

### Text Proposal

A text proposal is just a notice. You could send one with the following command:
```$xslt
iriscli gov submit-proposal --title=text_proposal --description=test_description --type=Text --deposit=20iris --fee=0.1iris  --from=bft --chain-id=fuxi-4000 
```
You could vote on it with the following command:
```$xslt
iriscli gov vote --proposal-id=6 --option=Yes --from=bft --chain-id=fuxi-4000 --fee=0.05iris --gas=20000 
```


### System Parameter Change Proposal

First, you could query the flexible parameters with the following command:
```
iriscli gov query-params --trust-node --module=gov 

```

The output is the shown below:
```$xslt
[
 "Gov/gov/DepositProcedure",
 "Gov/gov/TallyingProcedure",
 "Gov/gov/VotingProcedure"
]
```
It corresponds to the following fields in genesis file:

```json
gov: {
starting_proposalID: "1",
deposit_period: {
min_deposit: [
{
denom: "iris-atto",
amount: "10000000000000000000"
}
],
max_deposit_period: "50"
},
voting_period: {
voting_period: "50"
},
tallying_procedure: {
threshold: "1/2",
veto: "1/3",
governance_penalty: "1/100"
}
},
```


The default value of flexible parameters are shown above.

You could submit a parameter change proposal,
```$xslt
iriscli gov submit-proposal --title="update VotingProcedure" --description="test" --type="ParameterChange" --deposit=10iris --param='{"key":"Gov/gov/VotingProcedure","value":"{\"voting_period\": 250}","op":"update"}' --from=bft --chain-id=fuxi-4000  --fee=0.05iris --gas=20000
```
This proposal will change the voting period from default to 250.

### System Upgrade Proposal

You could read more about it in the [doc](../modules/upgrade/README.md)