# 通货膨胀

## 介绍

POW 共识网络的激励机制十分简明：一旦新的区块产生，那么区块的产生者将会获得一定数量的 token 和累积的交易费作为出块奖励。IRISnet 是 POS 区块链网络，这里的奖励生成过程跟 POW 区块链网络类似，但是奖励分配机制复杂很多。

POW 网络中，在每个区块生成期间，所有矿工竞争计算工作证明，最快计算出结果的矿工将成为赢家。实际上，所有失败的矿工都没有向优胜者矿工提供任何积极的帮助或合作，他们只是竞争对手。因此，向获胜者矿工授予所有奖励是合理的。但是，在 POS 区块链网络中，我们不能这样做。因为每个块生成过程都是所有验证人和委托人的协作，这意味着所有这些贡献者应该共享出块收益。至于如何将出块奖励分发给贡献者，我们将在 [distribution](distribution.md) 模块中详细解说。

在 IRISnet 网络中，奖励有两个来源，一个是区块中打包的交易的交易费；另一个是每个区块中增发的 token，我们把增发的这部分 token 称为**通胀**。这里，mint 模块负责**通胀**的计算，并把**通胀**的 token 添加到奖励池中。

## 计算通胀

### 区块时间

区块时间不是机器时间，因为不同机器的时间不可能完全相同。 他们或多或少一定会有一些偏差，这将导致不确定性。 这里的时间是指BFT时间。 有关详细说明，请参阅 [tendermint bft-time](https://github.com/tendermint/tendermint/blob/master/docs/spec/consensus/bft-time.md)。

### 通胀率

genesis 中指定的初始通胀率是 4%，这个值可以通过在 governance 中提交`参数修改`的提议来修改。相关步骤，请查阅 [governance](governance.md)。

### 通胀计算

通胀计算的公式如下：

```bash
blockCostTime  = (当前区块的BFT time) - (上一个区块的BFT time)
AnnualInflationAmount = inflationBasement * inflationRate
blockInflationAmount = AnnualInflationAmount * blockCostTime / (year)
```

`inflationBasement` 的值被定义在 genesis 文件中. 默认情况下，genesis 里面写入的值是 `2000000000iris`（20亿个 iris，`1 iris` 等于 `1*10^18 iris-atto`）。

假设 `blockCostTime` 是5000毫秒，通胀比例 `inflationRate` 是 `4%`，那么这个块增发的 token 数量是 `12675235125611580094iris-atto`（`12.675235125611580094iris`）

## 对用户的影响

通胀计算是一个自动过程，没有用户接口能直接干预此过程。每产生一个新的区块，就会增发一定数量的 token，loose tokens 的数量也会因此增加。

这里可以通过 staking 模块命令行和 restful api 来查询总的 `loose tokens` 的数量：

**`iris q staking pool`**

这个接口执行速度比较快，但是不能做默克尔证明，因此如果连接不上可信的全节点，请不要使用此接口。

```bash
iris q staking pool --node=<iris_node_url>
```

示例输出：

```bash
Pool
Loose Tokens: 1846663.900384156921391687
Bonded Tokens: 425182.329615843078608313
Token Supply: 2271846.230000000000000000
Bonded Ratio: 0.187152776500000000
```

**`/stake/pool`**

这个 restful api 的用法请参阅 LCD swagger 文档。

如何运行一个 LCD 节点，请参阅 [LCD 文档](../light-client/intro.md)。
