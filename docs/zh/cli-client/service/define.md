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