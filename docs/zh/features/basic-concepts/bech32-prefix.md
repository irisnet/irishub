# Bech32 on IRISnet

Bech32是由Pieter Wuille和Greg Maxwel提出的新比特币地址格式。除了比特币之外,bech32可以编码任何短二进制数据。在IRISnet里，私钥和地址可能指的是一些在网络中不同的角色，例如普通账户和验证人账户等。IRISnet设计使用Bech32地址格式来提供对数据鲁棒的完整性检查。  用户可读部分(human readable part) 使得用户更有效率的理解地址和阅读错误信息。


## 用户可读部分表

| HRP        | Definition |
| -----------|:-------------|
|faa|   IRISnet Account Address|
|fap|	IRISnet Account Public Key|
|fva|   IRISnet Validator's Operator Address|
|fvp|   IRISnet Validator's Operator Public Key|
|fca|   IRISnet Consensus Address|
|fcp|	IRISnet Consensus Public Key|

## 编码

不是所有IRISnet的用户地址都会以bech32的格式暴露出来。许多地址仍然是hex编码或者base64编码。 在bech32编码之前，首先要使用amino对其他地址私钥二进制表示进行编码。

## 例子

一旦创建一个新的账户，你将会看到以下信息:

```
NAME:	TYPE:	ADDRESS:						            PUBKEY:
test1	local	faa18ekc4dswwrh2a6lfyev4tr25h5y76jkpqsz7kl	fap1addwnpepqgxa40ww28uy9q46gg48g6ulqdzwupyjcwfumgfjpvz7krmg5mrnw6zv8uv
```

这意味着你创建了一个新账户地址 `faa18ekc4dswwrh2a6lfyev4tr25h5y76jkpqsz7kl`， 他的用户可读部分是 `faa`。他的公钥被密码成  `fap1addwnpepqgxa40ww28uy9q46gg48g6ulqdzwupyjcwfumgfjpvz7krmg5mrnw6zv8uv`， 他的用户可读部分是 `fap`。 