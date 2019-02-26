# Service User Guide

## 基本功能描述
IRIS Services（又名“iServices”）旨在对链下服务从定义、绑定（服务提供方注册）、调用到治理（分析和争端解决）的全生命周期传递，来跨越区块链世界和传统业务应用世界之间的鸿沟。 
IRIS-SDK通过增强的IBC处理逻辑来支持服务语义，以允许分布式商业服务在区块链互联网上可用。我们引入接口描述语言（[Interface description language](https://en.wikipedia.org/wiki/Interface_description_language)，
简称IDL）对服务进行标准化定义来满足跨语言的服务调用。目前支持的IDL语言为[protobuf](https://developers.google.com/protocol-buffers/)。该模块的主要功能点如下：
1. 服务定义
2. 服务绑定
3. 服务调用
4. 争议解决 (TODO)
5. 服务分析 (TODO)

### 系统参数
以下参数均可通过governance(./governance.md)修改

* `MinDepositMultiple`    服务绑定最小抵押金额的倍数
* `MaxRequestTimeout`     服务调用最大等待区块个数
* `ServiceFeeTax`         服务费的税收比例
* `SlashFraction`         惩罚百分比
* `ComplaintRetrospect`   可提起争议最大时长
* `ArbitrationTimeLimit`  争议解决最大时长

## 使用场景

### 服务定义

任何用户可以发起服务定义请求，在服务定义中，使用`protobuf`对该服务的方法，方法的输入、输出参数进行标准化定义。为了更好的支持服务的属性，IRISnet对`protobuf`进行了一些扩展，详细请参照[IDL文件扩展](#idl文件扩展)

```
# 创建服务定义
iriscli service define --chain-id=service-test  --from=x --fee=0.3iris --service-name=test-service --service-description=service-description --author-description=author-description --tags=tag1,tag2 --idl-content=<idl-content> --file=test.proto

# 查询服务定义
iriscli service definition --def-chain-id=service-test --service-name=test-service
```

### 服务绑定
在服务绑定中, 需要抵押一定数量的押金, 最小的抵押金额为该服务的服务费价格的`MinDepositMultiple`倍数。服务提供方可以随时更新他的服务绑定并调整服务价格，禁用、启用该服务绑定。如果想取回押金，需要禁用服务绑定并等待`ComplaintRetrospectParameter`+`ArbitrationTimelimitParameter`的周期。

```
# 服务绑定
iriscli service bind --chain-id=service-test  --from=x --fee=0.3iris --service-name=test-service --def-chain-id=service-test --bind-type=Local  --deposit=1000iris --prices=1iris --avg-rsp-time=10000 --usable-time=100

# 查询服务绑定
iriscli service binding --def-chain-id=service-test --service-name=test-service --bind-chain-id=service-test --provider=<your address>

# 查询服务绑定列表
iriscli service bindings --def-chain-id=service-test --service-name=test-service

# 服务绑定更新
iriscli service update-binding --chain-id=service-test  --from=x --fee=0.3iris --service-name=test-service --def-chain-id=service-test --bind-type=Local  --deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100

# 禁用服务绑定
iriscli service disable --chain-id=service-test  --from=x --fee=0.3iris --def-chain-id=service-test --service-name=test-service

# 开启服务绑定
iriscli service enable --chain-id=service-test  --from=x --fee=0.3iris --def-chain-id=service-test --service-name=test-service --deposit=1iris

# 取回押金
iriscli service refund-deposit --chain-id=service-test  --from=x --fee=0.3iris --def-chain-id=service-test --service-name=test-service
```

### 服务调用
服务消费者如果需要发起服务调用请求，需要支付服务提供方指定的服务费。服务提供方需要在`MaxRequestTimeout`定义的区块高度内响应该服务请求，如果超时未响应，将从服务提供方的该服务绑定押金中扣除`SlashFraction`比例的押金。如果服务调用被正常响应，系统从该次服务调用的服务费中将扣除`ServiceFeeTax`比例的系统税收，同时将剩余的服务费加入到服务提供方的收入池中。如果服务调用未及时响应，该次服务调用的服务费将退还到服务消费者的退费池中。服务提供方/消费者可以发起`withdraw-fees`/`refund-fees`交易取回自己在收入池/退费池中所有的token。

```
# 发起服务调用
iriscli service call --chain-id=test --from=node0 --fee=0.3iris --def-chain-id=test --service-name=test-service --method-id=1 --bind-chain-id=test --provider=faa1qm54q9ta97kwqaedz9wzd90cacdsp6mq54cwda --service-fee=1iris --request-data=434355

# 查询服务请求列表
iriscli service requests --def-chain-id=test --service-name=test-service --bind-chain-id=test --provider=faa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x

# 响应服务调用
iriscli service respond --chain-id=test --from=node0 --fee=0.3iris --request-chain-id=test --request-id=230-130-0 --response-data=abcd

# 查询服务响应
iriscli service response --request-chain-id=test --request-id=635-535-0

# 查询指定地址的服务费退款和收入
iriscli service fees [account address]

# 从服务费退款中退还所有费用
iriscli service refund-fees --chain-id=test --from=node0 --fee=0.3iris

# 从服务费收入中取回所有费用
iriscli service withdraw-fees --chain-id=test --from=node0 --fee=0.3iris
```

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

    > syntax = \"proto3\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output_privacy: NoPrivacy\n    //@Attribute output_cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n

* IDL文件参照

    [test.proto](./test.proto)
