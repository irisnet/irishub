# IRIS Command tool

## Introduction
`iristool` 现在包含debug。

## debug
简单调试的工具。

同时接受十六进制和base64格式，并提供相对应的响应。

通常情况下，日志中字节编码为十六进制，而在JSON中编码为base64。

### Usage

* Pubkeys 
下面得到相同的结果:

```bash
iristool debug pubkey TZTQnfqOsi89SeoXVnIw+tnFJnr4X8qVC0U8AsEmFk4=
iristool debug pubkey 4D94D09DFA8EB22F3D49EA17567230FAD9C5267AF85FCA950B453C02C126164E
```

* Txs
传入 hex/base64 tx 并返回完整的 JSON

```bash
iristool debug tx <hex or base64 transaction>
iristool debug tx lALZHnawCngqlySsCjgKFB8/ff4eQXErqiVk7Sm4/nrpwY+LEiAKCWlyaXMtYXR0bxITMjAwMDAwMDAwMDAwMDAwMDAwMBI4ChQfP33+HkFxK6olZO0puP566cGPixIgCglpcmlzLWF0dG8SEzIwMDAwMDAwMDAwMDAwMDAwMDASJQofCglpcmlzLWF0dG8SEjQwMDAwMDAwMDAwMDAwMDAwMBDQhgMabQom61rphyEC49U43CwWbrdmPS6djiJzj1P8S36rV/AFn70XlXu0tHESQApZyuWZB8oZnOyRO0Pk0fsmhCe9OXmsZ/JiSXCujKvdezBTqmRjlSq95Wqo8qoxMukLylhdlQF3GfkbW+PriBgg9x4=
```

返回
```bash
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

* 查询 addr

```bash
iristool debug addr iaa1rulhmls7g9cjh239vnkjnw870t5urrutsfwvmc
```
