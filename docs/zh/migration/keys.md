# Keys Migrate

irishub v0.16.x 的密钥文件（私钥）使用数据库存储。 新版本 v1.0.1 将提供一种存储用户私钥的新方法。为了支持将旧密钥文件迁移到新版本，提供了两种解决方案。

## 助记词

这种方式适用于具有助记词的用户。创建新帐户时，系统将为用户随机分配一个助记词，并使用该助记词恢复用户的私钥。 无论是 v0.16.x 或 v1.0.1 还是更高版本，助记词都保持不变。您可以在 `add` 命令中添加 `--recover` 使用助记词还原帐户，例如：

```bash
iris keys add n2 --recover
```

## Keystore

这种方法适用于丢失了助记词但保存了密钥的 db 文件或密钥文件的用户。irishub v0.16.x 的密钥文件的格式与以太坊的格式相似，并且 v1.0.1 也完全兼容。因此，用户可以使用密钥库导出旧的私钥，然后使用 irishub v1.0.1 版本导入密钥库以完成密钥迁移，操作流程如下：

**1. 通过 irishub v0.16.x 导出 keystore 文件**

```bash
iriscli keys export test1 --output-file=key.json --home ./iriscli_test 
```

输出：

```json
{
    "address":"iaa1k2j3ws7ghwl9qha36xdcmwuu4rend2yr9tw05q",
    "crypto":{
        "cipher":"aes-128-ctr",
        "ciphertext":"b5e586baf1126f982ee89ffa9fd23fc68e0a25e1d561d6d59896a0b4878a4d5f",
        "cipherparams":{
            "iv":"d02a7b40ce35b6e87f9a395850372bbc"
        },
        "kdf":"pbkdf2",
        "kdfparams":{
            "c":262144,
            "dklen":32,
            "prf":"hmac-sha256",
            "salt":"8c77a3a8a75a76da203b262e7fa0187bafbd2ab8bfd3b21ba99f88dcc550d1a6"
        },
        "mac":"4bdf3fd116a4b9d7eb8846d078399f41a6e271a80678ce8979e4fa86f793cdeb"
    },
    "id":"c63bdcd2-c470-4c9a-90eb-a4ef6d3d5937",
    "version":"1"
}
```

**2. 通过 irishub v1.0.1 to 导入 keystore 文件**

```bash
iris keys import n2 key.json --keyring-backend file 
```

**3. 验证导入的账户信息**

```bash
iris keys show n2 --keyring-backend file
```

输出：

```text
Enter keyring passphrase:
- name: n2
type: local
address: iaa1k2j3ws7ghwl9qha36xdcmwuu4rend2yr9tw05q
pubkey: iap1addwnpepqgrj4yshwmq7v7akp04empq9rrn6w26e8q6gpl7jkfjaexk93deq2pwa3m6
mnemonic: ""
threshold: 0
pubkeys: []
```

输出帐户地址与密钥库文件中的地址一致，且迁移成功。
