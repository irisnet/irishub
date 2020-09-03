# Bank

Bank模块用于管理本地帐户中的资产

## 可用命令

| 名称                                             | 描述                         |
| ------------------------------------------------ | ---------------------------- |
| [balances](#iris-query-bank-balances)            | 查询账户余额              |
| [total](#iris-query-bank-total)                  | 查询链上币的总数           |
| [send](#iris-tx-bank-send)                       | 创建、签名、广播一个转账交易 |

## 常见问题

### ERROR: decoding bech32 failed

```bash
iris bank account iaa1a0x4g8rqc90l3z9jh98x7mkd0w77e9q9r300h 
Error: decoding bech32 failed: checksum failed. Expected 9r300k, got 9r300h.
```

这表示该帐户地址拼写错误，请仔细检查该地址。


## iris query bank balances

该命令用于查询特定地址的余额信息。

```bash
iris query bank balances [address] [flags]
```

**标识：**

| 名称，速记   | 类型   | 必须 | 默认                  | 描述                                 |
| ------------ | ------ | -------- | --------------------- | ------------------------------------ |
| -h, --help   |        |          |                       | `balances`子命令的提示信息            |
| --denom      | string |          |                       | 要查询的指定余额面值          |


## iris query bank total

该命令用于查询链上币的总数。

```bash
iris query bank total [flags]
```

**标识：**

| 名称，速记   | 类型   | 必须 | 默认                  | 描述                                 |
| ------------ | ------ | -------- | --------------------- | ------------------------------------ |
| -h, --help   |        |          |                       | `balances`子命令的提示信息            |
| --denom      | string |          |                       | 要查询的指定余额面值          |

## iris tx bank send

发送令牌到另一个地址，此命令包括 `generate`，`sign` 和 `broadcast` 这些步骤。

```bash
iris tx bank send [from_key_or_address] [to_address] [amount] [flags]
```

**标识：**

| 名称，速记 | 类型   | 必须 | 默认 | 描述                           |
| ---------- | ------ | -------- | ---- | ------------------------------ |
| --amount   | string | Yes      |      | 要发送的通证数量，例如：10iris |
| --to       | string |          |      | 接收通证的Bech32编码地址       |

### 将通证发送到另一个地址

```bash
iris tx bank send [from_key_or_address] [to_address] [amount] [flags]
```

