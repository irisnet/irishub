# iriscli service update-binding 

## 描述

更新一个存在的服务绑定

## 用法

```
iriscli service update-binding [flags]
```

## 标志
| Name, shorthand       | Default                 | Description                                                                        | Required |
| --------------------- | ----------------------- | ---------------------------------------------------------------------------------- | -------- |
| --avg-rsp-time        |                         | [int]  服务平均返回时间的毫秒数表示                                                     | 是       |
| --bind-type           |                         | [string] 对服务是本地还是全局的设置，可选值Local/Global                                  | 是       |
| --def-chain-id        |                         | [string] 定义该服务的区块链ID                                                          | 是       |
| --deposit             |                         | [string] 绑定押金, 将会增加当前服务绑定押金                                               |          |
| --prices              |                         | [strings] 服务定价,按照服务方法排序的定价列表                                             |          |
| --service-name        |                         | [string] 服务名称                                                                    | 是       |
| --usable-time         |                         | [int] 每一万次服务调用的可用性的整数表示                                                  | 是       |
| -h, --help            |                         | 绑定更新命令帮助                                                                       |          |

## 示例

### 更新一个存在的服务绑定
```shell
iriscli service update-binding --chain-id=test-irishub --from=node0 --fee=0.004iris --service-name=test-service --def-chain-id=test-irishub --bind-type=Local --deposit=1iris --prices=1iris --avg-rsp-time=10000 --usable-time=100
```

运行成功以后，返回的结果如下:

```json
{
   "tags": {
     "action": "service-update-binding",
     "completeConsumedTxFee-iris-atto": "\"100840000000000\""
   }
 }
```

