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

通过委托人地址和验证人地址查询委托交易。

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

### 查询委托人的所有未绑定委托记录

```bash
iriscli stake unbonding-delegations <iaa...>
```

## iriscli stake unbonding-delegations-from

### 查询验证人的所有未绑定委托记录

```bash
iriscli stake unbonding-delegations-from <iva...>
```

## iriscli stake redelegations-from

查询验证人的所有转委托记录。

```bash
iriscli stake redelegations-from <validator-address> <flags>
```

### 查询验证人的所有转委托记录

```bash
iriscli stake redelegations-from <iva...>
```

## iriscli stake redelegation

通过委托人地址、原验证人地址、目标验证人地址查询转委托记录。

```bash
iriscli stake redelegation --address-validator-source=<source-validator-address> --address-validator-dest=<destination-validator-address> --address-delegator=<address-delegator> <flags>
```

**标志：**

| 名称，速记                 | 默认 | 描述               | 必须 |
| -------------------------- | ---- | ------------------ | ---- |
| --address-delegator        |      | 委托者bech地址     | 是   |
| --address-validator-dest   |      | 目标验证者bech地址 | 是   |
| --address-validator-source |      | 源验证者bech地址   | 是   |

### 查询转委托记录

```bash
iriscli stake redelegation --address-validator-source=<iva...> --address-validator-dest=<iva...> --address-delegator=<iaa...>
```

## iriscli stake redelegations

### 查询委托人的所有转委托记录

```bash
iriscli stake redelegations <iaa...>
```

## iriscli stake pool

### 查询当前权益池

```bash
iriscli stake pool
```

示例输出:

```bash
Pool:
  Loose Tokens:   1409493892.759816067399143966
  Bonded Tokens:  590526409.65743521209068061
  Token Supply:   2000020302.417251279489824576
  Bonded Ratio:   0.2952602076
```

## iriscli stake parameters

### 查询当前权益参数信息

```bash
iriscli stake parameters
```

示例输出:

```bash
Stake Params:
  stake/UnbondingTime:  504h0m0s
  stake/MaxValidators:  100
```

## iriscli stake signing-info

### 查询验证人签名信息

```bash
iriscli stake signing-info <iva...>
```

示例输出:

```bash
Signing Info
  Start Height:          0
  Index Offset:          3506
  Jailed Until:          1970-01-01 00:00:00 +0000 UTC
  Missed Blocks Counter: 0
```

## iriscli stake create-validator

发送交易申请成为验证人，并委托一定数量的iris到该验证人。

```bash
iriscli stake create-validator <flags>
```

**标志：**

| 名称，速记        | 类型   | 必须 | 默认  | 描述                                |
| ----------------- | ------ | ---- | ----- | ----------------------------------- |
| --amount          | string | 是   |       | 委托金额                            |
| --commission-rate | float  | 是   | 0.0   | 初始佣金比例                        |
| --details         | string |      |       | 验证人节点的详细信息                |
| --genesis-format  | bool   |      | false | 是否以genesis transaction的方式导出 |
| --identity        | string |      |       | 身份信息的签名                      |
| --ip              | string |      |       | 验证人节点的IP                      |
| --moniker         | string | 是   |       | 验证人节点名称                      |
| --pubkey          | string | 是   |       | Amino编码的验证人公钥               |
| --website         | string |      |       | 验证人节点的网址                    |

### 创建验证人

```bash
iriscli stake create-validator --chain-id=irishub --from=<key-name> --fee=0.3iris --pubkey=<validator-pubKey> --commission-rate=0.1 --amount=100iris --moniker=<validator-name>
```

:::tip
查看 [主网](../get-started/mainnet.md#升级为验证人节点) 说明以了解更多。
:::

## iriscli stake edit-validator

修改验证的的参数，包括佣金比率，验证人节点名称以及其他描述信息。

```bash
iriscli stake edit-validator <flags>
```

**标志：**

| 名称，速记        | 类型   | 必须 | 默认 | 描述               |
| ----------------- | ------ | ---- | ---- | ------------------ |
| --commission-rate | float  |      | 0.0  | 佣金比率           |
| --moniker         | string |      |      | 验证人名称         |
| --identity        | string |      |      | 身份签名           |
| --website         | string |      |      | 网址               |
| --details         | string |      |      | 验证人节点详细信息 |

### 编辑验证人信息

```bash
iriscli stake edit-validator --from=<key-name> --chain-id=irishub --fee=0.3iris --commission-rate=0.10 --moniker=<validator-name>
```

### 上传验证人头像

请参考 [如何将验证人的Logo上传到区块浏览器](../concepts/validator-faq.md#如何将验证人的logo上传到区块浏览器)。

## iriscli stake delegate

向验证人委托通证。

```bash
iriscli stake delegate --address-validator=<validator-address> <flags>
```

**标志：**

| 名称，速记          | 类型   | 必须 | 默认 | 描述       |
| ------------------- | ------ | ---- | ---- | ---------- |
| --address-validator | string | 是   |      | 验证人地址 |
| --amount            | string | 是   |      | 委托金额   |

```bash
iriscli stake delegate --chain-id=irishub --from=<key-name> --fee=0.3iris --amount=10iris --address-validator=<iva...>
```

## iriscli stake unbond

从验证人解委托通证。

```bash
iriscli stake unbond <flags>
```

**标志：**

| 名称，速记          | 类型   | 必须 | 默认 | 描述                    |
| ------------------- | ------ | ---- | ---- | ----------------------- |
| --address-validator | string | 是   |      | 验证人地址              |
| --shares-amount     | float  |      | 0.0  | 解绑的shares数量，正数  |
| --shares-percent    | float  |      | 0.0  | 解绑的比率，小于1的正数 |

用户必须指定解绑shares的数量，可使用`--shares-amount`或者`--shares-percent`指定。请勿同时使用这两个参数。

### 从验证人中解委托一定数量的shares

```bash
iriscli stake unbond --address-validator=<iva...> --shares-amount=10 --from=<key-name> --chain-id=irishub --fee=0.3iris
```

### 从验证人中解委托一定比例的shares

```bash
iriscli stake unbond --address-validator=<iva...> --shares-percent=0.1 --from=<key-name> --chain-id=irishub --fee=0.3iris
```

## iriscli stake redelegate

把某个委托的一部分或者全部从一个验证人转移到另外一个验证人。

:::tip
转委托没有`unbonding time`，所以你不会错过奖励。但是对每个验证人的转委托，在一个`unbonding time`区间内只能操作一次。
:::

```bash
iriscli stake redelegate <flags>
```

**标志：**

| 名称，速记                 | 类型   | 必须 | 默认 | 描述                          |
| -------------------------- | ------ | ---- | ---- | ----------------------------- |
| --address-validator-dest   | string | 是   |      | 目标验证人地址                |
| --address-validator-source | string | 是   |      | 源验证人地址                  |
| --shares-amount            | float  |      | 0.0  | 转移的shares数量，正数        |
| --shares-percent           | float  |      | 0.0  | 转移的shares比例，小于1的正数 |

用户必须指定解绑shares的数量，可使用`--shares-amount`或者`--shares-percent`指定。请勿同时使用这两个参数。

### 转委托一定数量shares到其他验证人

```bash
iriscli stake redelegate --chain-id=irishub --from=<key-name> --fee=0.3iris --address-validator-source=iva106nhdckyf996q69v3qdxwe6y7408pvyv3hgcms --address-validator-dest=iva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6sfqamll  --shares-amount=10
```

### 转委托一定比例shares到其他验证人

```bash
iriscli stake redelegate --chain-id=irishub --from=<key-name> --fee=0.3iris --address-validator-source=iva106nhdckyf996q69v3qdxwe6y7408pvyv3hgcms --address-validator-dest=iva1xpqw0kq0ktt3we5gq43vjphh7xcjfy6sfqamll  --shares-percent=0.1
```

## iriscli stake unjail

在PoS网络中，验证人的收益主要来自于staking抵押获利，但是若验证人不能保持在线，就会被当作一种作恶行为。系统会剥夺其作为验证人参与共识的资格。这样一来，其的状态会变成jailed，且投票权将立刻变为零。这种状态将持续一段时间。当jailed期结束，验证人节点的operator需要执行unjail操作来让节点的状态变为unjailed，再次成为验证人或者候选验证人。

```bash
iriscli stake unjail <flags>
```

### 解禁被监禁的验证人

```bash
iriscli stake unjail --from=<key-name> --fee=0.3iris --chain-id=irishub
```

### Validator still jailed, cannot yet be unjailed

如果执行解禁操作的tx报错 `Validator still jailed, cannot yet be unjailed`，意味着该验证人节点还在监禁期内，不能被解禁。您可以查询 [signing-info](#iriscli-stake-signing-info) 获取监禁结束时间。
