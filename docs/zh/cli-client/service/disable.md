# iriscli service disable 

## 描述

禁用一个可用的服务绑定

## 用法

```
iriscli service disable <flags>
```

## 特有标志

| Name, shorthand       | Default                 | Description                                                                        | Required |
| --------------------- | ----------------------- | ---------------------------------------------------------------------------------  | -------- |
| --def-chain-id        |                         | 定义该服务的区块链ID                                                         | 是       |
| --service-name        |                         | 服务名称                                                                   | 是       |

## 示例

### 禁用一个可用的服务绑定

```shell
iriscli service disable --chain-id=<chain-id>  --from=<key_name> --fee=0.3iris --def-chain-id=<service_define_chain_id> --service-name=<service>
```

运行成功以后，返回的结果如下:

```txt
Committed at block 588 (tx hash: 12333DCD222EB3620A2177E9A8E0D84E248BAE0D3BC445274E09A19096794A46, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 3123,
   "codespace": "",
   "tags": {
     "action": "service_disable"
   }
 })
```