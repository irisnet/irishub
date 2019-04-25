# Distribution User Guide

## Introduction 

This module is in charge of distributing collected transaction fee and inflated token to all validators and delegators. 
To reduce computation stress, a lazy distribution strategy is brought in. 
`lazy` means that the benefit won't be paid directly to contributors automatically. 
The contributors are required to explicitly send transactions to withdraw their benefit, otherwise, 
their benefit will be kept in the global pool. 

## Benefit

### Source

1.The first signer of transactions (Collected to feeCollector in DeliverTx)
2.Inflation tokens (Each hour Inflate stake token and add it to LooseTokens)

### Destination

1.Validator Operators: Self delegation benefit and delegation commission
2.Delegators: Delegation benefit
3.Community Funding: Community tax
4.Proposer Reward

If one validator is the proposer of current round, that validator (and their delegators) receives between 1% and 5% of the sum of fee rewards and inflated token as proposer reward.
It is calculated as:
```
 proposerReward = (TxFee + InflatedToken) * (0.01 + 0.04 * sumPowerPrecommitValidators / totalBondedTokens)
```

## Usage Scenario

### Set withdraw address

By default, the reward will be paid to the wallet address which send the delegation transaction.

The delegator could set a new wallet as reward paid address. To set another wallet(marked as `B`) as the paid address, delegator need to send another transaction from wallet `A`.

```bash
iriscli distribution set-withdraw-addr <address_of_wallet_B> --fee=0.3iris --from=<key_name_of_ wallet_A> --chain-id=<chain-id>
```  

Query withdraw address：

```bash
iriscli distribution withdraw-address <address_of_wallet_A> 
```
### Withdraw reward 

There are 3 ways to withdraw reward according to different scenarios

1.`WithdrawDelegationRewardsAll` : Withdraw all delegation reward

```bash
iriscli distribution withdraw-rewards --from=<key_name> --fee=0.3iris --chain-id=<chain-id>
```

2.`WithdrawDelegatorReward` : Only withdraw the self-delegation reward of from designated validator

```bash
iriscli distribution withdraw-rewards --only-from-validator=<validator_address>  --from=<key_name> --fee=0.3iris --chain-id=<chain-id>
```

3.`WithdrawValidatorRewardsAll` : Withdraw all delegation reward including commission benefit, only for validator

```bash
iriscli distribution withdraw-rewards --is-validator=true --from=<key_name> --fee=0.3iris --chain-id=<chain-id>
```

### Query reward token

There are 3 ways to query reward according to different scenarios

1.Execute `rewards` command. 

```bash
iriscli distribution rewards <delegator_address>
```

Output：
```bash
{
  "total": [
    {
      "denom": "iris-atto",
      "amount": "2035754787730512363646"
    }
  ],
  "delegations": [
    {
      "validator": <validator_address>,
      "reward": [
        {
          "denom": "iris-atto",
          "amount": "1052463556823086428786"
        }
      ]
    }
  ],
  "commission": [
    {
      "denom": "iris-atto",
      "amount": "983291230907425934859"
    }
  ]
}
```

2.Use `dry-run` mode (simulation only , tx won't be broadcast)

Execute command(validator only):
```bash
iriscli distribution withdraw-rewards --is-validator=true --from=node0 --dry-run --chain-id=irishub-stage --fee=0.3iris --commit
```

Output：`withdraw-reward-total`is your estimated inflation rewards
```bash
estimated gas = 16768
simulation code = 0
simulation log = Msg 0: 
simulation gas wanted = 50000
simulation gas used = 11179
simulation fee amount = 0
simulation fee denom = 
simulation tag action = withdraw_validator_rewards_all
simulation tag source-validator = iva1rulhmls7g9cjh239vnkjnw870t5urrut9cyrxl
simulation tag withdraw-reward-total = 2035775375047308887487iris-atto
simulation tag withdraw-address = iaa18cgtskr6cgqyyady8mumk05xk2g9c95qgw5556
simulation tag withdraw-reward-from-validator-iva1rulhmls7g9cjh239vnkjnw870t5urrut9cyrxl = 1052484144134629789682iris-atto
simulation tag withdraw-reward-commission = 983291230912679097804iris-atto
```

Execute command:
```bash
iriscli distribution withdraw-rewards --from=<key_name> --dry-run --chain-id=<chain-id> --fee=0.3iris --commit
```

Output：
```bash
estimated gas = 14329
simulation code = 0
simulation log = Msg 0: 
simulation gas wanted = 50000
simulation gas used = 9553
simulation fee amount = 0
simulation fee denom = 
simulation tag action = withdraw_delegation_rewards_all
simulation tag delegator = iaa1rulhmls7g9cjh239vnkjnw870t5urrutsfwvmc
simulation tag withdraw-reward-total = 1052472042330962430914iris-atto
simulation tag withdraw-address = iaa18cgtskr6cgqyyady8mumk05xk2g9c95qgw5556
simulation tag withdraw-reward-from-validator-iva1rulhmls7g9cjh239vnkjnw870t5urrut9cyrxl = 1052472042330962430914iris-atto
```