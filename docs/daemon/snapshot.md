---
order: 4
---

# Snapshot

## Introduction

IRIShub can snapshot the latest data of the full node, including blocks, consensus status, application status, and so on. For the block data, only the latest block information of the full node is reserved, and all the previous block information is discarded, so it is not suitable as a full node of the LCD connection. This feature is suitable for the following situations:

* Quickly start a new node and join the main network
  
* Free up disk space

## Usage

```bash
 iris snapshot <flags>
```

After the command is executed, a `data.bak` directory will be generated in the directory specified by `--tmp-dir`. Delete the old `data` directory, rename `data.bak` to `data`, and restart the node.

## Flags

| Name, shorthand | type   | Required | Default          | Description                                                        |
| --------------- | ------ | -------- | ---------------- | ------------------------------------------------------------------ |
| --tmp-dir       | string |          | Same as `--home` | Where the snapshot data is saved                                   |
| --home          | string |          | $HOME/.iris      | Specify the directory which stores node config and blockchain data |

## Examples

Snapshot current node's latest data

```bash
 iris snapshot --home=<path-to-your-home>
```
