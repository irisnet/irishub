# 收益分配

## 简介

该模块负责将收集的交易费和通胀的代币分发给所有验证人和委托人。为了减少计算压力，引入了一种被动分配策略。
`被动`意味着不会直接向贡献者自动支付利益。委托人或者验证人必须手动发送交易以取回其收益，否则，他们的收益将保留在全局收益池中。

## 收益

### 收益的来源

1. 交易所产生的交易费 `fee` (由交易的第一个签名者支付`fee`)
2. 通胀的代币 `inflation`   (目前IRISnet系统设置的通胀为4%每年，通证总量为`20亿`)

### 收益的去向

1. 验证人
2. 委托人
3. 社区基金
4. 出块奖励

:::tip
[计算公式](../concepts/general-concepts.md#staking-收益计算公式)
:::

## 使用场景

### 设置收益取回地址

默认情况下，收益将支付给发送委托交易的钱包。

委托人可以更改自己收益取回钱包的地址。将委托地址对应的钱包(标记为`A`)，希望收益取回钱包地址(标记为`B`)。

设置钱包B为取回收益的钱包：

```bash
iris tx distribution set-withdraw-addr [withdraw-addr] [flags]
```  

### 取回收益

根据取回场景的不同，有2种方式可以取回收益：

1.`withdraw-all-rewards` : 提取所有在外的委托收益）。

```bash
iris tx distribution withdraw-all-rewards [flags] --from=<key-name> --fees=0.3iris --chain-id=irishub
```

2.`withdraw-rewards` ：从指定验证人处提取委托收益。

```bash
iris tx distribution withdraw-rewards [validator-addr] [flags] --from=<key-name> --fees=0.3iris --chain-id=irishub
```

### 查询收益

查询委托人获得的所有奖励，可以选择为来自单个验证者的奖励。

```bash
iris query distribution rewards [delegator-addr] [validator-addr] [flags]
```

对于其它distribution相关的命令，请参考[stake-cli](../cli-client/distribution.md)
