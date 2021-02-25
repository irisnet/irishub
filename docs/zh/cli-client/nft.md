# NFT

`NFT`提供了将资产进行数字化的能力。通过该模块，每个链外资产将被建模为唯一的链上资产。

## 可用命令

| 名称                                     | 描述           |
| ---------------------------------------- | -------------- |
| [issue](#iris-tx-nft-issue)              | 发行资产       |
| [mint](#iris-tx-nft-mint)                | 增发资产       |
| [edit](#iris-tx-nft-edit)                | 编辑资产       |
| [transfer](#iris-tx-nft-transfer)        | 转让资产       |
| [burn](#iris-tx-nft-burn)                | 销毁资产       |
| [supply](#iris-query-nft-supply)         | 查询supply     |
| [owner](#iris-query-nft-owner)           | 通过owner查询  |
| [collection](#iris-query-nft-collection) | 查询collection |
| [denom](#iris-query-nft-denom)           | 查询denom      |
| [denoms](#iris-query-nft-denoms)         | 查询denoms     |
| [token](#iris-query-nft-token)           | 查询token      |

## iris tx nft issue

发行资产

```bash
iris tx nft issue [denom-id] [flags]
```

**标志：**

| 名称，速记 | 默认 | 描述              | 必须 |
| ---------- | ---- | ----------------- | ---- |
| --name     |      | denom名字         |      |
| --schema   |      | denom数据结构定义 |      |

## iris tx nft mint

增发资产

```bash
iris tx nft mint [denomID] [tokenID] [flags]
```

**标志：**

| 名称，速记  | 默认 | 描述               | 必须 |
| ----------- | ---- | ------------------ | ---- |
| --uri       |      | 链下token数据的URI |      |
| --recipient |      | nft接受者          |      |
| --name      |      | nft名字            |      |

## iris tx nft edit

编辑资产

```bash
iris tx nft edit [denomID] [tokenID] [flags]
```

**标志：**

| 名称，速记 | 默认 | 描述               | 必须 |
| ---------- | ---- | ------------------ | ---- |
| --uri      |      | 链下token数据的URI |      |
| --name     |      | nft名字            |      |

## iris tx nft transfer

转让资产

```bash
iris tx nft transfer [recipient] [denomID] [tokenID] [flags]
```

**标志：**

| 名称，速记 | 默认 | 描述               | 必须 |
| ---------- | ---- | ------------------ | ---- |
| --uri      |      | 链下token数据的URI |      |
| --name     |      | nft名字            |      |

## iris tx nft burn

销毁资产

```bash
iris tx nft burn [denomID] [tokenID] [flags]
```

## iris query nft

查询资产

### iris query nft supply

```bash
iris query nft supply [denomID]
iris query nft supply [denomID] --owner=<owner address>
```

### iris query nft owner

```bash
iris query nft owner [owner address]
iris query nft owner [owner address] --denom=<denomID>
```

### iris query nft collection

```bash
iris query nft collection [denomID]
```

### iris query nft denom

```bash
iris query nft denom [denomID]
```

### iris query nft denoms

```bash
iris query nft denoms
```

### iris query nft token

```bash
iris query nft token [denomID] [tokenID]
```
