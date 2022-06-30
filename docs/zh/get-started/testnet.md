---
order: 4
---

# 加入测试网

主网完成 IRIS Hub 1.0 升级后，**Nyancat** 测试网开始作为稳定的应用程序测试网运行，该测试网具有与主网相同的软件版本。IRISnet 的服务提供方可以在不需要运行 IRIShub 节点的情况下，在 Nyancat 测试网上开发其应用程序。

## 公共端点

- GRPC: 35.234.10.84:9090
- RPC: http://35.234.10.84:26657/
- REST: http://35.234.10.84:1317/swagger/



## 运行全节点

#### 从genesis开始运行节点

::提示
必须使用 irishub [v1.1.1](https://github.com/irisnet/irishub/releases/tag/v1.1.1)[ ](https://github.com/irisnet/irishub/releases/tag/v1.0.1) 初始化你的节点::

```bash
# 初始化节点
iris init <moniker> --chain-id=nyancat-8

# 下载公开的 config.toml 和 genesis.json
curl -o ~/.iris/config/config.toml https://github.com/irisnet/testnets/blob/master/nyancat/config/config.toml
curl -o ~/.iris/config/genesis.json https://raw.githubusercontent.com/irisnet/testnets/master/nyancat/config/genesis.json

# 启动节点（也可使用 nohup 或 systemd 等方式后台运行）
iris start
```



## 水龙头

欢迎加入我们的【[nyancat-faucet](https://discord.gg/Z6PXeTb5Mt)】频道申请测试通证

申请方法：在 [nyancat-faucet](https://discord.gg/Z6PXeTb5Mt) 频道中，发送：`$faucet <your_addr>`，每个 Discord 账号每 24 小时只可领取一次测试通证（NYAN）

## 浏览器

<https://nyancat.iobscan.io/>

## 社区

欢迎加入我们的社区进行讨论：[nyancat testnet](https://discord.gg/9cSt7MX2fn)

