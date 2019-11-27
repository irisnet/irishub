---
order: 2
---

# Command Line Tool

## Introduction

`iristool` is a tool for simple debugging.

Accept both hex and base64 formats and provide a readable response.

Usually, we encode bytes as hex in the logs, but as base64 in the JSON.

## Usage

### iristool debug pubkey

The following give the same result:

```bash
iristool debug pubkey TZTQnfqOsi89SeoXVnIw+tnFJnr4X8qVC0U8AsEmFk4=
iristool debug pubkey 4D94D09DFA8EB22F3D49EA17567230FAD9C5267AF85FCA950B453C02C126164E
  ```

### iristool debug tx

Input a hex/base64 tx and get back the full JSON

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

Convert raw bytes output (eg. [10 21 13 255]) to hex

```bash
iristool debug raw-bytes <raw-bytes>
iristool debug raw-bytes "[10 21 13 255]"
```

### iristool debug rand-secret

Generate a random secret

```bash
iristool debug rand-secret
```

### iristool debug hash-lock

Generate a hash lock with secret and timestamp (if privided)

```bash
iristool debug hash-lock <secret-hex64> <timestamp>
iristool debug hash-lock 10dfd779e15176f3b2867f0acc9e18d29f65ea4002957c632d1bea200b9b2915 1580000000
```
