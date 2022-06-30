# NFT

`NFT`提供了将资产进行数字化的能力。通过该模块，每个链外资产将被建模为唯一的链上资产。

## 可用命令

| 名称                                          | 描述              |
| --------------------------------------------- | ----------------- |
| [issue](#iris-tx-nft-issue)                   | 发行资产          |
| [transfer-denom](#iris-tx-nft-transfer-denom) | 转移nft分类所有权 |
| [mint](#iris-tx-nft-mint)                     | 增发资产          |
| [edit](#iris-tx-nft-edit)                     | 编辑资产          |
| [transfer](#iris-tx-nft-transfer)             | 转让资产          |
| [burn](#iris-tx-nft-burn)                     | 销毁资产          |
| [supply](#iris-query-nft-supply)              | 查询supply        |
| [owner](#iris-query-nft-owner)                | 通过owner查询     |
| [collection](#iris-query-nft-collection)      | 查询collection    |
| [denom](#iris-query-nft-denom)                | 查询denom         |
| [denoms](#iris-query-nft-denoms)              | 查询denoms        |
| [token](#iris-query-nft-token)                | 查询token         |

## iris tx nft issue

发行资产

```bash
iris tx nft issue [denom-id] [flags]
```

**标志：**

| 名称，速记          | 默认 | 描述                                                                                                           | 必须 |
| ------------------- | ---- | -------------------------------------------------------------------------------------------------------------- | ---- |
| --name              |      | denom名字                                                                                                      |      |
| --uri               |      | 用于补充的off-chain元数据的URI(JSON对象)                                                                       |      |
| --schema            |      | denom数据结构定义                                                                                              |      |
| --data              |      | 用于补充的off-chain元数据(JSON对象)                                                                            |      |
| --symbol            |      | 分类的简短名称                                                                                                 |      |
| --mint-restricted   |      | 此字段表示在此分类下是否有发行NFT的限制，true表示只有Denom的拥有者可以在此分类下发行NFT，false表示任何人可以   |      |
| --update-restricted |      | 此字段表示在此分类下是否有更新NFT的限制，true表示此分类下任何人不得更新NFT，false表示只有此NFT的拥有者可以更新 |      |

## iris tx nft transfer-denom

转移nft分类所有权

```bash
iris tx nft transfer-denom [recipient] [denom-id]
```

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
