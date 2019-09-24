# HTLC用户指南

## 概念

[Hash Time Locked Contract（HTLC）]() 是一个跨链原子交换协议。它以去中心化的方式完成用户间点对点的跨链交易。HTLC将跨链资产的交换过程分解为若干子过程。HTLC在协议上能够确保参与方之间的所有交换过程要么全部完成，要么一个都没发生。通过Iris Hub的HTLC功能, 用户能够完成与Bitcoin、Ethereum、LTC以及其他所有支持HTLC的链之间的资产交换。

## 交互过程

假定Alice用BTC交换Bob的IRIS。Alice与Bob的原子交换过程可以由以下几个步骤表述。

### 第0步

### 第1步

### 第2步

### 第3步

## Iris Hub HTLC 规范

### 创建HTLC消息

| **字段**       | **类型** | **描述**           |
|----------------|----------|------------------|
| receiver | Address | 接收者地址              |
| receiverOnOtherChain | string | 其他链上的接收者地址; 最大32字节, 十六进制表示 |
| amount | Coin | 欲交换的资产数量 |
| hashLock | string | 由secret(和timestamp, 如果提供)生成的sha256哈希值; 32字节, 十六进制表示 |
| timestamp | uint64 | 时间戳, 如果提供则参与hash生成; 精度为秒 |
| timeLock | uint64 | 过期区块数; [50, 25480] (大于5分钟, 小于48小时) |

### Claim HTLC消息

| **字段**       | **类型** | **描述**           |
|----------------|----------|------------------|
| hashLock | string | 创建HTLC时提供的hash lock |
| secret | string | 参与生成hash lock的随机数; 32字节, 十六进制表示 |

### 退回HTLC消息

| **字段**       | **类型** | **描述**           |
|----------------|----------|------------------|
| hashLock | string | 创建HTLC时提供的hash lock |

## 行为

### [创建HTLC]()
### [Claim HTLC]()
### [退回HTLC]()
### [查询HTLC]()
