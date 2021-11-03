# NFT

## 简介

`NFT`提供了将资产进行数字化的能力。通过该模块，每个链外资产将被建模为唯一的链上资产。

链上资产用 `ID` 进行标识，借助区块链安全、不可篡改的特性，资产的所有权将得到明确。资产在成员间的交易过程也将被公开地记录，以便于追溯以及争议处理。

资产的元数据（`metadata`）可以直接存储在链上，也可以将其在链外存储源的 `URI` 存储在链上。资产元数据按照特定的 [JSON Schema](https://JSON-Schema.org/) 进行组织。[这里](https://github.com/irisnet/irishub/blob/master/docs/zh/features/nft-metadata.json)是一个元数据 JSON Schema 示例。

资产在创建前需要发行，用以声明其抽象属性：

- _Denom_：即全局唯一的资产类别名
  
- _Denom ID_：Demon的全局唯一标识符 
  
- _Symbol_: 分类的简短名称

- _Mint-restricted_: 表示此分类下是否有发行NFT的限制，true表示只有Denom的拥有者可以在此分类下发行NFT，false表示任何人可以

- _Update-restricted_: 表示此分类下是否有更新NFT的限制，true表示此分类下任何人不得更新NFT，false表示只有此NFT的拥有者可以更新

- _元数据规范_：资产元数据应遵循的 JSON Schema

每一个具体的资产由以下元素描述：

- _Denom_：该资产的类别

- _ID_：资产的标识符，在此资产类别中唯一；此 ID 在链外生成

- _元数据_：包含资产具体数据的结构

- _元数据 URI_：当元数据存储在链外时，此 URI 表示其存储位置

## 功能

### 发行

指定资产 Denom ID（资产类别ID）、元数据 JSON Schema，即可发行资产。

`CLI`

```bash
iris tx nft issue <denom-id> --from=<key-name> --schema=<schema-content or path/to/schema.json> --symbol=<denom-symbol> --mint-restricted=<mint-restricted>  --update-restricted=<update-restricted> --chain-id=<chain-id> --fees=<fee>
```

### 转让NFT分类所有权

NFT分类拥有者可以转移NFT分类的所有权

`CLI`

```bash
iris tx nft transfer-denom <recipient> <denom-id>
```

### 增发

在发行资产之后即可增发（创建）该类型的具体资产。需指定资产 ID、接收者地址和URI。

`CLI`

```bash
iris tx nft mint <denom-id> <token-id> --uri=<uri> --recipient=<recipient> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

### 编辑

可对指定资产的元数据进行更新。

`CLI`

```bash
iris tx nft edit <denom-id> <token-id> --uri=<uri> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

### 转移

转移指定资产。

`CLI`

```bash
iris tx nft transfer <recipient-address> <denom-id> <token-id>
```

### 销毁

可以销毁已创建的资产。

`CLI`

```bash
iris tx nft burn <denom-id> <token-id> --from=<key-name> --chain-id=<chain-id> --fees=<fee>
```

### 查询指定的资产类别

根据 Denom ID查询资产类别信息。

`CLI`

```bash
iris q nft denom <denom-id>
```

### 查询所有资产类别信息

查询已发行的所有资产类别信息。

`CLI`

```bash
iris q nft denoms
```

### 查询指定类别资产的总量

根据 Denom ID查询资产总量；接受可选的 owner 参数。

`CLI`

```bash
iris q nft supply <denom-id> --owner=<owner>
```

### 查询指定账户的所有资产

查询某一账户所拥有的全部资产；可以指定 Denom ID参数。

`CLI`

```bash
iris q nft owner <address> --denom-id=<denom-id>
```

### 查询指定类别的所有资产

根据 Denom ID查询所有资产。

`CLI`

```bash
iris q nft collection <denom-id>
```

### 查询指定资产

根据 Denom ID以及 ID 查询具体资产。

`CLI`

```bash
iris q nft token <denom-id> <token-id>
```
