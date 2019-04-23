# iriscli distribution validator-distr-info

## 介绍

查询验证人的收益分配信息

## 用法

```
iriscli distribution validator-distr-info <validator_address> <flags>
```

打印帮助信息:

```
iriscli distribution validator-distr-info --help
```

## 示例

```
iriscli distribution validator-distr-info <validator_address>
```

执行结果示例
```json
{
  "operator_addr": "iva1e7wljxhz7u7xrh63xjlds8vcy047a47ejpnz7a",
  "fee_pool_withdrawal_height": "101290",
  "del_accum": {
    "update_height": "101290",
    "accum": "0.0000000000"
  },
  "del_pool": "0.0000000000000000000000000000iris",
  "val_commission": "12.8560369893449408111336573478iris"
}
```