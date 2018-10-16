# 参与到Genesis文件生成流程中


1. 每个希望成为验证人的参与者确保安装了对应版本的软件：iris v0.4.2

2. 执行gentx命令，获得一个node-id.json的文件。这个操作将默认生成一个余额为200IRIS的账户，该账户默认绑定100IRIS成为一个验证人候选人。

```
       iris init gen-tx --name=your_name --home=path_to_home --ip=Your_public_IP
```
   代码示例：
```
       iris init gen-tx --name=alice 
```

```
       {
        "app_message": {
          "secret": "similar spread grace kite security age pig easy always prize salon clip exhibit electric art abandon"
        },
        "gen_tx_file": {
          "node_id": "3385a8e3895b169eab3024079d987602b4d2b383",
          "ip": "192.168.1.7",
          "validator": {
            "pub_key": {
              "type": "AC26791624DE60",
              "value": "RDxXckkpTc35q9xlLNXjzUAov6xMkGJlwtWg2IqAkD8="
            },
            "power": 100,
            "name": ""
          },
          "app_gen_tx": {
            "name": "alice",
            "address": "8D3B5761BC2B9048E2A7745B14E62D51C82E0B7C",
            "pub_key": {
              "type": "AC26791624DE60",
              "value": "RDxXckkpTc35q9xlLNXjzUAov6xMkGJlwtWg2IqAkD8="
            }
          }
        }
       }
  ```
然后你可以发现在$IRISHOME/config目录下生成了一个gentx文件夹。里面存在一个gentx-node-ID.json文件。这个文件包含了如下信息：

   ```
       {
        "node_id": "3385a8e3895b169eab3024079d987602b4d2b383",
        "ip": "192.168.1.7",
        "validator": {
          "pub_key": {
            "type": "AC26791624DE60",
            "value": "RDxXckkpTc35q9xlLNXjzUAov6xMkGJlwtWg2IqAkD8="
          },
          "power": 100,
          "name": ""
        },
        "app_gen_tx": {
          "name": "alice",
          "address": "8D3B5761BC2B9048E2A7745B14E62D51C82E0B7C",
          "pub_key": {
            "type": "AC26791624DE60",
            "value": "RDxXckkpTc35q9xlLNXjzUAov6xMkGJlwtWg2IqAkD8="
          }
        }
       }
  ```
   validator字段对应了home/config下的节点信息

   `app_gen_tx`中说明了拥有这个节点的账户信息。这个账户的助记词就是刚刚的secret

3. 将上述提到的json文件以提交Pull Request的形式上传到`https://github.com/irisnet/testnets/tree/master/testnets/fuxi-3001/config/gentx`目录下：

   注意⚠️：json文中的IP改成公网IP




