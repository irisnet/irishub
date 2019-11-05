---
order: 4
---

# 导出区块状态

## 简介

IRIShub可以导出区块链状态并将其输出为json格式的字符串，该字符串可用作新区块链的创世文件。

默认情况下，IRIShub存储每10,000个块和最后100个块的快照。可以从任何现有快照高度导出区块链状态。

如果要从不存在的快照高度导出状态，则需要首先将区块链状态[重置](local-testnet.md#iris-reset)到指定的高度。

## 用法

```bash
 iris export <flags>
```

## 标识

| 名称，速记        | 类型   | 必须 | 默认值       | 描述                                                                               |
| ----------------- | ------ | ---- | ------------ | ---------------------------------------------------------------------------------- |
| --for-zero-height | bool   |      | false        | 导出数据之前做一些清理性的工作，如果不想以导出的数据启动一条新链，可以不加这个标识 |
| --height          | uint   |      | 0            | 从指定的高度导出，默认值为0表示导出当前高度状态                                    |
| --home            | string |      | $HOME/.iris  | 指定存储配置和区块链数据的目录                                                     |
| --output-file     | string |      | genesis.json | 存储导出状态的文件                                                                 |

## 示例

导出当前的区块链状态

```bash
 iris export --home=<path-to-your-home>
```

从特定高度导出区块链状态，该高度必须是现有快照高度

```bash
iris export --height 10000 --home=<path-to-your-home>
```

如果想导出指定高度的区块链状态，并且以这个状态启动一条新链，可以尝试这个命令

```bash
iris export --height 10000 --for-zero-height --home=<path-to-your-home>
```
