# 通货膨胀

## 介绍

POW共识网络的激励机制十分简明：一旦新的区块产生，那么区块的产生者将会获得一定数量的token和累积的交易费作为出块奖励。IRISnet是POS区块链网络，这里的奖励生成过程跟POW区块链网络类似，但是奖励分配机制复杂很多。

POW网络中，在每个区块生成期间，所有矿工竞争计算工作证明，最快计算出结果的矿工将成为赢家。实际上，所有失败的矿工都没有向优胜者矿工提供任何积极的帮助或合作，他们只是竞争对手。因此，向获胜者矿工授予所有奖励是合理的。但是，在POS区块链网络中，我们不能这样做。因为每个块生成过程都是所有验证人和委托人的协作，这意味着所有这些贡献者应该共享出块收益。至于如何将出块奖励分发给贡献者，我们将在[distribution](distribution.md)模块中详细解说。

在IRISnet网络中，奖励有两个来源，一个是区块中打包的交易的交易费；另一个是每个区块中增发的token，我们把增发的这部分token称为**通胀**。这里，mint模块负责**通胀**的计算，并把**通胀**的token添加到奖励池中。

## 计算通胀

### 区块时间

区块时间不是机器时间，因为不同机器的时间不可能完全相同。 他们或多或少一定会有一些偏差，这将导致不确定性。 这里的时间是指BFT时间。 有关详细说明，请参阅[tendermint bft-time](https://github.com/tendermint/tendermint/blob/master/docs/spec/consensus/bft-time.md)。

### 通胀率

genesis中指定的初始通胀率是4%，这个值可以通过在governance中提交`参数修改`的提议来修改。相关步骤, 请查阅[governance](governance.md)。

### 通胀计算

通胀计算的公式如下：

```bash
 blockCostTime  = (当前区块的BFT time) - (上一个区块的BFT time)
 AnnualInflationAmount = inflationBasement * inflationRate
 blockInflationAmount = AnnualInflationAmount * blockCostTime / (year)
```

`inflationBasement`的值被定义在genesis文件中. 默认情况下，genesis里面写入的值是 `2000000000iris`(20亿个iris, `1 iris`等于`1*10^18 iris-atto`)。

假设`blockCostTime`是5000毫秒， 通胀比例`inflationRate`是`4%`, 那么这个块增发的token数量是`12675235125611580094iris-atto` (`12.675235125611580094iris`)

## 对用户的影响

通胀计算是一个自动过程，没有用户接口能直接干预此过程。每产生一个新的区块，就会增发一定数量的token，loose tokens的数量也会因此增加。

这里有两个命令行接口和两个restful api来查询总的loose tokens的数量：

1. `iriscli stake pool`

    这个接口执行速度比较快，但是不能做默克尔证明，因此如果连接不上可信的全节点，请不要使用此接口。

    ```bash
    iriscli stake pool --node=<iris_node_url>
    ```

    示例输出:

    ```bash
    Pool
    Loose Tokens: 1846663.900384156921391687
    Bonded Tokens: 425182.329615843078608313
    Token Supply: 2271846.230000000000000000
    Bonded Ratio: 0.187152776500000000
    ```

2. `iriscli bank token-stats`

    如果不信任连接的全节点，请加上`--trust-node=false`这个标志。如果连接不上可信的全节点，这个接口十分必要。

    ```bash
    iriscli bank token-stats --trust-node=false --chain-id=<chain-id> --node=<iris_node_url>
    ```

    示例输出:

    ```bash
    TokenStats:
      Loose Tokens:  1864477.596384156921391687iris
      Burned Tokens:  177.59638iris
      Bonded Tokens:  425182.329615843078608313iris
    ```

3. `/stake/pool`和`/bank/token-stats`

    这两个restful api的用法请参阅LCD swagger文档。

    如何运行一个LCD节点，请参阅[LCD文档](../light-client/intro.md)。
