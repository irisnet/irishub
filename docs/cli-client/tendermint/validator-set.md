# iriscli tendermint validator-set

## Description

Get the full tendermint validator set at given height


## Usage

```
  Get the full tendermint validator set at given height


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
 iriscli tendermint validator-set 114360 --chain-id=fuxi-4000 --trust-node=true

```






