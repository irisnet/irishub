# NFT

## 简介

`NFT`提供了将资产进行数字化的能力。通过该模块，每个链外资产将被建模为唯一的链上资产。

链上资产用 `ID` 进行标识，借助区块链安全、不可篡改的特性，资产的所有权将得到明确。资产在成员间的交易过程也将被公开地记录，以便于追溯以及争议处理。

资产的元数据（`metadata`）可以直接存储在链上，也可以将其在链外存储源的 `URI` 存储在链上。资产元数据按照特定的 [JSON Schema](https://JSON-Schema.org/) 进行组织。[这里](./nft-metadata.json)是一个元数据 JSON Schema 示例。

资产在创建前需要发行，用以声明其抽象属性：

- _Denom_：即全局唯一的资产类别标识符

- _元数据规范_：资产元数据应遵循的 JSON Schema

每一个具体的资产由以下元素描述：

- _Denom_: 该资产的类别

- _ID_：资产的标识符，在此资产类别中唯一；此 ID 在链外生成

- _元数据_：包含资产具体数据的结构

- _元数据 URI_：当元数据存储在链外时，此 URI 表示其存储位置

## 功能

### 发行

指定资产 Denom（资产类别）、元数据 JSON Schema，即可发行资产。

`CLI`

```bash
iris tx nft issue <denom> --schema=<schema-content or path/to/schema.json>
```

### 增发

在发行资产之后即可增发（创建）该类型的具体资产。需指定资产 ID、接收者地址、元数据或其 URI。

`CLI`

```bash
iris tx nft mint <denom> <token-id> --recipient=<recipient-address> --uri=<token-uri> --data=<token-data>
```

### 编辑

可对指定资产的元数据进行更新。

`CLI`

```bash
iris tx nft edit <denom> <token-id> --uri=<token-uri> --data=<token-data>
```

### 转移

转移指定资产。

`CLI`

```bash
iris tx nft transfer <recipient-address> <denom> <token-id>
```

### 销毁

可以销毁已创建的资产。

`CLI`

```bash
iris tx nft burn <denom> <token-id>
```

### 查询指定的资产类别

根据 Denom 查询资产类别信息。

`CLI`

```bash
iris q nft denom <denom>
```

### 查询所有资产类别信息

查询已发行的所有资产类别信息。

`CLI`

```bash
iris q nft denoms
```

### 查询指定类别资产的总量

根据 Denom 查询资产总量；接受可选的 owner 参数。

`CLI`

```bash
iris q nft supply <denom> --owner=<owner>
```

### 查询指定账户的所有资产

查询某一账户所拥有的全部资产；可以指定 Denom 参数。

`CLI`

```bash
iris q nft owner <address> --denom=<denom>
```

### 查询指定类别的所有资产

根据 Denom 查询所有资产。

`CLI`

```bash
iris q nft collection <denom>
```

### 查询指定资产

根据 Denom 以及 ID 查询具体资产。

`CLI`

```bash
iris q nft token <denom> <token-id>
```
