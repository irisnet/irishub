# HTLC

[HTLC 模块](../features/htlc.md)提供了与其他链进行原子交换的相关功能。

在 HTLC 生命周期中有下面的几种状态：

- open：HTLC 是可申领的
- completed：HTLC 已经被申领
- refunded：HTLC 已退还

## 可用命令

| Name                           | Description                                          |
| ------------------------------ | ---------------------------------------------------- |
| [create](#iris-tx-htlc-create) | 创建 HTLC                                            |
| [claim](#iris-tx-htlc-claim)   | 将一个 OPEN 状态的 HTLC 中锁定的资金发放到收款人地址 |
| [htlc](#iris-query-htlc-htlc)  | 查询一个 HTLC 的详细信息                             |

## iris tx htlc create

创建一个 HTLC。

```bash
iris tx htlc create \
    --to=<recipient> \
    --receiver-on-other-chain=<receiver-on-other-chain> \
    --sender-on-other-chain=<sender-on-other-chain> \
    --amount=<amount> \
    --hash-lock=<hash-lock> \
    --secret=<secret> \
    --timestamp=<timestamp> \
    --time-lock=<time-lock> \
    --transfer=<true | false> \
    --from=mykey
```

**标识：**

| 名称，速记                | 类型   | 必须 | 默认  | 描述                                                                    |
| ------------------------- | ------ | ---- | ----- | ----------------------------------------------------------------------- |
| --to                      | string | 是   |       | Bech32 编码的收款人地址                                                 |
| --receiver-on-other-chain | string |      |       | 另一条链上的 HTLC 认领接收地址                                          |
| --receiver-on-other-chain | string |      |       | 另一条链上对应 HTLC 的创建者地址                                        |
| --amount                  | string | 是   |       | 要发送的金额                                                            |
| --secret                  | string |      |       | 用于生成 hash lock 的 secret，缺省将随机生成                            |
| --hash-lock               | string |      |       | 由 secret 和时间戳（如果提供）生成的 sha256 哈希，缺省将使用 secre 生成 |
| --time-lock               | string | 是   |       | 资金锁定的区块数                                                        |
| --timestamp               | uint   |      |       | 参与生成 hash lock 的 10 位时间戳（可选）                               |
| --transfer                | bool   |      | false | 是否是 HTLT 交易                                                        |

## iris tx htlc claim

将 HTLC 中锁定的资金发送到收款人地址。

```bash
iris tx htlc claim [id] [secret] [flags] --from=mykey
```

## iris query htlc htlc

查询一个 HTLC 的详细信息。

```bash
iris query htlc htlc <id>
```

## iris query htlc params

查询 HTLC 模块参数

```bash
iris query htlc params
```

## iris query htlc supplies

查询所有 HTLT 资产的 supply

```bash
iris query htlc supplies
```

## iris query htlc supply

查询一个 HTLT 资产的 supply

```bash
iris query htlc supply [denom]
```
