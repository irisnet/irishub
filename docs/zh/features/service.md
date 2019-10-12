# iService

## 简介

IRIS Services（又名“iService”）旨在对链下服务从定义、绑定（服务提供方注册）、调用到治理（分析和争端解决）的全生命周期传递，来跨越区块链世界和传统业务应用世界之间的鸿沟。

IRIS-SDK通过增强的IBC处理逻辑来支持服务语义，以允许分布式商业服务在区块链互联网上可用。我们引入接口描述语言（[Interface description language](https://en.wikipedia.org/wiki/Interface-description-language)，
简称IDL）对服务进行标准化定义来满足跨语言的服务调用。目前支持的IDL语言为[protobuf](https://developers.google.com/protocol-buffers/)。该模块的主要功能点如下：

1. 服务定义
2. 服务绑定
3. 服务调用
4. 争议解决 (TODO)
5. 服务分析 (TODO)

### 系统参数

以下参数均可通过[governance](governance.md)修改

* `MinDepositMultiple`    服务绑定抵押金额（相对于服务价格）最小的倍数
* `MaxRequestTimeout`     服务调用最大等待区块个数
* `ServiceFeeTax`         服务费的税收比例
* `SlashFraction`         惩罚百分比
* `ComplaintRetrospect`   可提起争议最大时长
* `ArbitrationTimeLimit`  争议解决最大时长

## 使用场景

### 服务定义

任何用户可以发起服务定义请求，在服务定义中，使用`protobuf`对该服务的方法，方法的输入、输出参数进行标准化定义。为了更好的支持服务的属性，IRISnet对`protobuf`进行了一些扩展，详细请参照[IDL文件扩展](#idl文件扩展)

```bash
# 创建服务定义
iriscli service define --chain-id=<chain-id>  --from=<key-name> --fee=0.6iris --gas=100000 --service-name=<service-name> --service-description=<service-description> --author-description=<author-description> --tags=<tag1>,<tag2> --idl-content=<idl-content> --file=</***/***.proto>

# 查询服务定义
iriscli service definition --def-chain-id=<def-chain-id> --service-name=<service-name>
```

### 服务绑定

在服务绑定中, 需要抵押一定数量的押金, 最小的抵押金额为该服务的服务费价格的`MinDepositMultiple`倍数。服务提供方可以随时更新他的服务绑定并调整服务价格，禁用、启用该服务绑定。如果想取回押金，需要禁用服务绑定并等待`ComplaintRetrospectParameter`+`ArbitrationTimelimitParameter`的周期。

```bash
# 服务绑定（抵押1000iris， 价格1iris， 平均响应时间10000毫秒， 服务可用性9999（10000次调用可用次数的整数表示））
iriscli service bind --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --service-name=<service-name> --def-chain-id=<def-chain-id> --bind-type=Local  --deposit=1000iris --prices=1iris --avg-rsp-time=10000 --usable-time=9999

# 查询服务绑定
iriscli service binding --def-chain-id=<def-chain-id> --service-name=<service-name> --bind-chain-id=<bind-chain-id> --provider=<provider-account-address>

# 查询服务绑定列表
iriscli service bindings --def-chain-id=<def-chain-id> --service-name=<service-name>

# 服务绑定更新
iriscli service update-binding --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --service-name=<service-name> --def-chain-id=<def-chain-id> --bind-type=Local  --deposit=1iris --prices=1iris,2iris --avg-rsp-time=10000 --usable-time=100

# 禁用服务绑定
iriscli service disable --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --def-chain-id=<def-chain-id> --service-name=<service-name>

# 开启服务绑定, 并追加抵押100iris
iriscli service enable --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --def-chain-id=<def-chain-id> --service-name=<service-name> --deposit=100iris

# 取回押金
iriscli service refund-deposit --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --def-chain-id=<def-chain-id> --service-name=<service-name>
```

### 服务调用

服务消费者如果需要发起服务调用请求，需要支付服务提供方指定的服务费。服务提供方需要在`MaxRequestTimeout`定义的区块高度内响应该服务请求，如果超时未响应，将从服务提供方的该服务绑定押金中扣除`SlashFraction`比例的押金，同时该次服务调用的服务费将退还到服务消费者的退费池中。如果服务调用被正常响应，系统从该次服务调用的服务费中将扣除`ServiceFeeTax`比例的系统税收，同时将剩余的服务费加入到服务提供方的收入池中。服务提供方/消费者可以发起`withdraw-fees`/`refund-fees`交易取回自己在收入池/退费池中所有的token。

```bash
# 发起服务调用
iriscli service call --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --def-chain-id=<def-chain-id> --service-name=<service-name> --method-id=1 --bind-chain-id=<bind-chain-id> --provider=<provider-account-address> --service-fee=1iris --request-data=<request-data>

# 查询服务请求列表
iriscli service requests --def-chain-id=<def-chain-id> --service-name=<service-name> --bind-chain-id=<bind-chain-id> --provider=<provider-account-address>

# 响应服务调用
iriscli service respond --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --request-chain-id=<request-chain-id> --request-id=<request-id (e.g.230-130-0)> --response-data=<response-data>

# 查询服务响应
iriscli service response --request-chain-id=<request-chain-id> --request-id=<request-id (e.g.230-130-0)>

# 查询指定地址的服务费退款和收入
iriscli service fees <account-address>

# 从服务费退款中退还所有费用
iriscli service refund-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris

# 从服务费收入中取回所有费用
iriscli service withdraw-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
```

## IDL文件扩展

在使用proto文件对服务的方法，输入、输出参数进行标准化定义时，可通过注释的方式增加method属性。

### 注释标准

* 使用 `//@Attribute 属性： 值`的方式添加在rpc方法上，即可将该属性添加为方法的属性。例如:

    > //@Attribute description: sayHello

### 目前支持的属性

* `description` 对该方法的描述
* `output-privacy` 是否对该方法的输出进行加密处理，{`NoPrivacy`,`PubKeyEncryption`}
* `output-cached` 是否对该方法的输出进行缓存，{`OffChainCached`,`NoCached`}

### IDL content参照

* IDL content参照

    > syntax = \"proto3\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output-privacy: NoPrivacy\n    //@Attribute output-cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n

* IDL文件参照

    [test.proto](https://github.com/irisnet/irishub/blob/master/docs/features/test.proto)
