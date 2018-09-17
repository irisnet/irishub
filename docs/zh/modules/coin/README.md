# Coin_Type

##  定义

coin_type 定义了iris-hub系统中代币的可用单位，只要是系统中已注册的coin_type类型，都可以使用该单位来进行交易。iris-hub中系统默认的代币为iris，iris存在以下几种可用单位:iris-milli,iris-micro,iris-nano,iris-pico,iris-femto,iris-atto。他们之间存在以下换算关系

```
1 iris = 10^3 iris-milli
1 iris = 10^6 iris-micro
1 iris = 10^9 iris-nano
1 iris = 10^12 iris-pico
1 iris = 10^15 iris-femto
1 iris = 10^18 iris-atto
```

## coin_type的数据模型

```golang
type CoinType struct {
	Name    string `json:"name"`
	MinUnit Unit   `json:"min_unit"`
	Units   Units  `json:"units"`
	Origin  Origin `json:"origin"`
	Desc    string `json:"desc"`
}
```

* Name :    代币名称，也是coin的主单位，例如iris
* MinUnit： coin_type的最小单位，系统中存在的代币都是以最小单位的形式存在，例如iris代币，在iris-hub中存储的单位是iris-atto。当用户发送交易到iris-hub中，使用的必须是该代币的最小单位。但是如果你使用的是iris-hub提供的命令行工具，你可以使用任何系统识别的单位，系统将自动转化为该代币对应的最小单位形式。比如如果你使用send命令转移1iris，命令行将在后端处理为10^18 iris-atto，使用交易hash查询到的交易详情，你也只会看到10^18 iris-atto。

## Unit结构定义

```golang
type Unit struct {
	Denom   string `json:"denom"`
	Decimal int    `json:"decimal"`
}
```

其中Denom定义为该单位的名称，Decimal定义为该单位支持的最大精度，例如iris-atto支持的最大精度为18
* Units：定义了coin_type下可用的一组单位
* Origin：定义了该coin_type的来源，取值:Native(系统内部，iris),External(系统外部,例如eth等),UserIssued(用户自定义)
* Desc：对该代币coin_type的描述

## 查询代币coin_type

如果想查询某种代币的coin_type配置，可以使用如下命令

```golang
iriscli coin types [coin_name]
```