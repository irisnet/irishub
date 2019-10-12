# iriscli service

Service模块允许在IRIS Hub中通过区块链定义、绑定、调用服务。[了解更多iService内容](../features/service.md)。

## 可用命令

| 名称                                              | 描述                           |
| ------------------------------------------------- | ------------------------------ |
| [define](#iriscli-service-define)                 | 创建一个新的服务定义           |
| [definition](#iriscli-service-definition)         | 查询服务定义                   |
| [bind](#iriscli-service-bind)                     | 创建一个新的服务绑定           |
| [binding](#iriscli-service-binding)               | 询服务绑定                     |
| [bindings](#iriscli-service-bindings)             | 查询服务绑定列表               |
| [update-binding](#iriscli-service-update-binding) | 更新一个存在的服务绑定         |
| [disable](#iriscli-service-disable)               | 禁用一个可用的服务绑定         |
| [enable](#iriscli-service-enable)                 | 启用一个不可用的服务绑定       |
| [refund-deposit](#iriscli-service-refund-deposit) | 取回所有押金                   |
| [call](#iriscli-service-call)                     | 调用服务方法                   |
| [requests](#iriscli-service-requests)             | 查询服务请求列表               |
| [respond](#iriscli-service-respond)               | 响应服务调用                   |
| [response](#iriscli-service-response)             | 查询服务响应                   |
| [fees](#iriscli-service-fees)                     | 查询指定地址的服务费退款和收入 |
| [refund-fees](#iriscli-service-refund-fees)       | 从服务费退款中退还所有费用     |
| [withdraw-fees](#iriscli-service-withdraw-fees)   | 从服务费收入中取回所有费用     |

## iriscli service define

创建一个新的服务定义。

```bash
iriscli service define <flags>
```

**标志：**

| 名称，速记            | 默认 | 描述                                     | 必须 |
| --------------------- | ---- | ---------------------------------------- | ---- |
| --service-description |      | 服务的描述                               |      |
| --author-description  |      | 服务创建者的描述                         |      |
| --service-name        |      | 服务名称                                 | 是   |
| --tags                |      | 该服务的关键字                           |      |
| --idl-content         |      | 对该服务描述的接口定义语言内容           |      |
| --file                |      | 对该服务描述的接口定义语言内容的文件路径 |      |

### define a service

```bash
iriscli service define --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --service-name=<service-name> --service-description=<service-description> --author-description=<author-description> --tags=tag1,tag2 --idl-content=<idl-content> --file=test.proto
```

如果文件项不是空的，将会替换Idl-content.  [IDL内容示例](#IDL内容示例)。

### IDL内容示例

* IDL内容示例

    > syntax = \\"proto3\\";\n\npackage helloworld;\n\n// The greeting service definition.\nservice Greeter {\n    //@Attribute description: sayHello\n    //@Attribute output-privacy: NoPrivacy\n    //@Attribute output-cached: NoCached\n    rpc SayHello (HelloRequest) returns (HelloReply) {}\n}\n\n// The request message containing the user's name.\nmessage HelloRequest {\n    string name = 1;\n}\n\n// The response message containing the greetings\nmessage HelloReply {\n    string message = 1;\n}\n

* IDL文件示例

    [test.proto](https://github.com/irisnet/irishub/blob/master/docs/features/test.proto)

## iriscli service definition

查询服务定义。

```bash
iriscli service definition <flags>
```

**标志：**

| 名称，速记     | 默认 | 描述                 | 必须 |
| -------------- | ---- | -------------------- | ---- |
| --def-chain-id |      | 定义该服务的区块链ID | 是   |
| --service-name |      | 服务名称             | 是   |

### 查询服务定义

查询指定Chain ID和服务名称的服务定义的详细信息。

```bash
iriscli service definition --def-chain-id=<service-define-chain-id> --service-name=<service-name>
```

## iriscli service bind

创建一个新的服务绑定。

```bash
iriscli service bind <flags>
```

**标志：**

| 名称，速记     | 默认 | 描述                                                | 必须 |
| -------------- | ---- | --------------------------------------------------- | ---- |
| --avg-rsp-time |      | 服务平均返回时间的毫秒数表示                        | 是   |
| --bind-type    |      | 对服务是本地还是全局的设置，可选值“Local”或“Global” | 是   |
| --def-chain-id |      | 定义该服务的区块链ID                                | 是   |
| --deposit      |      | 服务提供者的保证金                                  | 是   |
| --prices       |      | 服务定价，按照服务方法排序的定价列表                |      |
| --service-name |      | 服务名称                                            | 是   |
| --usable-time  |      | 每一万次服务调用的可用性的整数表示                  | 是   |

### 添加服务绑定到已存在的服务定义

在服务绑定中，你需要抵押`deposit`指定数量的IRIS，最小的抵押金额为该服务的服务费价格`price` * genesis中定义的`MinDepositMultiple`倍数。

```bash
iriscli service bind --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --service-name=<service-name> --def-chain-id=<service-define-chain-id> --bind-type=Local --deposit=1000iris --prices=1iris --avg-rsp-time=10000 --usable-time=9999
```

## iriscli service binding

查询服务绑定。

```bash
iriscli service binding <flags>
```

**标志：**

| 名称，速记      | 默认 | 描述                               | 必须 |
| --------------- | ---- | ---------------------------------- | ---- |
| --bind-chain-id |      | 绑定该服务的区块链ID               | 是   |
| --def-chain-id  |      | 定义该服务的区块链ID               | 是   |
| --provider      |      | 服务提供者的区块链地址(bech32编码) | 是   |
| --service-name  |      | 服务名称                           | 是   |

### 查询服务绑定

```bash
iriscli service binding --def-chain-id=<service-define-chain-id> --service-name=<service-name> --bind-chain-id=<service-bind-chain-id> --provider=<provider-address>
```

## iriscli service bindings

查询服务绑定列表。

```bash
iriscli service bindings <flags>
```

**标志：**

| 名称，速记     | 默认 | 描述                 | 必须 |
| -------------- | ---- | -------------------- | ---- |
| --def-chain-id |      | 定义该服务的区块链ID | 是   |
| --service-name |      | 服务名称             | 是   |

### Query service binding list

```bash
iriscli service bindings --def-chain-id=<chain-id> --service-name=<service-name>
```

## iriscli service update-binding

更新服务绑定。

```bash
iriscli service update-binding <flags>
```

**标志：**
| 名称，速记     | 默认 | 描述                                                | 必须 |
| -------------- | ---- | --------------------------------------------------- | ---- |
| --avg-rsp-time |      | 服务平均返回时间的毫秒数表示                        |      |
| --bind-type    |      | 对服务是本地还是全局的设置，可选值“Local”或“Global” |      |
| --def-chain-id |      | 定义该服务的区块链ID                                | 是   |
| --deposit      |      | 绑定押金，将会增加当前服务绑定押金                  |      |
| --prices       |      | 服务定价，按照服务方法排序的定价列表                |      |
| --service-name |      | 服务名称                                            | 是   |
| --usable-time  |      | 每一万次服务调用的可用性的整数表示                  |      |

### 更新一个存在的服务绑定

更新服务绑定，并且追加10iris的抵押。

```bash
iriscli service update-binding --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --service-name=<service-name> --def-chain-id=<service-define-chain-id> --bind-type=Local --deposit=10iris --prices=1iris --avg-rsp-time=10000 --usable-time=9999
```

## iriscli service disable

禁用一个可用的服务绑定。

```bash
iriscli service disable <flags>
```

**标志：**

| 名称，速记     | 默认 | 描述                 | 必须 |
| -------------- | ---- | -------------------- | ---- |
| --def-chain-id |      | 定义该服务的区块链ID | 是   |
| --service-name |      | 服务名称             | 是   |

### 禁用一个可用的服务绑定

```bash
iriscli service disable --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --def-chain-id=<service-define-chain-id> --service-name=<service>
```

## iriscli service enable

启用一个不可用的服务绑定。

```bash
iriscli service enable <flags>
```

**标志：**

| 名称，速记       | 默认 | 描述                               | 必须 |
| ---------------- | ---- | ---------------------------------- | ---- |
| --def-chain-id   |      | 定义该服务的区块链ID               | 是   |
| --deposit string |      | 绑定押金, 将会增加当前服务绑定押金 |      |
| --service-name   |      | 服务名称                           | 是   |

### 启用一个不可用的服务绑定

启用一个不可用的服务绑定，并且追加10iris的抵押。

```bash
iriscli service enable --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --def-chain-id=<service-define-chain-id> --service-name=<service-name> --deposit=10iris
```

## iriscli service refund-deposit

取回所有押金。

```bash
iriscli service refund-deposit <flags>
```

**标志：**

| 名称，速记     | 默认 | 描述                 | 必须 |
| -------------- | ---- | -------------------- | ---- |
| --def-chain-id |      | 定义该服务的区块链ID | 是   |
| --service-name |      | 服务名称             | 是   |

### 取回所有押金

取回抵押之前，必须先[禁用](#iriscli-service-disable)服务绑定。

```bash
iriscli service refund-deposit --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --def-chain-id=<service-define-chain-id> --service-name=<service-name>
```

## iriscli service call

调用服务方法。

```bash
iriscli service call <flags>
```

**标志：**

| 名称，速记      | 默认 | 描述                           | 必须 |
| --------------- | ---- | ------------------------------ | ---- |
| --def-chain-id  |      | 定义该服务的区块链ID           | 是   |
| --service-name  |      | 服务名称                       | 是   |
| --method-id     |      | 调用的服务方法ID               | 是   |
| --bind-chain-id |      | 绑定该服务的区块链ID           | 是   |
| --provider      |      | bech32编码的服务提供商账户地址 | 是   |
| --service-fee   |      | 服务调用支付的服务费           |      |
| --request-data  |      | hex编码的服务调用请求数据      |      |

### 发起一个服务调用请求

```bash
iriscli service call --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --def-chain-id=<service-define-chain-id> --service-name=<service-name> --method-id=1 --bind-chain-id=<service-bind-chain-id> --provider=<provider-address> --service-fee=1iris --request-data=<request-data>
```

## iriscli service requests

查询服务请求列表。

```bash
iriscli service requests <flags>
```

**标志：**

| 名称，速记      | 默认 | 描述                           | 必须 |
| --------------- | ---- | ------------------------------ | ---- |
| --def-chain-id  |      | 定义该服务的区块链ID           | 是   |
| --service-name  |      | 服务名称                       | 是   |
| --bind-chain-id |      | 绑定该服务的区块链ID           | 是   |
| --provider      |      | bech32编码的服务提供商账户地址 | 是   |

### 查询服务请求列表

```bash
iriscli service requests --def-chain-id=<service-define-chain-id> --service-name=<service-name> --bind-chain-id=<service-bind-chain-id> --provider=<provider-address>
```

## iriscli service respond

响应服务调用。

```bash
iriscli service respond <flags>
```

**标志：**

| 名称，速记         | 默认 | 描述                      | 必须 |
| ------------------ | ---- | ------------------------- | ---- |
| --request-chain-id |      | 发起该服务调用的区块链ID  | 是   |
| --request-id       |      | 该服务调用的ID            | 是   |
| --response-data    |      | hex编码的服务调用响应数据 |

### 响应一个服务调用

```bash
iriscli service respond --chain-id=<chain-id> --from=<key-name> --fee=0.3iris --request-chain-id=<request-chain-id> --request-id=<request-id> --response-data=<response-data>
```

:::tip
你可以从[服务调用](#iriscli-service-call)的返回结果中得到`request-id`。
:::

## iriscli service response

查询服务响应。

```bash
iriscli service response <flags>
```

**标志：**

| 名称，速记         | 默认 | 描述                     | 必须 |
| ------------------ | ---- | ------------------------ | ---- |
| --request-chain-id |      | 发起该服务调用的区块链ID | 是   |
| --request-id       |      | 该服务调用的ID           | 是   |

### 查询服务响应

```bash
iriscli service response --request-chain-id=<request-chain-id> --request-id=<request-id>
```

:::tip
你可以从[服务调用](#iriscli-service-call)的返回结果中得到`request-id`。
:::

## iriscli service fees

查询指定地址的服务费退款和收入。

### 查询服务费

```bash
iriscli service fees <service-provider-address>
```

示例输出:

```json
{
  "returned-fee": [
    {
      "denom": "iris-atto",
      "amount": "10000000000000000"
    }
  ],
  "incoming-fee": [
    {
      "denom": "iris-atto",
      "amount": "10000000000000000"
    }
  ]
}
```

## iriscli service refund-fees

从服务费退款中退还所有费用。

```bash
iriscli service refund-fees <flags>
```

### 从服务费退款中退还费用

```bash
iriscli service refund-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
```

## iriscli service withdraw-fees

从服务费收入中取回所有费用。

```bash
iriscli service withdraw-fees <flags>
```

### 从服务费收入中取回费用

```bash
iriscli service withdraw-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
```
