# 重置区块链状态

## 介绍

IRISnet支持重置区块链状态到任意高度，这里介绍重置区块链状态所用的命令为`iris reset`。

## 用法
```		
 iris reset <flags>
```
## 标志

 | Name，shorthand | type   | Required | Default     | Description                           |		
 | --------------- | -----  | -------- | ----------- | ------------------------------------- |		
 | --height        | uint   |          | 0           | 重置状态为特定高度(大于最新高度表示最新高度) |		
 | --home          | string |          | $HOME/.iris | 指定存储配置和区块链数据的目录             |		
 
## 示例
 
1. 重置区块链状态到区块100:
```
 iris reset --height 100 --home=<path_to_your_home>
```