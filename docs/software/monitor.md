# How to deploy IRIS Monitor

## Install IRIS Monitor 
Please refer to this [document](https://github.com/irisnet/irishub/blob/master/docs/get-started/Install-Iris.md) for deatiled instructions. Besides, please make sure your machine has these commands(bash, wc, ps) installed.

## Start IRIS Monitor

```
iristool monitor --validator-address=EAC535EC37EB3AE8D18C623BA4B4C8128BC082D2 \
--account-address=faa1nwpzlrs35nawthal6vz2rjr4k8xjvn7k8l63st \
--chain-id=<chain-id> --node=http://localhost:26657
```

Parameters：

- `validator-address`：hex encoded validator address
- `account-address`：bech32 encoded account address
- `chain-id`：blockchain id that you want to monitor
- `node`：listening address of the node that you want to monitor (The rpc url of a irishub node, default value is tcp://localhost:26657. If you want to monitor other irishub nodes, please change this value.)

Then you can visit `http://localhost:36660/` to see metrics data。

## Start Prometheus

### Edit Prometheus config file

You can visit [prometheus.yml](https://github.com/prometheus/prometheus/blob/master/documentation/examples/prometheus.yml) to download default `prometheus.yml`.

Then edit `prometheus.yml` and add `jobs` :

```yaml
      - job_name: fuxi-5000
          static_configs:
          - targets: ['localhost:36660']
            labels:
              instance: fuxi-5000
```

> Note：value of targets is ip:port which used by IRIS monitor 

### Start Prometheus

```bash
docker run -d --name=prometheus -p 9090:9090 -v ~/volumes/prometheus:/etc/prometheus prom/prometheus
```

> The above example, you can save `prometheus.yml` at `~/volumes/prometheus` on your host machine

You can visit `http://localhost:9090` to see prometheus data.

## Start Grafana

```
docker run -d --name=grafana -p 3000:3000 grafana/grafana
```

You can visit `http://localhost:3000/` to open grafana and create your own dashboard.

> Tips: The default username and password are both admin. We strongly recommend immediately changing your username & password after login