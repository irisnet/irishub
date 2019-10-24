# iriscli htlc

[HTLC模块](../features/htlc)提供了与其他链进行原子交换的相关功能。

## Available Commands

| Name                                   | Description                                          |
| -------------------------------------- | ---------------------------------------------------- |
| [create](#iriscli-htlc-create)         | 创建 HTLC                                            |
| [claim](#iriscli-htlc-claim)           | 将一个 OPEN 状态的 HTLC 中锁定的资金发放到收款人地址 |
| [refund](#iriscli-htlc-refund)         | 从过期的 HTLC 中取回退款                             |
| [query-htlc](#iriscli-htlc-query-htlc) | 查询一个 HTLC 的详细信息                             |

## iriscli htlc create

创建一个 HTLC。

```bash
iriscli htlc create --receiver=<receiver> --receiver-on-other-chain=<receiver-on-other-chain> --amount=<amount> --hash-lock=<hash-lock> --time-lock=<time-lock> --timestamp=<timestamp>
```

**标识：**

| 名称, 速记                | 类型     | 必须 | 默认 | 描述                                              |
| ------------------------- | -------- | ---- | ---- | ------------------------------------------------- |
| --receiver                | string   | 是   |      | 收款人地址                                        |
| --receiver-on-other-chain | bytesHex | 是   |      | 另一条链上的收款地址                              |
| --amount                  | string   | 是   |      | 要发送的金额                                      |
| --hash-lock               | bytesHex | 是   |      | 由 secret 和 时间戳（如果提供）生成的 sha256 哈希 |
| --time-lock               | string   | 是   |      | 资金锁定的区块数                                  |
| --timestamp               | uint     |      |      | 参与生成 hash lock 的10位时间戳（可选）           |

### 创建HTLC

```bash
iriscli htlc create \
--from=userX \
--receiver=faa1zx6n0jussc3lx0dk0rax6zsk80vgzyy7kyfud5 \
--receiver-on-other-chain=bb9188215a6112a6f7eb93e3e929197b3d44004cb691f95babde84cc18789364 \
--hash-lock=e8d4133e1a82c74e2746e78c19385706ea7958a0ca441a08dacfa10c48ce2561 \
--amount=10iris \
--time-lock=50 \
--timestamp=1580000000 \
--fee=0.3iris \
--chain-id=testNet \
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
| --hash-lock | bytesHex | 是   |      | T要发送锁定资金 HTLC 的 hash lock |
| --secret    | bytesHex | 是   |      | 用于生成 hash lock 的 secret      |

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
| --hash-lock | bytesHex | 是   |      | 要退款的 HTLC 的 hash Lock |

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
    Receiver:             iaa1zx6n0jussc3lx0dk0rax6zsk80vgzyy7kyfud5
    ReceiverOnOtherChain: 72656365697665724f6e4f74686572436861696e
    Amount:               10000000000000000000iris-atto
    Secret:               5f5f5f6162636465666768696a6b6c6d6e6f707172737475767778797a5f5f5f
    Timestamp:            0
    ExpireHeight:         107
    State:                completed
```
