# Service

Service模块允许在IRIS Hub中定义、绑定、调用服务。[了解更多iService内容](../features/service.md)。

## 可用命令

| 名称                                                    | 描述                                         |
| ------------------------------------------------------- | -------------------------------------------- |
| [define](#iris-tx-service-define)                       | 定义一个新的服务                             |
| [definition](#iris-query-service-definition)                | 查询服务定义                                 |
| [bind](#iris-tx-service-bind)                           | 绑定一个服务                                 |
| [binding](#iris-query-service-binding)                      | 查询服务绑定                                 |
| [bindings](#iris-query-service-bindings)                    | 查询服务绑定列表                             |
| [set-withdraw-addr](#iris-tx-service-set-withdraw-addr) | 设置服务提供者的提取地址                     |
| [withdraw-addr](#iris-query-service-withdraw-addr)          | 查询服务提供者的提取地址                     |
| [update-binding](#iris-tx-service-update-binding)       | 更新一个存在的服务绑定                       |
| [disable](#iris-tx-service-disable)                     | 禁用一个可用的服务绑定                       |
| [enable](#iris-tx-service-enable)                       | 启用一个不可用的服务绑定                     |
| [refund-deposit](#iris-tx-service-refund-deposit)       | 退还一个服务绑定的所有押金                   |
| [call](#iris-tx-service-call)                           | 发起服务调用                                 |
| [request](#iris-query-service-request)                      | 通过请求ID查询服务请求                       |
| [requests](#iris-query-service-requests)                    | 通过服务绑定或请求上下文查询服务请求列表     |
| [respond](#iris-tx-service-respond)                     | 响应服务请求                                 |
| [response](#iris-query-service-response)                    | 通过请求ID查询服务响应                       |
| [responses](#iris-query-service-responses)                  | 通过请求上下文ID和批次计数器查询服务响应列表 |
| [request-context](#iris-query-service-request-context)      | 查询请求上下文                               |
| [update](#iris-tx-service-update)                       | 更新请求上下文                               |
| [pause](#iris-tx-service-pause)                         | 暂停一个正在进行的请求上下文                 |
| [start](#iris-tx-service-start)                         | 启动一个暂停的请求上下文                     |
| [kill](#iris-tx-service-kill)                           | 终止请求上下文                               |
| [fees](#iris-query-service-fees)                            | 查询服务提供者的收益                         |
| [withdraw-fees](#iris-tx-service-withdraw-fees)         | 提取服务提供者的收益                         |
| [schema](#iris-query-service-schema)                        | 通过 schema 名称查询系统 schema              |

## iris tx service define

定义一个新的服务。

```bash
iris tx service define [flags]
```

**标志：**

| 名称，速记           | 默认 | 描述                            | 必须 |
| -------------------- | ---- | ------------------------------- | ---- |
| --name               |      | 服务名称                        | 是   |
| --description        |      | 服务的描述                      |      |
| --author-description |      | 服务创建者的描述                |      |
| --tags               |      | 服务的标签列表                  |      |
| --schemas            |      | 服务接口schemas的内容或文件路径 | 是   |

### 定义一个新的服务

```bash
iris tx service define --chain-id=irishub --from=<key-name> --fees=0.3iris 
--name=<service name> --description=<service description> --author-description=<author description>
--tags=tag1,tag2 --schemas=<schemas content or path/to/schemas.json>
```

### Schemas内容示例

```json
{"input":{"$schema":"http://json-schema.org/draft-04/schema#","title":"BioIdentify service input","description":"BioIdentify service input specification","type":"object","properties":{"id":{"description":"id","type":"string"},"name":{"description":"name","type":"string"},"data":{"description":"data","type":"string"}},"required":["id","data"]},"output":{"$schema":"http://json-schema.org/draft-04/schema#","title":"BioIdentify service output","description":"BioIdentify service output specification","type":"object","properties":{"data":{"description":"result data","type":"string"}},"required":["data"]}}
```

## iris query service definition

查询服务定义。

```bash
iris query service definition [service-name] [flags]
```

### 查询一个服务定义

查询指定服务定义的详细信息。

```bash
iris query service definition <service name>
```

## iris tx service bind

绑定一个服务。

```bash
iris tx service bind [flags]
```

**标志：**

| 名称，速记     | 默认 | 描述                                                         | 必须 |
| -------------- | ---- | ------------------------------------------------------------ | ---- |
| --service-name |      | 服务名称                                                     | 是   |
| --deposit      |      | 服务绑定的押金                                               | 是   |
| --pricing      |      | 服务定价内容或文件路径，是一个[Irishub Service Pricing JSON Schema](../features/service-pricing.md)实例 | 是   |
| --qos          |      | 最小响应时间                                                 | 是   |
| --provider     |      | 服务提供者地址

### 绑定一个存在的服务定义

抵押`deposit`应该满足最小抵押数量需求，最小抵押数量为`price` * `MinDepositMultiple` 和 `MinDeposit`中的最大值（`MinDepositMultiple`以及`MinDeposit`是可治理参数）。

```bash
iris tx service bind --chain-id=irishub --from=<key-name> --fees=0.3iris
--service-name=<service name> --deposit=10000iris --pricing=<pricing content or path/to/pricing.json> --qos=50
```

### Pricing内容示例

```json
{
    "price": "1iris"
}
```

## iris query service binding

查询服务绑定。

```bash
iris query service binding [service-name] [provider] [flags]
```

## iris query service bindings

查询服务绑定列表。

```bash
iris query service bindings [service-name] [flags]
```

### 查询服务绑定列表

```bash
iris query service bindings <service name> <owner address>
```

## iris tx service update-binding

更新服务绑定。

```bash
iris tx service update-binding [service-name] [provider-address] [flags]
```

**标志：**
| 名称，速记 | 默认 | 描述                                                         | 必须 |
| ---------- | ---- | ------------------------------------------------------------ | ---- |
| --deposit  |      | 增加的绑定押金，为空则不更新                                 |      |
| --pricing  |      | 服务定价内容或文件路径，是一个[Irishub Service Pricing JSON Schema](../features/service-pricing.md)实例，为空则不更新 |      |
| --qos      |      | 最小响应时间，为0则不更新                                    |      |

### 更新一个存在的服务绑定

更新服务绑定，追加 10 IRIS 的抵押。

```bash
iris tx service update-binding <service-name> <prvider-address>  --chain-id=irishub --from=<key-name> --fees=0.3iris --deposit=10iris
```

## iris tx service set-withdraw-addr

设置服务提供者的提取地址。

```bash
iris tx service set-withdraw-addr [withdrawal-address] [flags]
```

## iris query service withdraw-addr

查询服务提供者的提取地址。

```bash
iris query service withdraw-addr [provider] [flags]
```

## iris tx service disable

禁用一个可用的服务绑定。

```bash
iris tx service disable [service-name] [provider-address] [flags]
```

## iris tx service enable

启用一个不可用的服务绑定。

```bash
iris tx service enable [service-name] [provider-address] [flags]
```

**标志：**

| 名称，速记 | 默认 | 描述               | 必须 |
| ---------- | ---- | ------------------ | ---- |
| --deposit  |      | 启用绑定增加的押金 |      |

### 启用一个不可用的服务绑定

启用一个不可用的服务绑定，追加 10 IRIS 的抵押。

```bash
iris tx service enable <service name> <provider-address> --chain-id=irishub --from=<key-name> --fees=0.3iris --deposit=10iris
```

## iris tx service refund-deposit

从一个服务绑定中退还所有的押金。

```bash
iris tx service refund-deposit [service-name] [provider-address] [flags]
```

### 取回所有押金

取回抵押之前，必须先[禁用](#iris-tx-service-disable)服务绑定。

```bash
iris tx service refund-deposit <service name>  <provider-address> --chain-id=irishub --from=<key-name> --fees=0.3iris
```

## iris tx service call

发起服务调用。

```bash
iris tx service call [flags]
```

**标志：**

| 名称，速记        | 默认  | 描述                                                  | 必须 |
| ----------------- | ----- | ----------------------------------------------------- | ---- |
| --service-name     |       | 服务名称                                              | 是   |
| --providers       |       | 服务提供者列表                                        | 是   |
| --service-fee-cap |       | 愿意为单个请求支付的最大服务费用                      | 是   |
| --data            |       | 请求输入的内容或文件路径，是一个Input JSON Schema实例 | 是   |
| --timeout         |       | 请求超时                                              | 是   |
| --super-mode      | false | 签名者是否为超级用户                                  |      |
| --repeated        | false | 请求是否为重复性的                                    |      |
| --frequency       |       | 重复性请求的请求频率；默认为`timeout`值               |      |
| --total           |       | 重复性请求的请求总数，-1表示无限制                    |      |

### 发起一个服务调用请求

```bash
iris tx service call --chain-id=irishub --from=<key name> --fees=0.3iris --service-name=<service name> --providers=<provider list> --service-fee-cap=1iris --data=<request input or path/to/input.json> --timeout=100 --repeated --frequency=150 --total=100
```

### 请求输入示例

```json
{
    "id": "1",
    "name": "irisnet",
    "data": "facedata"
}
```

## iris query service request

通过请求ID查询服务请求。

```bash
iris query service request [request-id] [flags]
```

### 查询一个服务请求

```bash
iris query service request <request-id>
```

:::tip
你可以从[按高度获取区块信息](./tendermint.md#iris-query-tendermint-block)的结果中获取`request-id`。
:::

## iris query service requests

通过服务绑定或请请求上下文ID查询服务请求列表。

```bash
iris query service requests [service-name] [provider] | [request-context-id] [batch-counter] [flags]
```

### 查询服务绑定的活跃请求

```bash
iris query service requests <service name> <provider>
```

### 通过请求上下文ID和批次计数器查询服务请求列表

```bash
iris query service requests <request-context-id> <batch-counter>
```

## iris tx service respond

响应服务请求。

```bash
iris tx service respond [flags]
```

**标志：**

| 名称，速记   | 默认 | 描述                                                         | 必须 |
| ------------ | ---- | ------------------------------------------------------------ | ---- |
| --request-id |      | 欲响应请求的ID                                               | 是   |
| --result     |      | 响应结果的内容或文件路径, 是一个[Irishub Service Result JSON Schema](../features/service-result.md)实例 | 是   |
| --data       |      | 响应输出的内容或文件路径, 是一个Output JSON Schema实例       |      |

### 响应一个服务请求

```bash
iris tx service respond --chain-id=irishub --from=<key-name> --fees=0.3iris --request-id=<request-id> --result=<response result or path/to/result.json> --data=<response output or path/to/output.json>
```

:::tip
你可以从[按高度获取区块信息](./tendermint.md#iris-query-tendermint-block)的结果中获取`request-id`。
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

## iris query service response

通过请求ID查询服务响应。

```bash
iris query service response [request-id] [flags]
```

### 查询一个服务响应

```bash
iris query service response [request-id] [flags]
```

:::tip
你可以从[按高度获取区块信息](./tendermint.md#iris-query-tendermint-block)的结果中获取`request-id`。
:::

## iris query service responses

通过请求上下文ID以及批次计数器查询服务响应列表。

```bash
iris query service responses [request-context-id] [batch-counter] [flags]
```

### 根据指定的请求上下文ID以及批次计数器查询服务响应

```bash
iris query service responses <request-context-id> <batch-counter>
```

## iris query service request-context

查询请求上下文。

```bash
iris query service request-context [request-context-id] [flags]
```

### 查询一个请求上下文

```bash
iris query service request-context <request-context-id>
```

:::tip
你可以从[调用服务](#iris-tx-service-call)的结果中获取`request-context-id`
:::

## iris tx service update

更新请求上下文。

```bash
iris tx service update [request-context-id] [flags]
```

**标志：**

| 名称，速记        | 默认 | 描述                                           | 必须 |
| ----------------- | ---- | ---------------------------------------------- | ---- |
| --providers       |      | 服务提供者列表，为空则不更新                   |      |
| --service-fee-cap |      | 愿意为单个请求支付的最大服务费用，为空则不更新 |      |
| --timeout         |      | 请求超时，为0则不更新                          |      |
| --frequency       |      | 请求频率，为0则不更新                          |      |
| --total           |      | 请求总数，为0则不更新                          |      |

### 更新一个请求上下文

```bash
iris tx service update <request-context-id> --chain-id=irishub --from=<key name> --fees=0.3iris --providers=<provider list> --service-fee-cap=1iris --timeout=0 --frequency=150 --total=100
```

## iris tx service pause

暂停一个正在进行的请求上下文。

```bash
iris tx service pause [request-context-id] [flags]
```

### 暂停一个正在进行的请求上下文

```bash
iris tx service pause <request-context-id>
```

## iris tx service start

启动一个暂停的请求上下文。

```bash
iris tx service start [request-context-id] [flags]

## iris tx service kill

终止请求上下文。

```bash
iris tx service kill [request-context-id] [flags]
```

### 终止一个请求上下文

```bash
iris tx service kill <request-context-id>
```

## iris query service fees

查询服务提供者的收益。

```bash
iris query service fees [provider] [flags]
```

## iris tx service withdraw-fees

提取服务提供者的收益。

```bash
iris tx service withdraw-fees [flags]
```

## iris query service schema

通过 schema 名称查询系统 schema。有效的 schema 名称为 `pricing`（服务定价）和 `result`（响应结果）。

```bash
iris query service schema [schema-name] [flags]
```

### 查询服务定价 schema

```bash
iris query service schema pricing
```

### 查询响应结果 schema

```bash
iris query service schema result
```

