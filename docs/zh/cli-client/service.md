# iriscli service

Service模块允许在IRIS Hub中定义、绑定、调用服务。[了解更多iService内容](../features/service.md)。

## 可用命令

| 名称                                              | 描述                           |
| ------------------------------------------------- | ------------------------------ |
| [define](#iriscli-service-define)                 | 定义一个新的服务           |
| [definition](#iriscli-service-definition)         | 查询服务定义                   |
| [bind](#iriscli-service-bind)                     | 绑定一个服务        |
| [binding](#iriscli-service-binding)               | 查询服务绑定                     |
| [bindings](#iriscli-service-bindings)             | 查询服务绑定列表               |
| [set-withdraw-addr](#iriscli-service-set-withdraw-addr)             | 设置提取地址                              |
| [withdraw-addr](#iriscli-service-withdraw-addr)             | 查询提取地址                              |
| [update-binding](#iriscli-service-update-binding) | 更新一个存在的服务绑定         |
| [disable](#iriscli-service-disable)               | 禁用一个可用的服务绑定         |
| [enable](#iriscli-service-enable)                 | 启用一个不可用的服务绑定       |
| [refund-deposit](#iriscli-service-refund-deposit) | 取回所有押金                   |
| [call](#iriscli-service-call)                     | 调用服务                  |
| [request](#iriscli-service-request)             | 查询服务请求                          |
| [requests](#iriscli-service-requests)             | 查询服务请求列表               |
| [respond](#iriscli-service-respond)               | 响应一个服务请求                 |
| [response](#iriscli-service-response)             | 查询服务响应                   |
| [responses](#iriscli-service-responses)             | 查询服务响应列表                             |
| [request-context](#iriscli-service-request-context)             | 查询请求上下文                             |
| [update](#iriscli-service-update)             | 更新请求上下文                             |
| [pause](#iriscli-service-pause)             | 暂停运行的请求上下文                             |
| [start](#iriscli-service-start)             | 恢复暂停的请求上下文                             |
| [kill](#iriscli-service-kill)             | 终止请求上下文                            |
| [fees](#iriscli-service-fees)                     | 查询指定服务提供者的收益 |
| [withdraw-fees](#iriscli-service-withdraw-fees)   | 提取指定服务提供者的收益     |

## iriscli service define

定义一个新的服务。

```bash
iriscli service define <flags>
```

**标志：**

| 名称，速记            | 默认 | 描述                                     | 必须 |
| --------------------- | ---- | ---------------------------------------- | ---- |
| --name        |      | 服务名称                                 | 是   |
| --description |      | 服务的描述                               |      |
| --author-description  |      | 服务创建者的描述                         |      |
| --tags                |      | 服务的标签列表                        |      |
| --schemas         |      | 服务接口的schemas内容或路径           |      |

### define a service

```bash
iriscli service define --chain-id=<chain-id> --from=<key-name> --fee=0.3iris 
--name=<service name> --description=<service description> --author-description=<author description>
--tags=tag1,tag2 --schemas=<schemas content or path/to/schemas.json>
```

### Schemas内容示例


## iriscli service definition

查询服务定义。

```bash
iriscli service definition <service name>
```

**标志：**

| 名称，速记     | 默认 | 描述                 | 必须 |
| -------------- | ---- | -------------------- | ---- |
| --service-name |      | 服务名称             | 是   |

### 查询服务定义

查询指定服务定义的详细信息。

```bash
iriscli service definition <service name>
```

## iriscli service bind

绑定一个服务。

```bash
iriscli service bind <flags>
```

**标志：**

| 名称，速记     | 默认 | 描述                                                | 必须 |
| -------------- | ---- | --------------------------------------------------- | ---- |
| --service-name |      | 服务名称                                            | 是   |
| --deposit      |      | 服务绑定的押金                                  | 是   |
| --pricing       |      | 服务定价内容或路径，需符合Irishub Pricing JSON schema      |      |

### 绑定一个存在的服务定义

抵押`deposit`应该满足最小抵押数量需求，最小抵押数量为`price` * `MinDepositMultiple` 和 `MinDeposit`中的最大值（`MinDepositMultiple`以及`MinDeposit`是可治理参数）。

```bash
iriscli service bind --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
--service-name=<service name> --deposit=10000iris --pricing=<pricing>
```

### Pricing内容示例

## iriscli service binding

查询服务绑定。

```bash
iriscli service binding <service name> <provider>
```

### 查询服务绑定

```bash
iriscli service binding <service name> <provider>
```

## iriscli service bindings

查询服务绑定列表。

```bash
iriscli service bindings <service name>
```

### 查询服务绑定列表

```bash
iriscli service bindings <service name>
```

## iriscli service update-binding

更新服务绑定。

```bash
iriscli service update-binding <flags>
```

**标志：**
| 名称，速记     | 默认 | 描述                                                | 必须 |
| -------------- | ---- | --------------------------------------------------- | ---- |
| --service-name |      | 服务名称                                            | 是   |
| --deposit      |      | 增加的押金                  |      |
| --pricing      |      | 新的服务定价内容或路径，需符合Irishub Pricing JSON schema              |      |

### 更新一个存在的服务绑定

更新服务绑定，追加10iris的抵押。

```bash
iriscli service update-binding --chain-id=<chain-id> --from=<key-name> --fee=0.3iris 
--service-name=<service-name> --deposit=10iris
```

## iriscli service disable

禁用一个可用的服务绑定。

```bash
iriscli service disable <service name>
```

### 禁用一个可用的服务绑定

```bash
iriscli service disable <service name> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
```

## iriscli service enable

启用一个不可用的服务绑定。

```bash
iriscli service enable <service name> <flags>
```

**标志：**

| 名称，速记       | 默认 | 描述                               | 必须 |
| ---------------- | ---- | ---------------------------------- | ---- |
| --deposit |      | 增加的押金 |      |

### 启用一个不可用的服务绑定

启用一个不可用的服务绑定，并且追加10iris的抵押。

```bash
iriscli service enable <service name> --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --deposit=10iris
```

## iriscli service refund-deposit

取回所有押金。

```bash
iriscli service refund-deposit <service name>
```

### 取回所有押金

取回抵押之前，必须先[禁用](#iriscli-service-disable)服务绑定。

```bash
iriscli service refund-deposit <service name> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
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
