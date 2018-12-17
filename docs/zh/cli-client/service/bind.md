# iriscli service bind 

## 描述

创建一个新的服务绑定

## 用法

```
iriscli service bind [flags]
```

## 标志

| Name, shorthand       | Default                 | Description                                                                        | Required |
| --------------------- | ----------------------- | ---------------------------------------------------------------------------------- | -------- |
| --avg-rsp-time        |                         | [int]  服务平均返回时间的毫秒数表示                                                     | 是       |
| --bind-type           |                         | [string] 对服务是本地还是全局的设置，可选值Local/Global                                  | 是       |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID                                                          | 是       |
| --deposit             |                         | [string] 服务提供者的保证金                                                            | 是       |
| --prices              |                         | [strings] 服务定价,按照服务方法排序的定价列表                                             |          |
| --service-name        |                         | [string] 服务名称                                                                    | 是       |
| --usable-time         |                         | [int] 每一万次服务调用的可用性的整数表示                                                  | 是       |
| -h, --help            |                         | 绑定命令帮助                                                                          |          |

## 示例

### 添加服务绑定到已存在的服务定义
```shell
iriscli service bind --chain-id=test --from=node0 --fee=0.004iris --service-name=test-service --def-chain-id=test --bind-type=Local --deposit=1iris --prices=1iris --avg-rsp-time=10000 --usable-time=100
```

运行成功以后，返回的结果如下:

```txt
Committed at block 567 (tx hash: A48DBD217CBB843E72CC47B40F90CE7DEEEDD6437C86A74A2976ADC9F449A034, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 5893,
   "codespace": "",
   "tags": {
     "action": "service_bind"
   }
 })
```

