# 重置区块链状态

## 介绍

IRISnet支持重置区块链状态到任意高度，这里介绍重置区块链状态所用的命令为`iris reset`。

## 用法
```		
 iris reset <flags>
```
## 标志

 | Name，shorthand     | type   | Required | Default  | Description    |		
 | ------------------- | -----  | -------- | -------- | -------------- |		
 | --height            | int    | false    | 0        | Specify the height, default value is 0 which means to export the latest state |		
 | --home              | string | false    | $HOME/.iris       | Specify the directory which stores node config and blockchain data |		
 
## 示例
 
1. 重置区块链状态到区块100:
```		
 iris reset --height 100 --home=<path_to_your_home>
```