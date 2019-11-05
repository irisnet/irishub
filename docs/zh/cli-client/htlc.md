# iriscli htlc

[HTLC模块](../features/htlc)提供了与其他链进行原子交换的相关功能。

## Available Commands

| Name                                   | Description                                          |
| -------------------------------------- | ---------------------------------------------------- |
| [create](#iriscli-htlc-create)         | 创建HTLC                                            |
| [claim](#iriscli-htlc-claim)           | 将一个OPEN状态的HTLC中锁定的资金发放到收款人地址 |
| [refund](#iriscli-htlc-refund)         | 从过期的HTLC中取回退款                             |
| [query-htlc](#iriscli-htlc-query-htlc) | 查询一个HTLC的详细信息                             |

## iriscli htlc create

创建一个 HTLC。

```bash
iriscli htlc create --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --to=<to> --receiver-on-other-chain=<receiver-on-other-chain> --amount=<amount> --secret=<secret> --time-lock=<time-lock> --timestamp=<timestamp>
```

**标识：**

| 名称, 速记                | 类型     | 必须 | 默认 | 描述                                              |
| ------------------------- | -------- | ---- | ---- | ------------------------------------------------- |
| --to                      | string   | 是   |      | Bech32编码的收款人地址                                        |
| --receiver-on-other-chain | bytesHex |      |      | 另一条链上的HTLC认领接收地址                          |
| --amount                  | string   | 是   |      | 要发送的金额                                      |
| --secret                  | bytesHex |      |      | 用于生成Hash Lock的secret, 缺省将随机生成    |
| --hash-lock               | bytesHex | 是   |      | 由secret和时间戳（如果提供）生成的sha256哈希 |
| --time-lock               | string   | 是   |      | 资金锁定的区块数                                  |
| --timestamp               | uint     |      |      | 参与生成hash lock的10位时间戳（可选）           |

### 创建HTLC

```bash
iriscli htlc create \
--from=node0 \
--to=faa1zx6n0jussc3lx0dk0rax6zsk80vgzyy7kyfud5 \
--receiver-on-other-chain=0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826 \
--amount=10iris \
--secret=382aa2863398a31474616f1498d7a9feba132c4bcf9903940b8a5c72a46e4a41 \
--time-lock=50 \
--timestamp=1580000000 \
--fee=0.3iris \
--chain-id=test \
--commit
```

## iriscli htlc claim

将 HTLC 中锁定的资金发送到收款人地址。

```bash
iriscli htlc claim --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --hash-lock=<hash-lock> --secret=<secret>
```

**标识：**

| 名称, 速记  | 类型     | 必须 | 默认 | 描述                              |
| ----------- | -------- | ---- | ---- | --------------------------------- |
| --hash-lock | bytesHex | 是   |      | T要发送锁定资金HTLC的Hash Lock |
| --secret    | bytesHex | 是   |      | 用于生成hash lock的secret      |

### 将 HTLC 中锁定的资金发送到收款人地址

```bash
iriscli htlc claim \
--from=node0 \
--hash-lock=bae5acb11ad90a20cb07023f4bf0fcf4d38549feff486dd40a1fbe871b4aabdf \
--secret=382aa2863398a31474616f1498d7a9feba132c4bcf9903940b8a5c72a46e4a41 \
--fee=0.3iris \
--chain-id=test \
--commit
```

## iriscli htlc refund

从过期的 HTLC 中取回退款。

```bash
iriscli htlc refund --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --hash-lock=<hash-lock>
```

**标识：**

| 名称, 速记  | 类型     | 必须 | 默认 | 描述                       |
| ----------- | -------- | ---- | ---- | -------------------------- |
| --hash-lock | bytesHex | 是   |      | 要退款的HTLC的Hash Lock |

### 从过期的 HTLC 中取回退款

```bash
iriscli htlc refund \
--from=node0 \
--hash-lock=bae5acb11ad90a20cb07023f4bf0fcf4d38549feff486dd40a1fbe871b4aabdf \
--fee=0.3iris \
--chain-id=test \
--commit
```

## iriscli htlc query-htlc

查询一个 HTLC 的详细信息。

```bash
iriscli htlc query-htlc <hash-lock>
```

### 查询HTLC详细信息

```bash
iriscli htlc query-htlc bae5acb11ad90a20cb07023f4bf0fcf4d38549feff486dd40a1fbe871b4aabdf
```

执行完命令后，获得 HTLC 的详细信息如下。

```bash
HTLC:
        Sender:               faa1a2g4k9w3v2d2l4c4q5rvvu7ggjcrfnynvrpqze
        To:                   faa1zx6n0jussc3lx0dk0rax6zsk80vgzyy7kyfud5
        ReceiverOnOtherChain: 0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826
        Amount:               10iris
        Secret:               382aa2863398a31474616f1498d7a9feba132c4bcf9903940b8a5c72a46e4a41
        Timestamp:            1580000000
        ExpireHeight:         59
        State:                completed
```
