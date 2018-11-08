# How to deploy IRIShub Monitor

Please make sure that iris is installed in your computer and added to $PATH.You can see this page for insturcion https://github.com/irisnet/irishub. You also need /bin/bash, wc, ps to ensure the monitor work properly.

### Get Monitor Info

* Address of validator

```
-a= < hex-encoded-address-of-validator >
```

* Chain-ID
```
--chain-id= < id-of-blockchain >
```

* Node


```
--node= < localhost:26657 >
```
###  Start IRIS Monitor
Example:
```
irismon --account-address=faa1nwpzlrs35nawthal6vz2rjr4k8xjvn7k8l63st --address=EAC535EC37EB3AE8D18C623BA4B4C8128BC082D2 --chain-id=irishub-stage --node=http://localhost:26657&
```

then, you can visit http://localhost:36660/ to see the metrics page. 

### Start Prometheus

First, you need to edit the configuration file `prometheus.yml` of prometheus in `~/volumes/prometheus` folder. Add `jobs` :
```$xslt
- job_name: fuxi-4000

    static_configs:

    - targets: ['localhost:36660']

      labels:

        instance: fuxi-4000
```
Start prometheus service with Docker:
```
docker run -p 9090:9090 -v ~/volumes/prometheus:/etc/prometheus prom/prometheus 1>prometheus.log &
```

Then, you could see in your browser that there are some data available at port 36660.

### Start Grafana

You could start grafana with docker by running:
```$xslt
sudo docker run -p 3000:3000 grafana/grafana 1>grafana.log 2>grafana.error &
```

The default username and password are both admin. We strongly recommend immediately changing your username & password after login.

Then you could create your own  dashboard. 

###  Stop IRIS Monitor
```
kill -9 <irismon-process-id>
```