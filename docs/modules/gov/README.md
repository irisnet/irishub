# Governance

## Initialize the runtime environment of blockchain

```
rm -rf .iris
rm -rf .iriscli
iris init gen-tx --name=iris
iris init --gen-txs --chain-id=gov-test
iris start

```

## Proposal process

Here is an example of "parameter modification proposal". Other proposals do not need parameter "--params". For instance, sometimes we find there are a lot of useless proposals in the system, probably because the current minimum delegating amount is too small, so that many people submit meaningless ones. Here we can use "parameter modification proposal" to modify the system default minimum parameter of “delegating amount”. Firstly we need to know the value of key of this parameter before modification. The following command can be used to view:
```
iriscli params export gov
```
This command will export all the parameters that can be modified with the "parameter modification proposal". For instance, we get the following results:

```
[
  {
    "key": "gov/depositprocedure/deposit",
    "value": "10000000000000000000iris"
  },
  {
    "key": "gov/depositprocedure/maxDepositPeriod",
    "value": "10"
  },
  {
    "key": "gov/feeToken/gasPriceThreshold",
    "value": "20000000000"
  },
  {
    "key": "gov/tallyingprocedure/penalty",
    "value": "1/100"
  },
  {
    "key": "gov/tallyingprocedure/threshold",
    "value": "1/2"
  },
  {
    "key": "gov/tallyingprocedure/veto",
    "value": "1/3"
  },
  {
    "key": "gov/votingprocedure/votingPeriod",
    "value": "20"
  }
]

```
Each (key, value) is corresponding to a set of system-preset modifiable parameters, the specific meaning of which will be updated in future documents. The key of minimum delegating amount here is gov/depositprocedure/deposit. At present, value=10000000000000000000iris（This is the value after the decimal-binary conversion, refering to the “fee-token” module for details.）. To raise the threshold of proposals, it will be doubled to 20000000000000000000iris. The following command can be used:

```
iriscli gov submit-proposal --title="update MinDeposit" --description="test" --type="ParameterChange"
                            --deposit="9000000000000000000iris" 
                            --params='[{"key":"gov/depositprocedure/deposit","value":"20000000000000000001iris","op":"update"}]' 
                            --proposer=faa1pkunlumfyglqd9dgup0mwp66kjrp6y09twmuvd 
                            --from=iris 
                            --chain-id=gov-test 
                            --fee=400000000000000iris 
                            --gas=20000

```

Here I delegate 9000000000000000000iris, which is 1000000000000000000iris less than the minimum delegating amount, so the proposal has not been activated can not be voted. 10000000000000000000 more iris should be delegated in 10 (key:gov/depositprocedure/maxDepositPeriod) blocks (If it needs 5s to produce a block, it means you need to complete the delegation in 5 * 10s). The delegating command is as follows:

```
iriscli gov deposit --proposalID=1 
                    --depositer=faa1pkunlumfyglqd9dgup0mwp66kjrp6y09twmuvd 
                    --deposit=1000000000000000000iris   
                    --from=iris 
                    --chain-id=gov-test  
                    --fee=200000000000000iris 
                    --gas=20000

```
The above proposalID is the result from the first step. At this stage, we delegate 1000000000000000000iris tokens, which is exactly equal to the minimum delegating amount, so the proposal can be voted. And the proposer can send a voting request to each validators (currently only off-chain notification is available, but on-chain or monitoring notification will be implemented later). Then each validators can view the proposal first with the following command:

```
iriscli gov query-proposal --proposalID=1 
```

Later proposers can vote as they wish, here I vote Yes (option=Yes):

```
iriscli gov vote --proposalID=1 
                 --voter=faa1pkunlumfyglqd9dgup0mwp66kjrp6y09twmuvd 
                 --option=Yes  
                 --from=iris 
                 --chain-id=gov-test  
                 --fee=400000000000000iris 
                 --gas=20000
```

Notice that the maximum of waiting time is in 20 blocks during the voting period (key: gov/votingprocedure/votingPeriod). If the percentage of affirmative vote is still less than 50% during this period (key: gov/tallyingprocedure/threshold), the proposal will not be passed and the tokens delegated will also not be refunded (the validators haven't voted will be slashed, and 1/100 of the total tokens delegated currently will be deducted (key: gov/tallyingprocedure/penalty). This mechanism has not been implemented in current version ). Suppose there is only one validator. If I voted yes, the ratio of affirmative vote is 1>1/2 and the strong negative vote is 0<1/3 (key:gov/tallyingprocedure/veto), the proposal will be passed. After voting, the proposal  is automatically executed: (key: gov/depositprocedure/deposit, value: 10000000000000000000iris) is modified to (key: gov/depositprocedure/deposit, value: 20000000000000000000iris). Then we can verify this result and query the minimum delegating amount in current system:

```
iriscli iriscli params export gov/depositprocedure/deposit
```

This is the end of the governance process.
