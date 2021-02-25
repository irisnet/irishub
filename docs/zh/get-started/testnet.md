---
order: 4
---

# 加入测试网

主网完成 IRIS Hub 1.0 升级后，**Nyancat** 测试网开始作为稳定的应用程序测试网运行，该测试网具有与主网相同的软件版本。IRISnet 的服务提供方可以在不需要运行 IRIShub 节点的情况下，在 Nyancat 测试网上开发其应用程序。

## 公共端点

- GRPC：34.80.202.172:9090
- RPC：http://34.80.202.172:26657
- REST：http://34.80.202.172:1317/swagger/

## 运行节点

如果您想自行配置测试网节点而不是使用公共端点，可以参考加入[IRIS Hub 主网]（https://www.irisnet.org/docs/get-started/mainnet.html）步骤，除了：

### 创世文件

[点击下载](https://github.com/irisnet/testnets/raw/master/nyancat/config/genesis.json)

### 种子节点

在 `config.toml` 中添加以下 `seeds` 和 `persistent_peers`：

seeds：

```bash
07e58f179b2b7101b72f04248f542f67af8993bd@35.234.10.84:26656
```

persistent_peers：

```bash
bc77e49df0de4d70ab6f97f1e3a17bfb51a1ea7a@34.80.202.172:26656
```

### 水龙头

欢迎加入我们的【[Discord 社区](https://discord.gg/9cSt7MX2fn)】申请测试通证

## 浏览器

<https://nyancat.iobscan.io/>