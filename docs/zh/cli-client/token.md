# iris token

token模块用于管理你在IRIS Hub上发行的资产。

## 可用命令

| 名称                                            | 描述                     |
| ----------------------------------- | ------------------------ |
| [issue](#iris-tx-token-issue)       | 发行通证                 |
| [edit](#iris-tx-token-edit)         | 编辑通证                 |
| [transfer](#iris-tx-token-transfer) | 转让通证所有权           |
| [mint](#iris-tx-token-mint)         | 增发通证到指定账户       |
| [token](#iris-query-token-token)    | 查询通证                 |
| [tokens](#iris-query-token-tokens)  | 查询指定所有者的通证集合 |
| [fee](#iris-query-token-fee)        | 查询通证相关费用         |
| [params](#iris-query-token-params)  | 查询通证相关参数         |

## iris tx token issue

发行一个新通证。

```bash
iris tx token issue [flags]
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
iris tx token issue --name="Kitty Token" --symbol="kitty" --min-unit="kitty" --scale=0 --initial-supply=100000000000 --max-supply=1000000000000 --mintable=true --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

### 发送通证

您可以像[发送iris](./bank.md#iris-tx-bank-send)一样发送任何通证。

#### 发送通证

```bash
iris tx bank send [from_key_or_address] [to_address] [amount] [flags]
```

## iris tx token edit

编辑通证。

```bash
iris tx token edit [symbol] [flags]
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
iris tx token edit <symbol> --name="Cat Token" --max-supply=100000000000 --mintable=true --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

## iris tx token transfer

转让通证所有权。

```bash
iris tx token transfer [symbol] [flags]
```

**标识：**

| 名称，速记 | 类型   | 必须 | 默认 | 描述       |
| ---------- | ------ | ---- | ---- | ---------- |
| --to       | string | 是   |      | 接收人地址 |

### 转让通证所有者

```bash
iris tx token transfer <symbol> --to=<to> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

## iris tx token mint

增发通证到指定地址。

```bash
iris tx token mint [symbol] [flags]
```

**标识：**

| 名称，速记 | 类型   | 必须 | 默认 | 描述                                       |
| ---------- | ------ | ---- | ---- | ------------------------------------------ |
| --to       | string |      |      | 增发的通证的接收地址，默认为发起该交易地址 |
| --amount   | uint64 | 是   | 0    | 增发的数量                                 |

### 增发通证

```bash
iris tx token mint <symbol> --amount=<amount> --to=<to> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

## iris query token token

查询通证。

```bash
iris query token token [denom] [flags]
```

### 查询通证

```bash
iris query token token <denom>
```

## iris query token tokens

查询指定所有者的通证集合。所有者是可选的。

```bash
iris query token tokens [owner] [flags]
```

### 查询所有通证

```bash
iris query token tokens
```

### 查询指定所有者的通证

```bash
iris query token tokens <owner>
```

## iris query token fee

查询与通证相关的费用，包括通证发行和增发。

```bash
iris query token fee [symbol] [flags]
```

### 查询发行和增发通证的费用

```bash
iris query token fee kitty
```

## iris query token fee

查询与通证相关的参数，包括通证发行和增发。

```bash
iris query token fee [symbol] [flags]
```

### 查询发行和增发通证的参数

```bash
iris query token params kitty
```

