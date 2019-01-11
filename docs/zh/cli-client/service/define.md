# iriscli service define 

## 描述

创建一个新的服务定义

## 用法

```
iriscli service define [flags]
```
## 标志

| 命令，缩写       | 类型   | 是否必须 | 默认值                | 描述                                                         |
| ---------------- | ------ | -------- | --------------------- | ------------------------------------------------------------ |
| -h, --help       |        | 否       |                       | 打印帮助                                                     |
| --chain-id       | String | 否       |                       | tendermint 节点网络ID                                        |
| --account-number | int    | 否       |                       | 账户数字用于签名通证发送                                     |
| --async          |        | 否       | True                  | 异步广播传输信息                                             |
| --dry-run        |        | 否       |                       | 忽略--gas 标志 ，执行仿真传输，但不广播。                    |
| --fee            | String | 是       |                       | 设置传输需要的手续费                                         |
| --from           | String | 是       |                       | 用于签名的私钥名称                                           |
| --gas            | String | 否       | 20000                 | 每笔交易设定的gas限额; 设置为“simulate”以自动计算所需气体    |
| --gas-adjustment | Float  | 否       | 1                     | 调整因子乘以传输模拟返回的估计值; 如果手动设置气体限制，则忽略该标志 |
| --generate-only  |        | 否       |                       | 创建一个未签名的传输并写到标准输出中。                       |
| --indent         |        | 否       |                       | 在JSON响应中增加缩进                                         |
| --json           |        | 否       |                       | 以json格式返回输出                                           |
| --memo           | String | 否       |                       | 传输中的备注信息                                             |
| --print-response |        | 否       |                       | 返回传输响应 (仅仅当 async = false时有效)                    |
| --sequence       | Int    | 否       |                       | 等待签名传输的序列号。                                       |
| --ledger         | String | 否       |                       | 使用一个联网的分账设备                                       |
| --node           | String | 否       | tcp://localhost:26657 | <主机>:<端口> 链上的tendermint rpc 接口。                    |
| --trust-node     | String | 否       | True                  | 不验证响应的证明                                             |
| --commit         | String | 否     | True                  |是否等到交易有明确的返回值，如果是True，则忽略--async的内容|

## 特有标志

| Name, shorthand       | Default                 | Description                                                                       | Required |
| --------------------- | ----------------------- | --------------------------------------------------------------------------------- | -------- |
| --service-description |                         | [string] 服务的描述                                                                 |          |
| --author-description  |                         | [string] 服务创建者的描述                                                            |          |
| --service-name        |                         | [string] 服务名称                                                                   |   Yes    |
| --tags                |                         | [strings] 该服务的关键字                                                             |          |
| --idl-content         |                         | [string] 对该服务描述的接口定义语言内容                                                 |          |
| --file                |                         | [string] 对该服务描述的接口定义语言内容的文件路径                                         |          |

## 示例

### 创建一个新的服务定义
```shell
iriscli service define --chain-id=test  --from=node0 --fee=0.004iris --service-name=test-service --service-description=service-description --author-description=author-description --tags=tag1,tag2 --idl-content=<idl-content> --file=test.proto
```
如果文件项不是空的，将会替换Idl-content.  [IDL内容示例](#idl-content-example).

运行成功以后，返回的结果如下:

```txt
Committed at block 539 (tx hash: 9ED8B36F8DDA7745BF03E0F5271E55B6D0BC34B373BFCDB6B5BC78C502DAE032, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 7604,
   "codespace": "",
   "tags": {
     "action": "service_define"
   }
 })
```

### IDL内容示例
* IDL内容示例

    > syntax = \\"proto3\\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output_privacy: NoPrivacy\n    //@Attribute output_cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n

* IDL文件示例

    [test.proto](https://github.com/irisnet/irishub/blob/master/docs/features/test.proto)