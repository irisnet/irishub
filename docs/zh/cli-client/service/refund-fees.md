# iriscli service refund-fees 

## 描述

从服务费退款中退还所有费用

## 用法

```
iriscli service refund-fees <flags>
```

## 示例

### 从服务费退款中退还费用 

```shell
iriscli service refund-fees --chain-id=<chain-id> --from=<key_name> --fee=0.3iris
```

运行成功以后，返回的结果如下:

```txt
Committed at block 79 (tx hash: 1E3A690028116E0DF541A840BF5830598EAD4154F4374B2A4042911C27D68C64, response:
 {
   "code": 0,
   "data": null,
   "log": "Msg 0: ",
   "info": "",
   "gas_wanted": 200000,
   "gas_used": 3912,
   "codespace": "",
   "tags": {
     "action": "service_refund_fees"
   }
 })
```

