# 导出区块链状态

## 介绍

这里介绍一种能导出区块链状态，并以json格式返回给用户。如果把返回的json字符串保存到一个json文件里，那么这个json文件可以作为一个新区块链网络的创世块。导出区块链状态所用的命令为`iris export`

## 用法

```
iris export [flags]
```

## 标志

| 名称，速记          | 类型   | 是否必填 | 默认值   | 介绍    |
| ------------------- | -----  | -------- | -------- | -------------- |
| --for-zero-height   | bool   | false    | false    | 导出数据之前做一些清理性的工作，如果不想以导出的数据启动一条新链，可以不加这个标志 |
| --height            | int    | false    | -1       | 高度 |
| --home              | string | false    | $HOME/.iris | 指定存储配置和区块链数据的目录 |

## 示例

1. 导出最新的区块链状态:
```
iris export
```
2. 导出高度10000的区块链状态
```
iris export --height 10000
```
3. 如果想导出105000高度的区块链状态，并且以这个状态启动一条新链，可以尝试这个命令
```
iris export --height=20010 --for-zero-height --home=[your_home]
```
可能会遇到的错误
```
ERROR: error exporting state: failed to load rootMultiStore: wanted to load target 105000 but only found up to 100000
```
默认情况下，节点启动时设置的历史数据删除策略是`syncable`。这就意味着此节点只会保存最近100块的区块链状态，对于更老的状态则每10000块保留一个，比如高度为10000,20000和30000的状态会被永久保留。这种情况下用户需要在节点上发起replay的操作：
```
iris start --home=[your_home] --replay_height=105000
```
执行上面这个命令，节点会从高度100000开始replay到105000，这时105000的状态会被重建出来，然后再次尝试导出高度为105000的区块链状态：
```
iris export --height=20010 --for-zero-height --home=[your_home]
```