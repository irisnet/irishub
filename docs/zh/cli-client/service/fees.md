# iriscli service fees 

## Description

查询指定地址的服务费退款和收入

## Usage

```
iriscli service fees [account address]
```

## Flags

| Name, shorthand       | Default                 | Description                                                                                                                                           | Required |
| --------------------- | ----------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| -h, --help            |                         | 服务费查询命令帮助                                                                                                                                         |          |

## Examples

### 查询服务费
```shell
iriscli service fees faa1f02ext9duk7h3rx9zm7av0pnlegxve8ne5vw6x
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

