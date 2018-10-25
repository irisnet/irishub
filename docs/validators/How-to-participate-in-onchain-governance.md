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

* Community Members: They could come up with interesting for improvement of network



## Governance Process
![workflow](../pics/flow.jpg)

The governance process is divided in a few steps that are outlined below:

* Proposal submission: Proposal is submitted to the blockchain with a deposit.
* Vote: Once deposit reaches a certain value (MinDeposit), proposal is confirmed and vote opens. Bonded Atom holders can then send TxGovVote transactions to vote on the proposal.
* If the proposal is passed, then there sill be some changes. 


For `System Parameter Change` proposals, the following parameters are eligible for change:

* Mininimum Depost: `100IRIS`
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

To add deposit to some proposal, you could execute this command to add `10IRIS` to the proposal's deposit:

```
iriscli gov deposit --proposalID=1 --depositer=<account> --deposit=1000000000000000000iris   --from=<name>  --chain-id=fuxi-3001  --fee=400000000000000iris --gas=20000  --node=http://localhost:36657 
```
## How to query proposals?

You could also use IRISplorer to query all proposal. 

## How to vote a proposal?

In the current version of governance module, you have the following choices for each proposal:
* Yes
* No
* NoWithVeto
* Abstien

You could put one of the choices in the `--option` field. 

To vote for a proposal, you need to get the correct `<proposal_id>`.You could execute the following command to vote on proposal with ID = 1:
```
iriscli  vote --from=jerry --voter=<account> --proposalID=1 --option=Yes --chain-id=fuxi-3001
   --fee=2000000000000000iris --gas=20000  --node=http://localhost:36657
```

## How to get more information about a proposal?

You could use the following command to get the first proposal:  
```
iriscli gov query-proposal --proposalID=1 --chain-id=fuxi-3001 --node=http://localhost:26657

```

