# Gov/Iparam User Guide

## Basic Function Description

1. On-chain governance proposals on text
2. On-chain governance proposals on parameter change
3. On-chain governance proposals on software upgrade （unavailable)

## Interactive process

### governance process

1. Any users can deposit some tokens to initiate a proposal. Once deposit reaches a certain value `min_deposit`, enter voting period, otherwise it will remain in the deposit period. Others can deposit the proposals on the deposit period. Once the sum of the deposit reaches `min_deposit`, enter voting period. However, if the block-time exceeds `max_deposit_period` in the deposit period, the proposal will be closed.
2. The proposals which enter voting period only can be voted by validators and delegators. The vote of a delegator who hasn't vote will be the same as his validator's vote, and the vote of a delegator who has voted will be remained. The votes wil be tallyed when reach `voting_period'.
3. More details about voting for proposals:
[CosmosSDK-Gov-spec](https://github.com/cosmos/cosmos-sdk/blob/develop/docs/spec/governance/overview.md)

## Usage Scenario
### Create an environment

```
rm -rf iris                                                                         
rm -rf .iriscli
iris init gen-tx --name=x --home=iris
iris init --gen-txs --chain-id=gov-test -o --home=iris
iris start --home=iris
```

### Usage scenario of parameter change

Scenario 1：Change the parameters through the command lines

```
# Query parameters can be changed by the modules'name in gov 
iriscli gov query-params --module=gov --trust-node

# Results
[
 "Gov/gov/DepositProcedure",
 "Gov/gov/TallyingProcedure",
 "Gov/gov/VotingProcedure"
]

# Query parameters can be modified by "key”
iriscli gov query-params --key=Gov/gov/DepositProcedure --trust-node

# Results
{"key":"Gov/gov/DepositProcedure","value":"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":10}","op":""}

# Send proposals, return changed parameters
echo 1234567890 | iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange" --deposit="10iris"  --param='{"key":"Gov/gov/DepositProcedure","value":"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":20}","op":"update"}' --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# Deposit for a proposal
echo 1234567890 | iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# Vote for a proposal
echo 1234567890 | iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# Query the state of a proposal
iriscli gov query-proposal --proposal-id=1 --trust-node

```

Scenario 2: Change the parameters by the files

```
# Export profiles
iriscli gov pull-params --path=iris --trust-node

# Query profiles' info
cat iris/config/params.json                                                         
{
  "gov": {
    "Gov/gov/DepositProcedure": {
      "min_deposit": [
        {
          "denom": "iris-atto",
          "amount": "10000000000000000000"
        }
      ],
      "max_deposit_period": "10"
    },
    "Gov/gov/VotingProcedure": {
      "voting_period": "10"
    },
    "Gov/gov/TallyingProcedure": {
      "threshold": "1/2",
      "veto": "1/3",
      "governance_penalty": "1/100"
    }
  }
}
# Modify profiles (TallyingProcedure的governance_penalty)
vi iris/config/params.json                                                            
{
  "gov": {
    "Gov/gov/DepositProcedure": {
      "min_deposit": [
        {
          "denom": "iris-atto",
          "amount": "10000000000000000000"
        }
      ],
      "max_deposit_period": "10"
    },
    "Gov/gov/VotingProcedure": {
      "voting_period": "10"
    },
    "Gov/gov/TallyingProcedure": {
      "threshold": "1/2",
      "veto": "1/3",
      "governance_penalty": "20/100"
    }
  }
}

# Change the parameters through files, return changed parameters
echo 1234567890 | iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange" --deposit="10iris"  --path=iris --key=Gov/gov/TallyingProcedure --op=update --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# Deposit for a proposal
echo 1234567890 | iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# Vote for a proposal
echo 1234567890 | iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=20000

# Query the state of a proposal
iriscli gov query-proposal --proposal-id=1 --trust-node
```

## CLI Command Details

### Basic method of gov modules

```
# Text proposals
iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="Text" --deposit="10iris" --from=x --chain-id=gov-test --fee=0.05iris --gas=20000
```

* `--title`       The title of a proposal
* `--description` The description of a proposal
* `--type`        The type of a proposal {'Text','ParameterChange','SoftwareUpgrade'}
* `--deposit`     The number of the tokens deposited
* The basic text proposals are as below

```
iriscli gov deposit --proposal-id=1 --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=20000
```

* `--propsal-id` The ID of the proposal deposited
* `--deposit`    The number of the tokens deposited

```
iriscli gov vote --proposal-id=1 --option=Yes  --from=x --chain-id=gov-test --fee=0.05iris --gas=20000
```

* `--proposal-id` The ID of the proposal in voting period
* `--option`      Vote option{'Yes'-agree,'Abstain'-abstain,'No'-disagree,'nowithVeto'-strongly disagree }


```
# Query the state of a proposal
iriscli gov query-proposal --proposal-id=1 --trust-node
```

* `--proposal-id` Query the ID of a proposal



### The proposals on parameters modification

```
# Query parameters can be modified by the modules'name in gov 
iriscli gov query-params --module=gov --trust-node
```

* `--module` Query the list of "key" of the parameters can be changed in the module


```
# Query the parameters can be modified by "key"
iriscli gov query-params --key=Gov/gov/DepositProcedure --trust-node
```

* `--key` Query the parameter corresponding to the "key"

```
# Export profiles
iriscli gov pull-params --path=iris --trust-node
```

* `--path` The folder of node initialization



```
# Modify the parameters through the command lines
iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange" --deposit="10iris"  --param='{"key":"Gov/gov/DepositProcedure","value":"{\"min_deposit\":[{\"denom\":\"iris-atto\",\"amount\":\"10000000000000000000\"}],\"max_deposit_period\":20}","op":"update"}' --from=x --chain-id=gov-test --fee=0.05iris --gas=20000
```

* `--param` The details of changed parameters （get parameters through query-params, modify it and then add "update" on the "op", more details in usage scenarios）
* Other fields' proposals are similar with text proposal

```
# Change the parameters through files, return modified parameters
echo 1234567890 | iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange" --deposit="10iris"  --path=iris --key=Gov/gov/TallyingProcedure --op=update --from=x --chain-id=gov-test --fee=0.05iris --gas=20000
```

* `--path` The folder of node initialization 
* `--key`  The key of the parameter to be modified
* `--op`   The type of changed parameters; only 'update' is implemented at present
* Other fields' proposals are similar with text proposal

### Proposals on software upgrade 

## Basic parameters

```
# DepositProcedure（The parameters in deposit period）
"Gov/gov/DepositProcedure": {
    "min_deposit": [
    {
        "denom": "iris-atto",
        "amount": "10000000000000000000"
    }
    ],
    "max_deposit_period": "10"
}
```

* Parameters can be changed
* The key of parameters:"Gov/gov/DepositProcedure"
* `min_deposit[0].denom`  The minimum tokens deposited are counted by iris-atto.
* `min_deposit[0].amount` The number of minimum tokens and the default scope：10iris,（1iris，200iris）
* `max_deposit_period`    Window period for repaying deposit, default :10， scope（0，1）     

```
# VotingProcedure（The parameters in voting period）
"Gov/gov/VotingProcedure": {
    "voting_period": "10"
},
```

* Parameters can be changed   
* `voting_perid`  Window period for vote, default：10, scope（20，20000）
   
```
# TallyingProcedure (The parameters in Tallying period)    
"Gov/gov/TallyingProcedure": {
    "threshold": "1/2",
    "veto": "1/3",
    "governance_penalty": "1/100"
}
``` 
  
* Parameters can be changed
* `veto` default: 1/3, scope（0，1）
* `threshold` default 1/2, scope（0，1）
* `governance_penalty` The default ratio of slashing tokens of validators who didn't vote: 1/100, scope（0，1）
*  Vote rules: If the ratio of voting power of "strongly disagree" over "veto", the proposal won't be passed. If the ratio of voting_power of "agree" over "veto", the proposal won't be passed. Otherwise, it will be passed.

