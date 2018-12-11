# iriscli tendermint

## Description

Tendermint state querying subcommands

## Usage

```
iriscli tendermint [command]
```

## Available Commands

| Name, shorthand | Description        |
| --------------- | -------------------------- |
| [tx](tx.md)            |  Matches this txhash over all committed blocks           |  
| txs   | Search for all transactions that match the given tags  |                            
| [block](block.md)| 	Get verified data for a the block at given height    |   
| [validator-set](validator-set.md) | Get the full tendermint validator set at given height  |   

## Flags

|Name, shorthand|Description|
|---			|---		|
|--help,-h		|help for tendermint|


# iris tendermint

## Description

This command will return some Tendermint related info to you. 

## Usage

```shell
iris tendermint [command]
```

## Available Commands

| Name                    | Description                                                                                  |
| ----------------------- | -------------------------------------------------------------------------------------------- |
| [show-node-id](show-node-id.md) | show node id |
| [show-validator](show-validator.md) | show validator |
| [show-address](show-address.md)           |     show address                                                 |

## Flags

| Name, shorthand | Default | Description   | Required |
| --------------- | ------- | ------------- | -------- |
| --help, -h      |         | Help for keys |          |

## Global Flags

| Name, shorthand | Default        | Description                            | Required |
| --------------- | -------------- | -------------------------------------- | -------- |
| --encoding, -e  | hex            | [string] Binary encoding (hex|b64|btc) |          |
| --home          | $HOME/.iriscli | [string] Directory for config and data |          |
| --output, -o    | text           | [string] Output format (text|json)     |          |
| --trace         |                | Print out full stack trace on errors   |          |
