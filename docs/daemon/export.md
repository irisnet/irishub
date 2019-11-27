---
order: 4
---

# Export Blockchain State

## Introduction

IRIShub can export the blockchain state and output to a json-format string which can be used as the genesis file of a new blockchain.

By default, IRIShub stores snapshots of every 10,000 blocks and the last 100 blocks. You can export the blockchain state from any existing snapshot height.

If you want to export the state from a nonexisting snapshot height, you need to [reset](local-testnet.md#iris-reset) the blockchain state to the specified height first.

## Usage

```bash
 iris export <flags>
```

## Flags

| Name, shorthand   | type   | Required | Default      | Description                                                  |
| ----------------- | ------ | -------- | ------------ | ------------------------------------------------------------ |
| --for-zero-height | bool   |          | false        | Do some clean up work before exporting state. If you want to use the exported state to start a new blockchain, please add this flag. Otherwise, just leave out it |
| --height          | uint   |          | 0            | Export state from a particular height, default value is 0 which means to export the latest state |
| --home            | string |          | $HOME/.iris  | Specify the directory which stores node config and blockchain data |
| --output-file     | string |          | genesis.json | Target file to save exported state                           |

## Examples

Export the current blockchain state

```bash
 iris export --home=<path-to-your-home>
```

Export blockchain state from a particular height, the height must be an existing snapshot height

```bash
iris export --height 10000 --home=<path-to-your-home>
```

If you want to export the blockchain state from a particular height and use the exported state as genesis state of another blockchain

```bash
iris export --height 10000 --for-zero-height --home=<path-to-your-home>
```
