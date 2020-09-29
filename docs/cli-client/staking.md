# Staking

Staking module provides a set of subcommands to query staking state and send staking transactions.

## Available Commands

| Name                                                                         | Description                                                                                   |
| ---------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| [validator](#iris-query-staking-validator)                                   | Query a validator                                                                             |
| [validators](#iris-query-staking-validators)                                 | Query for all validators                                                                      |
| [delegation](#iris-query-staking-delegation)                                 | Query a delegation based on address and validator address                                     |
| [delegations](#iris-query-staking-delegations)                               | Query all delegations made from one delegator                                                 |
| [delegations-to](#iris-query-staking-delegations-to)                         | Query all delegations to one validator                                                        |
| [unbonding-delegation](#iris-query-staking-unbonding-delegation)             | Query an unbonding-delegation record based on delegator and validator address                 |
| [unbonding-delegations](#iris-query-staking-unbonding-delegations)           | Query all unbonding-delegations records for one delegator                                     |
| [unbonding-delegations-from](#iris-query-staking-unbonding-delegations-from) | Query all unbonding delegatations from a validator                                            |
| [redelegations-from](#iris-query-staking-redelegations-from)                 | Query all outgoing redelegatations from a validator                                           |
| [redelegation](#iris-query-staking-redelegation)                             | Query a redelegation record based on delegator and a source and destination validator address |
| [redelegations](#iris-query-staking-redelegations)                           | Query all redelegations records for one delegator                                             |
| [pool](#iris-query-staking-pool)                                             | Query the current staking pool values                                                         |
| [params](#iris-query-staking-params)                                         | Query the current staking parameters information                                              |
| [historical-info](#iris-query-staking-historical-info)                       | Query historical info at given height                                                         |
| [create-validator](#iris-tx-staking-create-validator)                        | Create new validator initialized with a self-delegation to it                                 |
| [edit-validator](#iris-tx-staking-edit-validator)                            | Edit existing validator account                                                               |
| [delegate](#iris-tx-staking-delegate)                                        | Delegate liquid tokens to an validator                                                        |
| [unbond](#iris-tx-staking-unbond)                                            | Unbond shares from a validator                                                                |
| [redelegate](#iris-tx-staking-redelegate)                                    | Redelegate illiquid tokens from one validator to another                                      |

## iris query staking validator

### Query a validator by validator address

```bash
iris query staking validator <iva...>
```

## iris query staking validators

### Query all validators

```bash
iris query staking validators
```

## iris query staking delegation

Query a delegation based on delegator address and validator address.

```bash
iris query staking delegation [delegator-addr] [validator-addr]
```

### Query a delegation

```bash
iris query staking delegation <iaa...> <iva...>
```

Example Output:

```bash
Delegation:
  Delegator:  iaa13lcwnxpyn2ea3skzmek64vvnp97jsk8qrcezvm
  Validator:  iva15grv3xg3ekxh9xrf79zd0w077krgv5xfzzunhs
  Shares:     1.0000000000000000000000000000
  Height:     26
```

## iris query staking delegations

Query all delegations delegated from one delegator.

```bash
iris query staking delegations [delegator-address] [flags]
```

### Query all delegations of a delegator

```bash
iris query staking delegations <iaa...>
```

## iris query staking delegations-to

Query all delegations to one validator.

```bash
iris query staking delegations-to [validator-address] [flags]
```

### Query all delegations to one validator

```bash
iris query staking delegations-to <iva...>
```

Example Output:

```bash
Delegation:
  Delegator:  iaa13lcwnxpyn2ea3skzmek64vvnp97jsk8qrcezvm
  Validator:  iva1yclscskdtqu9rgufgws293wxp3njsesxxlnhmh
  Shares:     100.0000000000000000000000000000
  Height:     0
Delegation:
  Delegator:  iaa1td4xnefkthfs6jg469x33shzf578fed6n7k7ua
  Validator:  iva1yclscskdtqu9rgufgws293wxp3njsesxxlnhmh
  Shares:     1.0000000000000000000000000000
  Height:     26
```

## iris query staking unbonding-delegation

Query an unbonding-delegation record based on delegator and validator address.

```bash
iris query staking unbonding-delegation [delegator-addr] [validator-addr] [flags]```
```

### Query an unbonding delegation record

```bash
iris query staking unbonding-delegation <iaa...> <iva...>
```

## iris query staking unbonding-delegations

### Query all unbonding delegations records of a delegator

```bash
iris query staking unbonding-delegations <iaa...>
```

## iris query staking unbonding-delegations-from

### Query all unbonding delegations from a validator

```bash
iris query staking unbonding-delegations-from <iva...>
```

## iris query staking redelegations-from

Query all outgoing redelegations of a validator

```bash
iris query staking redelegations-from [validator-address] [flags]
```

### Query all outgoing redelegatations of a validator

```bash
iris query staking redelegations-from <iva...>
```

## iris query staking redelegation

Query a redelegation record based on delegator and source validator address and destination validator address.

```bash
iris query staking redelegation [delegator-addr] [src-validator-addr] [dst-validator-addr] [flags]
```

### Query a redelegation record

```bash
iris query staking redelegation <iaa...> <iva...> <iva...>
```

## iris query staking redelegations

### Query all redelegations records of a delegator

```bash
iris query staking redelegations <iaa...>
```

## iris query staking pool

### Query the current staking pool values

```bash
iris query staking pool
```

Example Output:

```bash
Pool:
  Loose Tokens:   1409493892.759816067399143966
  Bonded Tokens:  590526409.65743521209068061
  Token Supply:   2000020302.417251279489824576
  Bonded Ratio:   0.2952602076
```

## iris query staking params

### Query the current staking parameters information

```bash
iris query staking params
```

## iris query staking historical-info

### Query historical info at given height

```bash
iris query staking historical-info <height>
```

## iris tx staking create-validator

Send a transaction to apply to be a validator and delegate a certain amount of iris to it.

```bash
iris tx staking create-validator [flags]
```

**Flags:**

| Name, shorthand              | type   | Required | Default | Description                                                                                      |
| ---------------------------- | ------ | -------- | ------- | ------------------------------------------------------------------------------------------------ |
| --amount                     | string | Yes      |         | Amount of coins to bond                                                                          |
| --commission-rate            | float  | Yes      | 0.0     | The initial commission rate percentage                                                           |
| --commission-max-rate        | float  |          | 0.0     | The maximum commission rate percentage                                                           |
| --commission-max-change-rate | float  |          | 0.0     | The maximum commission change rate percentage (per day)                                          |
| --min-self-delegation        | string |          |         | The minimum self delegation required on the validator                                            |
| --details                    | string |          |         | Optional details                                                                                 |
| --genesis-format             | bool   |          | false   | Export the transaction in gen-tx format; it implies --generate-only                              |
| --identity                   | string |          |         | Optional identity signature (ex. UPort or Keybase)                                               |
| --ip                         | string |          |         | Node's public IP. It takes effect only when used in combination with                             |
| --node-id                    | string |          |         | The node's ID                                                                                    |
| --moniker                    | string | Yes      |         | Validator name                                                                                   |
| --pubkey                     | string | Yes      |         | Go-Amino encoded hex PubKey of the validator. For Ed25519 the go-amino prepend hex is 1624de6220 |
| --website                    | string |          |         | Optional website                                                                                 |
| --security-contact           | string |          |         | The validator's (optional) security contact email                                                |

### Create a validator

```bash
iris tx staking create-validator --chain-id=irishub --from=<key-name> --fees=0.3iris --pubkey=<validator-pubKey> --commission-rate=0.1 --amount=100iris --moniker=<validator-name>
```

:::tip
Follow the [Mainnet](../get-started/mainnet.md#create-validator) instructions to learn more.
:::

## iris tx staking edit-validator

Edit an existing validator's settings, such as commission rate, name, etc.

```bash
iris tx staking edit-validator [flags]
```

**Flags:**

| Name, shorthand       | type   | Required | Default | Description                                           |
| --------------------- | ------ | -------- | ------- | ----------------------------------------------------- |
| --commission-rate     | float  |          | 0.0     | Commission rate percentage                            |
| --moniker             | string |          |         | Validator name                                        |
| --identity            | string |          |         | Optional identity signature (ex. UPort or Keybase)    |
| --website             | string |          |         | Optional website                                      |
| --details             | string |          |         | Optional details                                      |
| --security-contact    | string |          |         | The validator's (optional) security contact email     |
| --min-self-delegation | string |          |         | The minimum self delegation required on the validator |

### Edit validator information

```bash
iris tx staking edit-validator --from=<key-name> --chain-id=irishub --fees=0.3iris --commission-rate=0.10 --moniker=<validator-name>
```

### Upload validator avatar

Please refer to [How to upload my validator's logo to the Explorers](../concepts/validator-faq.md#how-to-upload-my-validator-s-logo-to-the-explorers)

## iris tx staking delegate

Delegate tokens to a validator.

```bash
iris tx staking delegate [validator-addr] [amount] [flags]
```

```bash
iris tx staking delegate <iva...> <amount> --chain-id=irishub --from=<key-name> --fees=0.3iris
```

## iris tx staking unbond

Unbond tokens from a validator.

```bash
iris tx staking unbond [validator-addr] [amount] [flags]
```

### Unbond some tokens from a validator

```bash
iris tx staking unbond <iva...> 10iris --from=<key-name> --chain-id=irishub --fees=0.3iris
```

## iris tx staking redelegate

Transfer delegation from one validator to another.

:::tip
There is no `unbonding time` during the redelegation, so you will not miss the rewards. But you can only redelegate once per validator, until a period (= `unbonding time`) exceed.
:::

```bash
iris tx staking redelegate [src-validator-addr] [dst-validator-addr] [amount] [flags]
```

### Redelegate some tokens to another validator

```bash
iris tx staking redelegate <iva...> <iva...> 10iris --chain-id=irishub --from=<key-name> --fees=0.3iris
```
