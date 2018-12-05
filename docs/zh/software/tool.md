# IRIS Command tool

## Introduction
`iristool` 现在包含 [monitor](./monitor.md) 和 debug.

## debug
用于简单调试的简单工具。

我们尝试同时接受十六进制和base64格式，并提供有用的响应。

注意，我们通常在日志中将字节编码为十六进制，而在JSON中编码为base64。

### Usage

* Pubkeys 
下面得到相同的结果:

```bash
iristool debug pubkey TZTQnfqOsi89SeoXVnIw+tnFJnr4X8qVC0U8AsEmFk4=
iristool debug pubkey 4D94D09DFA8EB22F3D49EA17567230FAD9C5267AF85FCA950B453C02C126164E
```

* Txs
传入 hex/base64 tx 并返回完整的 JSON

```bash
iristool debug tx [hex or base64 transaction]
```

* Hack
这是一个带有样板的命令，用于将Go作为脚本语言来攻击现有的Iris状态。

如果你运行 
```
iristool debug hack $HOME/.iris
```
在那个状态, 它将对状态历史进行二分查找来发现何时违背了状态不变量。
