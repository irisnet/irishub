# IRISLCD

## 介绍

IRISLCD是IRIShub的轻节点。与IRIShub全节点不同，它不会存储所有块并执行所有交易，这意味着它只需要最小的带宽，计算和存储资源。在不信任模式下，它将跟踪验证人集的演变，并要求完整节点返回共识证明和默克尔证明。除非具有超过2/3投票权的验证者集体执行拜占庭行为，否则IRISLCD证明验证算法可以检测所有潜在的恶意数据，这意味着IRISLCD节点可以提供与完整节点相同的安全性。

irislcd的默认主文件夹是`$HOME/.irislcd`。一旦IRISLCD启动，它将创建两个目录：`keys`和`trust-base.db`。密钥存储db位于`keys`中。`trust-base.db`存储所有可信验证器集和其他验证相关文件。

当IRISLCD在不信任模式下启动时，它将检查`trust-base.db`是否为空。如果是，那么它将获取最新块作为其信任基础并将其保存在`trust-base.db`下。IRISLCD节点无条件信任这个区块。所有查询证明都将已这个区块为基础进行验证，这意味着IRISLCD只能验证之后高度上的区块和交易。如果要查询较低高度的交易和区块，请以信任模式启动IRISLCD。有关详细的验证算法介绍，请参阅[tendermint lite](https://github.com/tendermint/tendermint/blob/master/docs/tendermint-core/light-client-protocol.md)。

## 使用

有关如何启动IRISLCD，请参阅[lcd_start](../light-client/README.md)
