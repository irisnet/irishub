# 如何部署monitor

确保已经安装了iris等工具，系统中需要有/bin/bash、wc、ps等命令。 你可以参考这个页面来安装iris工具:https://github.com/irisnet/irishub

## 获得修改运行参数

* 验证人地址的hex编码

```
-a=<你的验证人地址的hex编码>
```

* 确定Chain-ID

```
--chain-id=<你要监控的网络的ID>
```

* 确定监控的节点


```
--node=<你要监控的节点监听的rpc端口(默认为26657)>
```

## 启动IRIS Monitor

```
irismon --account-address=faa1nwpzlrs35nawthal6vz2rjr4k8xjvn7k8l63st --address=EAC535EC37EB3AE8D18C623BA4B4C8128BC082D2 --chain-id=irishub-stage --node=http://localhost:36657&
```


### 启动Prometheus

从以下链接下载默认配置文件prometheus.yml：
* https://github.com/prometheus/prometheus/blob/master/documentation/examples/prometheus.yml
复制到`~/volumes/prometheus`目录下。

在配置文件`prometheus.yml`中添加 `jobs` :

```$xslt
      - job_name: fuxi-4000
      
          static_configs:
      
          - targets: ['localhost:36660']
      
            labels:
      
              instance: fuxi-4000
```

通过Docker启动Prometheus后可以在本地36660端口看到监控数据。
```
docker run -p 9090:9090 -v ~/volumes/prometheus:/etc/prometheus prom/prometheus 1>prometheus.log &
```
### 启动Grafana

通过Docker启动Grafana
```$xslt
sudo docker run -p 3000:3000 grafana/grafana 1>grafana.log 2>grafana.error &
```

接下来就可以访问localhost:3000来查看grafana监控。打开网页后使用默认用户名admin，默认密码admin登录。建议登录之后立即修改密码。

5. 关闭监控
```$xslt
killl -9 <process-ID>
```