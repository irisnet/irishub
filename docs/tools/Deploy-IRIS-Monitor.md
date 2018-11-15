# How to deploy IRIS Monitor

Please make sure that iris is installed in your computer and added to $PATH.You can see this page for insturcion https://github.com/irisnet/irishub. You also need /bin/bash, wc, ps to ensure the monitor work properly.

## Start IRIS Monitor

```
irismon --address=EAC535EC37EB3AE8D18C623BA4B4C8128BC082D2 --account-address=faa1nwpzlrs35nawthal6vz2rjr4k8xjvn7k8l63st --chain-id=irishub-stage --node=http://localhost:36657
```

Parameters：

- `address`：hex encoded validator address
- `account-address`：bech32 encoded account address
- `chain-id`：blockchain id that you want to monitor
- `node`：listening address of the node that you want to monitor ("tcp://localhost:26657" by default, you should not change this if you didn't modify your rpc port)

Then you can visit `http://localhost:36660/` to see metrics data。

## Start Prometheus

### Edit Prometheus config file

You can visit [prometheus.yml](https://github.com/prometheus/prometheus/blob/master/documentation/examples/prometheus.yml) to download default `prometheus.yml` and save it.

Then edit `prometheus.yml` and add `jobs` :

```yaml
      - job_name: fuxi-4000
          static_configs:
          - targets: ['localhost:36660']
            labels:
              instance: fuxi-4000
```

> Note：value of targets is ip:port which used by IRIS monitor 

### Start Prometheus

```bash
docker run -d --name=prometheus -p 9090:9090 -v ~/volumes/prometheus:/etc/prometheus prom/prometheus
```

> The above example, the path of `prometheus.yml` is `~/volumes/prometheus` on host machine

You can visit `http://localhost:9090` to see prometheus data.

## Start Grafana

```
docker run -d --name=grafana -p 3000:3000 grafana/grafana
```

You can visit `http://localhost:3000/` to open grafana and create your own dashboard.

> Tips: The default username and password are both admin. We strongly recommend immediately changing your username & password after login.