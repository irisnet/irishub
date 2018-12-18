# iriscli service enable 

## 描述

启用一个不可用的服务绑定

## 用法

```
iriscli service enable [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                       | Required |
| --------------------- | ----------------------- | --------------------------------------------------------------------------------- | -------- |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID                                                         | 是       |
| --deposit string      |                         | [string] 绑定押金, 将会增加当前服务绑定押金                                             |          |
| --service-name        |                         | [string] 服务名称                                                                   | 是       |
| -h, --help            |                         | 启用命令帮助                                                                         |          |

## 示例

### 启用一个不可用的服务绑定
```shell
iriscli service enable --chain-id=test  --from=node0 --fee=0.004iris --def-chain-id=test --service-name=test-service
```

运行成功以后，返回的结果如下:

```txt
Committed at block 599 (tx hash: 87AD0FE53A665114EB3CF40505821C63952F56E9E7EF844481167C1D7B026432, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 5447,
   "codespace": "",
   "tags": {
     "action": "service_enable"
   }
 })
```