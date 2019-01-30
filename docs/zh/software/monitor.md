# 如何部署 IRIS Monitor

## Metrics

IRISnet可以报告和提供Prometheus metrics，Prometheus收集器可以使用这些指标。

默认情况下禁用此功能。

要启用Prometheus metrics，请在配置文件中设置`instrumentation.prometheus=true`。默认情况下，Metrics将在26660端口下的/metrics提供。可以在配置文件中修改服务地址(请参阅 instrumentation.prometheus\_listen\_addr).

### List of available metrics

The following metrics are available:

| Name | Type | Tags | Description |
| ---- | ---- | ---- | ----------- |
| stake_bonded_token | Gauge | validator_address | 验证人被绑定的token总数 |
| stake_loosen_token | Gauge |                   | 未被绑定的token总数 |
| stake_burned_token | Gauge |                   | 销毁的token总数 |
| stake_slashed_token | Counter | validator_address | 验证人被惩罚的token总数 |
| stake_jailed        | Gauge | validator_address | 验证人监禁状态，0（未监禁）或1（被监禁） |
| stake_power         | Gauge | validator_address | 验证人投票权 |
| distribution_community_tax  | Gauge |  | 社区资金池的token总数 |
| upgrade_upgrade  | Gauge |  | 是否需要安装新软件，0（否）或1（是） |
| upgrade_signal  | Gauge | validator_address, version | 验证人是否运行了新版本软件，0（否）或1（是）|
| service_active_requests  | Gauge |  | 活跃的请求数 |

IRISnet metrics同时也包括了tendermint metrics，访问[tendermint metrics](https://github.com/irisnet/tendermint/blob/irisnet/v0.27.3-iris/docs/tendermint-core/metrics.md) 获取更多信息。

## 启动 IRIS Monitor

```
iristool monitor --validator-address=EAC535EC37EB3AE8D18C623BA4B4C8128BC082D2 \
--account-address=faa1nwpzlrs35nawthal6vz2rjr4k8xjvn7k8l63st \
--chain-id=<chain-id> --node=http://localhost:36657
```

参数说明：

- `validator-address`：要监测的验证人地址（hex编码）
- `account-address`：要监测的账户地址（bech32 编码）
- `chain-id`：要监测的链 id
- `node`：要监控的节点地址（默认为 tcp://localhost:26657）

启动之后, 你可以通过 `http://localhost:36660/` 看到 Metrics 数据页面。

## 启动 Prometheus

### 编辑配置文件

从 [https://github.com/prometheus/prometheus/blob/master/documentation/examples/prometheus.yml](https://github.com/prometheus/prometheus/blob/master/documentation/examples/prometheus.yml) 下载默认配置文件到本地：

在配置文件 `prometheus.yml` 中添加以下 `jobs` :

```yaml
      - job_name: fuxi-5000
          static_configs:
          - targets: ['localhost:36660']
            labels:
              instance: fuxi-5000
```

> Note：targets 配置项的值为 IRIS Monitor 启动后所占用的 ip 和 port。 

### 启动 Prometheus

```bashg
docker run -d --name=prometheus -p 9090:9090 -v ~/volumes/prometheus:/etc/prometheus prom/prometheus
```

将编辑好的配置文件 `prometheus.yml` 放在宿主机的目录下并映射到容器中。
例如在上例中配置文件位于宿主机的 `~/volumes/prometheus` 中。


## 启动 Grafana

```
docker run -d --name=grafana -p 3000:3000 grafana/grafana
```

接下来就可以访问 `http://localhost:3000/` 来查看 grafana 监控。

> Tips: 打开网页后使用默认用户名 admin，默认密码 admin 登录。建议登录之后立即修改密码。
