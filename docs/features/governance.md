# Gov User Guide

## Basic Function Description

1. On-chain governance proposals on parameter change
2. On-chain governance proposals on software upgrade 
3. On-chain governance proposals on software halt
4. On-chain governance proposals on tax usage

## Interactive process

### Proposal Level

Specific Proposal for different levels：
- Critical：`SoftwareUpgrade`, `SystemHalt`
- Important：`ParameterChange`
- Normal：`TxTaxUsage`

`SoftwareUpgrade Proposal` and `SystemHalt Proposal` can only be submitted by the profiler.

Different levels correspond to different parameters：

| GovParams | Critical | Important | Normal |Range|
| ------ | ------ | ------ | ------|------| 
| govDepositProcedure/MinDeposit | 4000 iris | 2000 iris | 1000 iris |[10iris,10000iris]|
| govDepositProcedure/MaxDepositPeriod | 24 hours | 24 hours | 24 hours |[20s,3d]|
| govVotingProcedure/VotingPeriod | 72 hours | 60 hours | 48 hours |[20s,3d]|
| govVotingProcedure/MaxProposal | 1 | 5 | 2 |Critical==1, other(1,)|
| govTallyingProcedure/Participation | 7/8 | 5/6 | 3/4 |(0,1)|
| govTallyingProcedure/Threshold | 6/7 | 4/5 | 2/3 |(0,1)|
| govTallyingProcedure/Veto | 1/3 | 1/3 | 1/3 |(0,1)|
| govTallyingProcedure/Penalty | 0.0009 | 0.0007 | 0.0005 |(0,1)|


* `MinDeposit`  The minimum of  deposit
* `MaxDepositPeriod`  Window period for deposit
* `VotingPeriod` Window period for voting
* `MaxProposal` The maximum number of proposal that can exist at the same time
* `Penalty`   The proportion of the slash
* `Veto`  the power of Veto / all voted power
* `Threshold` the power of Yes / all voted power
* `Participation` all voted power / total voting power

### Deposit Procedure

The proposer at least deposit more the 30% amount of `MinDeposit` to submit a proposal, when the total deposit amount exceeds `MinDeposit`, the proposal enter the voting procedure. If the time exceeds `MaxDepositPeriod` and the total deposit has not yet exceeded `MinDeposit`, the proposal will be deleted and the full deposit won't be refunded. It is not allowed to deposit a proposal which is in voting procedure.

### Voting Procedure
Only the validator can vote , and they can't vote twice for one proposal. The voting options are `Yes` , `Abstain` , `No` , `NoWithVeto` .

### Tallying Procedure

There are three tallying results: `PASS`，`REJECT`，`REJECTVETO`。

On the premise that the `voting_power of all voters` / `total voting_power of the system` exceeds `participation`,if the ratio of `NoWithVeto` voting power to all voters' voting power over `veto`, the result is `REJECTVETO`. Then if the ratio of `Yes` voting power to all voter's voting power over `threshold`, the result is `PASS`. Otherwise, the result is `REJECT`. 
 

### Burning Mechanism

Whether the proposal is passed or not passed, 20% `Deposit` will be burned for the cost of governance. The remaining `Deposit` will be returned. But if the result of proposal is `REJECTVETO`,  all `Deposit` will be burned.

### Slashing Mechanism

The validator should be slashed according to the proportion of `Penalty` if he fails to vote for a proposal.

If the validator quit from validator set in voting procedure, then he wont be slashed.

## Usage Scenario

### Usage scenario of parameter change

Change the parameters through the command lines

```
# Query parameters can be changed by the modules'name in gov 
iriscli gov query-params --module=mint

# Result
mint/Inflation=0.0400000000

# Query parameters can be modified by "key”
iriscli gov query-params --module=mint --key=mint/Inflation

# Results
mint/Inflation=0.0400000000

# Send proposal for parameters change
iriscli gov submit-proposal --title=<title> --description=<description> --type=ParameterChange --deposit=8iris  --param="mint/Inflation=0.0000000000" --from=<key_name> --chain-id=<chain-id> --fee=0.3iris --commit

# Deposit for a proposal
iriscli gov deposit --proposal-id=<proposal-id> --deposit=1000iris --from=<key_name> --chain-id=<chain-id> --fee=0.3iris --commit

# Vote for a proposal 
iriscli gov vote --proposal-id=<proposal-id> --option=Yes --from=<key_name> --chain-id=<chain-id> --fee=0.3iris --commit

# Query the state of a proposal
iriscli gov query-proposal --proposal-id=<proposal-id>
```

### Proposals on community funds usage
There are three usages, `Burn`, `Distribute` and `Grant`. `Burn` means burning tokens from community funds. `Distribute` and `Grant` will transfer tokens to the destination trustee's account from community funds.

```shell
# Submit Burn usage proposal
iriscli gov submit-proposal --title="burn tokens 5%" --description=<description> --type="TxTaxUsage" --usage="Burn" --deposit="10iris"  --percent=0.05 --from=<key_name> --chain-id=<chain-id> --fee=0.3iris --commit

# Submit Distribute usage proposal
iriscli gov submit-proposal --title="distribute tokens 5%" --description="test" --type="TxTaxUsage" --usage="Distribute" --deposit="10iris"  --percent=0.05 --dest-address=<dest-address (only trustees)> --from=<key_name> --chain-id=<chain-id> --fee=0.3iris --commit

# Submit Grant usage proposal
iriscli gov submit-proposal --title="grant tokens 5%" --description="test" --type="TxTaxUsage" --usage="Grant" --deposit="10iris"  --percent=0.05 --dest-address=<dest-address (only trustees)> --from=<key_name> --chain-id=<chain-id> --fee=0.3iris --commit
```

### Proposals on system halting

Sending this proposal which can terminate the system, the node will be closed after systemHaltHeight (= proposal height + systemHaltPeriod), and only `query-only` mode is available after re-starting the node.

```
# submit the SystemHaltProposal
iriscli gov submit-proposal --title=<title> --description=<description> --type=SystemHalt --deposit=10iris --fee=0.3iris --from=<key_name> --chain-id=<chain-id> --commit
```

### Proposals on software upgrade

Detail in [Upgrade](upgrade.md)

