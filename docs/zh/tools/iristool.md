---
order: 2
---

# 命令行工具

## 介绍

`iristool`是用于简单调试的工具。

接受十六进制和base64格式，并提供可读的响应。

通常，我们在日志中将字节编码为十六进制，但在JSON中将字节编码为base64。

## 用法

### iristool debug pubkey

下面得到相同的结果:

```bash
iristool debug pubkey TZTQnfqOsi89SeoXVnIw+tnFJnr4X8qVC0U8AsEmFk4=
iristool debug pubkey 4D94D09DFA8EB22F3D49EA17567230FAD9C5267AF85FCA950B453C02C126164E
  ```

### iristool debug tx

传入hex/base64编码的tx返回完整的JSON

```bash
iristool debug tx <hex or base64 transaction>
iristool debug tx lALZHnawCngqlySsCjgKFB8/ff4eQXErqiVk7Sm4/nrpwY+LEiAKCWlyaXMtYXR0bxITMjAwMDAwMDAwMDAwMDAwMDAwMBI4ChQfP33+HkFxK6olZO0puP566cGPixIgCglpcmlzLWF0dG8SEzIwMDAwMDAwMDAwMDAwMDAwMDASJQofCglpcmlzLWF0dG8SEjQwMDAwMDAwMDAwMDAwMDAwMBDQhgMabQom61rphyEC49U43CwWbrdmPS6djiJzj1P8S36rV/AFn70XlXu0tHESQApZyuWZB8oZnOyRO0Pk0fsmhCe9OXmsZ/JiSXCujKvdezBTqmRjlSq95Wqo8qoxMukLylhdlQF3GfkbW+PriBgg9x4=
```

returns

```json
{
  "type": "irishub/bank/StdTx",
  "value": {
    "msg": [
      {
        "type": "irishub/bank/Send",
        "value": {
          "inputs": [
            {
              "address": "iaa1rulhmls7g9cjh239vnkjnw870t5urrutsfwvmc",
              "coins": [
                {
                  "denom": "iris-atto",
                  "amount": "2000000000000000000"
                }
              ]
            }
          ],
          "outputs": [
            {
              "address": "iaa1rulhmls7g9cjh239vnkjnw870t5urrutsfwvmc",
              "coins": [
                {
                  "denom": "iris-atto",
                  "amount": "2000000000000000000"
                }
              ]
            }
          ]
        }
      }
    ],
    "fee": {
      "amount": [
        {
          "denom": "iris-atto",
          "amount": "400000000000000000"
        }
      ],
      "gas": "50000"
    },
    "signatures": [
      {
        "pub_key": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AuPVONwsFm63Zj0unY4ic49T/Et+q1fwBZ+9F5V7tLRx"
        },
        "signature": "ClnK5ZkHyhmc7JE7Q+TR+yaEJ705eaxn8mJJcK6Mq917MFOqZGOVKr3laqjyqjEy6QvKWF2VAXcZ+Rtb4+uIGA==",
        "account_number": "0",
        "sequence": "3959"
      }
    ],
    "memo": ""
  }
}
```

### iristool debug addr

```bash
iristool debug addr iaa1rulhmls7g9cjh239vnkjnw870t5urrutsfwvmc
  ```

returns

```bash
Address (Hex): 1F3F7DFE1E41712BAA2564ED29B8FE7AE9C18F8B
Bech32 Acc: iaa1rulhmls7g9cjh239vnkjnw870t5urrutsfwvmc
Bech32 Val: iva1rulhmls7g9cjh239vnkjnw870t5urrut9cyrxl
Bech32 Cons: ica1rulhmls7g9cjh239vnkjnw870t5urrutvqntyu
```

### iristool debug raw-bytes

将原始字节输出(如[10 21 13 255])转化为hex编码

```bash
iristool debug raw-bytes <raw-bytes>
iristool debug raw-bytes "[10 21 13 255]"
```

### iristool debug rand-secret

产生随机密钥

```bash
iristool debug rand-secret
```

### iristool debug hash-lock

生成带有密钥和时间戳（如果提供）的哈希锁

```bash
iristool debug hash-lock <secret-hex64> <timestamp>
iristool debug hash-lock 10dfd779e15176f3b2867f0acc9e18d29f65ea4002957c632d1bea200b9b2915 1580000000
```
