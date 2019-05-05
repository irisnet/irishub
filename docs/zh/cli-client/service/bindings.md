# iriscli service bindings

## 描述

查询服务绑定列表

## 用法

```
iriscli service bindings <flags>
```

## 标志

| Name, shorthand | Default                    | Description                                            | Required |
| --------------- | -------------------------- | ------------------------------------------------------ | -------- |
| --def-chain-id  |                            | 定义该服务的区块链ID                             | 是        |
| --service-name  |                            | 服务名称                                       | 是        |
| --help, -h      |                            | 查询绑定列表命令帮助                                       |          |

## 示例

### 查询服务绑定列表

```shell
iriscli service bindings --def-chain-id=<chain-id> --service-name=<service-name>
```

运行成功以后，返回的结果如下:

```json
[{
	"def_name": "test-service",
	"def_chain_id": "test",
	"bind_chain_id": "test",
	"provider": "iaa1ydhmma8l4m9dygsh7l08fgrwka6yczs0se0tvs",
	"binding_type": "Local",
	"deposit": [{
		"denom": "iris-atto",
		"amount": "1000000000000000000000"
	}],
	"price": [{
		"denom": "iris-atto",
		"amount": "1000000000000000000"
	}],
	"level": {
		"avg_rsp_time": "10000",
		"usable_time": "100"
	},
	"available": true,
	"disable_height": "0"
}]
```