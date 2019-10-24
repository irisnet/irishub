# 链上随机数

::: warning
此模块处于 Beta 阶段，使用前请评估风险
:::

## 简介

该规范描述了IRISHub上的随机数的使用方法和适用范围。目前该功能处于 Beta 阶段，使用前请自行评估风险。

## 概念

### 适用范围

**适用于**应用层获取基于区块链产生的随机数，比如博彩、游戏等

**不适用于**私钥、区块链共识算法等

### PRNG

伪随机数生成器（PRNG），也称为确定性随机数生成器（DRBG），是用于生成近似于随机数序列特性的数字序列算法。 -- 维基百科

IRISHub 在 Beta 阶段使用此算法生成随机数。

一方面，我们通过区块链生成的多个指标作为“因子”来计算随机数，使得此随机数公开透明，方便验证；

随机数“因子”具体包含以下指标：

- 上一个区块的 Hash：区块 Hash 的生成，取决于该区块的多方面因素，比如区块高度、交易数量、时间戳等等，因此区块 Hash 本身就具有一定的不可预测性
- 当前区块的时间戳：区块时间戳采用 BFT 时间，即根据验证人的权重，使用上一个区块中每一个Precommit的时间，加权计算出来的分布式时间戳（毫秒级别），也具有一定的不可预测性 [[BFT Time](https://tendermint.com/docs/spec/consensus/bft-time.html#bft-time)]
- 请求随机数的账户地址：主要是为了实现不同人在同一个区块高度得到不同的随机数

由于区块 Hash 和 BFT 时间的计算都是基于上一个区块的信息，为了避免请求随机数之前可以预先计算结果，所以另一方面，我们通过“未来区块”，加强随机数的不可预测性。

***但是，不可预测不代表不可操纵，比如一个新区块的提议者，可以选择性的打包交易，可以选择性的接受Precommit，以此来影响区块 Hash 和 BFT 时间***

#### 计算公式

```bash
seed = sha256(timestamp + int(sha256(blockhash)) / timestamp + int(sha256(consumer)) / timestamp)
rand = seed mod 10^20 / 10^20
```

### TRNG

硬件随机数生成器（HRNG）或真随机数生成器（TRNG），是通过物理过程（而不是通过算法）生成随机数的设备。 -- 维基百科

鉴于 PRNG 存在的安全风险，我们拟计划在下一个版本通过预言机引入 TRNG，以提高随机数的安全性

## 操作

- [请求随机数](../cli-client/rand.md#iriscli-rand-request-rand)
- [查询随机数](../cli-client/rand.md#iriscli-rand-query-rand)
- [查询随机数队列](../cli-client/rand.md#iriscli-rand-query-queue)
