# iriscli keys add

## 描述

创建一个新密钥，或通过助记词导入已有密钥

## 使用方式

```
iriscli keys add <name> <flags>
```

## 标志

| 名称, 速记       | 默认值     | 描述                                                              | 是否必须  |
| --------------- | --------- | ----------------------------------------------------------------- | -------- |
| --account       |           | HD推导的账号                                              |          |
| --dry-run       |           | 忽略--gas标志并进行本地的交易仿真                                     |          |
| --help, -h      |           | 查询命令帮助                                                       |          |
| --index         |           | HD推导的索引号                                            |          |
| --ledger        |           | 使用ledger设备                                               |          |
| --no-backup     |           | 不输出助记词（如果其他人正在看着操作终端）                              |          |
| --recover       |           | 提供助记词以恢复现有密钥而不是新建                                     |          |
| --type, -t      | secp256k1 | 私钥类型 (secp256k\|ed25519)                              |          |

## 例子

### 创建密钥

```shell
iriscli keys add MyKey
```

执行命令后，系统会要求你输入密钥密码，注意：密码必须至少为8个字符。

```txt
Enter a passphrase for your key:
Repeat the passphrase:
```

之后，你已经完成了创建新密钥的工作，但请记住备份你的助记词短语，如果你不慎忘记密码或丢失了密钥，这是唯一能恢复帐户的方法。

```txt
NAME:	TYPE:	ADDRESS:						PUBKEY:
MyKey	local	iaa1mmsm487rqkgktl2qgrjad0z3yaf9n8t5ee8f3x	iap1addwnpepq2g0u7cnxp5ew0yhqep8j4rth5ugq8ky7gjmunk8tkpze95ss23akexx3tn
**Important** write this seed phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

oval green shrug term already arena pilot spirit jump gain useful symbol hover grid item concert kiss zero bleak farm capable peanut snack basket
```

上面24个单词只是助记词的示例，**不要**在生产环境中使用。

### 通过助记词恢复密钥

如果你忘记了密码或丢失了密钥，或者你想在其他地方使用密钥，则可以通过助记词短语来恢复。

```txt
iriscli keys add MyKey --recover
```

系统会要求你输入并确认密钥的新密码，然后输入助记词。这样就能恢复你的密钥。

```txt
Enter a passphrase for your key:
Repeat the passphrase:
Enter your recovery seed phrase:
```

