# Distribution

## Introduction

This module is in charge of distributing collected transaction fee and inflated token to all validators and delegators. To reduce computation stress, a lazy distribution strategy is brought in. `lazy` means that the benefit won't be paid directly to contributors automatically. The contributors are required to explicitly send transactions to withdraw their benefit, otherwise, their benefit will be kept in the global pool.

## Benefit

### Source

1. The first signer of transactions (Collected to feeCollector in DeliverTx)
2. Inflation tokens (Each hour Inflate stake token and add it to LooseTokens)

### Destination

1. Validator
2. Delegator
3. Community Tax
4. Proposer Reward

:::tip
[Calculation Formula](../concepts/general-concepts.md#staking-rewards-calculation-formula)
:::

## Usage Scenario

### Set withdraw address

By default, the reward will be paid to the wallet address which send the delegation transaction.

The delegator could set a new wallet as reward paid address. To set another wallet(marked as `B`) as the paid address, delegator need to send another transaction from wallet `A`.

```bash
iriscli distribution set-withdraw-addr <address-of-wallet-B> --fee=0.3iris --from=<key-name-of- wallet-A> --chain-id=<chain-id>
```  

Query withdraw address:

```bash
iriscli distribution withdraw-address <address-of-wallet-A>
```

### Withdraw rewards

There are 3 ways to withdraw rewards according to different scenarios

- `WithdrawDelegationRewardsAll` : Withdraw all delegation reward

  ```bash
  iriscli distribution withdraw-rewards --from=<key-name> --fee=0.3iris --chain-id=<chain-id>
  ```

- `WithdrawDelegatorReward` : Only withdraw the self-delegation reward of from designated validator

  ```bash
  iriscli distribution withdraw-rewards --only-from-validator=<validator-address>  --from=<key-name> --fee=0.3iris --chain-id=<chain-id>
  ```

- `WithdrawValidatorRewardsAll` : Withdraw all delegation reward including commission benefit, only for validator

  ```bash
  iriscli distribution withdraw-rewards --is-validator=true --from=<key-name> --fee=0.3iris --chain-id=<chain-id>
  ```

### Query reward token

There are 2 ways to query rewards according to different scenarios

- Execute `rewards` query command.

  ```bash
  iriscli distribution rewards <delegator-address>
  ```

  Example Output：

  ```bash
  Total:        270.33761964714393479iris
  Delegations:  
    validator: iva1q7602ujxxx0urfw7twm0uk5m7n6l9gqsgw4pqy, reward: 2.899411557255275253iris
  Commission:   267.438208089888659537iris
  ```

- Use `dry-run` mode (simulation only , tx won't be broadcasted)

Execute command(validator only):

```bash
iriscli distribution withdraw-rewards --is-validator=true --from=node0 --dry-run --chain-id=irishub-stage --fee=0.3iris --commit
```

Example Output：`withdraw-reward-total`is your estimated inflation rewards

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
