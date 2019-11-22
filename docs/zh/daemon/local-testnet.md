---
order: 3
---

# 本地测试网

出于测试或开发目的，您可能需要运行本地测试网。

## 单节点测试网

**需求:**

- [安装iris](../get-started/install.md)

:::tip
对于以下示例，我们全部使用默认的[主目录](intro.md#主目录)
:::

### iris init

初始化genesis.json文件，它将帮助你启动网络

```bash
iris init --chain-id=testing --moniker=testing
```

### 创建一个钱包

创建一个钱包作为您的验证人帐户

```bash
iriscli keys add MyValidator
```

### iris add-genesis-account

将该钱包地址添加到genesis文件中的genesis.app_state.accounts数组中

:::tip
此命令使您可以设置通证数量。确保此帐户有iris，这是IRISnet上唯一的质押通证
:::

```bash
iris add-genesis-account $(iriscli keys show MyValidator --address) 100000000iris
```

### iris gentx

生成创建验证人的交易。gentx存储在`~/.iris/config/`中

```bash
iris gentx --name MyValidator
```

### iris collect-gentxs

将生成的质押交易添加到创世文件

```bash
iris collect-gentxs
```

### iris start

现在可以启动`iris`了

```bash
iris start
```

### iris unsafe-reset-all

可以使用此命令来重置节点，包括本地区块链数据库，地址簿文件，并将priv_validator.json重置为创世状态。

当本地区块链数据库以某种方式中断和无法同步或参与共识时，这是有用的。

```bash
iris unsafe-reset-all
```

### iris reset

与[iris unsafe-reset-all](#iris-unsafe-reset-all)不同，此命令允许将节点的区块链状态重置为指定的高度，因此可以更快地修复区块链数据库。

```bash
# e.g. reset the blockchain state to height 100
iris reset --height 100
```

还有一个修复区块链数据库的方式，如果在主网上出现 `Wrong Block.Header.AppHash`的错误，请确认您使用的是正确的[主网版本](../get-started/install.md#最新版本)，然后通过以下方式重新启动节点：

```bash
iris start --replay-last-block
```

### iris tendermint

查询可以在p2p连接中使用的唯一节点ID，例如在[config.toml](intro.md#cnofig-toml)中`seeds`和`persistent_peers`的格式`<node-id>@ip:26656`。

节点ID存储在[node_key.json](intro.md#node_key-json)中。

```bash
iris tendermint show-node-id
```

 查询[Tendermint Pubkey](../concepts/validator-faq.md#tendermint-密钥)，用于[identify your validator](../cli-client/stake.md#iriscli-stake-create-validator),并将用于在共识过程中签署Pre-vote/Pre-commit。

[Tendermint Key](../concepts/validator-faq.md#tendermint-密钥)存储在[priv_validator.json](intro.md#priv_validator-json)中，创建验证人后，请一定要记得[备份](../concepts/validator-faq.md#如何备份验证人节点)。

```bash
iris tendermint show-validator
```

查询bech32前缀验证人地址

```bash
iris tendermint show-address
```

### iris export

请参阅[导出区块状态](export.md)。

## 多节点测试网

**前提:**

- [安装 iris](../get-started/install.md)
- [安装 jq](https://stedolan.github.io/jq/download/)
- [安装 docker](https://docs.docker.com/engine/installation/)
- [安装 docker-compose](https://docs.docker.com/compose/install/)

### 构建和初始化

```bash
# Work from the irishub repo
cd $GOPATH/src/github.com/irisnet/irishub

# Build the linux binary in ./build
make build_linux

# Quick init a 4-node testnet configs
make testnet_init
```

`make testnet_init`将调用`iris testnet`命令在`build/nodecluster`目录下生成4个节点的测试网配置文件。

```bash
$ tree -L 3 build/nodecluster/
build/nodecluster/
├── gentxs
│   ├── node0.json
│   ├── node1.json
│   ├── node2.json
│   └── node3.json
├── node0
│   ├── iris
│   │   ├── config
│   │   └── data
│   └── iriscli
│       ├── key_seed.json
│       └── keys
├── node1
│   ├── iris
│   │   ├── config
│   │   └── data
│   └── iriscli
│       └── key_seed.json
├── node2
│   ├── iris
│   │   ├── config
│   │   └── data
│   └── iriscli
│       └── key_seed.json
└── node3
    ├── iris
    │   ├── config
    │   └── data
    └── iriscli
        └── key_seed.json
```

### 启动

```bash
make testnet_start
```

该命令将使用ubuntu:16.04的docker镜像创建4个节点的测试网。下表列出了每个节点的端口：

| Node      | P2P Port | RPC Port |
| --------- | -------- | -------- |
| irisnode0 | 26656    | 26657    |
| irisnode1 | 26659    | 26660    |
| irisnode2 | 26661    | 26662    |
| irisnode3 | 26663    | 26664    |

要更新二进制文件，只需重新构建它并重新启动节点即可：

```bash
make build_linux testnet_start
```

### 停止

停止所有正在运行的节点：

```bash
make testnet_stop
```

### 重置

要停止所有正在运行的节点并将网络重置为创世状态：

```bash
make testnet_unsafe_reset
```

### 清理

要停止所有正在运行的节点并删除`build/`目录中的所有文件：

```bash
make testnet_clean
```
