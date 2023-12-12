# Distribution

## Summary

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
iris tx distribution set-withdraw-addr [withdraw-addr] [flags] --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```  

### Withdraw rewards

There are 2 ways to withdraw rewards according to different scenarios

- `withdraw-all-rewards` : Withdraw all delegations rewards for a delegator

```bash
iris tx distribution withdraw-all-rewards [flags] --from=<key-name> --fees=0.3iris --chain-id=irishub
```

- `withdraw-rewards` : Withdraw rewards from a given validator address

```bash
iris tx distribution withdraw-rewards [validator-addr] [flags] --from=<key-name> --fees=0.3iris --chain-id=irishub
```

### Query reward token

Query all rewards earned by a delegator, optionally restrict to rewards from a single validator.

```bash
iris query distribution rewards [delegator-addr] [validator-addr] [flags]
```

For other distribution commands, please refer to [distribution cli client](../cli-client/distribution.md)
