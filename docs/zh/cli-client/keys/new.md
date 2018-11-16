# iriscli keys new

## 描述

使用交互式命令派生新的私钥，该命令将提示您输入每个输入。可选择指定bip39助记符，bip39密码短语以进一步保护助记符，和bip32 HD路径，以获得特定帐户。密钥将以给定名称存储并使用给定的密码加密。这里唯一需要的输入是加密密码。

## 使用方式

```
iriscli keys new <name> [flags]
```

## 标志

| 名称, 速记       | 默认值             | 描述                                                            | 是否必须  |
| --------------- | ----------------- | --------------------------------------------------------------- | -------- |
| --bip44-path    | 44'/118'/0'/0/0   | BIP44 path from which to derive a private key                   |          |
| --default       |                   | Skip the prompts and just use the default values for everything |          |
| --help, -h      |                   | help for add                                                    |          |
| --ledger        |                   | Store a local reference to a private key on a Ledger device     |          |

## 例子

### Create a new key by the specified method

```shell
iriscli keys new MyKey
```

执行命令后，系统提示你输入 bip44 路径, 直接回车的默认值是 44'/118'/0'/0/0.

```txt
> -------------------------------------
> Enter your bip44 path. Default is 44'/118'/0'/0/0
```

然后你会被要求输入你的bip39助记词，或者直接敲回车键生成一个。注意保存好助记词。

```txt
> Enter your bip39 mnemonic, or hit enter to generate one.
```

直接按回车键生成bip39助记符，会显示一个新提示，要求您输入bip39密码。

```txt
> -------------------------------------
> Enter your bip39 passphrase. This is combined with the mnemonic to derive the seed
> Most users should just hit enter to use the default, ""
```

你也可以直接回车键跳过，然后你会收到一个输入并确认密码的提示：

```txt
> -------------------------------------
> Enter a passphrase to encrypt your key to disk:
> Repeat the passphrase:
```

设置密码完成这些步骤，你就创建了一个自定义的密钥。