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

:::tip
Please stop your node before executing the command.
:::

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

## FAQ

### What if every node in the network runs from a snapshot and a new node tries to catch up from genesis

If that happens, a new node will be unable to sync from scratch, but it can use a snapshot to catch up faster.

Presumably not all nodes will delete historical data, such as nodes for explorers and wallets. And we, the IRIS Foundation will keep all the data too, also we can offer a full data snapshot. And we encourage snapshot service providers could provide both the minimal and the full snapshot :)

### But this will have a slow speed of download probably

Correct, but most people don't need the full data, they can download a latest snapshot to sync up much faster than before. If they want, they can also download the full data snapshot too.

### Can I snapshot the validator node

Yes, but we wouldn't recommend that for the time being.