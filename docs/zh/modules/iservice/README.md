# IService User Guide

## 基本功能描述

1. 服务定义
2. 服务绑定 (TODO)
3. 服务调用 (TODO)
4. 争议解决 (TODO)
5. 服务分析 (TODO)

## 交互流程

### 服务定义流程

1. 任何用户可以发起服务定义请求，在服务定义中，使用 [protobuf](https://developers.google.com/protocol-buffers/) 对该服务的方法，方法的输入、输出参数进行标准化定义。

## 使用场景
### 创建使用环境

```
rm -rf iris
rm -rf .iriscli
iris init gen-tx --name=x --home=iris
iris init --gen-txs --chain-id=service-test -o --home=iris
iris start --home=iris
```

### 服务定义

```
# 服务定义
iriscli iservice define --chain-id=service-test  --from=x --fee=0.004iris --name=test-service --service-description=service-description --author-description=author-description --tags "tag1 tag2" --broadcast=Broadcast --idl-content=<idl-content> --file=test.proto

# 结果
Committed at block 1040 (tx hash: 58FD40B739F592F5BD9B904A661B8D7B19C63FA9, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:13601 Tags:[{Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[247 102 151 120 200 0]}] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "completeConsumedTxFee-iris-atto": "159740000000000"
   }
}

# 查询服务定义
iriscli iservice definition --name=test-service --chain-id=service-test

```

## 命令详情

```
iriscli iservice define --chain-id=service-test  --from=x --fee=0.004iris --name=test-service --service-description=service-description --author-description=author-description --tags "tag1 tag2" --broadcast=Broadcast --idl-content=<idl-content> --file=test.proto
```

* `--name`  该iService服务的名称
* `--service-description`  该iService服务的描述
* `--author-description`  该iService服务创建者的描述. 可选
* `--tags`  该iService服务的关键字
* `--broadcast`  此服务消息广播类型{`Broadcast`,`Unicast`}
* `--idl-content`  对该iService服务的methods的标准化定义内容
* `--file`  可使用文件代替idl-content，当该项不为空时，覆盖`idl-content`内容

```
iriscli iservice definition --name=test-service --chain-id=service-test
```

* `--chain-id` 定义该iservice服务的区块链ID
* `--name`  iService服务的名称

## IDL文件扩展
在使用proto文件对服务的方法，输入、输出参数进行标准化定义时，可通过注释的方式增加method属性。

### 注释标准
* 使用 `//@Attribute 属性： 值`的方式添加在rpc方法上，即可将该属性添加为方法的属性。例如: 
> //@Attribute description: sayHello

### 目前支持的属性
* `description` 对该方法的描述
* `output_privacy` 是否对该方法的输出进行加密处理，{`NoPrivacy`,`PubKeyEncryption`}
* `output_cached` 是否对该方法的输出进行缓存，{`OffChainCached`,`NoCached`}

### IDL content参照
* IDL content参照
`syntax = \"proto3\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output_privacy: NoPrivacy\n    //@Attribute output_cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n`

* IDL文件参照
[test.proto](./test.proto)