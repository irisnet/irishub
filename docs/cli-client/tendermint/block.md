# iriscli tendermint block

## Description

Get verified data for a the block at given height. If no height is specified, the latest height will be used as default.


## Usage

```
  iriscli tendermint block [height] [flags]
```
or 
```
  iriscli tendermint block [flags]
```

## Flags

| Name, shorthand | Default                    |Description                                                             | Required     |
| --------------- | -------------------------- | --------------------------------------------------------- | -------- |
| --chain-id    |     | Chain ID of Tendermint node   | yes     |
| --node string     |   tcp://localhost:26657                         | Node to connect to (default "tcp://localhost:26657")  |                                     
| --help, -h      |       | 	help for block|    |
| --trust-node    |              true         | Trust connected full node (don't verify proofs for responses)     |          |

## Examples

### Get block at height 114263

```shell
iriscli tendermint block 114263  --chain-id=fuxi-4000 --trust-node=true

```

### Get the latest block

```shell
iriscli tendermint block  --chain-id=fuxi-4000 --trust-node=true

```







