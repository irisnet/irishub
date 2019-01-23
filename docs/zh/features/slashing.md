# Slash

## 基本功能

通过对一些作恶的验证人进行惩罚，更好地维持网络的健康成长，保持网络的活跃度。

主要有三种惩罚类型：

1. 惩罚验证人长时间不在线
2. 惩罚验证人的double sign行为
3. 惩罚验证人hack代码，使用极小的代价构造垃圾数据上链。

## Long Downtime

在固定时间窗口`SignedBlocksWindow`内，验证人的缺席出块数目比重大于`MinSignedPerWindow`，则以`SlashFractionDowntime`比例惩罚验证人的绑定的token,并jail验证人。直到jail时间超过`DowntimeJailDuration`，才能通过unjail命令解除jail。

### 参数

* `SignedBlocksWindow` 默认值: 20000
* `MinSignedPerWindow` 默认值: 0.5
* `DowntimeJailDuration` 默认值: 2天
* `SlashFractionDowntime` 默认值: 0.005

## Double Sign

执行区块时, 收到某验证人对同一高度同一Round不同区块都进行签名的作恶证据，如果作恶的时间距当前区块时间小于`MaxEvidenceAge`，则以`SlashFractionDoubleSign`比例惩罚验证人的绑定的token,并jail验证人。直到jail时间超过`DoubleSignJailDuration`，才能通过unjail命令解除jail。

### 参数

* `MaxEvidenceAge` 默认值: 1天
* `DoubleSignJailDuration` 默认值: 5天
* `SlashFractionDoubleSign`默认值: 0.01

## Propoer Censorship

如果节点在执行区块过程中，检测到其中交易只要没有通过`txDecoder`, `validateTx`, `validateBasicTxMsgs`, 则以`SlashFractionCensorship`比例惩罚验证人的绑定的token, 并jail验证人。直到jail时间超过`CensorshipJailDuration`, n。

* `txDecode` 对Tx的反序列化
* `validateTx` 对Tx的大小限制
* `validateBasicTxMsgs` 对tx中msg的基本检查

### 参数

* `CensorshipJailDuration` 默认值: 7天
* `SlashFractionCensorship` 默认值: 0.02

## slash 命令

### unjail

如果validator被jail，并且jail的时间已经过去，则可以通过以下命令unjail。

```
iriscli stake unjail --from=<key name> --fee=0.004iris --chain-id=<chain-id>
```
