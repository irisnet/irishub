# 资产管理

## 简介

该规范描述了IRISHub上的资产管理。任何人都可以在IRISHub上发布新资产，或者通过On-Chain Governance将现有资产与任何其他区块链挂钩。

## 概念

### 资产

#### 内部资产（用户发行资产）

IRISHub 允许个人和公司创建和发行他们自己的资产，用于他们可以想象的任何场景。内部资产的潜在用例是数不胜数的。一方面，可用作存放在客户手机上的门票，以通过音乐会的入口。另一方面，它们可用于众筹，权益跟踪甚至以股票形式出售公司股权。
想要创建新的内部资产，您需要做的仅仅是执行一行命令，为您的资产定义初始化参数，例如总量，符号，描述等。然后，您可以发送一些您自己发行的 Token 到任何人的账户，就像 iris 转账一样简单。
作为该资产的所有者，您无需处理区块链的任何技术细节，例如分布式共识算法，区块链开发或集成，而且也不需要运行任何挖矿设备或服务器。

### 费用

#### 相关参数

| name              | Type | Default   | Description                     |
| ----------------- | ---- | --------- | ------------------------------- |
| AssetTaxRate      | Dec  | 0.4       | 资产税率，即Community Tax的比例 |
| IssueTokenBaseFee | Coin | 60000iris | 发行FT的基准费用                |
| MintTokenFeeRatio | Dec  | 0.1       | 增发FT的费率(相对于发行费用)    |

注：以上参数均为可治理参数

#### 发行 Fungible Token 费用

- 基准费用：发行FT所需的基本费用，即FT Symbol长度为最小(3)时的费用
- 费用因子计算公式：(ln(len({symbol}))/ln3)^4
- 发行FT费用计算公式：IssueFTBaseFee/费用因子；结果取整到iris（大于1时四舍五入，小于等于1时取值为1）

#### 增发 Fungible Token 费用

- 增发FT费率：相对于发行FT时的费率
- 增发FT费用计算公式：发行FT费用 * MintFTFeeRatio；结果取整到iris（大于1时四舍五入，小于等于1时取值为1）

#### 费用扣除方式

- Community Tax：资产相关的操作费用一部分将作为Community Tax，比例由AssetTaxRate决定。
- Burned：剩余部分将被销毁

## 操作

- **资产**

  - [发行资产](../cli-client/asset.md#iriscli-asset-token-issue)

    - [发行原生资产](../cli-client/asset.md#发行通证)

    - [转账](../cli-client/asset.md#发送通证)

  - [查询资产列表](../cli-client/asset.md#iriscli-asset-token-tokens)

  - [编辑资产信息](../cli-client/asset.md#iriscli-asset-token-edit)

  - [增发](../cli-client/asset.md#iriscli-asset-token-mint)

  - [销毁](../cli-client/bank.md#iriscli-bank-burn)

  - [转让所有权](../cli-client/asset.md#iriscli-asset-token-transfer)

- **费用**

  - [查询资产发行和增发费用](../cli-client/asset.md#查询发行和增发通证的费用)