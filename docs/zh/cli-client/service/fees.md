# iriscli service fees 

## 描述

查询指定地址的服务费退款和收入

## 用法

```
iriscli service fees <service_provider_address>
```

## 标志

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| -h, --help            |                         | 服务费查询命令帮助                                                                                                                                         |          |

## 示例

### 查询服务费
```shell
iriscli service fees <service_provider_address>
```

运行成功以后，返回的结果如下:
```json
{
  "returned_fee": [
    {
      "denom": "iris-atto",
      "amount": "10000000000000000"
    }
  ],
  "incoming_fee": [
    {
      "denom": "iris-atto",
      "amount": "10000000000000000"
    }
  ]
}
```

