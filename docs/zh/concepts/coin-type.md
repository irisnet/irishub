# 通证单位

## 定义

Coin type 定义了IRIShub系统中代币的可用单位，只要是系统中已注册的coin-type类型，都可以使用该单位来进行交易。IRIShub中系统默认的代币为iris，iris存在以下几种可用单位： `iris-milli`，`iris-micro`，`iris-nano`，`iris-pico`，`iris-femto` 和 `iris-atto`。他们之间存在以下换算关系：

```toml
1 iris = 10^3 iris-milli
1 iris = 10^6 iris-micro
1 iris = 10^9 iris-nano
1 iris = 10^12 iris-pico
1 iris = 10^15 iris-femto
1 iris = 10^18 iris-atto
```

## 数据结构

### CoinType

```go
type CoinType struct {
    Name    string `json:"name"`
    MinUnit Unit   `json:"min_unit"`
    Units   Units  `json:"units"`
    Origin  Origin `json:"origin"`
    Desc    string `json:"desc"`
}
```

### Unit

```go
type Unit struct {
    Denom   string `json:"denom"`
    Decimal int    `json:"decimal"`
}
```

* `Name`:  代币名称，也是coin的主单位，例如`iris`
* `MinUnit`: coin-type 的最小单位，系统中存在的代币都是以最小单位的形式存在，例如 `iris` 代币，在IRIShub中存储的单位是 `iris-atto`。当用户发送交易到IRIShub中，使用的必须是该代币的最小单位。但是如果你使用的是IRIShub提供的命令行工具，你可以使用任何系统识别的单位，系统将自动转化为该代币对应的最小单位形式。比如如果你使用send命令转移1iris，命令行将在后端处理为10^18 iris-atto，使用交易hash查询到的交易详情，你也只会看到10^18 iris-atto。

* `Denom`:定义为该单位的名称，Decimal定义为该单位支持的最大精度，例如iris-atto支持的最大精度为18
* `Units`：定义了coin-type下可用的一组单位
* `Origin`：定义了该coin-type的来源，取值:Native(系统内部，iris),External(系统外部,例如eth等),UserIssued(用户自定义)
* `Desc`：对该代币coin-type的描述

## 查询 CoinType

如果想查询某种代币的CoinType配置，可以使用如下命令：

```bash
iriscli bank coin-type <coin-name>
```

比如查询 `iris` 的 CoinType，可以执行以下命令：

```bash
iriscli bank coin-type iris
```

示例输出:

```bash
CoinType:
  Name:     iris
  MinUnit:  iris-atto: 18
  Units:    iris: 0,  iris-milli: 3,  iris-micro: 6,  iris-nano: 9,  iris-pico: 12,  iris-femto: 15,  iris-atto: 18
  Origin:   native
  Desc:     IRIS Network
```
