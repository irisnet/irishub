# Slash

## 功能介绍

收集验证人异常的行为，并根据异常行为的类型实施相应的惩罚机制。

验证人异常的行为主要有以下三种：

1. 验证人节点长期不参与网络共识
2. 对同一次共识过程多次投票，并且这些投票相互矛盾
3. 验证人节点通过打包不合法的交易进入区块来扰乱网络共识

## 惩罚机制

1. 根据当前验证人拥有的voting power，计算验证人节点所绑定的token数量。
2. 惩罚作恶验证人一定比例的token，并把其踢出验证人集合；同时禁止此验证人在一段时间内再次进入验证人集合，这个过程被称为jail验证人。
3. 对于不同类型的异常行为，采用不同的惩罚比例和jail的时间。
4. 惩罚细则：
	1. 如果当前验证人token总数为A，惩罚比例为B，那么对此验证人最多惩罚的token的数量为A*B。
	2. 如果在当前高度有处于unbonding阶段的unbonding delegation和redelegation,并且unbonding delegation和redelegation的创建高度小于作恶实施高度，则按比例B先惩罚这两部分的token
	3. 对unbonding delegation和redelegation总共惩罚的token数量为S。如果S小于A*B，则惩罚验证人的token梳理为`A*B-S`。否则不惩罚绑定在验证人上的token

## 长时间不参与网络共识

在固定时间窗口`SignedBlocksWindow`内，验证人的缺席出块数目比重大于`MinSignedPerWindow`，则以`SlashFractionDowntime`比例惩罚验证人的绑定的token,并jail验证人。直到jail时间超过`DowntimeJailDuration`，才能通过unjail命令解除jail。

### 参数

* `SignedBlocksWindow` 默认值: 20000
* `MinSignedPerWindow` 默认值: 0.5
* `DowntimeJailDuration` 默认值: 2天
* `SlashFractionDowntime` 默认值: 0.005

## 恶意投票

执行区块时, 收到某验证人对同一高度同一Round不同区块都进行签名的作恶证据，如果作恶的时间距当前区块时间小于`MaxEvidenceAge`，则以`SlashFractionDoubleSign`比例惩罚验证人的绑定的token,并jail验证人。直到jail时间超过`DoubleSignJailDuration`，才能通过unjail命令解除jail。

### 参数

* `MaxEvidenceAge` 默认值: 1天
* `DoubleSignJailDuration` 默认值: 5天
* `SlashFractionDoubleSign`默认值: 0.01

## 打包不合法的交易

如果节点在执行区块过程中，检测到其中交易只要没有通过`txDecoder`, `validateTx`, `validateBasicTxMsgs`, 则以`SlashFractionCensorship`比例惩罚验证人的绑定的token, 并jail验证人。直到jail时间超过`CensorshipJailDuration`，才能通过unjail命令解除jail。

* `txDecode` 对Tx的反序列化
* `validateTx` 对Tx的大小限制
* `validateBasicTxMsgs` 对tx中msg的基本检查

### 参数

* `CensorshipJailDuration` 默认值: 7天
* `SlashFractionCensorship` 默认值: 0.02

## 相关用户接口

### unjail

如果validator被jail，并且jail的时间已经过去，则可以通过以下命令unjail。

```
iriscli stake unjail --from=<key name> --fee=0.004iris --chain-id=<chain-id>
```
