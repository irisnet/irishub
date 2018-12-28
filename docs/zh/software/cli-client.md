# IRIS Command Line Client

## 介绍

iriscli是参与IRISnet网络的客户端。IRISnet的用户都可以通过iriscli来发送各种不同的交易或者进行各种查询。

## iriscli目录

iriscli客户端的默认目录是`$HOME/.iriscli`，主要用来保存配置文件和数据。 IRISnet `key` 的数据就保存在iriscli的HOME目录下。也可以通过`--home`来指定客户端的HOME目录。

## iriscli --node

tendermint节点的rpc地址,交易和查询的消息都发送到监听这个端口的进程。默认是`tcp://localhost:26657`，也可以通过`--node`指定rpc地址。

## iriscli config命令

iriscli config命令可以交互式地配置一些默认参数，例如chain-id，home，fee 和 node。

## 手续费和Gas

iriscli发送交易可以通过`--fee`指定手续费和`--gas`指定Gas（默认值是200000）。Gas除以手续费就是Gas Price，Gas Price不能小于区块链设定的最小值。执行完整个交易以后剩余的手续费会返还给用户。可以设置`--gas="simulate"`, 它可以通过仿真运行估算出交易大致所需要消耗的Gas，并且乘以一个由`--gas-adjustment`指定的系数得到最终的Gas，作为这次交易的Gas。最后交易才会被广播出去。

## dry-run模式

iriscli默认关闭dry-run模式。如果想打开dry-run模式，就可以使用`--dry-run`。它和simulate处理逻辑类似，可以计算出需要消耗的Gas，但是之后它不会广播给全节点，直接返回并打印此次消耗的Gas。

例子：使用dry-run模式发送命令

```
iriscli gov submit-proposal --title="ABC" --description="test" --type=Text --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --dry-run
```

返回：

```
estimated gas = 8370
```

## 异步模式

iriscli发送的交易默认是同步模式。同步模式指的是发送交易，然后会堵塞，直到收到回复或者超时，整个命令才结束。如果想打开异步模式，就可以使用`--async`。发送交易之后，会立马返回交易的hash。

## generate-only

`generate-only`默认是关闭的，但可以使能`--generate-only`，然后会打印命令行生成的未签名交易。

例子：使能generate-only生成未签名交易

```
iriscli gov submit-proposal --title="ABC" --description="test" --type=Text --deposit=1iris --from=x --chain-id=gov-test --fee=0.05iris --gas=200000 --generate-only
```

返回：

```json
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "cosmos-sdk/MsgSubmitProposal",
        "value": {
          "title": "ABC",
          "description": "test",
          "proposal_type": "Text",
          "proposer": "faa1k47r0nxd6ec8n6sc6tzvk2053u4eff0vx99755",
          "initial_deposit": [
            {
              "denom": "iris-atto",
              "amount": "1000000000000000000"
            }
          ],
          "Param": {
            "key": "",
            "value": "",
            "op": ""
          }
        }
      }
    ],
    "fee": {
      "amount": [
        {
          "denom": "iris-atto",
          "amount": "50000000000000000"
        }
      ],
      "gas": "200000"
    },
    "signatures": null,
    "memo": ""
  }
}

```

## trust-node

trust-node默认为true。当trust-node是true时， iriscli的客户端只查询数据并不对数据进行默克尔证明。你也可以通过`--trust-node=false`, 对查询得到的数据进行默克尔证明。
