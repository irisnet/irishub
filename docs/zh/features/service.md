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

为响应用户请求，服务提供者发送一个由 _结果_ 对象和可选的输出对象组成的响应。结果code等于200时输出对象必须提供。_结果_ 对象必须符合此[schema](service-result.json)。下面是一个示例：

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

通过创建对现有服务定义的绑定, 任何用户都可以提供相应的服务。绑定主要由四个部分组成：_提供者地址_、_定价_、_押金_ 以及 _服务质量_。

### 提供者地址

提供者地址是 _服务提供者_ （即链外服务进程）用来监听请求的一个端点。在服务提供者能够接受并处理服务请求之前，其运营者或所有者必须为它创建一个链上地址，并且发起一个 `绑定` 交易将这个地址关联到相关的服务定义。

为调用一个服务，用户或消费者通过发起一个请求交易向一个有效服务绑定的提供者地址发起请求；服务提供者检测和处理这个请求，并且通过一个响应交易发送处理结果。

### 定价

定价指定服务提供者如何对其提供的服务收费。定价必须符合此[schema](service-pricing.json)。下面是一个示例：

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
      "end_time": "2020-06-30T23:59:59Z",
      "discount": 0.9
    }
  ]
}
```

服务提供者能选择接受 `iris` 以外的 tokens 作为服务费用，例如 `0.03link`。价格是消费者从多个提供相同服务的提供者中遴选的一个考虑因素。

### 押金

运营一个服务提供者意味着重要的服务责任，因此创建服务绑定需要一定数量的押金。押金数量必须大于 _押金阈值_，该值为 `MinDepositMultiple * price` 与 `MinDeposit` 两者中的最大值。如果服务提供者未能在超时之前响应请求，则其绑定押金的一小部分，即 `SlashFraction * deposit` 将被罚没和销毁。如果押金降至阈值以下，服务绑定将被暂时禁用，直到其所有者增加足够的押金重新激活。

> **_提示：_**  `service/MinDepositMultiple`、`service/MinDeposit` 和 `service/SlashFraction`是可以通过链上[治理](governance.md)更改的系统参数。

### 服务质量

服务质量承诺是根据提供者将服务响应发送回区块链所需的平均区块数来声明的。这是消费者选择潜在提供者时考虑的另一个因素。

### 命令

服务绑定可以由其所有者随时更新，以调整定价、增加押金或者改变 QoS；也可以被禁用和重新启用。如果服务提供者的所有者不想再提供服务，则需要禁用绑定并等待一段时间，才能取回押金。

```bash
# 创建服务绑定
iriscli service bind <service-name> <provider-address> <deposit> <min-response-time> <pricing-json or path/to/pricing.json>

# 更新服务绑定
iriscli service update-binding <service-name> --deposit=<added-deposit> --min-resp-time=<min-response-time> --pricing=<pricing-json or path/to/pricing.json>

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

### 请求上下文

消费者通过创建一个 _请求上下文_ 来指定如何调用一个服务。_请求上下文_ 像智能合约一样自动产生实际的请求。_请求上下文_ 由一些参数组成，可以大致分为如下五组：

#### 目标和输入

* _服务名_：要调用的目标服务的名称
* _输入数据_：符合目标服务输入 schema 的 json 格式数据

#### 提供者过滤

* _提供者列表_：逗号分隔的候选服务提供者的地址列表
* _服务费上限_：消费者愿意为每次调用支付的最大服务费用
* _超时_：消费者为接收响应愿意等待的区块数

#### 响应处理

* _模块_：包含回调函数的模块名称
* _响应阈值_：为调用回调函数所需接收的最小响应数

> **_提示：_** 这两个参数不能从 CLI 和 API 设置；它们只对使用 iService 的其他模块可用，比如 [oracle](oracle.md) 和 [random](random.md)。

#### 重复性

* _重复_：指示请求上下文是否可重复的一个布尔标志
* _频率_：重复调用批次之间的区块间隔数
* _总数_: 重复调用批次的总数，负数表示无限



```bash
```
