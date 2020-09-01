# Bank模块

## 简介

该模块主要用于账户之间转账、查询账户余额，同时提供了通用的离线签名与交易广播方法。

## 使用场景

1. 账户查询

    可以通过账户地址查询该账户的余额。

    ```bash
    iris q bank account <account-address>
    ```

2. 账户间转账

    该命令包括了交易“构造，签名，广播”的所有操作。 如从账户A转账10iris给账户B：

    ```bash
    iris tx bank send [A] [B] [10iris] --fees=0.3iris --from=<key-name-of-wallet-A> --chain-id=irishub
    ```

    IRISnet支持多种代币流通，将来IRISnet可以在一个交易中包含多种代币交换。

3. 交易签名

    为了提高账户安全性，IRISnet支持交易离线签名保护账户的私钥。在任意交易中，使用参数--generate-only可以构建一个未签名的交易。使用转账交易作为示例：

    ```bash
    iris tx bank send [from_key_or_address] [to_address] [amount] --amount=10iris --fees=0.3iris --from=<key-name-of-wallet-A> --generate-only
    ```

    以上命令将构建一未签名交易：

    ```json
    {
      "type": "auth/StdTx",
      "value": {
        "msg": [ "txMsg" ],
        "fee": "fee",
        "signatures": null,
        "memo": ""
      }
    }
    ```

    将结果保存到文件`<file>`。

    对上述的离线交易进行签名：

    ```bash
    iris tx sign <file> --chain-id=irishub --from=<key-name>
    ```

    将返回已签名的交易：

    ```json
    {
      "type": "auth/StdTx",
      "value": {
        "msg": [ "txMsg" ],
        "fee": "fee",
        "signatures": [
          {
            "pub_key": {
              "type": "tendermint/PubKeySecp256k1",
              "value": "A+qXW5isQDb7blT/KwEgQHepji8RfpzIstkHpKoZq0kr"
            },
            "signature": "5hxk/R81SWmKAGi4kTW2OapezQZpp6zEnaJbVcyDiWRfgBm4Uejq8+CDk6uzk0aFSgAZzz06E014UkgGpelU7w==",
            "account_number": "0",
            "sequence": "11"
          }
        ],
        "memo": ""
      }
    }
    ```

    将结果保存到文件`<file>`。

4. 广播交易

    广播离线产生的已签名的交易，在这里，你可以使用上面的sign命令生成的交易。当然，您可以通过任何方法生成已签名的交易，例如：[irisnet-crypto](https://github.com/irisnet/irisnet-crypto)。

    ```bash
    iris tx broadcast <file>
    ```

    该交易将在IRIShub中广播并执行。
