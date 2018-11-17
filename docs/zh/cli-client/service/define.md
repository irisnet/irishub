# iriscli service define 

## 描述

创建一个新的服务定义

## 用法

```
iriscli service define [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                       | Required |
| --------------------- | ----------------------- | --------------------------------------------------------------------------------- | -------- |
| --service-description |                         | [string] 服务的描述                                                                 |          |
| --author-description  |                         | [string] 服务创建者的描述                                                            |          |
| --service-name        |                         | [string] 服务名称                                                                   |   Yes    |
| --tags                |                         | [strings] 该服务的关键字                                                             |          |
| --idl-content         |                         | [string] 对该服务描述的接口定义语言内容                                                 |          |
| --file                |                         | [string] 对该服务描述的接口定义语言内容的文件路径                                         |          |
| -h, --help            |                         | 服务定义命令帮助                                                                    |          |
| --account-number      |                         | [int] 用来签名交易的AccountNumber                                                     |          |
| --async               |                         | 异步广播交易                                                                         |          |
| --chain-id            |                         | [string] tendermint节点的链ID                                                       | 是       |
| --dry-run             |                         | 忽略--gas标志并进行本地的交易仿真                                                       |          |
| --fee                 |                         | [string] 支付的交易费用                                                              | 是       |
| --from                |                         | [string] 用来签名的私钥名                                                            | 是       |
| --from-addr           |                         | [string] 指定generate-only模式下的from address                                       |          |
| --gas string          | 200000                  | 单次交易的gas上限; 设置为"simulate"将自动计算相应的阈值                                   |          |
| --gas-adjustment      | 1                       | [float] 这个调整因子将乘以交易仿真返回的估计值; 如果手动设置了gas的上限，这个标志将被忽略       |          |
| --generate-only       |                         | 构建一个未签名交易并将其打印到标准输出                                                   |          |
| --indent              |                         | 在JSON响应中添加缩进                                                                 |          |
| --json                |                         | 输出将以json格式返回                                                                 |          |
| --ledger              |                         | 使用连接的硬件记账设备                                                                |          |
| --memo                |                         | [string] 发送交易的备忘录                                                            |          |
| --node                | tcp://localhost:26657   | [string] [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口>                   |          |
| --print-response      |                         | 返回交易响应 (当且仅当同步模式下使用))                                                  |          |
| --sequence            |                         | [int] 用来签名交易的sequence number                                                 |          |
| --trust-node          | true                    | 关闭响应结果校验                                                                    |          |

## 例子

### 创建一个新的服务定义
```shell
iriscli service define --chain-id=test  --from=node0 --fee=0.004iris --service-name=test-service --service-description=service-description --author-description=author-description --tags=tag1,tag2 --idl-content=<idl-content> --file=test.proto
```
如果文件项不是空的，将会替换Idl-content.  [IDL内容示例](#idl-content-example).

运行成功以后，返回的结果如下:

```txt
Password to sign with 'node0':
Committed at block 65 (tx hash: 663B676E453F91BFCDC87B0308910501DD14DF79C88390FC15E06C4CC9612422, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:7968 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 100 101 102 105 110 101] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 49 53 57 51 54 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "service-define",
     "completeConsumedTxFee-iris-atto": "\"159360000000000\""
   }
 }
```

### IDL内容示例
* IDL内容示例

    > syntax = \\"proto3\\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output_privacy: NoPrivacy\n    //@Attribute output_cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n

* IDL文件示例

    [test.proto](../../../features/test.proto)

