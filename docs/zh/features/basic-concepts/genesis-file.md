# Genesis

Genesis文件是整个网络初始化的基础。它包含了创建创世区块的大部分信息(比如ChainID，共识的参数),初始化账户余额，每个模块的参数，validators信息。
在genesis.json文件中指明了各个账户的余额信息。账户余额数据是需要通过链下共识获得的。这个过程可能会依赖其他区块链的结果，或者一个代币产生事件。

## Basic State

* **genesis_time** Genesis文件创建时间
* **chain_id**     区块链的ID

## Consensus Params

* **block_size** 
  * `max_bytes` 区块大小限制
  * `max_gas`  区块最大Gas数量，默认值为-1表示无限制。如果累积的Gas超出Gas限制，该交易和之后的交易将执行失败。
* **evidence**   区块内作恶证据的生命周期

## App State

* **accounts** 初始化账户信息

* **auth** 与系统相关的参数 
  * params 详情见[Gov-Params](gov-params.md)

* **stake** 与抵押共识相关的参数
  * pool   当前绑定的token总数
  * validators 当前每个验证人的信息
  * params 详情见[Gov-Params](gov-params.md)
  
* **mint**  与通货膨胀相关的参数
  * minter 初始化状态
  * params 详情见[Gov-Params](gov-params.md)
  
* **distribution** 与分配收益有关的参数
  * params 详情见[Gov-Params](gov-params.md)

* **gov**  与链上治理相关的参数
  * params 详情见[Governance.md](../governance.md)
  

* **upgrade** 与升级相关的参数
  * GenesisVersion 初始化当前版本信息

* **slashing** 与惩罚validator相关的参数
  * params 详情见[Gov-Params](gov-params.md)

* **service**  与service相关的参数
  * params 详情见[Gov-Params](gov-params.md)
    
* **guardian** 与guardian相关的参数
  * `profilers` profiler列表
  * `trustees` trustee列表
  
## Gentxs

Gentxs中包含了创世区块中创建validators的交易集合。Cosmos SDK通过一个 gen-tx可以有效的对genesis文件的产生进行管控。 gen-tx即Genesis Transaction 是经过签名的交易数据，在区块链启动的过程中，这些交易将被执行，然后确定初始验证人集合。
通过提交gen-txs交易，代币的持有者证明了自己愿意参与到这个区块链网络的启动流程中，并且愿意承担潜在的损失。