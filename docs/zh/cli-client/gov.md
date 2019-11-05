# iriscli gov

该模块提供了[链上治理](../features/governance.md)的基本功能。

## 可用命令

| 名称                                            | 描述                                              |
| ----------------------------------------------- | ------------------------------------------------- |
| [query-proposal](#iriscli-gov-query-proposal)   | 查询单个提案的详细信息                            |
| [query-proposals](#iriscli-gov-query-proposals) | 按条件查询提案                                    |
| [query-vote](#iriscli-gov-query-vote)           | 查询投票                                          |
| [query-votes](#iriscli-gov-query-votes)         | 按条件查询投票                                    |
| [query-deposit](#iriscli-gov-query-deposit)     | 查询抵押详情                                      |
| [query-deposits](#iriscli-gov-query-deposits)   | 查询所有抵押                                      |
| [query-tally](#iriscli-gov-query-tally)         | 查询提案的统计信息                                |
| [submit-proposal](#iriscli-gov-submit-proposal) | 提交提案以及初始化抵押金额                        |
| [deposit](#iriscli-gov-deposit)                 | 为有效的提案抵押通证                              |
| [vote](#iriscli-gov-vote)                       | 为有效的提案投票，选项：Yes/No/NoWithVeto/Abstain |

## iriscli gov query-proposal

查询提案的详细信息。

```bash
iriscli gov query-proposal <flags>
```

**标识：**

| 名称, 速记    | 类型 | 必须 | 默认 | 描述     |
| ------------- | ---- | -------- | ---- | -------- |
| --proposal-id | uint | 是       |      | 提案的Id |

### 查询一个提议信息

```bash
iriscli gov query-proposal --chain-id=irishub --proposal-id=<proposal-id>
```

## iriscli gov query-proposals

按条件查询提案。

```bash
iriscli gov query-proposals <flags>
```

**标识：**

| 名称, 速记  | 类型    | 必须 | 默认 | 描述                                  |
| ----------- | ------- | -------- | ---- | ------------------------------------- |
| --depositor | Address |          |      | 按抵押人地址过滤提案                  |
| --limit     | uint    |          |      | 限制返回提案的个数， 默认返回所有提案 |
| --status    | string  |          |      | 按状态过滤提案 (passed / rejected)    |
| --voter     | Address |          |      | 按投票人地址过滤提案                  |

### 查询所有提案

```bash
iriscli gov query-proposals --chain-id=irishub
```

### 按条件查询提案

```bash
iriscli gov query-proposals --chain-id=irishub --limit=3 --status=passed --depositor=<iaa...>
```

## iriscli gov query-vote

查询投票信息。

```bash
iriscli gov query-vote <flags>
```

**标识：**

| 名称, 速记    | 类型    | 必须 | 默认 | 描述       |
| ------------- | ------- | -------- | ---- | ---------- |
| --proposal-id | uint    | 是       |      | 提案的Id   |
| --voter       | Address | 是       |      | 投票人地址 |

### 查询一个投票信息

```bash
iriscli gov query-vote --chain-id=irishub --proposal-id=<proposal-id> --voter=<iaa...>
```

## iriscli gov query-votes

查询提案的所有投票信息。

```bash
iriscli gov query-votes <flags>
```

**标识：**

| 名称, 速记    | 类型 | 必须 | 默认 | 描述     |
| ------------- | ---- | -------- | ---- | -------- |
| --proposal-id | uint | 是       |      | 提案的Id |

### 查询指定提案的所有投票信息

```bash
iriscli gov query-votes --chain-id=irishub --proposal-id=<proposal-id>
```

## iriscli gov query-deposit

查询指定提案的抵押信息。

```bash
iriscli gov query-deposit <flags>
```

**标识：**

| 名称, 速记    | 类型    | 必须 | 默认 | 描述       |
| ------------- | ------- | -------- | ---- | ---------- |
| --proposal-id | uint    | 是       |      | 提案的Id   |
| --depositor   | Address | 是       |      | 抵押人地址 |

### 查询指定人的抵押信息

```bash
iriscli gov query-deposit --chain-id=irishub --proposal-id=<proposal-id> --depositor=<iaa...>
```

## iriscli gov query-deposits

查询指定提案的所有抵押信息。

```bash
iriscli gov query-deposits <flags>
```

**标识：**

| 名称, 速记    | 类型 | 必须 | 默认 | 描述     |
| ------------- | ---- | -------- | ---- | -------- |
| --proposal-id | uint | 是       |      | 提案的Id |

### 查询指定提案的所有抵押信息

```bash
iriscli gov query-deposits --chain-id=irishub --proposal-id=<proposal-id>
```

## iriscli gov query-tally

查询提案投票的统计信息。

```bash
iriscli gov query-tally <flags>
```

**标识：**

| 名称, 速记    | 类型 | 必须 | 默认 | 描述     |
| ------------- | ---- | -------- | ---- | -------- |
| --proposal-id |  uint    | 是       |      | 提案的Id |

### 查询提案投票的统计信息

```bash
iriscli gov query-tally --chain-id=irishub --proposal-id=<proposal-id>
```

## iriscli gov submit-proposal

提交提案以及初始化抵押金额。

```bash
iriscli gov submit-proposal <flags>
```

**标识：**

| 名称, 速记               | 类型   | 必须 | 默认  | 描述                                                                                           |
| ------------------------ | ------ | -------- | ----- | ---------------------------------------------------------------------------------------------- |
| --deposit                | Coin   | 是       |       | 初始抵押金额(至少最小抵押金额的30% of)                                                         |
| --description            | string | 是       |       | 提案的描述信息                                                                                 |
| --param                  | string |          |       | 提案修改的参数，例如`mint/Inflation=0.050`                                                     |
| --title                  | string | 是       |       | 提案的标题                                                                                     |
| --type                   | string | 是       |       | 提案的类型（PlainText/Parameter/SoftwareUpgrade/SoftwareHalt/CommunityTaxUsage/TokenAddition） |
| --version                | uint   |          | 0     | 升级的版本                                                                                     |
| --software               | string |          |       | 新软件的地址                                                                                   |
| --switch-height          | uint   |          | 0     | 软件升级过程中，切换的区块高度                                                                 |
| --threshold              | string |          | "0.8" | 软件升级的升级信号阈值                                                                         |
| --token-canonical-symbol | string |          |       | 外部通证的源符号                                                                               |
| --token-symbol           | string |          |       | 通证符号。 创建后，将无法修改                                                                  |
| --token-name             | string |          |       | 通证名称                                                                                       |
| --token-decimal          | uint   |          |       | 通证的最大精度，最大值为18                                                                     |
| --token-min-unit-alias   | string |          |       | 通证最小单位别名                                                                               |
| --token-initial-supply   | uint64 |          |       | 通证初始总量                                                                                   |

:::tip
提案人必须至少抵押[MinDeposit](../features/governance.md#提议级别)的30％才能提交提案。
:::

### 提交参数修改提案

:::tip
[有哪些参数可以在线修改？](../concepts/gov-params.md)
:::

**唯一必须参数：** `--param`

```bash
iriscli gov submit-proposal --chain-id=irishub --title=<proposal-title> --description=<proposal-description> --from=<key-name> --fee=0.3iris --deposit=2000iris --type=Parameter --param='mint/Inflation=0.050'
```

### 提交软件升级提案

**必须参数：** `--software`, `--version`, `--switch-height`, `--threshold`

```bash
iriscli gov submit-proposal --chain-id=irishub --title=<proposal-title> --description=<proposal-description> --from=<key-name> --fee=0.3iris --deposit=2000iris --type=SoftwareUpgrade --software=https://github.com/irisnet/irishub/tree/v0.15.1 --version=2 --switch-height=8000 --threshold=0.8
```

### 提交通证添加提案

**所需参数：**

- 必选: `--token-symbol`, `--token-canonical-symbol`, `--token-name`
- 可选: `--token-decimal`, `--token-min-unit-alias`

```bash
iriscli gov submit-proposal --chain-id=irishub --title=<proposal-title> --description=<proposal-description> --from=<key-name> --fee=1iris --deposit=2000iris --type=TokenAddition --token-symbol=btc --token-canonical-symbol=btc --token-name=Bitcoin --token-decimal=18 --token-min-unit-alias=satoshi
```

## iriscli gov deposit

为有效的提案抵押通证。

```bash
iriscli gov deposit <flags>

```

**标识：**

| 名称, 速记    | 类型 | 必须 | 默认 | 描述           |
| ------------- | ---- | -------- | ---- | -------------- |
| --deposit     | Coin | 是       |      | 抵押的通证金额 |
| --proposal-id | uint | 是       |      | 提案Id         |

### 为有效的提案抵押通证

当总抵押金额超过[MinDeposit](../features/governance.md#提议级别)时，提案将进入投票阶段。

```bash
iriscli gov deposit --chain-id=irishub --proposal-id=<proposal-id> --deposit=50iris --from=<key-name> --fee=0.3iris
```

## iriscli gov vote

为有效的提案投票，选项：Yes/No/NoWithVeto/Abstain。

:::tip
[No VS NoWithVeto](../features/governance.md#销毁机制)

在投票期内，只有验证人和委托人可以对提案进行投票。
:::

```bash
iriscli gov vote <flags>
```

**标识：**

| 名称, 速记    | 类型   | 必须 | 默认 | 描述                            |
| ------------- | ------ | -------- | ---- | ------------------------------- |
| --option      | string | 是       |      | 选项：Yes/No/NoWithVeto/Abstain |
| --proposal-id | uint   | 是       |      | 提案Id                          |
