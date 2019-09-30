---
order: 2
---

# 命令

## 描述

IRIS守护程序命令允许您初始化，启动，重置节点或生成创世文件等。

你可以通过创建[Local Testnet](local-testnet.md)来熟悉这些命令。

## 用法

```bash
iris <command>
```

## 可用命令

| Name                                                             | Description                                                                                                     |
| ---------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| [init](local-testnet.md#iris-init)                               | Initialize private validator, p2p, genesis, and application configuration files                                 |
| [add-genesis-account](local-testnet.md#iris-add-genesis-account) | Add genesis account to genesis.json                                                                             |
| [gentx](local-testnet.md#iris-gentx)                             | Generate a genesis tx carrying a self delegation                                                                |
| [collect-gentxs](local-testnet.md#iris-collect-gentxs)           | Collect genesis txs and output a genesis.json file                                                              |
| [start](local-testnet.md#iris-start)                             | Run the full node                                                                                               |
| [unsafe-reset-all](local-testnet.md#iris-unsafe-reset-all)       | Resets the blockchain database, removes address book files, and resets priv_validator.json to the genesis state |
| [tendermint](local-testnet.md#iris-tendermint)                   | Tendermint subcommands                                                                                          |
| [testnet](local-testnet.md#build-and-init)                       | Initialize files for a Irishub testnet                                                                          |
| [reset](local-testnet.md#iris-reset)                             | Reset app state to the specified height                                                                         |
| [export](export.md)                                              | Export state to JSON                                                                                            |
| version                                                          | Show executable binary version                                                                                  |

## 全局标识

| 名称，速记  | 默认值       | 描述                                               | 必须 | 类型   |
| ----------- | ------------ | -------------------------------------------------- | ---- | ------ |
| -h, --help  |              | Help for iris                                      |      |        |
| --home      | /$HOME/.iris | Directory for config and data                      |      | String |
| --log_level | \*:info      | Log level (default "main:info,state:info,*:error") |      | String |
| --trace     |              | Print out full stack trace on errors               |      |        |
