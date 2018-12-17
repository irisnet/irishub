# iriscli tendermint validator-set

## Description

Get the full tendermint validator set at given height. If no height is specified, the latest height will be used as default.


## Usage

```
  iriscli tendermint validator-set [height] [flags]
```
or
```
  iriscli tendermint validator-set
```

## Flags

| Name, shorthand | Default                    |Description                                                             | Required     |
| --------------- | -------------------------- | --------------------------------------------------------- | -------- |
| --chain-id    |     | Chain ID of Tendermint node   | yes     |
| --node string     |   tcp://localhost:26657                         | Node to connect to (default "tcp://localhost:26657")  |                                     
| --help, -h      |       | 	help for block|    |
| --trust-node    |              true         | Trust connected full node (don't verify proofs for responses)     |          |

## Examples

### Get validator-set at height 114360

```shell
 iriscli tendermint validator-set 114360 --chain-id=fuxi-4000 --trust-node=true

```

### Get the latest validator-set

```shell
 iriscli tendermint validator-set --chain-id=fuxi-4000 --trust-node=true

```






