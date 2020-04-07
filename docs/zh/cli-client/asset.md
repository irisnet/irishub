# iriscli asset

Asset模块用于管理你在IRIS Hub上发行的资产。

## 可用命令

| 名称                                            | 描述                     |
| ----------------------------------------------- | ------------------------ |
| [token issue](#iriscli-asset-token-issue)       | 发行通证                 |
| [token edit-token](#iriscli-asset-token-edit)   | 编辑通证                 |
| [token transfer](#iriscli-asset-token-transfer) | 转让通证所有权           |
| [token mint](#iriscli-asset-token-mint)         | 增发通证到指定账户       |
| [token token](#iriscli-asset-token-token)       | 查询通证                 |
| [token tokens](#iriscli-asset-token-tokens)     | 查询指定所有者的通证集合 |
| [token fee](#iriscli-asset-token-fee)           | 查询通证相关费用         |

## iriscli asset token issue

发行一个新通证。

```bash
iriscli asset token issue [flags]
```

**标识：**

| 名称，速记       | 类型    | 必须 | 默认          | 描述                                                               |
| ---------------- | ------- | ---- | ------------- | ------------------------------------------------------------------ |
| --name           | string  | 是   |               | 通证的名称，限制为32个unicode字符，例如"IRIS Network"              |
| --symbol         | string  | 是   |               | 通证的符号，长度在3到8之间，字母数字字符，以字符开始，不区分大小写 |
| --initial-supply | uint64  | 是   |               | 此通证的初始供应。 增发前的数量不应超过1000亿。                    |
| --max-supply     | uint64  |      | 1000000000000 | 通证上限，总供应不能超过最大供应。 增发前的数量不应超过1万亿       |
| --scale          | uint8   | 是   |               | 通证最多可以有18位小数                                             |
| --min-unit       | string  |      |               | 最小单位别名                                                       |
| --mintable       | boolean |      | false         | 首次发行后是否可以增发此通证                                       |

### 发行通证

```bash
iriscli asset token issue --symbol="kitty" --name="Kitty Token" --initial-supply=100000000000 --max-supply=1000000000000 --scale=0 --mintable=true --fee=1iris --chain-id=irishub --from=<key-name> --commit
```

### 发送通证

您可以像[发送iris](./bank.md#iriscli-bank-send)一样发送任何通证。

#### 发送通证

```bash
iriscli bank send --from=<key-name> --to=<address> --amount=10kitty --fee=0.3iris --chain-id=irishub --commit
```

## iriscli asset token edit

编辑通证。

```bash
iriscli asset token edit [symbol] [flags]
```

**标识：**

| 名称，速记   | 类型   | 必须 | 默认  | 描述                          |
| ------------ | ------ | ---- | ----- | ----------------------------- |
| --name       | string |      |       | 通证名称，例如：IRIS Network  |
| --max-supply | uint   |      | 0     | 通证的最大供应量              |
| --mintable   | bool   |      | false | 通证是否可以增发，默认为false |

`max-supply` 不得少于当前的总供应量。

### 编辑通证

```bash
iriscli asset token edit kitty --name="Cat Token" --max-supply=100000000000 --mintable=true --from=<key-name> --chain-id=irishub --fee=0.3iris --commit
```

## iriscli asset token transfer

转让通证所有权。

```bash
iriscli asset token transfer [symbol] [flags]
```

**标识：**

| 名称，速记 | 类型   | 必须 | 默认 | 描述       |
| ---------- | ------ | ---- | ---- | ---------- |
| --to       | string | 是   |      | 接收人地址 |

### 转让通证所有者

```bash
iriscli asset token transfer kitty --to=<new-owner-address> --from=<key-name> --chain-id=irishub --fee=0.3iris --commit
```

## iriscli asset token mint

增发通证到指定地址。

```bash
iriscli asset token mint [symbol] [flags]
```

**标识：**

| 名称，速记 | 类型   | 必须 | 默认 | 描述                                       |
| ---------- | ------ | ---- | ---- | ------------------------------------------ |
| --to       | string |      |      | 增发的通证的接收地址，默认为发起该交易地址 |
| --amount   | uint64 | 是   | 0    | 增发的数量                                 |

### 增发通证

```bash
iriscli asset token mint kitty --amount=1000000 --from=<key-name> --chain-id=irishub --fee=0.3iris
```

## iriscli asset token token

查询通证。

```bash
iriscli asset token token [symbol] [flags]
```

### 查询通证

```bash
iriscli asset token token kitty
```

## iriscli asset token tokens

查询指定所有者的通证集合。所有者是可选的。

```bash
iriscli asset token tokens [owner] [flags]
```

### 查询所有通证

```bash
iriscli asset token tokens
```

### 查询指定所有者的通证

```bash
iriscli asset token tokens <owner>
```

## iriscli asset token fee

查询与通证相关的费用，包括通证发行和增发。

```bash
iriscli asset token fee [symbol] [flags]
```

### 查询发行和增发通证的费用

```bash
iriscli asset token fee kitty
```
