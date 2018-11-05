# 参与到Genesis文件生成流程中


1. 每个希望成为验证人的参与者确保安装了对应版本的软件：iris v0.7.0

2. 先创建账户,再执行gentx命令，获得一个gentx-node-ID.json的文件。这个操作将默认生成一个余额为150IRIS的账户，该账户默认绑定100IRIS成为一个验证人候选人。

```
iriscli keys add your_name
iris gentx --name=your_name --home=<path_to_home> --ip=Your_public_IP
```

代码示例：
   
```
iriscli keys add alice
iris gentx --name=alice --home=iris --chain-id=irishub-stage --ip=1.1.1.1
```
然后你可以发现在$IRISHOME/config目录下生成了一个gentx文件夹。里面存在一个gentx-node-ID.json文件。这个文件包含了如下信息：

```
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "cosmos-sdk/MsgCreateValidator",
        "value": {
          "Description": {
            "moniker": "chenggedexiaokeai.local",
            "identity": "",
            "website": "",
            "details": ""
          },
          "Commission": {
            "rate": "0.1000000000",
            "max_rate": "0.2000000000",
            "max_change_rate": "0.0100000000"
          },
          "delegator_address": "faa1cf25tf4pfjdhkzx8lqnkajlse6jcpm2fyw4yme",
          "validator_address": "fva1cf25tf4pfjdhkzx8lqnkajlse6jcpm2f3lltx7",
          "pubkey": {
            "type": "tendermint/PubKeyEd25519",
            "value": "/JvLFsvyMgm2ND4QgN4JKyLxhL42dVgat67383Q+mPY="
          },
          "delegation": {
            "denom": "iris-atto",
            "amount": "100000000000000000000"
          }
        }
      }
    ],
    "fee": {
      "amount": null,
      "gas": "200000"
    },
    "signatures": [
      {
        "pub_key": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AtfNRj0zYvffAQG+iad6SScfdl29ag9G3EI0JDSwKJmy"
        },
        "signature": "BwTejBceK4M+3LzmNl62jVFUr9wVv//UO7iI/yWi5KFoez9eY43HSlaZJf+3rnKLjosn2tD79EIw55BJ6SbYzQ==",
        "account_number": "0",
        "sequence": "0"
      }
    ],
    "memo": "0eb02fdabb96923ac1e855ac012a5a624793264a@1.1.1.1:26656"
  }
}
```

`msg` 是创建验证人节点的交易

3. 将上述提到的json文件以提交Pull Request的形式上传到`https://github.com/irisnet/testnets/tree/master/testnets/fuxi-5000/config/gentx`目录下：

   注意⚠️：json文中的IP改成公网IP




