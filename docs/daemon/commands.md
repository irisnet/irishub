---
order: 2
---

# Commands

## Introduction

IRIS Daemon Commands allow you to init, start, reset a node, or generate a genesis file, etc.

You can get familiar with these commands by creating a [Local Testnet](local-testnet.md).

## Usage

```bash
iris <command>
```

## Available Commands

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

## Global Flags

| Name,shorthand | Default      | Description                                        | Required | Type   |
| -------------- | ------------ | -------------------------------------------------- | -------- | ------ |
| -h, --help     |              | Help for iris                                      |          |        |
| --home         | /$HOME/.iris | Directory for config and data                      |          | String |
| --log_level    | \*:info      | Log level (default "main:info,state:info,*:error") |          | String |
| --trace        |              | Print out full stack trace on errors               |          |        |
