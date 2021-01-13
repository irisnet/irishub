# Bank

Bank模块用于管理本地帐户中的资产。

## 可用命令

| 名称                                  | 描述                         |
| ------------------------------------- | ---------------------------- |
| [balances](#iris-query-bank-balances) | 查询账户余额                 |
| [total](#iris-query-bank-total)       | 查询链上币的总数             |
| [send](#iris-tx-bank-send)            | 创建、签名、广播一个转账交易 |

## iris query bank balances

该命令用于查询指定账户的余额信息。

```bash
iris query bank balances [address] [flags]
```

**标识：**

| 名称，速记 | 类型   | 必须 | 默认 | 描述                       |
| ---------- | ------ | ---- | ---- | -------------------------- |
| -h, --help |        |      |      | `balances`子命令的提示信息 |
| --denom    | string |      |      | 要查询的指定余额面值       |
| --count-total   |        |          |         | 所有查询的余额记录总数 |

## iris query bank total

该命令用于查询链上币的总数。

```bash
iris query bank total [flags]
```

**标识：**

| 名称，速记 | 类型   | 必须 | 默认 | 描述                       |
| ---------- | ------ | ---- | ---- | -------------------------- |
| -h, --help |        |      |      | `balances`子命令的提示信息 |
| --denom    | string |      |      | 要查询的指定余额面值       |

## iris tx bank send

发送代币到另一个地址，此命令包括 `generate`，`sign` 和 `broadcast` 这些步骤。

```bash
iris tx bank send [from_key_or_address] [to_address] [amount] [flags]
```

**标识：**

| 名称，速记 | 类型 | 必须 | 默认 | 描述                       |
| ---------- | ---- | ---- | ---- | -------------------------- |
| -h, --help |      |      |      | `balances`子命令的提示信息 |
