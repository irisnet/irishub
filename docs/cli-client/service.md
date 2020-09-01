# iris tx service

Service module allows you to define, bind, invoke services on the IRIS Hub. [Read more about iService](../features/service.md).

## 可用命令

| Name                                                    | Description                                                 |
| ------------------------------------------------------- | ----------------------------------------------------------- |
| [define](#iris-tx-service-define)                       | Define a new service                                        |
| [definition](#iris-q-service-definition)               | Query a service definition                                  |
| [bind](#iris-tx-service-bind)                           | Bind a service                                              |
| [binding](#iris-q-service-binding)                     | Query a service binding                                     |
| [bindings](#iris-q-service-bindings)                   | Query all bindings of a service definition                  |
| [set-withdraw-addr](#iris-tx-service-set-withdraw-addr) | Set a withdrawal address for a provider                     |
| [withdraw-addr](#iris-q-service-withdraw-addr)         | Query the withdrawal address of a provider                  |
| [update-binding](#iris-tx-service-update-binding)       | Update an existing service binding                          |
| [disable](#iris-tx-service-disable)                     | Disable an available service binding                        |
| [enable](#iris-tx-service-enable)                       | Enable an unavailable service binding                       |
| [refund-deposit](#iris-tx-service-refund-deposit)       | Refund all deposit from a service binding                   |
| [call](#iris-tx-service-call)                           | Initiate a service call                                     |
| [request](#iris-q-service-request)                     | Query a request by the request ID                           |
| [requests](#iris-q-service-requests)                   | Query active requests by the service binding or request context ID   |
| [respond](#iris-tx-service-respond)                     | Respond to a service request                                |
| [response](#iris-q-service-response)                   | Query a response by the request ID                          |
| [responses](#iris-q-service-responses)                 | Query active responses by the request context ID and batch counter |
| [request-context](#iris-q-service-request-context)     | Query a request context                                     |
| [update](#iris-tx-service-update)                       | Update a request context                                    |
| [pause](#iris-tx-service-pause)                         | Pause a running request context                             |
| [start](#iris-tx-service-start)                         | Start a paused request context                              |
| [kill](#iris-tx-service-kill)                           | Terminate a request context                                 |
| [fees](#iris-q-service-fees)                           | Query the earned fees of a provider                         |
| [withdraw-fees](#iris-tx-service-withdraw-fees)         | Withdraw the earned fees of a provider                      |
| [schema](#iris-q-service-schema)         | Query the system schema by the schema name       |

## iris tx service define

Define a new service.

```bash
iris tx service define [flags]
```

**Flags:**

| Name, shorthand      | Default | Description                               | Required |
| -------------------- | ------- | ----------------------------------------- | -------- |
| --name               |         | Service name                              | Yes      |
| --description        |         | Service description                       |          |
| --author-description |         | Service author description                |          |
| --tags               |         | Service tags                              |          |
| --schemas            |         | Content or file path of service interface schemas  | Yes      |

### define a service

```bash
iris tx service define --chain-id=irishub --from=<key-name> --fees=0.3iris 
--name=<service name> --description=<service description> --author-description=<author description>
--tags=tag1,tag2 --schemas=<schemas content or path/to/schemas.json>
```

### Schemas content example

```json
{"input":{"$schema":"http://json-schema.org/draft-04/schema#","title":"BioIdentify service input","description":"BioIdentify service input specification","type":"object","properties":{"id":{"description":"id","type":"string"},"name":{"description":"name","type":"string"},"data":{"description":"data","type":"string"}},"required":["id","data"]},"output":{"$schema":"http://json-schema.org/draft-04/schema#","title":"BioIdentify service output","description":"BioIdentify service output specification","type":"object","properties":{"data":{"description":"result data","type":"string"}},"required":["data"]}}
```

## iris q service definition

Query a service definition.

```bash
iris q service definition [service-name] [flags]
```

### Query a service definition

Query the detailed info of the service definition with the specified service name.

```bash
iris q service definition <service name>
```

## iris tx service bind

Bind a service.

```bash
iris tx service bind [flags]
```

**Flags:**

| Name, shorthand | Default | Description                                                  | Required |
| --------------- | ------- | ------------------------------------------------------------ | -------- |
| --service-name  |         | Service name                                                 | Yes      |
| --deposit       |         | Deposit of the binding                                       | Yes      |
| --pricing       |         | Pricing content or file path, which is an instance of [Irishub Service Pricing JSON Schema](../features/service-pricing.json) | Yes      |
| --qos           |         | Minimum response time                                        | Yes      |
| --provider      |         | provider address, default to the owner                       |          |

### Bind an existing service definition

The deposit needs to satisfy the minimum deposit requirement, which is the maximal one between `price` * `MinDepositMultiple` and `MinDeposit`(`MinDepositMultiple` and `MinDeposit` are the system parameters, which can be modified through the governance).

```bash
iris tx service bind --chain-id=irishub --from=<key-name> --fees=0.3iris
--service-name=<service name> --deposit=10000iris --pricing=<pricing content or path/to/pricing.json> --min-resp-time=50
```

### Pricing content example

```json
{
    "price": "1iris"
}
```

## iris q service binding

Query a service binding.

```bash
iris q service binding <service name> <provider>
```

### Query a service binding

```bash
iris q service binding [service-name] [provider] [flags]
```

## iris q service bindings

Query all bindings of a service definition.

```bash
iris q service bindings [service-name] [flags]
```

### Query service binding list

```bash
iris q service bindings <service name> --owner=<address>
```

## iris tx service update-binding

Update a service binding.

```bash
iris tx service update-binding [service-name] [provider-address] [flags]
```

**Flags:**

| Name, shorthand | Default | Description                                                                         | Required |
| --------------- | ------- | ----------------------------------------------------------------------------------- | -------- |
| --deposit       |         | Deposit added for the binding, not updated if empty                                                     |          |
| --pricing       |         | Pricing content or file path, which is an instance of [Irishub Service Pricing JSON Schema](../features/service-pricing.md), not updated if empty |          |
| --qos |         | Minimum response time, not updated if set to 0 |  |

### Update an existing service binding

The following example updates the service binding with the additional 10 IRIS deposit

```bash
iris tx service update-binding <service-name> --chain-id=irishub --from=<key-name> --fees=0.3iris --deposit=10iris
```

## iris tx service set-withdraw-addr

Set a withdrawal address for a provider.

```bash
iris tx service set-withdraw-addr [withdrawal-address] [flags]
```

### Set a withdrawal address

```bash
iris tx service set-withdraw-addr <withdrawal address> --chain-id=irishub --from=<key-name> --fees=0.3iris
```

## iris q service withdraw-addr

Query the withdrawal address of a provider.

```bash
iris q service withdraw-addr [provider] [flags]
```

### Query the withdrawal address of a provider

```bash
iris q service withdraw-addr <provider>
```

## iris tx service disable

Disable an available service binding.

```bash
iris tx service disable [service-name] [provider-address] [flags]
```

### Disable an available service binding

```bash
iris tx service disable <service name> --chain-id=irishub --from=<key-name> --fees=0.3iris
```

## iris tx service enable

Enable an unavailable service binding.

```bash
iris tx service enable [service-name] [provider-address] [flags]
```

**Flags:**

| Name, shorthand | Default | Description                            | Required |
| --------------- | ------- | -------------------------------------- | -------- |
| --deposit       |         | deposit added for enabling the binding |          |

### Enable an unavailable service binding

The following example enables an unavailable service binding with the additional 10 IRIS deposit.

```bash
iris tx service enable <service name> --chain-id=irishub --from=<key-name> --fees=0.3iris --deposit=10iris
```

## iris tx service refund-deposit

Refund all deposits from a service binding.

```bash
iris tx service refund-deposit [service-name] [provider-address] [flags]
```

### Refund all deposits from an unavailable service binding

Before refunding, you should [disable](#iris-tx-service-disable) the service binding first.

```bash
iris tx service refund-deposit <service name> --chain-id=irishub --from=<key-name> --fees=0.3iris
```

## iris tx service call

Initiate a service call.

```bash
iris tx service call [flags]
```

**Flags:**

| Name, shorthand   | Default | Description                                                          | Required |
| ----------------- | ------- | -------------------------------------------------------------------- | -------- |
| --service-name    |         | Service name                                                         | Yes      |
| --providers       |         | Provider list to request                                             | Yes      |
| --service-fee-cap |         | Maximum service fee to pay for a single request                      | Yes      |
| --data            |         | Content or file path of the request input, which is an Input JSON Schema instance | Yes      |
| --timeout         |         | Request timeout                                                      |   Yes       |
| --super-mode      | false   | Indicate if the signer is a super user                               |
| --repeated        | false   | Indicate if the reqeust is repetitive                                |          |
| --frequency       |         | Request frequency when repeated, default to `timeout`                |          |
| --total           |         | Request count when repeated, -1 means unlimited                      |          |

### Initiate a service invocation request

```bash
iris tx service call --chain-id=irishub --from=<key name> --fees=0.3iris --service-name=<service name>
--providers=<provider list> --service-fee-cap=1iris --data=<request input or path/to/input.json> --timeout=100 --repeated --frequency=150 --total=100
```

### Input example

```json
{
    "id": "1",
    "name": "irisnet",
    "data": "facedata"
}
```

## iris q service request

Query a request by the request ID.

```bash
iris q service request [request-id] [flags]
```

### Query a service request

```bash
iris q service request <request-id>
```

:::tip
You can retrieve the `request-id` in the result of [tendermint block](./tendermint.md#iris-q-tendermint-block)
:::

## iris q service requests

Query active requests by the service binding or request context ID.

```bash
iris q service requests [service-name] [provider] | [request-context-id] [batch-counter] [flags]
```

### Query active requests of a service binding

```bash
iris q service requests <service name> <provider>
```

### Query service requests by the request context ID and batch counter

```bash
iris q service requests <request-context-id> <batch-counter>
```

## iris tx service respond

Respond to a service request.

```bash
iris tx service respond [flags]
```

**Flags:**

| Name, shorthand | Default | Description                                                             | Required |
| --------------- | ------- | ----------------------------------------------------------------------- | -------- |
| --request-id    |         | ID of the request to respond to                                         | Yes      |
| --result        |         | Content or file path of the response result, which is an instance of [Irishub Service Result JSON Schema](../features/service-result.json)  | Yes      |
| --data          |         | Content or file path of the response output, which is an Output JSON Schema instance |          |

### Respond to a service request

```bash
iris tx service respond --chain-id=irishub --from=<key-name> --fees=0.3iris
--request-id=<request-id> --result=<response result or path/to/result.json> --data=<response output or path/to/output.json>
```

:::tip
You can retrieve the `request-id` in the result of [tendermint block](./tendermint.md#iriscli-tendermint-block)
:::

### Result example

```
{
    "code": 200,
    "message": ""
}
```

### Output example

```json
{
    "data": "userdata"
}
```

## iris q service response

Query a service response.

```bash
iris q service response [request-id] [flags]
```

### Query a service response

```bash
iris q service response <request-id>
```

:::tip
You can retrieve the `request-id` in the result of [tendermint block](./tendermint.md#iris-q-tendermint-block)
:::

## iris q service responses

Query active responses by the request context ID and batch counter.

```bash
iris q service responses [request-context-id] [batch-counter] [flags]
```

### Query responses by the request context ID and batch counter

```bash
iris q service responses <request-context-id> <batch-counter>
```

## iris q service request-context

Query a request context.

```bash
iris q service request-context [request-context-id] [flags]
```

### Query a request context

```bash
iris q service request-context <request-context-id>
```

:::tip
You can retrieve the `request-context-id` in the result of [service call](#iris-tx-service-call)
:::

## iris tx service update

Update a request context.

```bash
iris tx service update [request-context-id] [flags]
```

**Flags:**

| Name, shorthand   | Default | Description                                                           | Required |
| ----------------- | ------- | --------------------------------------------------------------------- | -------- |
| --providers       |         | Provider list to request, not updated if empty                        |          |
| --service-fee-cap |         | Maximum service fee to pay for a single request, not updated if empty |          |
| --timeout         |         | Request timeout, not updated if set to 0                              |          |
| --frequency       |         | Request frequency, not updated if set to 0                            |          |
| --total           |         | Request count, not updated if set to 0                                |          |

### Update a request context

```bash
iris tx service update <request-context-id> --chain-id=irishub --from=<key name> --fees=0.3iris
--providers=<provider list> --service-fee-cap=1iris --timeout=0 --frequency=150 --total=100
```

## iris tx service pause

Pause a running request context.

```bash
iris tx service pause [request-context-id] [flags]
```

### Pause a running request context

```bash
iris tx service pause <request-context-id>
```

## iris tx service start

Start a paused request context.

```bash
iris tx service start [request-context-id] [flags]
```

### Start a paused request context

```bash
iris tx service start <request-context-id>
```

## iris tx service kill

Terminate a request context.

```bash
iris tx service kill [request-context-id] [flags]
```

### Kill a request context

```bash
iris tx service kill <request-context-id>
```

## iris q service fees

Query the earned fees of a provider.

```bash
iris q service fees [provider] [flags]
```

### Query service fees

```bash
iris q service fees <provider>
```

## iris tx service withdraw-fees

Withdraw the earned fees of a provider.

```bash
iris tx service withdraw-fees [provider-address] [flags]
```

### Withdraw the earned fees

```bash
iris tx service withdraw-fees [provider-address] --chain-id=irishub --from=<key-name> --fees=0.3iris
```

## iris q service schema

Query the system schema by the schema name, only pricing and result allowed.

```bash
iris q service schema [schema-name] [flags]
```

### Query the service pricing schema

```bash
iris q service schema pricing
```

### Query the response result schema

```bash
iris q service schema result
```

