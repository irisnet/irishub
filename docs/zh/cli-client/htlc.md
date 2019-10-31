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
iriscli htlc create --to=<receiver> --receiver-on-other-chain=<receiver-on-other-chain> --amount=<amount> --hash-lock=<hash-lock> --time-lock=<time-lock> --timestamp=<timestamp>
```

**标识：**

| 名称, 速记                | 类型     | 必须 | 默认 | 描述                                              |
| ------------------------- | -------- | ---- | ---- | ------------------------------------------------- |
| --to                      | string   | 是   |      | Bech32编码的收款人地址                                        |
| --receiver-on-other-chain | bytesHex |      |      | 交换对方在另一条链上的地址                           |
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
--secret=382aa2863398a31474616f1498d7a9feba132c4bcf9903940b8a5c72a46e4a41 \
--receiver-on-other-chain=0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826 \
--amount=10iris \
--time-lock=50 \
--timestamp=1580000000 \
--fee=0.3iris \
--chain-id=test \
--commit
```

## iriscli htlc claim

将 HTLC 中锁定的资金发送到收款人地址。

```bash
iriscli htlc claim --hash-lock=<hash-lock> --secret=<secret>
```

**标识：**

| 名称, 速记  | 类型     | 必须 | 默认 | 描述                              |
| ----------- | -------- | ---- | ---- | --------------------------------- |
| --hash-lock | bytesHex | 是   |      | T要发送锁定资金HTLC的Hash Lock |
| --secret    | bytesHex | 是   |      | 用于生成hash lock的secret      |

### 将 HTLC 中锁定的资金发送到收款人地址

```bash
iriscli htlc claim \
--from=userX \
--hash-lock=f054e34abd9ccc3cab12a5b797b8e9c053507f279e7e53fb3f9f44d178c94b20 \
--secret=5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f \
--fee=0.3iris \
--chain-id=testNet \
--commit
```

## iriscli htlc refund

从过期的 HTLC 中取回退款。

```bash
iriscli htlc refund --hash-lock=<hash-lock>
```

**标识：**

| 名称, 速记  | 类型     | 必须 | 默认 | 描述                       |
| ----------- | -------- | ---- | ---- | -------------------------- |
| --hash-lock | bytesHex | 是   |      | 要退款的HTLC的Hash Lock |

### 从过期的 HTLC 中取回退款

```bash
iriscli htlc refund \
--from=userX \
--hash-lock=f054e34abd9ccc3cab12a5b797b8e9c053507f279e7e53fb3f9f44d178c94b20 \
--fee=0.3iris \
--chain-id=testNet \
--commit
```

## iriscli htlc query-htlc

查询一个 HTLC 的详细信息。

```bash
iriscli htlc query-htlc <hash-lock>
```

### 查询HTLC详细信息

```bash
iriscli htlc query-htlc f054e34abd9ccc3cab12a5b797b8e9c053507f279e7e53fb3f9f44d178c94b20
```

执行完命令后，获得 HTLC 的详细信息如下。

```bash
HTLC:
    Sender:               iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym
    To:             iaa1zx6n0jussc3lx0dk0rax6zsk80vgzyy7kyfud5
    ReceiverOnOtherChain: 0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826
    Amount:               10000000000000000000iris-atto
    Secret:               5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f
    Timestamp:            0
    ExpireHeight:         107
    State:                completed
```
