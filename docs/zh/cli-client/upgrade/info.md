# iriscli upgrade info

## 描述

查询软件版本信息和升级模块的信息

## 用法

```
iriscli upgrade info
```

打印帮助信息:

```
iriscli upgrade info --help
```
## 标志

| 名称, 速记       | 默认值                     | 描述                                                        | 必需     |
| --------------- | -------------------------- | ----------------------------------------------------------------- | -------- |

## 例子

查询当前版本信息

```
iriscli upgrade info 
```

然后它会打印如下内容：

```
{
  "current_proposal_id": "0",
  "current_proposal_accept_height": "-1",
  "version": {
    "Id": "0",
    "ProposalID": "0",
    "Start": "0",
    "ModuleList": [
      {
        "Start": "0",
        "End": "9223372036854775807",
        "Handler": "bank",
        "Store": [
          "acc"
        ]
      },
      {
        "Start": "0",
        "End": "9223372036854775807",
        "Handler": "stake",
        "Store": [
          "stake",
          "acc",
          "mint",
          "distr"
        ]
      },
      .......
    ]
  }
}
```
