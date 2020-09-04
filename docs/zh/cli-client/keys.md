# Keys

Keys模块用于管理IRIS的Tendermint本地密钥库（钱包）。

## 可用命令

| 名称                               | 描述                                                         |
| ---------------------------------- | ------------------------------------------------------------ |
| [add](#iris-keys-add)              | 创建新密钥，加密后存入磁盘  |
| [delete](#iris-keys-delete)        | 删除指定的密钥                                                                         |
| [export](#iris-keys-export)        | 将密钥导出为json文件                                                              |
| [import](#iris-keys-import)        | 从本地钥匙库导入私钥                        |
| [list](#iris-keys-list)            | 列出所有密钥                                                                                |
| [migrate](#iris-keys-migrate)      | 从钥匙库将密钥信息迁移到另一个钥匙库                   |
| [mnemonic](#iris-keys-mnemonic)    | 从一些输入熵中计算出bip32助记词 |
| [parse](#iris-keys-parse)          | 将地址从hex解析为bech32，反之亦然                       |
| [show](#iris-keys-show)            | 按名称或地址检索密钥信息                                                             |

## iris keys add

获得一个密钥并写入磁盘。

```bash
iris keys add <key-name> [flags]
```

**标志：**

| 名称，速记           | 默认      | 描述                                               | 必须 |
| -------------------- | --------- | -------------------------------------------------- | ---- |
| --multisig           |           | 构造和存储multisig公钥                        |          |
| --multisig-threshold |     1     | 指定多签密钥最少签名数                                    |          |
| --nosort             |   false   | 传递给--multisig的密钥按其提供的顺序获取 |          |
| --pubkey             |           | 解析bech32格式的公钥并将其保存到磁盘           |          |
| --interactive        |   false   | 交互式提示用户输入BIP39密码和助记词       |          |
| --ledger             |   false   | 对本地设备上的密钥存储的引用       |          |
| --recover            |   false   | 提供种子短语来恢复现有密钥           |          |
| --no-backup          |   false   | 不要打印出种子短语（如果其他人正在观看终端） |          |
| --dry-run            |   false   | 执行操作，但不要向本地密钥库添加密钥               |          |
| --hd-path            |           | 手动HD路径导出（覆盖BIP44 配置）                |          |
| --coin-type          |    118    | HD衍生的硬币类型编号                                |          |
| --account            |     0     | 高清衍生账号                                 |          |
| --index              |     0     | HD派生的地址索引号                            |          |
| --algo               |  secp256k | 生成密钥的密钥签名算法                        |          |


### 创建密钥

```bash
iris keys add MyKey
```

执行该命令后输入并确认密码，将生成一个新的密钥。密码至少8个字符。

:::warning
**重要**

写下助记词并保存在安全的地方！如果你不慎忘记密码或丢失了密钥，这是恢复账户的唯一方法。
:::

### 通过助记词恢复密钥

如果你忘记了密码或丢失了密钥，或者你想在其他地方使用密钥，可以通过助记词短语来恢复密钥。

```bash
iris keys add MyKey --recover
```

系统会要求你输入密钥的新密码并确认，然后输入助记词。然后你将得到恢复的密钥。

```bash
Enter a passphrase for your key:
Repeat the passphrase:
Enter your recovery seed phrase:
```

### 创建多签密钥

以下例子为创建一个包含3个子密钥的多签密钥，且指定最小签名数为2。只有交易签名数大于等于2时，该交易才会被广播。

```bash
iris keys add <multisig-keyname> --multisig-threshold=2 --multisig=<signer-keyname-1>,<signer-keyname-2>,<signer-keyname-3>
```

:::tip
`<signer-keyname>` 可以为“local/offline/ledger”类型，但不允许为“multi”类型。

如果你没有子密钥的所有许可，则可先请求获取pubkeys并以此创建“offline”密钥，然后你将可以创建多签密钥。

其中“offline”类型密钥可以通过“iris keys add --pubkey”命令创建。
:::

如何使用多签密钥签名和广播交易，请参阅 [multisig](tx.md#iris-tx-multisign)

## iris keys delete

根据名称删除本地密钥

```bash
iris keys delete <key-name> [flags]
```

**标志：**

| 名称, 速记  | 默认  | 描述                                   | 必须 |
| ----------- | ----- | -------------------------------------- | ---- |
| --force, -f | false | 无密码强制删除秘钥                     |      |
| --yes, -y   | false | 删除离线和ledger引用密钥时跳过确认提示 |      |

### 删除一个本地密钥

```bash
iris keys delete MyKey
```

## iris keys export

将密钥keystore信息导出为json文件。

```bash
iris keys export <key-name> [flags]
```

### 导出keystore

```bash
iris keys export Mykey --output-file=<path-to-keystore>
```

## iris keys import

将密钥导入到本地钥匙库

### 将密钥导入到本地钥匙库

```bash
iris keys import <name> <keyfile> [flags]
```

## iris keys list

返回此密钥管理器存储的所有密钥的名称、类型、地址和公钥列表。

**标志：**

| 名称，速记        | 默认     | 描述                  | 必须 |
| --------------- | ------- | -------------------- | -------- |
| --list-name     |         | List names only |          |

### 列出所有密钥

```bash
iris keys list
```

## iris keys migrate

将钥匙库中的钥匙信息迁移到新钥匙库

**标志：**

| 名称，速记        | 默认     | 描述                                                                     | 必须 |          
| --------------- | ------- | ------------------------------------------------------------------------ | -------- |
| --dry-run       |         | Run migration without actually persisting any changes to the new Keybase |          |

### 迁移钥匙信息

```bash
iris keys migrate [flags]
```

## iris keys mnemonic

通过读取系统熵来创建24个单词组成的bip39助记词（也称为种子短语）。如果需要传递自定义的熵，请使用`unsafe-entropy`模式。

```bash
iris keys mnemonic [flags]
```

**标志：**

| 名称, 速记       | 默认 | 描述                                     | 必须 |
| ---------------- | ---- | ---------------------------------------- | ---- |
| --unsafe-entropy |      | 提示用户提供自定义熵，而不是通过系统生成 |      |

### 创建助记词

```bash
iris keys mnemonic
```

执行上述命令后你将得到24个单词组成的助记词，例如：

```bash
beauty entire blue tape ordinary fix rotate learn smart tiger dolphin cycle cigar dish alcohol slab bachelor vital design consider paper panther mad eternal
```

## iris keys parse

将十六进制的密钥地址和指纹转换为标准输出并打印为bech32 cosmos前缀格式，反之亦然。

### 格式化打印密钥和指纹

```bash
iris keys parse <hex-or-bech32-address> [flags]
```

## iris keys show

查询本地密钥的详细信息。

```bash
iris keys show <key-name> [flags]
```

**标志：**

| 名称, 速记           | 默认 | 描述                                 | 必须 |
| -------------------- | ---- | ------------------------------------ | ---- |
| --address            |false | 仅输出地址（覆盖 --output）          |      |
| --bech               | acc  | 在账本设备中输出地址                 |      |
| --device             |false | 查询命令帮助                         |      |
| --multisig-threshold | 1    | N个必需签名中的K个       |      |
| --pubkey             |false | 仅输出公钥 覆盖 --output）           |      |

### 查询指定密钥

```bash
iris keys show MyKey
```

执行命令后将会显示以下信息：

```bash
- name: Mykey
  type: local
  address: iaa1tulwx2hwz4dv8te6cflhda64dn0984harlzegw
  pubkey: iap1addwnpepq24rufap6u0sysqcpgsfzqhw3x8nfkhqhtmpgqt0369rlyqcg0vkgwzc4k0
  mnemonic: ""
  threshold: 0
  pubkeys: []
```

### 查询验证人operator地址

如果某个地址已绑定成为验证人operator（用于创建验证人的地址），则可以使用`--bech val`获取其`iva`前缀的operator地址和`ivp`前缀的公钥地址：

```bash
iris keys show MyKey --bech val
```

示例输出：

```bash
- name: Mykey
  type: local
  address: iva1tulwx2hwz4dv8te6cflhda64dn0984hakwgk4f
  pubkey: ivp1addwnpepq24rufap6u0sysqcpgsfzqhw3x8nfkhqhtmpgqt0369rlyqcg0vkgd8e6zy
  mnemonic: ""
  threshold: 0
  pubkeys: []
```
