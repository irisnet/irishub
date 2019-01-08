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

`SoftwareUpgradeProposal` and `SystemHaltProposal` can only be submitted by the profiler. `TxTaxUsage` can only be submitted by `trustee`.

Different levels correspond to different parameters：

| GovParams | Critical | Important | Normal |Range|
| ------ | ------ | ------ | ------|------| 
| govDepositProcedure/MinDeposit | 4000 iris | 2000 iris | 1000 iris |[10iris,10000iris]|
| govDepositProcedure/MaxDepositPeriod | 24 hours | 24 hours | 24 hours |[20s,3d]|
| govVotingProcedure/VotingPeriod | 72 hours | 60 hours | 48 hours |[20s,3d]|
| govVotingProcedure/MaxProposal | 1 | 2 | 1 |Critial==1, other(1,)|
| govTallyingProcedure/Participation | 6/7 | 5/6 | 3/4 |(0,1)|
| govTallyingProcedure/Threshold | 5/6 | 4/5 | 2/3 |(0,1)|
| govTallyingProcedure/Veto | 1/3 | 1/3 | 1/3 |(0,1)|
| govTallyingProcedure/Penalty | 0.0009 | 0.0007 | 0.0005 |(0,1)|


* `MinDeposit`  The minimum of  deposit
* `MaxDepositPeriod`  Window period for deposit
* `VotingPeriod` Window period for voting
* `MaxProposal` The maximum number of proposal that can exist at the same time
* `Penalty`   The proportion of the slash
* `Veto` 
* `Threshold` 
* `Participation` 

### Deposit Procedure

The submitted proposal has the deposit and when the deposit exceeds `MinDeposit`, it can enter the voting procedure. If the proposal exceeds `MaxDepositPeriod` and has not yet entered the voting procedure, the proposal will be deleted and the full deposit will be refunded. It is not possible to deposit the proposal which has been in  the voting procedure .

### Voting Procedure
Only the validator can vote once, and the vote cannot be repeated. The voting options are: `Yes` , `Abstain` , `No` , `NoWithVeto` .

### Tallying Procedure

There are three tallying results: `PASS`，`REJECT`，`REJECTVETO`。

 Under the premise that the ratio of all voters' voting power to the total voting power in system is more than `participation`, if the ratio of `NoWithVeto` voting power  to all voters' voting power over `veto`, the result is `REJECTVETO`. Then if the ratio of `Yes` voting power  to all voter's voting power  over `threshold`, the result is `REJECT`. Otherwise, the result is `REJECT`.

### Burning Mechanism

If the proposal is passed or not passed, 20% Deposit will be burned. As the cost of the governance, the remaining `Deposit` will be returned in proportion. But if the result of proposal is `REJECTVETO`, destroy all `Deposit`.

### Slashing Mechanism

If the proposal enters the voting procedure, the account is a validator and then the proposal enters the tallying procedure, he is still a validator, but if he does not vote, he will be slashed according to the proportion of `Penalty`.

## Usage Scenario

### Usage scenario of parameter change

Change the parameters through the command lines

```
# Query parameters can be changed by the modules'name in gov 
iriscli gov query-params --module=mint

# Results
iriscli gov query-params --module=mint

# Query parameters can be modified by "key”
iriscli gov query-params --module=mint --key=mint/Inflation

# Results
iriscli gov query-params --module=mint

# Send proposals, return changed parameters
iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange" --deposit=8iris  --param mint/Inflation=0.0000000000 --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# Deposit for a proposal
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# Vote for a proposal 
echo 1234567890 | iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# Query the state of a proposal
iriscli gov query-proposal --proposal-id=1 

```

### Proposals on community funds usage
There are three usages, `Burn`, `Distribute` and `Grant`. `Burn` means burning tokens from community funds. `Distribute` and `Grant` will transfer tokens to the destination trustee's account from community funds and then trustee will distribute or grant these tokens to others.
```shell
# Submit Burn usage proposal
iriscli gov submit-proposal --title="burn tokens 5%" --description="test" --type="TxTaxUsage" --usage="Burn" --deposit="10iris"  --percent=0.05 --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# Submit Distribute usage proposal
iriscli gov submit-proposal --title="distribute tokens 5%" --description="test" --type="TxTaxUsage" --usage="Distribute" --deposit="10iris"  --percent=0.05 --dest-address=[destnation-address] --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# Submit Grant usage proposal
iriscli gov submit-proposal --title="grant tokens 5%" --description="test" --type="TxTaxUsage" --usage="Grant" --deposit="10iris"  --percent=0.05 --dest-address=[destnation-address] --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# Deposit for a proposal
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# Vote for a proposal
iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --commit

# Query the state of a proposal
iriscli gov query-proposal --proposal-id=1
```

### Proposals on system halting

Sending this proposal which can terminate the system, the node will be closed after the proposed systemHaltHeight (= proposal height + systemHaltPeriod), and then start to enter the query-only mode.

```
# submit the SystemHaltProposal
iriscli gov submit-proposal  --title=test_title --description=test_description --type=SystemHalt --deposit=10iris --fee=0.005iris --from=x1 --chain-id=gov-test --commit

# Deposit for a proposal
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Vote for a proposal
iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Query the state of a proposal
iriscli gov query-proposal --proposal-id=1 
```

### Proposals on software upgrade

Detail in [Upgrade](upgrade.md)

