---
order: 3
---

# 加入主网

:::tip
需要先 [安装 iris](install.md)
:::

## 运行全节点

### Start node from genesis

:::tip
必须使用 irishub [v1.0.1](https://github.com/irisnet/irishub/releases/tag/v1.0.1) 初始化你的节点
:::

```bash
# 初始化节点
iris init <moniker> --chain-id=irishub-1

# 下载主网公开的 config.toml 和 genesis.json
curl -o ~/.iris/config/config.toml https://raw.githubusercontent.com/irisnet/mainnet/master/config/config.toml
curl -o ~/.iris/config/genesis.json https://raw.githubusercontent.com/irisnet/mainnet/master/config/genesis.json

# 启动节点（也可使用 nohup 或 systemd 等方式后台运行）
iris start
```

接下来，你的节点将执行所有链升级过程。在每次升级之间，你必须使用特定的版本同步区块。不用担心在升级高度使用旧版本，节点会自动停止。

| 提案 | 起始高度 | 升级高度 | irishub 版本 |
| -------- | ------------ | -------------- | ----- |
| genesis  |  9146456     |  9593205  | [v1.0.1](https://github.com/irisnet/irishub/releases/tag/v1.0.1) |
| [#1](https://irishub.iobscan.io/#/ProposalsDetail/1)  |  9593206     |    | [v1.1.0](https://github.com/irisnet/irishub/releases/tag/v1.1.0), [v1.1.1](https://github.com/irisnet/irishub/releases/tag/v1.1.1) |
| [#8](https://irishub.iobscan.io/#/ProposalsDetail/8)  |  12393048     | 12534300 | [v1.2.0](https://github.com/irisnet/irishub/releases/tag/v1.2.0), [v1.2.1](https://github.com/irisnet/irishub/releases/tag/v1.2.1) |
| [#11](https://irishub.iobscan.io/#/ProposalsDetail/11)  |  14166918     |  14301916  | [v1.3.0](https://github.com/irisnet/irishub/releases/tag/v1.3.0) |

:::tip
您可能会看到一些连接错误，这没关系，P2P网络正在尝试查找可用的连接

可以添加几个[社区公开节点](https://github.com/irisnet/mainnet/blob/master/config/community-peers.md) 到`config.toml`中的`persistent_peers`。

如果您在不需要历史数据的情况下要快速启动节点并加入 IRIS Hub，可以考虑使用 [state_sync](./state-sync.md) 功能快速启动节点。
:::

## 升级为验证人节点

### 创建钱包

您可以[创建新的钱包](../cli-client/keys.md#创建密钥)或[导入现有的钱包](../cli-client/keys.md#通过助记词恢复密钥)，然后从交易所或其他任何地方转入一些IRIS到您刚刚创建的钱包中：

```bash
# 创建一个新钱包
iris keys add <key-name>
```

:::warning
在安全的地方备份好助记词！如果您忘记密码，这是恢复帐户的唯一方法。
:::

### 确认节点同步状态

```bash
# 可以使用此命令安装 jq
# apt-get update && apt-get install -y jq

# 如果输出为 false，则表明您的节点已经完成同步
iris status | jq .sync_info.catching_up
```

### 创建验证人

只有节点已完成同步时，才可以运行以下命令将您的节点升级为验证人：

```bash
iris tx staking create-validator \
    --pubkey=$(iris tendermint show-validator) \
    --moniker=<your-validator-name> \
    --amount=<amount-to-be-delegated, e.g. 10000iris> \
    --min-self-delegation=1 \
    --commission-max-change-rate=0.1 \
    --commission-max-rate=0.1 \
    --commission-rate=0.1 \
    --gas=100000 \
    --fees=0.6iris \
    --chain-id=irishub-1 \
    --from=<key-name>
```

:::warning
**重要**

一定要备份好 home（默认为〜/.iris/）目录中的 `config` 目录！如果您的服务器磁盘损坏或您准备迁移服务器，这是恢复验证人的唯一方法。
:::

如果以上命令没有出现错误，则您的节点已经是验证人或候选人了（取决于您的Voting Power是否在前100名中）

阅读更多：

- 概念
  - [基础概念](../concepts/general-concepts.md)
  - [验证人问答](../concepts/validator-faq.md)
- 验证人安全
  - [哨兵节点 (DDOS 防护)](../concepts/sentry-nodes.md)
  - [密钥管理](../tools/kms.md)

## 水龙头

前往 Stakely 开发上线的水龙头申请 IRISnet 主网通证。

具体申请方法请参见页面提示：https://stakely.io/faucet/irisnet-iris
