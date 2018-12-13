# Gov User Guide

## Basic Function Description

1. On-chain governance proposals on text
2. On-chain governance proposals on parameter change
3. On-chain governance proposals on software upgrade 

## Interactive process

### governance process

1. Any users can deposit some tokens to initiate a proposal. Once deposit reaches a certain value `min_deposit`, enter voting period, otherwise it will remain in the deposit period. Others can deposit the proposals on the deposit period. Once the sum of the deposit reaches `min_deposit`, enter voting period. However, if the block-time exceeds `max_deposit_period` in the deposit period, the proposal will be closed.
2. The proposals which enter voting period only can be voted by validators and delegators. The vote of a delegator who hasn't vote will be the same as his validator's vote, and the vote of a delegator who has voted will be remained. The votes wil be tallyed when reach `voting_period'.
3. Our tally have a limit on participation, Other details about voting for proposals:
[CosmosSDK-Gov-spec](https://github.com/cosmos/cosmos-sdk/blob/v0.26.0/docs/spec/governance/overview.md)

## Usage Scenario

### Usage scenario of parameter change

Scenario 1：Change the parameters through the command lines

```
# Query parameters can be changed by the modules'name in gov 
iriscli gov query-params --module=gov --trust-node

# Results
[
"Gov/govDepositProcedure",
"Gov/govTallyingProcedure",
"Gov/govVotingProcedure"
]

# Query parameters can be modified by "key”
iriscli gov query-params --key=Gov/govDepositProcedure --trust-node

# Results
{"key":"Gov/govDepositProcedure","value":"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":172800000000000}","op":""}

# Send proposals, return changed parameters
iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange" --deposit="10iris"  --param='{"key":"Gov/govDepositProcedure","value":"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"20000000000000000000\"}],\"max_deposit_period\":172800000000000}","op":"update"}}' --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Deposit for a proposal
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Vote for a proposal
echo 1234567890 | iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Query the state of a proposal
iriscli gov query-proposal --proposal-id=1 --trust-node

```

Scenario 2: Change the parameters by the files

```
# Export profiles
iriscli gov pull-params --path=iris --trust-node

# Query profiles' info
cat iris/config/params.json                                              {
"gov": {
"Gov/govDepositProcedure": {
"min_deposit": [
{
"denom": "iris-atto",
"amount": "10000000000000000000"
}
],
"max_deposit_period": "172800000000000"
},
"Gov/govVotingProcedure": {
"voting_period": "10000000000"
},
"Gov/govTallyingProcedure": {
"threshold": "0.5000000000",
"veto": "0.3340000000",
"participation": "0.6670000000"
}
}

# Modify profiles (TallyingProcedure的governance_penalty)
vi iris/config/params.json                                               {
"gov": {
"Gov/govDepositProcedure": {
"min_deposit": [
{
"denom": "iris-atto",
"amount": "10000000000000000000"
}
],
"max_deposit_period": "172800000000000"
},
"Gov/govVotingProcedure": {
"voting_period": "10000000000"
},
"Gov/govTallyingProcedure": {
"threshold": "0.5000000000",
"veto": "0.3340000000",
"participation": "0.4990000000"
}
}

# Change the parameters through files, return changed parameters
iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange" --deposit="10iris"  --path=iris --key=Gov/govTallyingProcedure --op=update --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# Deposit for a proposal
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# Vote for a proposal
iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# Query the state of a proposal
iriscli gov query-proposal --proposal-id=1 --trust-node
```

### Proposals on software upgrade

Detail in [Upgrade](upgrade.md)

## Basic parameters

```
# DepositProcedure（The parameters in deposit period）
"Gov/govDepositProcedure": {
"min_deposit": [
{
"denom": "iris-atto",
"amount": "10000000000000000000"
}
],
"max_deposit_period": "172800000000000"
}
```

* Parameters can be changed
* The key of parameters:"Gov/gov/DepositProcedure"
* `min_deposit[0].denom`  The minimum tokens deposited are counted by iris-atto.
* `min_deposit[0].amount` The number of minimum tokens and the default scope：1000iris,（1iris，10000iris）
* `max_deposit_period`    Window period for repaying deposit, default:172800000000000ns==2Days， scope（20s，3Day）     

```
# VotingProcedure（The parameters in voting period）
"Gov/govVotingProcedure": {
"voting_period": "10000000000"
}
```

* Parameters can be changed   
* `voting_perid`  Window period for vote, default:172800000000000ns==2Days, scope（20s，3Days）

```
# TallyingProcedure (The parameters in Tallying period)    
"Gov/govTallyingProcedure": {
"threshold": "0.5000000000",
"veto": "0.3340000000",
"participation": "0.6670000000"
}
``` 

* Parameters can be changed
* `veto` default: 0.334, scope（0，1）
* `threshold` default: 0.5, scope（0，1）
* `governance_penalty` default: 0.667, scope（0，1）
*  Vote rules:If the ratio of all voters' `voting_power` to the total 'voting_power' in system less than “participation”, the proposal won't be passed. If the ratio of strongly opposed `voting_power` to all voters' `voting_power` more than “veto”, the proposal won't be passed. Then if the ratio of approved `voting_power` to all voter's `voting_power` except abstentions over “threshold”, the proposal will be passed. Otherwise, N/A.


### Proposals on community funds usage
There are three usages, `Burn`, `Distribute` and `Grant`. `Burn` means burning tokens from community funds. `Distribute` and `Grant` will transfer tokens to the destination trustee's account from community funds and then trustee will distribute or grant these tokens to others.
```shell
# Submit Burn usage proposal
iriscli gov submit-proposal --title="burn tokens 5%" --description="test" --type="TxTaxUsage" --usage="Burn" --deposit="10iris"  --percent=0.05 --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Submit Distribute usage proposal
iriscli gov submit-proposal --title="distribute tokens 5%" --description="test" --type="TxTaxUsage" --usage="Distribute" --deposit="10iris"  --percent=0.05 --dest-address=[destnation-address] --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Submit Grant usage proposal
iriscli gov submit-proposal --title="grant tokens 5%" --description="test" --type="TxTaxUsage" --usage="Grant" --deposit="10iris"  --percent=0.05 --dest-address=[destnation-address] --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Deposit for a proposal
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Vote for a proposal
iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=200000

# Query the state of a proposal
iriscli gov query-proposal --proposal-id=1
```