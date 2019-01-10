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

## 例子

查询当前版本信息

```
iriscli upgrade info 
```

然后它会打印当前生效的协议版本信息和当前正在准备升级的协议信息：

```
{
  "version": {
    "ProposalID": "1",
    "Success": true,
    "Protocol": {
      "version": "1",
      "software": "https://github.com/irisnet/irishub/tree/v0.10.0",
      "height": "30"
    }
  },
  "upgrade_config": {
    "ProposalID": "3",
    "Definition": {
      "version": "2",
      "software": "https://github.com/irisnet/irishub/tree/v0.10.1",
      "height": "8000"
    }
  }
}
```
