# iriscli stake

Stake模块提供了一系列查询staking状态和发送staking交易的命令。

## Available Commands

| 名称                                                                    | 描述                                                         |
| ----------------------------------------------------------------------- | ------------------------------------------------------------ |
| [validator](#iriscli-stake-validator)                                   | 查询某个验证者                                               |
| [validators](#iriscli-stake-validators)                                 | 查询所有的验证者                                             |
| [delegation](#iriscli-stake-delegation)                                 | 基于委托者地址和验证者地址的委托查询                         |
| [delegations](#iriscli-stake-delegations)                               | 基于委托者地址的所有委托查询                                 |
| [delegations-to](#iriscli-stake-delegations-to)                         | 查询在某个验证人上的所有委托                                 |
| [unbonding-delegation](#iriscli-stake-unbonding-delegation)             | 基于委托者地址和验证者地址的unbonding-delegation记录查询     |
| [unbonding-delegations](#iriscli-stake-unbonding-delegations)           | 基于委托者地址的所有unbonding-delegation记录查询             |
| [unbonding-delegations-from](#iriscli-stake-unbonding-delegations-from) | 基于验证者地址的所有unbonding-delegation记录查询             |
| [redelegations-from](#iriscli-stake-redelegations-from)                 | 基于某一验证者的所有转委托查询                               |
| [redelegation](#iriscli-stake-redelegation)                             | 基于委托者地址，原验证者地址和目标验证者地址的转委托记录查询 |
| [redelegations](#iriscli-stake-redelegations)                           | 基于委托者地址的所有转委托记录查询                           |
| [pool](#iriscli-stake-pool)                                             | 查询最新的权益池                                             |
| [parameters](#iriscli-stake-parameters)                                 | 查询最新的权益参数信息                                       |
| [signing-info](#iriscli-stake-signing-info)                             | 查询验证者签名信息                                           |
| [create-validator](#iriscli-stake-create-validator)                     | 以自委托的方式创建一个新的验证者                             |
| [edit-validator](#iriscli-stake-edit-validator)                         | 编辑已存在的验证者信息                                       |
| [delegate](#iriscli-stake-delegate)                                     | 委托一定代币到某个验证者                                     |
| [unbond](#iriscli-stake-unbond)                                         | 从指定的验证者解绑一定的股份                                 |
| [redelegate](#iriscli-stake-redelegate)                                 | 转委托一定的token从一个验证者到另一个验证者                  |
| [unjail](#iriscli-stake-unjail)                                         | 恢复之前由于宕机被惩罚的验证者的身份                         |

## iriscli stake validator

### 通过地址查询验证人

```bash
iriscli stake validator <iva...>
```

## iriscli stake validators

### 查询所有验证人

```bash
iriscli stake validators
```

## iriscli stake delegation

通过委托人和验证人地址查询委托交易。

```bash
iriscli stake delegation --address-validator=<address-validator> --address-delegator=<address-delegator>
```

**标志：**

| 名称，速记          | 默认 | 描述           | 必须 |
| ------------------- | ---- | -------------- | ---- |
| --address-delegator |      | 委托人bech地址 | 是   |
| --address-validator |      | 验证人bech地址 | 是   |

### 查询委托交易

```bash
iriscli stake delegation --address-validator=<iva...> --address-delegator=<iaa...>
```

示例输出:

```bash
Delegation:
  Delegator:  iaa13lcwnxpyn2ea3skzmek64vvnp97jsk8qrcezvm
  Validator:  iva15grv3xg3ekxh9xrf79zd0w077krgv5xfzzunhs
  Shares:     1.0000000000000000000000000000
  Height:     26
```

## iriscli stake delegations

查询某个委托人发起的所有委托记录。

```bash
iriscli stake delegations <delegator-address> <flags>
```

### 查询某个委托人发起的所有委托记录

```bash
iriscli stake delegations <iaa...>
```

## iriscli stake delegations-to

查询某个验证人接受的所有委托。

```bash
iriscli stake delegations-to <validator-address> <flags>
```

### 查询某个验证人接受的所有委托

```bash
iriscli stake delegations-to <iva...>
```

示例输出:

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

通过委托人与验证人地址查询unbonding-delegation记录。

```bash
iriscli stake unbonding-delegation --address-delegator=<delegator-address> --address-validator=<validator-address> <flags>
```

**标志：**

| 名称，速记          | 默认 | 描述           | 必须 |
| ------------------- | ---- | -------------- | ---- |
| --address-delegator |      | 委托人bech地址 | 是   |
| --address-validator |      | 验证人bech地址 | 是   |

### 查询unbonding-delegation记录

```bash
iriscli stake unbonding-delegation --address-delegator=<iaa...> --address-validator=<iva...>
```

## iriscli stake unbonding-delegations

### 通过委托人与验证人地址查询所有unbonding-delegation记录

```bash
iriscli stake unbonding-delegations <iaa...>
```

## iriscli stake unbonding-delegations-from

### 查询所有unbonding-delegation-from记录

```bash
iriscli stake unbonding-delegations-from <iva...>
```

## iriscli stake redelegations-from

### Query all outgoing redelegations of a validator

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

**标志：**

| 名称，速记                 | 默认 | 描述                                      | 必须 |
| -------------------------- | ---- | ----------------------------------------- | ---- |
| --address-delegator        |      | Bech address of the delegator             | 是   |
| --address-validator-dest   |      | Bech address of the destination validator | 是   |
| --address-validator-source |      | Bech address of the source validator      | 是   |

### Query a redelegation record

```bash
iriscli stake redelegation --address-validator-source=<iva...> --address-validator-dest=<iva...> --address-delegator=<iaa...>
```

## iriscli stake redelegations

Query all redelegations records of a delegator.

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

**标志：**

| 名称，速记        | 类型   | 必须 | 默认  | 描述                                                                                             |
| ----------------- | ------ | ---- | ----- | ------------------------------------------------------------------------------------------------ |
| --amount          | string | 是   |       | Amount of coins to bond                                                                          |
| --commission-rate | float  | 是   | 0.0   | The initial commission rate percentage                                                           |
| --details         | string |      |       | Optional details                                                                                 |
| --genesis-format  | bool   |      | false | Export the transaction in gen-tx format; it implies --generate-only                              |
| --identity        | string |      |       | Optional identity signature (ex. UPort or Keybase)                                               |
| --ip              | string |      |       | Node's public IP. It takes effect only when used in combination with                             |
| --moniker         | string | 是   |       | Validator name                                                                                   |
| --pubkey          | string | 是   |       | Go-Amino encoded hex PubKey of the validator. For Ed25519 the go-amino prepend hex is 1624de6220 |
| --website         | string |      |       | Optional website                                                                                 |

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

**标志：**

| 名称，速记        | 类型   | 必须 | 默认 | 描述                                               |
| ----------------- | ------ | ---- | ---- | -------------------------------------------------- |
| --commission-rate | float  |      | 0.0  | Commission rate percentage                         |
| --moniker         | string |      |      | Validator name                                     |
| --identity        | string |      |      | Optional identity signature (ex. UPort or Keybase) |
| --website         | string |      |      | Optional website                                   |
| --details         | string |      |      | Optional details                                   |

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

**标志：**

| 名称，速记          | 类型   | 必须 | 默认 | 描述                          |
| ------------------- | ------ | ---- | ---- | ----------------------------- |
| --address-validator | string | 是   |      | Bech address of the validator |
| --amount            | string | 是   |      | Amount of coins to bond       |

```bash
iriscli stake delegate --chain-id=irishub --from=<key-name> --fee=0.3iris --amount=10iris --address-validator=<iva...>
```

## iriscli stake unbond

Unbond tokens from a validator.

```bash
iriscli stake unbond <flags>
```

**标志：**

| 名称，速记          | 类型   | 必须 | 默认 | 描述                                                                                                |
| ------------------- | ------ | ---- | ---- | --------------------------------------------------------------------------------------------------- |
| --address-validator | string | 是   |      | Bech address of the validator                                                                       |
| --shares-amount     | float  |      | 0.0  | Amount of source-shares to either unbond or redelegate as a positive integer or decimal             |
| --shares-percent    | float  |      | 0.0  | Percent of source-shares to either unbond or redelegate as a positive integer or decimal >0 and <=1 |

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

**标志：**

| 名称，速记                 | 类型   | 必须 | 默认 | 描述                                                                                                |
| -------------------------- | ------ | ---- | ---- | --------------------------------------------------------------------------------------------------- |
| --address-validator-dest   | string | 是   |      | Bech address of the destination validator                                                           |
| --address-validator-source | string | 是   |      | Bech address of the source validator                                                                |
| --shares-amount            | float  |      | 0.0  | Amount of source-shares to either unbond or redelegate as a positive integer or decimal             |
| --shares-percent           | float  |      | 0.0  | Percent of source-shares to either unbond or redelegate as a positive integer or decimal >0 and <=1 |

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

That means your validator is still in jail period, you can query the [signing-info](#iriscli-stake-signing-info) for the jail end time:
