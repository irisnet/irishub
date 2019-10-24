# iriscli keys

Keys模块用于管理IRIS的Tendermint本地密钥库（钱包）。

## 可用命令

| 名称                               | 描述                                                         |
| ---------------------------------- | ------------------------------------------------------------ |
| [add](#iriscli-keys-add)           | 创建新密钥、从助记词导入已有密钥，或从备份的keystore导入秘钥 |
| [list](#iriscli-keys-list)         | 列出所有密钥                                                 |
| [show](#iriscli-keys-show)         | 根据给定的名称显示密钥信息                                   |
| [export](#iriscli-keys-export)     | 将密钥导出为json文件                                         |
| [delete](#iriscli-keys-delete)     | 删除指定的密钥                                               |
| [update](#iriscli-keys-update)     | 更改用于保护私钥的密码                                       |
| [mnemonic](#iriscli-keys-mnemonic) | 通过读取系统熵来创建bip39助记词，也可以称为种子短语          |
| [new](#iriscli-keys-new)           | 使用交互式命令派生新的私钥，该命令将对你的每个输入作出提示   |

## iriscli keys add

创建一个新的密钥（钱包），或通过助记词/密钥库导入已有密钥。

```bash
iriscli keys add <key-name> <flags>
```

**标志：**

| 名称，速记           | 默认      | 描述                                               | 必须 |
| -------------------- | --------- | -------------------------------------------------- | ---- |
| --account            |           | HD推导的账号                                       |      |
| --dry-run            |           | 本地模拟交易，不会向本地密钥库中添加密钥           |      |
| --help, -h           |           | 查询命令帮助                                       |      |
| --index              |           | HD推导的索引号                                     |      |
| --ledger             |           | 在Ledger设备上存储私钥的本地引用                   |      |
| --no-backup          |           | 不输出助记词（如果有其他人正在观察你在终端的操作） |      |
| --recover            |           | 提供助记词以恢复现有密钥而不是新建                 |      |
| --keystore           |           | 从备份的keystore导入秘钥                           |      |
| --multisig           |           | 创建多签密钥                                       |      |
| --multisig-threshold |           | 指定多签密钥最少签名数                             |      |
| --type, -t           | secp256k1 | 私钥类型（secp256k\ed25519）                       |      |

### 创建密钥

```bash
iriscli keys add MyKey
```

执行该命令后输入并确认密码，将生成一个新的密钥。密码至少8个字符。

:::warning
**重要**

写下助记词并保存在安全的地方！如果你不慎忘记密码或丢失了密钥，这是恢复账户的唯一方法。
:::

### 通过助记词恢复密钥

如果你忘记了密码或丢失了密钥，或者你想在其他地方使用密钥，可以通过助记词短语来恢复密钥。

```bash
iriscli keys add MyKey --recover
```

系统会要求你输入密钥的新密码并确认，然后输入助记词。然后你将得到恢复的密钥。

```bash
Enter a passphrase for your key:
Repeat the passphrase:
Enter your recovery seed phrase:
```

### 通过keystore导入秘钥

```bash
iriscli keys add Mykey --recover --keystore=<path-to-keystore>
```

### 创建多签密钥

以下例子为创建一个包含3个子密钥的多签密钥，且指定最小签名数为2。只有交易签名数大于等于2时，该交易才会被广播。

```bash
iriscli keys add <multisig-keyname> --multisig-threshold=2 --multisig=<signer-keyname-1>,<signer-keyname-2>,<signer-keyname-3>
```

:::tip
`<signer-keyname>` 可以为“local/offline/ledger”类型，但不允许为“multi”类型。

如果你没有子密钥的所有许可，则可先请求获取pubkeys并以此创建“offline”密钥，然后你将可以创建多签密钥。

其中“offline”类型密钥可以通过“iriscli keys add --pubkey”命令创建。
:::

如何使用多签密钥签名和广播交易，请参阅 [multisig](tx.md#iriscli-tx-multisign)

## iriscli keys list

返回此密钥管理器存储的所有密钥的名称、类型、地址和公钥列表。

### 列出所有密钥

```bash
iriscli keys list
```

## iriscli keys show

查询本地密钥的详细信息。

```bash
iriscli keys show <key-name> <flags>
```

**标志：**

| 名称, 速记           | 默认 | 描述                                 | 必须 |
| -------------------- | ---- | ------------------------------------ | ---- |
| --address            |      | 仅输出地址（覆盖 --output）          |      |
| --bech               | acc  | 密钥的Bech32前缀编码（acc/val/cons） |      |
| --help, -h           |      | 查询命令帮助                         |      |
| --multisig-threshold | 1    | K out of N required signatures       |      |
| --pubkey             |      | 仅输出公钥 覆盖 --output）           |      |

### 查询指定密钥

```bash
iriscli keys show MyKey
```

执行命令后将会显示以下信息：

```bash
NAME:    TYPE:    ADDRESS:                                      PUBKEY:
MyKey    local    iaa1kkm4w5pvmcw0e3vjcxqtfxwqpm3k0zak83e7nf    iap1addwnpepq0gsl90v9dgac3r9hzgz53ul5ml5ynq89ax9x8qs5jgv5z5vyssskzc7exa
```

### 查询验证人operator地址

如果某个地址已绑定成为验证人operator（用于创建验证人的地址），则可以使用`--bech val`获取其`iva`前缀的operator地址和`ivp`前缀的公钥地址：

```bash
iriscli keys show MyKey --bech val
```

示例输出：

```bash
NAME:    TYPE:    ADDRESS:                                      PUBKEY:
MyKey    local    iva12nda6xwpmp000jghyneazh4kkgl2tnzyx7trze    ivp1addwnpepqfw52vyzt9xgshxmw7vgpfqrey30668g36f9z837kj9dy68kn2wxqm8gtmk
```

## iriscli keys export

将密钥keystore信息导出为json文件。

```bash
iriscli keys export <key-name> <flags>
```

**标志：**

| 名称, 速记    | 默认 | 描述                 | 必须 |
| ------------- | ---- | -------------------- | ---- |
| --output-file |      | keystore文件导出路径 |      |

### 导出keystore

```bash
iriscli keys export Mykey --output-file=<path-to-keystore>
```

## iriscli keys delete

根据名称删除本地密钥

```bash
iriscli keys delete <key-name> <flags>
```

**标志：**

| 名称, 速记  | 默认  | 描述                                   | 必须 |
| ----------- | ----- | -------------------------------------- | ---- |
| --force, -f | false | 无密码强制删除秘钥                     |      |
| --yes, -y   | false | 删除离线和ledger引用密钥时跳过确认提示 |      |

### 删除一个本地密钥

```bash
iriscli keys delete MyKey
```

## iriscli keys update

更改用于保护私钥的密码。

### 修改本地密钥的密码

```bash
iriscli keys update MyKey
```

## iriscli keys mnemonic

通过读取系统熵来创建24个单词组成的bip39助记词（也称为种子短语）。如果需要传递自定义的熵，请使用`unsafe-entropy`模式。

```bash
iriscli keys mnemonic <flags>
```

**标志：**

| 名称, 速记       | 默认 | 描述                                     | 必须 |
| ---------------- | ---- | ---------------------------------------- | ---- |
| --unsafe-entropy |      | 提示用户提供自定义熵，而不是通过系统生成 |      |

### 创建助记词

```bash
iriscli keys mnemonic
```

执行上述命令后你将得到24个单词组成的助记词，例如：

```bash
police possible oval milk network indicate usual blossom spring wasp taste canal announce purpose rib mind river pet brown web response sting remain airport
```

## iriscli keys new

:::warning
**Deprecated**
:::

使用交互式命令派生新的私钥，该命令将对每个输入进行提示。

指定bip39助记词、用于保护助记词的bip39密码、bip32 HD路径，以派生特定帐户。 密钥将以给定名称存储并用给定密码加密。 唯一需要的输入是加密密码。

```bash
iriscli keys new <key-name> <flags>
```

**标志：**

| 名称, 速记   | 默认            | 描述                             | 必须 |
| ------------ | --------------- | -------------------------------- | ---- |
| --bip44-path | 44'/118'/0'/0/0 | 从中导出私钥的BIP44路径          |      |
| --default    |                 | 跳过提示，只使用默认值           |      |
| --help, -h   |                 | 查询命令帮助                     |      |
| --ledger     |                 | 在Ledger设备上存储私钥的本地引用 |      |

### 通过指定方法创建新密钥

```bash
iriscli keys new MyKey
```

执行命令后，系统提示输入bip44路径, 默认值是 44'/118'/0'/0/0.

```bash
> -------------------------------------
> Enter your bip44 path. Default is 44'/118'/0'/0/0
```

然后将要求输入你的bip39助记词，或者直接按回车键生成。

```bash
> Enter your bip39 mnemonic, or hit enter to generate one.
```

直接按回车键生成bip39助记词，然后会提示要求输入bip39密码。

```bash
> -------------------------------------
> Enter your bip39 passphrase. This is combined with the mnemonic to derive the seed
> Most users should just hit enter to use the default, ""
```

你也可以直接按回车键跳过，然后将提示输入密码并确认：

```bash
> -------------------------------------
> Enter a passphrase to encrypt your key to disk:
> Repeat the passphrase:
```

完成以上步骤，你将得到一个新创建的密钥。
