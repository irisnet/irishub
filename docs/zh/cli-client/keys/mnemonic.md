# iriscli keys mnemonic

## 描述

通过读取系统熵来创建24个单词组成的bip39助记词（也称为种子短语）。如果需要传递自定义的熵，请使用`unsafe-entropy`参数。

## 使用方式

```
iriscli keys mnemonic <flags>
```

## 标志

| 名称, 速记        | 默认值     | 描述                                                                          | 是否必须  |
| ---------------- | --------- | ----------------------------------------------------------------------------- | -------- |
| --help, -h       |           | 查询命令帮助                                                                   |          |
| --unsafe-entropy |           | 提示用户提供自己的熵，而不是依赖于系统生成                                          |          |

## 例子

### 创建助记词

```shell
iriscli keys mnemonic
```

执行命令就可以得到24个单词组成的助记词。为了安全考虑，请注意保存，比如将单词手抄纸并将纸张妥善保存。

```txt
police possible oval milk network indicate usual blossom spring wasp taste canal announce purpose rib mind river pet brown web response sting remain airport
```

### 使用`unsafe-entropy`模式

此模式创建固定的助记词

```shell
root@ubuntu16:~# iriscli keys mnemonic --unsafe-entropy

WARNING: Generate at least 256-bits of entropy and enter the results here:
<input_your_own_entropy_string>
Input length: 128 [y/n]:y
-------------------------------------
wine hire tongue weasel air puzzle claim pole curtain taste box learn exchange where become inside blur tragic suffer fruit hole transfer race unit
```

