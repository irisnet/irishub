---
order: 4
---

# 加入测试网

我们有3个测试网：Fuxi、Nyancat 和 Bifrost。

自从主网启动以来，**Fuxi 测试网** 开始作为稳定的应用程序测试网运行，该测试网具有与主网相同的版本，因此 IRISnet 的服务提供商可以在不运行节点和 LCD 实例的情况下在 IRIShub 上开发和测试其应用程序。

但是，在 IRIShub 的新版本升级到主网之前，我们还需要和验证人一起进行升级和新功能相关的测试，这是 **Nyancat 测试网**所关注的重点。新的验证人也可以使用 Nyancat 测试网来熟悉验证人相关操作。

IRISnet 的 DeFi 和跨链 **Bifrost 测试网** 将与 Cosmos 的 Stargate 测试网同时运行，以测试 IRIS Hub 新增的 DeFi 功能、与 Cosmos SDK 新版本的集成以及通过 IBC 协议进行跨链协作 。

## 安装

我们使用不同的bech32前缀来区分主网和测试网，您所需要做的就是在[构建或安装](install.md) iris二进制文件之前在[irishub](https://github.com/irisnet/irishub) 源码根目录中执行以下命令：

```bash
source scripts/setTestEnv.sh # 用于构建和安装测试网版本
```

## Fuxi 测试网

Fuxi 测试网不支持以全节点的方式连接，您可以使用公共RPC和LCD来开发和测试您的应用程序：

- RPC：<http://rpc.testnet.irisnet.org:80>

- LCD：<https://lcd.testnet.irisnet.org/swagger-ui/>

- 区块浏览器：<https://testnet.irisplorer.io>

## Nyancat 测试网

请参考[如何加入nyancat测试网](https://github.com/irisnet/testnets/tree/master/nyancat#how-to-join-nyancat-testnet)

## Bifrost 测试网

请参考[如何加入bifrost测试网](https://github.com/irisnet/testnets/blob/master/bifrost/README.md)
