# iriscli stake

Stake module provides a set of subcommands to query staking state and send staking transactions.

## Available Commands

| Name                                                                    | Description                                                                                   |
| ----------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| [validator](#iriscli-stake-validator)                                   | Query a validator                                                                             |
| [validators](#iriscli-stake-validators)                                 | Query for all validators                                                                      |
| [delegation](#iriscli-stake-delegation)                                 | Query a delegation based on address and validator address                                     |
| [delegations](#iriscli-stake-delegations)                               | Query all delegations made from one delegator                                                 |
| [delegations-to](#iriscli-stake-delegations-to)                         | Query all delegations to one validator                                                        |
| [unbonding-delegation](#iriscli-stake-unbonding-delegation)             | Query an unbonding-delegation record based on delegator and validator address                 |
| [unbonding-delegations](#iriscli-stake-unbonding-delegations)           | Query all unbonding-delegations records for one delegator                                     |
| [unbonding-delegations-from](#iriscli-stake-unbonding-delegations-from) | Query all unbonding delegatations from a validator                                            |
| [redelegations-from](#iriscli-stake-redelegations-from)                 | Query all outgoing redelegatations from a validator                                           |
| [redelegation](#iriscli-stake-redelegation)                             | Query a redelegation record based on delegator and a source and destination validator address |
| [redelegations](#iriscli-stake-redelegations)                           | Query all redelegations records for one delegator                                             |
| [pool](#iriscli-stake-pool)                                             | Query the current staking pool values                                                         |
| [parameters](#iriscli-stake-parameters)                                 | Query the current staking parameters information                                              |
| [signing-info](#iriscli-stake-signing-info)                             | Query a validator's signing information                                                       |
| [create-validator](#iriscli-stake-create-validator)                     | Create new validator initialized with a self-delegation to it                                 |
| [edit-validator](#iriscli-stake-edit-validator)                         | Edit existing validator account                                                               |
| [delegate](#iriscli-stake-delegate)                                     | Delegate liquid tokens to an validator                                                        |
| [unbond](#iriscli-stake-unbond)                                         | Unbond shares from a validator                                                                |
| [redelegate](#iriscli-stake-redelegate)                                 | Redelegate illiquid tokens from one validator to another                                      |
| [unjail](#iriscli-stake-unjail)                                         | Unjail validator previously jailed for downtime                                               |

## iriscli stake validator

### Query a validator by validator address

```bash
iriscli stake validator <iva...>
```

## iriscli stake validators

### Query all validators

```bash
iriscli stake validators
```

## iriscli stake delegation

Query a delegation based on delegator address and validator address.

```bash
iriscli stake delegation --address-validator=<address-validator> --address-delegator=<address-delegator>
```

**Flags:**

| Name, shorthand     | Default | Description                   | Required |
| ------------------- | ------- | ----------------------------- | -------- |
| --address-delegator |         | Bech address of the delegator | Yes      |
| --address-validator |         | Bech address of the validator | Yes      |

### Query a delegation

```bash
iriscli stake delegation --address-validator=<iva...> --address-delegator=<iaa...>
```

Example Output:

```bash
Delegation:
  Delegator:  iaa13lcwnxpyn2ea3skzmek64vvnp97jsk8qrcezvm
  Validator:  iva15grv3xg3ekxh9xrf79zd0w077krgv5xfzzunhs
  Shares:     1.0000000000000000000000000000
  Height:     26
```

## iriscli stake delegations

Query all delegations delegated from one delegator.

```bash
iriscli stake delegations <delegator-address> <flags>
```

### Query all delegations of a delegator

```bash
iriscli stake delegations <iaa...>
```

## iriscli stake delegations-to

Query all delegations to one validator.

```bash
iriscli stake delegations-to <validator-address> <flags>
```

### Query all delegations to one validator

```bash
iriscli stake delegations-to <iva...>
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

## iriscli stake unbonding-delegation

Query an unbonding-delegation record based on delegator and validator address.

```bash
iriscli stake unbonding-delegation --address-delegator=<delegator-address> --address-validator=<validator-address> <flags>
```

**Flags:**

| Name, shorthand     | Default | Description                   | Required |
| ------------------- | ------- | ----------------------------- | -------- |
| --address-delegator |         | Bech address of the delegator | Yes      |
| --address-validator |         | Bech address of the validator | Yes      |

### Query an unbonding delegation record

```bash
iriscli stake unbonding-delegation --address-delegator=<iaa...> --address-validator=<iva...>
```

## iriscli stake unbonding-delegations

### Query all unbonding delegations records of a delegator

```bash
iriscli stake unbonding-delegations <iaa...>
```

## iriscli stake unbonding-delegations-from

### Query all unbonding delegatations from a validator

```bash
iriscli stake unbonding-delegations-from <iva...>
```

## iriscli stake redelegations-from

Query all outgoing redelegations of a validator

```bash
iriscli stake redelegations-from <validator-address> <flags>
```

### Query all outgoing redelegatations of a validator

```bash
iriscli stake redelegations-from <iva...>
```

## iriscli stake redelegation

Query a redelegation record based on delegator and source validator address and destination validator address.

```bash
iriscli stake redelegation --address-validator-source=<source-validator-address> --address-validator-dest=<destination-validator-address> --address-delegator=<address-delegator> <flags>
```

**Flags:**

| Name, shorthand            | Default | Description                               | Required |
| -------------------------- | ------- | ----------------------------------------- | -------- |
| --address-delegator        |         | Bech address of the delegator             | Yes      |
| --address-validator-dest   |         | Bech address of the destination validator | Yes      |
| --address-validator-source |         | Bech address of the source validator      | Yes      |

### Query a redelegation record

```bash
iriscli stake redelegation --address-validator-source=<iva...> --address-validator-dest=<iva...> --address-delegator=<iaa...>
```

## iriscli stake redelegations

### Query all redelegations records of a delegator

```bash
iriscli stake redelegations <iaa...>
```

## iriscli stake pool

### Query the current staking pool values

```bash
iriscli stake pool
```

Example Output:

```bash
Pool:
  Loose Tokens:   1409493892.759816067399143966
  Bonded Tokens:  590526409.65743521209068061
  Token Supply:   2000020302.417251279489824576
  Bonded Ratio:   0.2952602076
```

## iriscli stake parameters

### Query the current staking parameters information

```bash
iriscli stake parameters
```

Example Output:

```bash
Stake Params:
  stake/UnbondingTime:  504h0m0s
  stake/MaxValidators:  100
```

## iriscli stake signing-info

### Query a validator's signing information

```bash
iriscli stake signing-info <iva...>
```

Example Output:

```bash
Signing Info
  Start Height:          0
  Index Offset:          3506
  Jailed Until:          1970-01-01 00:00:00 +0000 UTC
  Missed Blocks Counter: 0
```

## iriscli stake create-validator

Send a transaction to apply to be a validator and delegate a certain amount of iris to it.

```bash
iriscli stake create-validator <flags>
```

**Flags:**

| Name, shorthand   | type   | Required | Default | Description                                                                                      |
| ----------------- | ------ | -------- | ------- | ------------------------------------------------------------------------------------------------ |
| --amount          | string | Yes      |         | Amount of coins to bond                                                                          |
| --commission-rate | float  | Yes      | 0.0     | The initial commission rate percentage                                                           |
| --details         | string |          |         | Optional details                                                                                 |
| --genesis-format  | bool   |          | false   | Export the transaction in gen-tx format; it implies --generate-only                              |
| --identity        | string |          |         | Optional identity signature (ex. UPort or Keybase)                                               |
| --ip              | string |          |         | Node's public IP. It takes effect only when used in combination with                             |
| --moniker         | string | Yes      |         | Validator name                                                                                   |
| --pubkey          | string | Yes      |         | Go-Amino encoded hex PubKey of the validator. For Ed25519 the go-amino prepend hex is 1624de6220 |
| --website         | string |          |         | Optional website                                                                                 |

### Create a validator

```bash
iriscli stake create-validator --chain-id=irishub --from=<key-name> --fee=0.3iris --pubkey=<validator-pubKey> --commission-rate=0.1 --amount=100iris --moniker=<validator-name>
```

:::tip
Follow the [Mainnet](../get-started/mainnet.md#create-validator) instructions to learn more.
:::

## iriscli stake edit-validator

Edit an existing validator's settings, such as commission rate, name, etc.

```bash
iriscli stake edit-validator <flags>
```

**Flags:**

| Name, shorthand   | type   | Required | Default | Description                                        |
| ----------------- | ------ | -------- | ------- | -------------------------------------------------- |
| --commission-rate | float  |          | 0.0     | Commission rate percentage                         |
| --moniker         | string |          |         | Validator name                                     |
| --identity        | string |          |         | Optional identity signature (ex. UPort or Keybase) |
| --website         | string |          |         | Optional website                                   |
| --details         | string |          |         | Optional details                                   |

### Edit validator information

```bash
iriscli stake edit-validator --from=<key-name> --chain-id=irishub --fee=0.3iris --commission-rate=0.10 --moniker=<validator-name>
```

### Upload validator avatar

Please refer to [How to upload my validator's logo to the Explorers](../concepts/validator-faq.md#how-to-upload-my-validator-s-logo-to-the-explorers)

## iriscli stake delegate

Delegate tokens to a validator.

```bash
iriscli stake delegate --address-validator=<validator-address> <flags>
```

**Flags:**

| Name, shorthand     | type   | Required | Default | Description                   |
| ------------------- | ------ | -------- | ------- | ----------------------------- |
| --address-validator | string | Yes      |         | Bech address of the validator |
| --amount            | string | Yes      |         | Amount of coins to bond       |

```bash
iriscli stake delegate --chain-id=irishub --from=<key-name> --fee=0.3iris --amount=10iris --address-validator=<iva...>
```

## iriscli stake unbond

Unbond tokens from a validator.

```bash
iriscli stake unbond <flags>
```

**Flags:**

| Name, shorthand     | type   | Required | Default | Description                                                                                         |
| ------------------- | ------ | -------- | ------- | --------------------------------------------------------------------------------------------------- |
| --address-validator | string | Yes      |         | Bech address of the validator                                                                       |
| --shares-amount     | float  |          | 0.0     | Amount of source-shares to either unbond or redelegate as a positive integer or decimal             |
| --shares-percent    | float  |          | 0.0     | Percent of source-shares to either unbond or redelegate as a positive integer or decimal >0 and <=1 |

Users must specify the unbond amount. There two options can do this: `--shares-amount` or `--shares-percent`. Keep in mind, don't specify both of them.

### Unbond amounts of shares from a validator

```bash
iriscli stake unbond --address-validator=<iva...> --shares-amount=10 --from=<key-name> --chain-id=irishub --fee=0.3iris
```

### Unbond percentage of shares from a validator

```bash
iriscli stake unbond --address-validator=<iva...> --shares-percent=0.1 --from=<key-name> --chain-id=irishub --fee=0.3iris
```

## iriscli stake redelegate

Transfer delegation from one validator to another.

:::tip
There is no `unbonding time` during the redelegation, so you will not miss the rewards. But you can only redelegate once per validator, until a period (= `unbonding time`) exceed.
:::

```bash
iriscli stake redelegate <flags>
```

**Flags:**

| Name, shorthand            | type   | Required | Default | Description                                                                                         |
| -------------------------- | ------ | -------- | ------- | --------------------------------------------------------------------------------------------------- |
| --address-validator-dest   | string | Yes      |         | Bech address of the destination validator                                                           |
| --address-validator-source | string | Yes      |         | Bech address of the source validator                                                                |
| --shares-amount            | float  |          | 0.0     | Amount of source-shares to either unbond or redelegate as a positive integer or decimal             |
| --shares-percent           | float  |          | 0.0     | Percent of source-shares to either unbond or redelegate as a positive integer or decimal >0 and <=1 |

Users must specify the redelegation token amount. There two options can do this: `--shares-amount` or `--shares-percent`. Keep in mind, don't specify both of them.

### Redelegate amounts of shares to another validator

```bash
iriscli stake redelegate --chain-id=irishub --from=<key-name> --fee=0.3iris --address-validator-source=iva106nhdckyf996q69v3qdxwe6y7408pvyv3hgcms --address-validator-dest=iva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6sfqamll  --shares-amount=10
```

### Redelegate percentage of shares to another validator

```bash
iriscli stake redelegate --chain-id=irishub --from=<key-name> --fee=0.3iris --address-validator-source=iva106nhdckyf996q69v3qdxwe6y7408pvyv3hgcms --address-validator-dest=iva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6sfqamll  --shares-percent=0.1
```

## iriscli stake unjail

In Proof-of-Stake blockchain, validators will get block provisions by staking their token. But if they failed to keep online, they will be punished by slashing a small portion of their staked tokens. The offline validators will be removed from the validator set and put into jail, which means their voting power is zero. During the jail period, these nodes are not even validator candidates. Once the jail period ends, they can send `unjail` transactions to free themselves and become validator candidates again.

```bash
iriscli stake unjail <flags>
```

### Unjail a jailed validator

```bash
iriscli stake unjail --from=<key-name> --fee=0.3iris --chain-id=irishub
```

### Validator still jailed, cannot yet be unjailed

That means your validator is still in jail period, you can query the [signing-info](#iriscli-stake-signing-info) for the jail end time.
