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
| [set-withdraw-addr](#iriscli-service-set-withdraw-addr)             | 设置服务提供者的提取地址                              |
| [withdraw-addr](#iriscli-service-withdraw-addr)             | 查询服务提供者的提取地址                              |
| [update-binding](#iriscli-service-update-binding) | 更新一个存在的服务绑定         |
| [disable](#iriscli-service-disable)               | 禁用一个可用的服务绑定         |
| [enable](#iriscli-service-enable)                 | 启用一个不可用的服务绑定       |
| [refund-deposit](#iriscli-service-refund-deposit) | 退还一个服务绑定的所有押金                   |
| [call](#iriscli-service-call)                     | 调用服务                  |
| [request](#iriscli-service-request)             | 通过请求ID查询服务请求                          |
| [requests](#iriscli-service-requests)             | 通过服务绑定或请求上下文查询服务请求列表              |
| [respond](#iriscli-service-respond)               | 响应服务请求                 |
| [response](#iriscli-service-response)             | 通过请求ID查询服务响应                   |
| [responses](#iriscli-service-responses)             | 通过请求上下文ID和批次计数器查询服务响应列表                            |
| [request-context](#iriscli-service-request-context)             | 查询请求上下文                             |
| [update](#iriscli-service-update)             | 更新请求上下文                             |
| [pause](#iriscli-service-pause)             | 暂停一个正在进行的请求上下文                             |
| [start](#iriscli-service-start)             | 启动一个暂停的请求上下文                             |
| [kill](#iriscli-service-kill)             | 终止请求上下文                            |
| [fees](#iriscli-service-fees)                     | 查询服务提供者的收益 |
| [withdraw-fees](#iriscli-service-withdraw-fees)   | 提取服务提供者的收益     |

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
| --schemas         |      | 服务接口的schemas内容或路径           | 是   |

### 定义一个新的服务

```bash
iriscli service define --chain-id=<chain-id> --from=<key-name> --fee=0.3iris 
--name=<service name> --description=<service description> --author-description=<author description>
--tags=tag1,tag2 --schemas=<schemas content or path/to/schemas.json>
```

### Schemas内容示例

```json
{"input":{"$schema":"http://json-schema.org/draft-04/schema#","title":"BioIdentify service input","description":"BioIdentify service input specification","type":"object","properties":{"id":{"description":"id","type":"string"},"name":{"description":"name","type":"string"},"data":{"description":"data","type":"string"}},"required":["id","data"]},"output":{"$schema":"http://json-schema.org/draft-04/schema#","title":"BioIdentify service output","description":"BioIdentify service output specification","type":"object","properties":{"data":{"description":"result data","type":"string"}},"required":["data"]}}
```

## iriscli service definition

查询服务定义。

```bash
iriscli service definition <service name>
```

**标志：**

| 名称，速记     | 默认 | 描述                 | 必须 |
| -------------- | ---- | -------------------- | ---- |
| --service-name |      | 服务名称             | 是   |

### 查询一个服务定义

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

```json
{
    "price": [
        {
            "denom": "iris-atto",
            "amount": "1000000000000000000"
        }
    ]
}
```

## iriscli service binding

查询服务绑定。

```bash
iriscli service binding <service name> <provider>
```

### 查询一个服务绑定

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
| --deposit      |      | 增加的绑定押金                  |      |
| --pricing      |      | 服务定价内容或路径，需符合Irishub Pricing JSON schema              |      |

### 更新一个存在的服务绑定

更新服务绑定，追加 10 IRIS 的抵押。

```bash
iriscli service update-binding --chain-id=<chain-id> --from=<key-name> --fee=0.3iris 
--service-name=<service-name> --deposit=10iris
```

## iriscli service set-withdraw-addr

设置服务提供者的提取地址

```bash
iriscli service set-withdraw-addr <withdrawal address>
```

### 设置一个提取地址

```bash
iriscli service set-withdraw-addr <withdrawal address> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
```

## iriscli service withdraw-addr

查询服务提供者的提取地址

```bash
iriscli service withdraw-addr <provider>
```

### 查询一个服务提供者的提取地址

```bash
iriscli service withdraw-addr <provider>
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
| --deposit |      | 启用绑定增加的押金 |      |

### 启用一个不可用的服务绑定

启用一个不可用的服务绑定，追加 10 IRIS 的抵押。

```bash
iriscli service enable <service name> --chain-id=<chain-id>  --from=<key-name> --fee=0.3iris --deposit=10iris
```

## iriscli service refund-deposit

从一个服务绑定中退还所有的押金。

```bash
iriscli service refund-deposit <service name>
```

### 取回所有押金

取回抵押之前，必须先[禁用](#iriscli-service-disable)服务绑定。

```bash
iriscli service refund-deposit <service name> --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
```

## iriscli service call

调用服务。

```bash
iriscli service call <flags>
```

**标志：**

| 名称，速记      | 默认 | 描述                           | 必须 |
| --------------- | ---- | ------------------------------ | ---- |
| --name  |      | 服务名称                       | 是   |
| --providers     |         | 服务提供者列表                             | 是      |
| --service-fee-cap |         | 愿意为单个请求支付的最大服务费用        | 是     |
| --data      |         | 请求的输入，是一个Input JSON schema实例 | 是      |
| --timeout | | 请求超时 | |
| --super-mode| false | 签名者是否为超级用户 |
| --repeated   |    false     | 请求是否为重复性的                |          |
| --frequency   |         | 重复性请求的请求频率；默认为`timeout`值              |          |
| --total  |         | 重复性请求的请求总数，-1表示无限制    |          |

### 发起一个服务调用请求

```bash
iriscli service call --chain-id=<chain-id> --from=<key name> --fee=0.3iris --service-name=<service name>
--providers=<provider list> --service-fee-cap=1iris --data=<request data> --timeout=100 --repeated --frequency=150 --total=100
```

### 请求输入示例

```json
{
    "id": "1",
    "name": "irisnet",
    "data": "facedata"
}
```

## iriscli service request

通过请求ID查询服务请求。

```bash
iriscli service request <request-id>
```

### 查询一个服务请求

```bash
iriscli service request <request-id>
```

:::tip
你可以从[按高度获取区块信息](#iriscli-tendermint-block)的结果中获取`request-id`。
:::

## iriscli service requests

通过服务绑定或请请求上下文ID查询服务请求列表。

```bash
iriscli service requests [<service name> <provider>] | [<request-context-id> <batch-counter>]
```

### 查询服务绑定的活跃请求

```bash
iriscli service requests <service name> <provider>
```

### 通过请求上下文ID和批次计数器查询服务请求列表

```bash
iriscli service requests <request-context-id> <batch-counter>
```

## iriscli service respond

响应服务请求。

```bash
iriscli service respond <flags>
```

**标志：**

| 名称，速记         | 默认 | 描述                      | 必须 |
| ------------------ | ---- | ------------------------- | ---- |
| --request-id       |      | 欲响应请求的ID            | 是   |
| --result | | 服务响应的结果, 是一个Result JSON schema实例  | 是 |
| --data    |      | 服务响应的输出, 是一个Output JSON schema实例 | |

### 响应一个服务请求

```bash
iriscli service respond --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
--request-id=<request-id> --result=<response result> --data=<response output>
```

:::tip
你可以从[按高度获取区块信息](#iriscli-tendermint-block)的结果中获取`request-id`。
:::

### 响应结果示例

```json
{
    "code": 200,
    "message": ""
}
```

### 响应输出示例

```json
{
    "data": "userdata"
}
```

## iriscli service response

通过请求ID查询服务响应。

```bash
iriscli service response <request-id>
```

### 查询一个服务响应

```bash
iriscli service response <request-id>
```

:::tip
你可以从[按高度获取区块信息](#iriscli-tendermint-block)的结果中获取`request-id`。
:::

## iriscli service responses

通过请求上下文ID以及批次计数器查询服务响应列表。

```bash
iriscli service responses <request-context-id> <batch-counter>
```

### 根据指定的请求上下文ID以及批次计数器查询服务响应

```bash
iriscli service responses <request-context-id> <batch-counter>
```

## iriscli service request-context

查询请求上下文。

```bash
iriscli service request-context <request-context-id>
```

### 查询一个请求上下文

```bash
iriscli service request-context <request-context-id>
```

:::tip
你可以从[调用服务](#iriscli-service-call)的结果中获取`request-context-id`
:::

## iriscli service update

更新请求上下文。

```bash
iriscli service update <request-context-id> <flags>
```

**标志：**

| 名称，速记         | 默认 | 描述                      | 必须 |
| --------------- | ------- | -------------------------------------------------- | -------- |
| --providers     |         | 服务提供者列表，为空则不更新                            | Yes      |
| --service-fee-cap |         | 愿意为单个请求支付的最大服务费用, 为空则不更新       |      |
| --timeout | | 请求超时, 为0则不更新 | |
| --frequency   |         | 请求频率, 为0则不更新           |          |
| --total  |         | 请求总数, 为0则不更新    |          |

### 更新一个请求上下文

```bash
iriscli service update <request-context-id> --chain-id=<chain-id> --from=<key name> --fee=0.3iris
--providers=<provider list> --service-fee-cap=1iris --timeout=0 --frequency=150 --total=100
```

## iriscli service pause

暂停一个正在进行的请求上下文。

```bash
iriscli service pause <request-context-id>
```

### 暂停一个正在进行的请求上下文

```bash
iriscli service pause <request-context-id>
```

## iriscli service start

启动一个暂停的请求上下文。

```bash
iriscli service start <request-context-id>
```

### 启动一个暂停的请求上下文

```bash
iriscli service start <request-context-id>
```

## iriscli service kill

终止请求上下文。

```bash
iriscli service kill <request-context-id>
```

### 终止一个请求上下文

```bash
iriscli service kill <request-context-id>
```

## iriscli service fees

查询服务提供者的收益。

```bash
iriscli service fees <provider>
```

### 查询服务提供者的收益

```bash
iriscli service fees <provider>
```

## iriscli service withdraw-fees

提取服务提供者的收益。

```bash
iriscli service withdraw-fees <flags>
```

### 提取服务提供者的收益

```bash
iriscli service withdraw-fees --chain-id=<chain-id> --from=<key-name> --fee=0.3iris
```
