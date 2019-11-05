# iriscli distribution

The distribution module allows you to manage your [Staking Rewards](../concepts/general-concepts.md#staking-rewards).

## Available Subommands

| Name                                                         | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| [withdraw-address](#iriscli-distribution-withdraw-address)   | Query withdraw address                                       |
| [rewards](#iriscli-distribution-rewards)                     | Query all the rewards of validator or delegator              |
| [set-withdraw-address](#iriscli-distribution-set-withdraw-addr) | Change withdraw address                                      |
| [withdraw-rewards](#iriscli-distribution-withdraw-rewards)   | withdraw rewards for either: all-delegations, a delegation, or a validator |

## iriscli distribution withdraw-address

Query the withdraw address of a delegator

```bash
iriscli distribution withdraw-address <delegator-address> <flags>
```

### Query withdraw address

```bash
iriscli distribution withdraw-address <delegator-address>
```

If the delegator did not specify the withdraw address other than himself, the query result will be empty.

## iriscli distribution rewards

Query all the rewards of a validator or a delegator

```bash
iriscli distribution rewards <address> <flags>
```

### Query rewards

```bash
iriscli distribution rewards <iaa...>
```

Output:

```bash
Total:        270.33761964714393479iris
Delegations:  
  validator: iva..., reward: 2.899411557255275253iris
  validator: iva..., reward: 2.899411557255275253iris
  validator: iva..., reward: 2.899411557255275253iris
Commission:   267.438208089888659537iris
```

## iriscli distribution set-withdraw-addr

Set another address to receive the rewards instead of using the delegator address

```bash
iriscli distribution set-withdraw-addr <withdraw-address> <flags>
```

### Set withdraw address

```bash
iriscli distribution set-withdraw-addr <iaa...> --from=<key-name> --fee=0.3iris --chain-id=irishub
```

## iriscli distribution withdraw-rewards

Withdraw rewards to the withdraw-address(default is the delegator address, you can set to another address via [set-withdraw-addr](#iriscli-distribution-set-withdraw-addr))

```bash
iriscli distribution withdraw-rewards <flags>
```

**Flags:**

| Name, shorthand       | type   | Required | Default  | Description                                                         |
| --------------------- | -----  | -------- | -------- | ------------------------------------------------------------------- |
| --only-from-validator | string |          |          | Only withdraw from this validator address (in bech)                 |
| --is-validator        | bool   |          | false    | Also withdraw validator's commission                                |

:::tip
Do not specify the above 2 flags together
:::

### Withdraw delegation rewards from a specified validator

```bash
iriscli distribution withdraw-rewards --only-from-validator=<validator-address> --from=<key-name> --fee=0.3iris --chain-id=irishub
```

### Withdraw all delegation rewards

```bash
iriscli distribution withdraw-rewards --from=<key-name> --fee=0.3iris --chain-id=irishub
```

### Validator withdraws all delegation rewards and commission rewards

```bash
iriscli distribution withdraw-rewards --is-validator=true --from=<key-name> --fee=0.3iris --chain-id=irishub
```
