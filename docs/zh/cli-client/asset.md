# iriscli asset

Asset模块用于管理你在IRIS Hub上发行的资产。

## 可用命令

| 名称                                                            | 描述                       |
| --------------------------------------------------------------- | -------------------------- |
| [create-gateway](#iriscli-asset-create-gateway)                 | 创建网关                   |
| [edit-gateway](#iriscli-asset-edit-gateway)                     | 编辑网关                   |
| [transfer-gateway-owner](#iriscli-asset-transfer-gateway-owner) | 转让网关所有权             |
| [issue-token](#iriscli-asset-issue-token)                       | 发行通证                   |
| [edit-token](#iriscli-asset-edit-token)                         | 编辑通证                   |
| [transfer-token-owner](#iriscli-asset-transfer-token-owner)     | 转让通证所有权             |
| [mint-token](#iriscli-asset-mint-token)                         | 增发通证到指定账户         |
| [query-token](#iriscli-asset-query-token)                       | 查询指定通证信息           |
| [query-tokens](#iriscli-asset-query-tokens)                     | 查询符合条件的一组通证信息 |
| [query-gateway](#iriscli-asset-query-gateway)                   | 根据别名查询指定网关信息   |
| [query-gateways](#iriscli-asset-query-gateways)                 | 根据所有者查询网关信息     |
| [query-fee](#iriscli-asset-query-fee)                           | 查询资产相关费率           |

## iriscli asset create-gateway

创建网关用于管理外部资产。

```bash
iriscli asset create-gateway <flags>
```

**标识：**

| 名称, 速记 | 类型   | 必须 | 默认 | 描述                                                   |
| ---------- | ------ | -------- | ---- | ------------------------------------------------------ |
| --moniker  | string | 是       |      | 唯一名称，大小在3到8之间，以字母开头，后跟字母数字字符 |
| --identity | string |          |      | 可选，身份签名，最大长度为128（例如UPort或Keybase）    |
| --details  | string |          |      | 可选，详细信息，最大长度为280                          |
| --website  | string |          |      | 可选，网站，最大长度为128                              |

### 创建网关

```bash
iriscli asset create-gateway --moniker=cats --identity=<pgp-id> --details="Cat Tokens" --website="www.example.com" --from=<key-name> --chain-id=irishub --fee=0.3iris
```

## iriscli asset edit-gateway

根据指定别名编辑网关。

```bash
iriscli asset edit-gateway <flags>
```

**标识：**

| 名称, 速记 | 类型   | 必须 | 默认 | 描述                                                 |
| ---------- | ------ | -------- | ---- | ---------------------------------------------------- |
| --moniker  | string | 是       |      | 唯一，别名，大小在3到8之间，以字母开头，后跟字母数字 |
| --identity | string |          |      | 可选，身份签名，最大长度为128                        |
| --details  | string |          |      | 可选，详细信息，最大长度为280                        |
| --website  | string |          |      | 可选，网站，最大长度为128                            |

### 编辑网关

```bash
iriscli asset edit-gateway --moniker=cats --identity=<pgp-id> --details="Cat Tokens" --website="http://www.example.com" --from=<key-name> --chain-id=irishub --fee=0.3iris --commit
```

## iriscli asset transfer-gateway-owner

将网关的所有权转让给新所有者。

```bash
iriscli asset transfer-gateway-owner <flags>
```

**标识：**

| 名称, 速记 | 类型    | 必须 | 默认 | 描述                 |
| ---------- | ------- | -------- | ---- | -------------------- |
| --moniker  | string  | 是       |      | 网关的唯一别名 |
| --to       | Address | 是       |      | 网关的接收者     |

### 转让网关所有权

```bash
iriscli asset transfer-gateway-owner --moniker=cats --to=<new-owner-address> --from=<key-name> --chain-id=irishub --fee=0.3iris --commit
```

## iriscli asset issue-token

此命令用于在IRIS Hub上发行新通证。

```bash
iriscli asset issue-token <flags>
```

**标识：**

| 名称, 速记         | 类型    | 必须 | 默认          | 描述                                                         |
| ------------------ | ------- | -------- | ------------- | ------------------------------------------------------------ |
| --family           | string  | 是       | fungible      | 通证类型: fungible, non-fungible (不支持)                    |
| --source           | string  |          | native        | 通证来源: native, gateway                                    |
| --name             | string  | 是       |               | 通证的名称，限制为32个unicode字符，例如"IRIS Network"        |
| --gateway          | string  |          |               | 网关的唯一名称，当`source`是`gateway`时必需                  |
| --symbol           | string  | 是       |               | 通证的符号，长度在3到8之间，字母数字字符，不区分大小写       |
| --canonical-symbol | string  |          |               | 当`source`是`gateway`时，它用于标识其原始链上的符号          |
| --min-unit-alias   | string  |          |               | 最小单位的别名                                               |
| --initial-supply   | uint64  | 是       |               | 此通证的初始供应。 增发前的数量不应超过1000亿。              |
| --max-supply       | uint64  |          | 1000000000000 | 通证上限，总供应不能超过最大供应。 增发前的金额不应超过1万亿 |
| --decimal          | uint8   | 是       |               | 通证最多可以有18位小数                                       |
| --mintable         | boolean |          | false         | 首次发行后是否可以增发此通证                               |

### 发行原生通证

```bash
iriscli asset issue-token --family=fungible --source=native --name="Kitty Token" --symbol=kitty --initial-supply=100000000000 --max-supply=1000000000000 --decimal=0 --mintable=true --fee=1iris --from=<key-name> --commit
```

### 发行网关通证

#### 创建一个网关

在此示例之前，需要创建一个名为`cats`的网关, [详细信息](#iriscli-asset-create-gateway)。

```bash
iriscli asset create-gateway --moniker=cats --identity=<identity> --details=<details> --website=<website> --from=<key-name> --commit
```

#### 发行一个网关通证

```bash
iriscli asset issue-token --family=fungible --source=gateway --gateway=cats --canonical-symbol=cat --name="Kitty Token" --symbol=kitty --initial-supply=100000000000 --max-supply=1000000000000 --decimal=0 --mintable=true  --fee=1iris --from=<key-name> --commit
```

### 发送通证

您可以像[发送iris](./bank.md#iriscli-bank-send)一样发送任何通证。

#### 发送原生通证

```bash
iriscli bank send --from=<key-name> --to=<address> --amount=10kitty --fee=0.3iris --chain-id=irishub
```

#### 发送网关通证

```bash
iriscli bank send --from=<key-name> --to=<address> --amount=10cats.kitty --fee=0.3iris --chain-id=irishub
```

## iriscli asset edit-token

编辑通证信息。

```bash
iriscli asset edit-token <token-id> <flags>
```

**标识：**

| 名称, 速记         | 类型   | 必须 | 默认  | 描述                          |
| ------------------ | ------ | -------- | ----- | ----------------------------- |
| --name             | string |          |       | 通证名称，例如：IRIS Network  |
| --canonical-symbol | string |          |       | 网关或外部通证的源符号        |
| --min-unit-alias   | string |          |       | 通证最小单位别名              |
| --max-supply       | uint   |          | 0     | 通证的最大供应量              |
| --mintable         | bool   |          | false | 通证是否可以增发，默认为false |

`max-supply` 只能减少，且不得少于当前的总供应量。

### 编辑通证

```bash
iriscli asset edit-token cat --name="Cat Token" --canonical-symbol="cat" --min-unit-alias=kitty --max-supply=100000000000 --mintable=true --from=<key-name> --chain-id=irishub --fee=0.3iris --commit
```

## iriscli asset transfer-token-owner

转让通证所有权。

```bash
iriscli asset transfer-token-owner <token-id> <flags>
```

**标识：**

| 名称, 速记 | 类型   | 必须 | 默认 | 描述       |
| ---------- | ------ | -------- | ---- | ---------- |
| --to       | string | 是       |      | 接收人地址 |

### 转让通证所有者

```bash
iriscli asset transfer-token-owner kitty --to=<new-owner-address> --from=<key-name> --chain-id=irishub --fee=0.3iris --commit
```

## iriscli asset mint-token

资产所有者可以直接将通证增发到指定地址。

```bash
iriscli asset mint-token <token-id> <flags>
```

**标识：**

| 名称, 速记 | 类型   | 必须 | 默认 | 描述                               |
| ---------- | ------ | -------- | ---- | ---------------------------------- |
| --to       | string |          |      | 增发的通证的接收地址，默认为发起该交易地址 |
| --amount   | uint64 | 是       | 0    | 增发的金额                         |

### 增发通证

```bash
iriscli asset mint-token kitty --amount=1000000 --from=<key-name> --chain-id=irishub --fee=0.3iris
```

## iriscli asset query-token

查询在IRIS Hub发行的通证。

```bash
iriscli asset query-token <token-id>
```

### 全局唯一令牌ID生成规则

- 当`Source`为`native`时：ID = [Symbol]，例如`iris`

- 当`Source`为`external`时：ID = x.[Symbol]，例如`x.btc`

- 当`Source`为`gateway`时：ID = [Gateway].[Symbol]，例如`cats.kitty`

### 查询名为`kitty`的原生通证

```bash
iriscli asset query-token kitty
```

### 查询名为`kitty`的网关通证`cats`

```bash
iriscli asset query-token cats.kitty
```

### 查询名为`btc`的外部通证

```bash
iriscli asset query-token x.btc
```

## iriscli asset query-tokens

根据条件查询在IRIS Hub上发行的通证的集合。

```bash
iriscli asset query-tokens <flags>
```

**标识：**

| 名称, 速记 | 类型   | 必须 | 默认 | 描述                                        |
| ---------- | ------ | -------- | ---- | ------------------------------------------- |
| --source   | string |          |      | 枚举值: native / gateway / external         |
| --gateway  | string |          |      | 网关的唯一名称，当`source`是`gateway`时需要 |
| --owner    | string |          |      | 通证的所有者                                |

### Query Rules

- 当`source`为`native`
  - `gateway` 将被忽略
  - `owner` 为可选
- 当`source`为`gateway`
  - `gateway` 必须
  - `owner` 将被忽略（因为网关通证全部归网关所有）
- 当`source`为`external`
  - `gateway` 和 `owner` 被忽略
- 当`gateway`不为空
  - `source` 可选

### 查询所有通证

```bash
iriscli asset query-tokens
```

### 查询所有原生通证

```bash
iriscli asset query-tokens --source=native
```

### 查询网关名称为”cats“的所有通证

```bash
iriscli asset query-tokens --gateway=cats
```

### 查询指定所有者的所有通证

```bash
iriscli asset query-tokens --owner=<address>
```

## iriscli asset query-gateway

查询指定别名的网关。

```bash
iriscli asset query-gateway <flags>
```

**标识：**

| 名称, 速记 | 类型   | 必须 | 默认 | 描述                                               |
| ---------- | ------ | -------- | ---- | -------------------------------------------------- |
| --moniker  | string | 是       |      | 唯一名称，大小在3到8之间，以字母开头，后跟字母数字 |

### 查询网关

```bash
iriscli asset query-gateway --moniker cats
```

## iriscli asset query-gateways

查询指定所有者的所有网关。

```bash
iriscli asset query-gateways <flags>
```

**标识：**

| 名称, 速记 | 类型    | 必须 | 默认 | 描述               |
| ---------- | ------- | -------- | ---- | ------------------ |
| --owner    | Address |          |      | 要查询的所有者地址 |

### 查询网关列表

```bash
iriscli asset query-gateways --owner=<owner-address>
```

## iriscli asset query-fee

查询与资产相关的费用，包括网关创建，通证发行和增发。

```bash
iriscli asset query-fee <flags>
```

**标识：**

| 名称, 速记 | 类型   | 必须 | 默认 | 描述                                           |
| ---------- | ------ | -------- | ---- | ---------------------------------------------- |
| --gateway  | string |          |      | 网关名称，用于查询网关费用                     |
| --token    | string |          |      | 通证ID，用于查询通证费用 |

### 查询创建网关的费用

```bash
iriscli asset query-fee --gateway=cats
```

### 查询发行和增发原生通证的费用

```bash
iriscli asset query-fee --token=kitty
```

### 查询发行和增发网关通证的费用

```bash
iriscli asset query-fee --token=cats.kitty
```
