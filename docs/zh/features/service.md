# Service

> **_提示：_** 本文档中显示的命令仅供说明。有关命令的准确语法，请参阅[cli docs](../client/service.md)。

## 简介

IRIS服务（又称iService）旨在弥合区块链和传统应用之间的鸿沟。它规范化了链外服务的定义和绑定（提供者注册），促进了调用以及与这些服务的交互，并能调解服务治理过程（分析和争议解决）。

## 服务定义

### 服务接口 schema

任何用户都可以在区块链上定义服务。服务的接口即输入和输出需要使用[JSON Schema](https://JSON-Schema.org/)来指定。下面是一个示例：

```json
{
  "input": {
    "$schema": "http://json-schema.org/draft-04/schema#",
    "title": "service-def-input-example",
    "description": "Schema for a service input example",
    "type": "object",
    "properties": {
      "base": {
        "description": "base token denom",
        "type": "string"
      },
      "quote": {
        "description": "quote token denom",
        "type": "string"
      }
    },
    "required": [
      "base",
      "quote"
    ]
  },
  "output": {
    "$schema": "http://json-schema.org/draft-04/schema#",
    "title": "service-def-output-example",
    "description": "Schema for a service output example",
    "type": "object",
    "properties": {
      "price": {
        "description": "price",
        "type": "number"
      }
    },
    "required": [
      "price"
    ]
  }
}
```

### 服务结果 schema

服务提供者为响应用户请求（包含一个输入对象），发送回一个由结果对象和可选的输出对象组成的响应。结果code等于200时输出对象必须提供。结果对象必须符合此[schema](service-result.json)，下面是一个示例：

```json
{
  "result" : {
    "code": 400,
    "message": "user input out of range"
  }
}
```

一旦服务定义就绪，就可以通过执行以下命令将其发布到区块链：

```bash
# 创建服务定义
iriscli service define <service-name> <schemas-json or path/to/schemas.json> --description=<service-description> --author-description=<author-description> --tags=<tag1,tag2,...>

# 查询服务定义
iriscli service definition <service-name>
```

## 服务绑定

任何人通过创建对现有服务定义的绑定，就可以提供相应的服务。绑定主要由三个组件组成：提供者地址（执行`绑定`交易的账户地址）、定价和押金。

### 提供者地址

消费者应该能够向目标提供者地址发起服务请求（输入），并能获取从该地址返回的响应（输出）。

### 定价

定价必须符合此[schema](service-pricing.json)。下面是一个定价示例：

```json
{
  "price": "0.1iris",
  "promotions_by_time": [
    {
      "start_time": "2020-01-01T00:00:00Z",
      "end_time": "2020-03-31T23:59:59Z",
      "discount": 0.7
    },
    {
      "start_time": "2020-04-01T00:00:00Z",
      "end_time": "2019-06-30T23:59:59Z",
      "discount": 0.9
    }
  ]
}
```

### 押金

创建服务绑定需要一定的押金用于服务承诺。押金必须大于 _押金阈值_，该值由`max(DepositMultiple*price,MinDeposit)`得出。如果服务提供者未能在超时之前响应请求，则其绑定押金的一小部分，即`SlashFraction*deposit`将被扣除。如果押金降至阈值以下，服务绑定将被暂时禁用，直到其所有者增加足够的押金重新激活。

> **_提示：_**  `service/DepositMultiple`、`service/MinDeposit`和`service/SlashFraction`是可以通过链上[治理](governance.md)更改的系统参数。

### 生命周期

服务绑定可以由其所有者随时更新，以调整定价或增加押金；也可以被禁用和重新启用。如果服务绑定所有者不想再提供服务，则需要禁用绑定并等待一段时间，然后才能取回押金。

```bash
# 创建服务绑定
iriscli service bind <service-name> <deposit> <pricing-json or path/to/pricing.json>

# 更新服务绑定
iriscli service update-binding <service-name> --deposit=<added-deposit> --pricing=<pricing-json or path/to/pricing.json>

# 设置收益提取地址
iriscli service set-withdraw-addr <withdrawal-address>

# 提取收益到指定的提取地址
iriscli service withdraw-fees

# 启用服务绑定
iriscli service enable <service-name> <added-deposit>

# 禁用服务绑定
iriscli service disable <service-name>

# 取回服务绑定的押金
iriscli service refund-deposit <service-name>

# 受托人提取服务税到指定地址
iriscli service withdraw-tax <destination-address> <withdrawal-amount>

# 查询服务绑定
iriscli service binding <service-name> <provider-address>

# 查询一个服务的绑定列表
iriscli service bindings <service-name>

# 查询服务提供者的收益提取地址
iriscli service withdraw-addr <provider-address>

# 查询服务提供者的收益
iriscli service fees <provider-address>

# 查询系统 schemas（有效的 schema 名称: pricing, result）
iriscli service schema <schema-name>
```

## 服务调用

服务消费者如果需要发起服务调用请求，需要支付服务提供方指定的服务费。服务提供方需要在MaxRequestTimeout定义的区块高度内响应该服务请求，如果超时未响应，将从服务提供方的该服务绑定押金中扣除SlashFraction比例的押金，同时该次服务调用的服务费将退还到服务消费者的退费池中。如果服务调用被正常响应，系统从该次服务调用的服务费中将扣除ServiceFeeTax比例的系统税收，同时将剩余的服务费加入到服务提供方的收入池中。服务提供方/消费者可以发起withdraw-fees/refund-fees交易取回自己在收入池/退费池中所有的token。

```bash
```
