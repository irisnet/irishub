# HTLC

## 概念

[Hash Time Locked Contract（HTLC）](https://en.bitcoin.it/wiki/Hash_Time_Locked_Contracts) 是一个跨链原子交换协议。它以去中心化的方式完成用户间点对点的跨链交易。HTLC将跨链资产的交换过程分解为若干子过程。HTLC在协议上能够确保参与方之间的所有交换过程要么全部完成，要么一个都没发生。通过IRIS Hub的HTLC功能, 用户能够完成与Bitcoin、Ethereum、LTC以及其他所有支持HTLC的链之间的资产交换。

## 交互过程

假定Alice用BTC交换Bob的IRIS。Alice与Bob的原子交换过程可以由以下几个步骤表述。

### 第0步

  Alice and Bob 在链外对此次交换达成协议。这个过程通常在一个第三方交易平台完成。我们假设Alice是一个maker，而Bob是一个taker。这个协议包括BTC与IRIS的兑换率、彼此在对方链上的地址、由一个秘密值secret生成的哈希锁以及一个合适的过期时间。我们这里假定secret由Bob持有。

### 第1步

  Bob发送一个创建HTLC的交易到IRIS Hub，这个交易将通过协商的哈希锁锁定指定数量的IRIS到特定的地址。

### 第2步

  上述交易将在IRIS Hub上触发一个特定事件，监听工具（通常是钱包）或者平台将通知Alice。在验证这个事件符合协议之后，Alice使用相同的哈希锁在Bitcoin上发起一个HTLC创建交易。特别地，这个HTLC将使用一个相较于Bob提供的值足够小的时间锁。

### 第3步

  Bob通过类似方式接收到Bitcoin上的交易事件，确认之后在Alice提供的时间锁过期之前用secret去claim Bitcoin上锁定的BTC。

### 第4步

  当Bob在Bitcoin上成功claim之后，secret将被披露。通过这个secret，Alice将在IRIS Hub上完成相同的claim操作。这个操作也需要在Bob提供的过期时间之前完成。

## IRIS Hub HTLC 规范

### 创建HTLC消息

| **字段**             | **类型** | **描述**                                                     |
| -------------------- | -------- | ------------------------------------------------------------ |
| receiver             | Address  | 接收者地址                                                   |
| receiverOnOtherChain | string   | 另一条链上的HTLC认领接收地址，最大128个字符                        |
| amount               | Coins    | 欲交换的资产数量                                             |
| hashLock             | string   | 由secret(和timestamp, 如果提供)生成的sha256哈希值; 32字节, 十六进制表示 |
| timestamp            | uint64   | 时间戳, 如果提供则参与hash生成; 精度为秒                     |
| timeLock             | uint64   | 过期区块数; [50, 25480] (大于5分钟, 小于48小时)              |

### 认领HTLC消息

| **字段** | **类型** | **描述**                                        |
| -------- | -------- | ----------------------------------------------- |
| hashLock | string   | 创建HTLC时提供的hash lock                       |
| secret   | string   | 参与生成hash lock的随机数; 32字节, 十六进制表示 |

### 退款HTLC消息

| **字段** | **类型** | **描述**                  |
| -------- | -------- | ------------------------- |
| hashLock | string   | 创建HTLC时提供的hash lock |

## 操作

- [创建 HTLC](../cli-client/htlc.md#iriscli-htlc-create)
- [认领 HTLC](../cli-client/htlc.md#iriscli-htlc-claim)
- [退款 HTLC](../cli-client/htlc.md#iriscli-htlc-refund)
- [查询 HTLC](../cli-client/htlc.md#iriscli-htlc-query-htlc)
