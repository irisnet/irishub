# iriscli distribution rewards

## 介绍

查询验证人或委托人的所有收益

## 用法

```
iriscli distribution rewards <address> [flags]
```

打印帮助信息：
```
iriscli distribution rewards --help
```

## 示例

```
iriscli distribution rewards <address>
```
执行结果示例：
```json
{
  "total_rewards": [
    {
      "denom": "iris-atto",
      "amount": "235492744310548933957"
    }
  ],
  "delegations": [
    {
      "validator": "iva1q7602ujxxx0urfw7twm0uk5m7n6l9gqsgw4pqy",
      "reward": [
        {
          "denom": "iris-atto",
          "amount": "2148801198916275"
        }
      ]
    }
  ],
  "commission": [
    {
      "denom": "iris-atto",
      "amount": "235490595509350017682"
    }
  ]
}
```