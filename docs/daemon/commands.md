# Commands

## Description

IRIS Daemon Commands allow you to init, start, reset a node, or generate a genesis file, etc.

## Usage

```bash
iris <command>
```

## Available Commands

| Name                                        | Description                                                  |
| ------------------------------------------- | ------------------------------------------------------------ |
| [init](#init)                               | Initialize private validator, p2p, genesis, and application configuration files |
| [gentx](#gentx)                             | Generate a genesis tx carrying a self delegation             |
| [add-genesis-account](#add-genesis-account) | Add genesis account to genesis.json                          |
| [testnet](#testnet)                         | Initialize files for a Irishub testnet                       |
| [collect-gentxs](#collect-gentxs)           | Collect genesis txs and output a genesis.json file           |
| [start](#start)                             | Run the full node                                            |
| [unsafe-reset-all](#unsafe-reset-all)       | Resets the blockchain database, removes address book files, and resets priv_validator.json to the genesis state |
| [tendermint](#tendermint)                   | Tendermint subcommands                                       |
| [reset](#reset)                             | Reset app state to the specified height                      |
| [export](#export)                           | Show executable binary version                               |
| [version](#version)                         | Query the asset related fees                                 |

## Global Flags

| Name,shorthand | Default         | Description                                        | Required | Type   |
| -------------- | --------------- | -------------------------------------------------- | -------- | ------ |
| -h, --help     |                 | Help for iris                                      | False    |        |
| --home         | /$HOME/.iris    | Directory for config and data                      | False    | String |
| --log_level    | \*:info         | Log level (default "main:info,state:info,*:error") | False    | String |
| --trace        |                 | Print out full stack trace on errors               | False    |        |
