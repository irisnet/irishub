# Tendermint subcommand of iris and iriscli

## iriscli tendermint

### Description

Tendermint state querying subcommands

### Usage

```
iriscli tendermint <subcommand>
```

### Available Subcommands

| Name, shorthand | Description        |
| --------------- | -------------------------- |
| [tx](tx.md)     |  Matches this txhash over all committed blocks           |  
| [txs](txs.md)   | Search for all transactions that match the given tags  |                        
| [block](block.md)| 	Get verified data for a the block at given height    |   
| [validator-set](validator-set.md) | Get the full tendermint validator set at given height  |   


## iris tendermint

### Description

This command will return some tendermint related info to you. 

### Usage

```shell
iris tendermint <subcommand>
```

### Available Commands

| Name                    | Description                                                                                  |
| ----------------------- | -------------------------------------------------------------------------------------------- |
| [show-node-id](show-node-id.md) | show node id |
| [show-validator](show-validator.md) | show validator |
| [show-address](show-address.md) |     show address  |

### Global Flags

| Name, shorthand | Default        | Description                            | Required |
| --------------- | -------------- | -------------------------------------- | -------- |
| --encoding, -e  | hex            | Binary encoding (hex\b64\btc) |          |
| --home          | $HOME/.iris    | Directory for config and data |          |
| --output, -o    | text           | Output format (text\json)     |          |
| --trace         |                | Print out full stack trace on errors   |          |
