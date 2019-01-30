# How to deploy IRIS Monitor

## Metrics

IRISnet can report and serve the Prometheus metrics, which in their turn can be consumed by Prometheus collector(s).

This functionality is disabled by default.

To enable the Prometheus metrics, set `instrumentation.prometheus=true` in your config file. Metrics will be served under /metrics on 26660 port by default. Listen address can be changed in the config file (see instrumentation.prometheus\_listen\_addr).

### List of available metrics

The following metrics are available:

| Name | Type | Tags | Description |
| ---- | ---- | ---- | ----------- |
| stake_bonded_token | Gauge | validator_address | Total bonded token by validator |
| stake_loosen_token | Gauge |                   | Total loosen token |
| stake_burned_token | Gauge |                   | Total burned token |
| stake_slashed_token | Counter | validator_address | Total slashed token by validator |
| stake_jailed        | Gauge | validator_address | Jailed status by validator, either 0 (not jailed) or 1 (jailed) |
| stake_power         | Gauge | validator_address | Voting power by validator |
| distribution_community_tax  | Gauge |  | Total token of community funds pool |
| upgrade_upgrade  | Gauge |  | Whether new software needs to be installed, either 0 (no) or 1 (yes) |
| upgrade_signal  | Gauge | validator_address, version | Whether validator have run the new version software, either 0 (no) or 1 (yes)|
| service_active_requests  | Gauge |  | Number of active requests |

IRISnet metrics also contains tendermint metrics, Visit [tendermint metrics](https://github.com/irisnet/tendermint/blob/irisnet/v0.27.3-iris/docs/tendermint-core/metrics.md) for more information.

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