# iriscli tx

签名或广播交易

## 可用命令

| 名称                               | 描述                     |
| ---------------------------------- | ------------------------ |
| [sign](#iriscli-tx-sign)           | 签名生成的离线交易文件   |
| [broadcast](#iriscli-tx-broadcast) | 广播一个已签名交易到网络 |
| [multisig](#iriscli-tx-multisign)  | 用多个账户为同一交易签名 |

## iriscli tx sign

签名生成的离线交易文件。该文件由`generate-only`标志生成。

```bash
iriscli tx sign <file> <flags>
```

### 标志

| 命令，速记   | 类型   | 必须 | 默认  | 描述                                           |
| ------------ | ------ | ---- | ----- | ---------------------------------------------- |
| -h，--help   |        |      |       | 打印帮助                                       |
| --append     | bool   | 是   | true  | 将签名附加到现有签名                           |
| --name       | string | 是   |       | 用于签名的私钥名称                             |
| --offline    | bool   | 是   | false | 离线模式                                       |
| --print-sigs | bool   | 是   | false | 打印必须签署交易的地址和已签名的地址，然后退出 |

### 生成离线交易

:::tip
任意类型的离线交易都可以通过追加`--generate-only`标志来生成
:::

下面示例中使用Transfer交易：

```bash
iriscli bank send --to=<address> --amount=10iris --from=<key-name> --fee=0.3iris --chain-id=irishub --generate-only > unsigned.json
```

`unsigned.json` 看起来是这样的：

```json
{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"iaa19aamjx3xszzxgqhrh0yqd4hkurkea7f646vaym","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"4000000000000000"}],"gas":"200000"},"signatures":null,"memo":""}}
```

### 签名离线交易

```bash
iriscli tx sign unsigned.json --name=<key-name> > signed.tx
```

`signed.json` 看起来是这样的：

```json
{"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/Send","value":{"inputs":[{"address":"iaa106nhdckyf996q69v3qdxwe6y7408pvyvyxzhxh","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}],"outputs":[{"address":"iaa1893x4l2rdshytfzvfpduecpswz7qtpstevr742","coins":[{"denom":"iris-atto","amount":"10000000000000000000"}]}]}}],"fee":{"amount":[{"denom":"iris-atto","amount":"40000000000000000"}],"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"Auouudrg0P86v2kq2lykdr97AJYGHyD6BJXAQtjR1gzd"},"signature":"sJewd6lKjma49rAiGVfdT+V0YYerKNx6ZksdumVCvuItqGm24bEN9msh7IJ12Sil1lYjqQjdAcjVCX/77FKlIQ==","account_number":"0","sequence":"3"}],"memo":"test"}}
```

签名之后，`signed.json`中的`signature`字段将不再为空。

现在准备[广播这个已签名交易](#iriscli-tx-broadcast)到IRIS Hub。

## iriscli tx broadcast

这个命令用于广播已离线签名的交易到网络。

### 广播离线签名的交易

```bash
iriscli tx broadcast signed.json --chain-id=irishub
```

## iriscli tx multisign

用多个账户为一个交易签名。这个交易只有在签名数满足multisig-threshold时才可以广播。

```bash
iriscli tx multisign <file> <key-name> <[signature]...> <flags>
```

### 用多签密钥创建离线交易

:::tip
没有多签密钥？[创建一个](keys.md#创建多签密钥)
:::

```bash
iriscli bank send --to=<address> --amount=10iris --fee=0.3iris --chain-id=irishub --from=<multisig-keyname> --generate-only > unsigned.json
```

### 签名多签交易

#### 查询多签地址

```bash
iriscli keys show <multisig-keyname>
```

#### 签名`unsigned.json`

假定multisig-threshold是2，我们使用2个签名者签名`unsigned.json`

用signer-1签名：

```bash
iriscli tx sign unsigned.json --name=<signer-keyname-1> --chain-id=irishub --multisig=<multisig-address> --signature-only > signed-1.json
```

用signer-2签名：

```bash
iriscli tx sign unsigned.json --name=<signer-keyname-2> --chain-id=irishub --multisig=<multisig-address> --signature-only > signed-2.json
```

#### 合并签名

合并所有签名到 `signed.json`

```bash
iriscli tx multisign --chain-id=irishub unsigned.json <multisig-keyname> signed-1.json signed-2.json > signed.json
```

现在可以[广播](#iriscli-tx-broadcast)这个已签名交易了。
