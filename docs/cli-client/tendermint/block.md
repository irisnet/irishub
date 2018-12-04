# iriscli tendermint block

## Description

Get verified data for a the block at given height


## Usage

```
  iriscli tendermint block [height] [flags]


```

## Flags

| Name, shorthand | Default                    |Description                                                             | Required     |
| --------------- | -------------------------- | --------------------------------------------------------- | -------- |
| --chain-id    |     | Chain ID of Tendermint node   | yes     |
| --node string     |   tcp://localhost:26657                         | Node to connect to (default "tcp://localhost:26657")  |                                     
| --help, -h      |       | 	help for block|    |
| --trust-node    |              true         | Trust connected full node (don't verify proofs for responses)     |          |

## Examples

### tx

```shell
iriscli tendermint block 114263  --chain-id=fuxi-4000 --trust-node=true

```







