# iris distribution

The distribution module allows you to manage your [Staking Rewards](../concepts/general-concepts.md#staking-rewards).

## Available Subcommands

| Name                                                                                      | Description                                                  |
| ----------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- |
| [commission](#iris-query-distribution-commission)                                         | Query distribution validator commission                                                                                |
| [community-pool](#iris-query-distribution-community-pool)                                 | Query the amount of coins in the community pool                                                                   |
| [params](#iris-query-distribution-params)                                                 | Query distribution params                                                                                   |
| [rewards](#iris-query-distribution-rewards)                                               | Query all distribution delegator rewards or rewards from a particular validator  |
| [slashes](#iris-query-distribution-slashes)                                               | Query distribution validator slashes.                                                                                   |
| [validator-outstanding-rewards](#iris-tx-distribution-validator-outstanding-rewards)      | Query distribution outstanding (un-withdrawn) rewards for a validator and all their delegations                                                                                   |
| [fund-community-pool](#iris-tx-distribution-fund-community-pool)                          | Funds the community pool with the specified amount                                                         |
| [set-withdraw-addr](#iris-tx-distribution-set-withdraw-addr)                              | Set the withdraw address for rewards associated with a delegator address                                                                                   |
| [withdraw-all-rewards](#iris-tx-distribution-withdraw-all-rewards)                        | Withdraw all rewards for a single delegator                                                                                   |
| [withdraw-rewards](#iris-tx-distribution-withdraw-rewards)                                | Withdraw rewards from a given delegation address,and optionally withdraw validator commission if the delegation address given is a validator operator  |

## iris query distribution commission

Query validator commission rewards from delegators to that validator.

```bash
iris query distribution commission [validator] [flags]
```

## iris query distribution community-pool

Query all coins in the community pool which is under Governance control.

```bash
iris query distribution community-pool [flags]
```

## iris query distribution params

Query distribution params.

```bash
 iris query distribution params [flags]
```

## iris query distribution rewards

Query all rewards earned by a delegator, optionally restrict to rewards from a single validator.

```bash
iris query distribution rewards [delegator-addr] [validator-addr] [flags]
```

## iris query distribution slashes

Query all slashes of a validator for a given block range.

```bash
iris query distribution slashes [validator] [start-height] [end-height] [flags]
```

## iris query distribution validator-outstanding-rewards

Query distribution outstanding (un-withdrawn) rewards for a validator and all their delegations.

```bash
iris query distribution validator-outstanding-rewards [validator] [flags]
```
## iris tx distribution fund-community-pool

Funds the community pool with the specified amount.

```bash
iris tx distribution fund-community-pool [amount] [flags] [validator-addr] [flags]
```
## iris tx distribution set-withdraw-addr

Set the withdraw address for rewards associated with a delegator address.

```bash
iris tx distribution set-withdraw-addr [withdraw-addr] [flags]
```

## iris tx distribution withdraw-all-rewards

Withdraw all rewards for a single delegator.

```bash
iris tx distribution withdraw-all-rewards [flags]
```

## iris tx distribution withdraw-rewards

Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator.

```bash
iris tx distribution withdraw-rewards [validator-addr] [flags]
```
