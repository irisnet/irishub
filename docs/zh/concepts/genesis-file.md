---
order: 4
---

# Genesis 文件

Genesis 文件 (~/.iris/config/genesis.json) 是整个网络初始化的基础，其中包含有关创建Genesis区块的大部分信息（例如ChainID，共识参数，应用程序状态），初始化帐户余额，模块级别的参数，以及验证人信息。
Genesis 文件用于设置IRIS网络的初始参数，在Genesis文件上建立健全的社会共识对于启动网络至关重要。

账户余额数据是需要通过链下共识获得的，这个过程可能会依赖其他区块链的数据，或者一个代币产生事件。

## 基础参数

* **genesis_time** 区块链启动时间
* **chain_id**     区块链ID
* **initial_height** 区块链初始块高

## 共识参数

* **block**
  * `max_bytes` 区块大小限制。
  * `max_gas` 区块最大Gas数量，默认值为-1表示无限制。如果累积的Gas超出Gas限制，该交易和之后的交易将执行失败。
  * `time_iota_ms` 连续块之间的最小时间增量（以毫秒为单位）。
* **evidence** 区块内作恶证据的生命周期
  * `max_age_num_blocks` 最长证据年龄，以区块为单位。
  * `max_age_duration`  最长证据年龄，以时间为单位。
  * `max_bytes` 可以在单个块中提交的总存证的最大大小（以字节为单位）。
* **validator**  验证人信息
  * `pub_key_types` 验证人使用的公钥类型。

## 应用参数

* **auth** 系统相关的参数

* **bank** 账户相关参数

* **capability** capability模块相关参数

* **coinswap** 流动性池相关参数

* **crisis** 危机模块相关参数

* **distribution** 收益分配相关的参数

* **evidence** evidence模块相关参数

* **genutil** 生成工具模块相关参数

* **gov**  链上治理相关的参数

* **guardian** 特殊账户相关的参数

* **htlc** 哈希时间锁定合约相关参数

* **ibc** ibc模块相关参数

* **nft** 非同质化通证相关参数

* **mint**  通货膨胀相关的参数

* **oracle**  oracle模块相关参数

* **params**  参数模块相关参数

* **random**  随机数相关参数

* **record** 存证相关参数

* **service** 服务相关的参数

* **slashing**  惩罚机制相关的参数

* **staking**  抵押和共识相关的参数

* **token**  资产相关的参数

* **transfer**  ibc transfer模块相关参数

* **upgrade** 升级相关的参数

* **vesting** 授权模块相关参数

可治理参数详见 [治理参数](gov-params.md)

## Gentxs

Gentxs中包含了创世区块中创建validators的交易集合。Cosmos SDK通过一个`gen-tx`可以有效的对genesis文件的产生进行管控。`gen-tx`即Genesis Transaction是经过签名的交易数据，在区块链启动的过程中，这些交易将被执行，然后确定初始验证人集合。
通过提交`gen-tx`交易，代币的持有者证明了自己愿意参与到这个区块链网络的启动流程中，并且愿意承担潜在的损失。
