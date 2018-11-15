# 如何部署 IRIS Monitor

确保已经安装了iris等工具，系统中需要有/bin/bash、wc、ps等命令。 你可以参考这个页面来安装iris工具: https://github.com/irisnet/irishub

## 启动 IRIS Monitor

```
irismon --address=EAC535EC37EB3AE8D18C623BA4B4C8128BC082D2 \
--account-address=faa1nwpzlrs35nawthal6vz2rjr4k8xjvn7k8l63st \
--chain-id=irishub-stage --node=http://localhost:36657
```

参数说明：

- `address`：要监测的验证人地址（hex编码）
- `account-address`：要监测的账户地址（bech32 编码）
- `chain-id`：要监测的链 id
- `node`：要监控的节点地址（默认为 tcp://localhost:26657）

启动之后, 你可以通过 `http://localhost:36660/` 能看到 Metrics 数据页面。

## 启动 Prometheus

### 编辑配置文件

从 [https://github.com/prometheus/prometheus/blob/master/documentation/examples/prometheus.yml](https://github.com/prometheus/prometheus/blob/master/documentation/examples/prometheus.yml) 下载默认配置文件到本地：

在配置文件 `prometheus.yml` 中添加以下 `jobs` :

```yaml
      - job_name: fuxi-4000
          static_configs:
          - targets: ['localhost:36660']
            labels:
              instance: fuxi-4000
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