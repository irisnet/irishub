# iriscli service definition

## 描述

查询服务定义

## 用法

```
iriscli service definition [flags]
```

## 标志

| Name, shorthand | Default                    | Description                                            | Required |
| --------------- | -------------------------- | ------------------------------------------------------ | -------- |
| --def-chain-id  |                            | [string] 定义该服务的区块链ID                              | 是        |
| --service-name  |                            | [string] 服务名称                                        | 是        |
| --help, -h      |                            | 查询定义命令帮助                                           |          |
| --chain-id      |                            | [string] tendermint节点的链ID                            |          |
| --height        | 最近可证明区块高度            | [int] 查询的区块高度                                       |          |
| --indent        |                            | 在JSON格式的应答中添加缩进                                  |          |
| --ledger        |                            | 使用连接的硬件记账设备                                      |          |
| --node          | tcp://localhost:26657      | [string] tendermint节点开启的远程过程调用接口\<主机>:\<端口>  |          |
| --trust-node    | true                       | 关闭响应结果校验                                           |          |


## 例子

### 查询服务定义

```shell
iriscli service definition --def-chain-id=test --service-name=test-service
```

运行成功以后，返回的结果如下:

```json
{
  "SvcDef": {
    "name": "test-service",
    "chain_id": "test",
    "description": "service-description",
    "tags": [
      "tag1",
      "tag2"
    ],
    "author": "faa1ydhmma8l4m9dygsh7l08fgrwka6yczs0gkfnvd",
    "author_description": "author-description",
    "idl_content": "syntax = \"proto3\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output_privacy: NoPrivacy\n    //@Attribute output_cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n"
  },
  "methods": [
    {
      "id": "1",
      "name": "SayHello",
      "description": "sayHello",
      "output_privacy": "NoPrivacy",
      "output_cached": "NoCached"
    }
  ]
}
```