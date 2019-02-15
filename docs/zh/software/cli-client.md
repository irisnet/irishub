# IRIS Command Line Client

## 介绍

iriscli是参与IRISnet网络的客户端。IRISnet的用户可通过iriscli来发送各种不同的交易或者进行查询。

## iriscli目录

iriscli客户端的默认工作目录是`$HOME/.iriscli`，主要用来保存配置文件和数据。 IRISnet `key` 的数据就保存在iriscli的HOME目录下。也可以通过`--home`来指定客户端的HOME目录。

## iriscli --node

--node用来指定所连接iris节点的rpc地址，交易和查询的消息都发送到监听这个端口的iris进程。默认是`tcp://localhost:26657`，也可以通过`--node`指定rpc地址。

## iriscli config命令

iriscli config 命令可以交互式地配置一些公共参数的默认值，例如chain-id，home，fee 和 node。完成配置后，后续的iriscli命令可以省略对这些flag参数的指定。

## Fee和Gas

iriscli发送交易可以通过`--fee`指定交易费和`--gas`指定Gas（默认值是20000）。交易费除以Gas就是Gas Price，Gas Price不能小于区块链设定的最小值。执行完整个交易以后剩余的交易费会返还给用户。可以设置`--gas="simulate"`, 它可以通过仿真运行估算出交易大致所需要消耗的Gas，并且乘以一个由`--gas-adjustment`指定的系数得到最终的Gas，作为这次交易的Gas。最后交易才会被广播出去。

## dry-run模式

iriscli默认关闭dry-run模式。如果想打开dry-run模式，就可以使用`--dry-run`。它和simulate处理逻辑类似，可以计算出需要消耗的Gas，但是之后它不会广播给全节点，直接返回并打印此次消耗的Gas。

例子：使用dry-run模式发送命令

```
iriscli gov submit-proposal --title="ABC" --description="test" --type="ParameterChange" --deposit=600iris --param='mint/Inflation=0.050' --from=x --chain-id=test --fee=0.5iris --dry-run
```

返回：

```
estimated gas = 5398
```

## 交易发送模式

async：不对交易进行任何验证，立即返回交易的hash
sync：对交易进行合法性验证（交易格式和签名），返回验证结果和交易hash，交易在网络中等待被打包出块
commit：等待交易被打包上链再返回交易完整执行结果，将堵塞请求，直到收到回复或者超时，整个命令才结束

iriscli发送的交易默认是sync模式。如果想用其他模式发送交易，可以使用`--async` 或 `--commit`。

## generate-only

`generate-only`默认是关闭的，但可以使能`--generate-only`，打印命令行生成未签名的交易。

例子：使能generate-only以生成未签名的交易

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
