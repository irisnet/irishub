### Running in command line
这个命令会在宿主机后台启动一个iris进程，启动前需要确认宿主机已安装iris、iriscli，并确认`26656`和`26657`端口没有被占用，
命令会删除$HOME/.iris*目录，请注意做好备份
```
nohup sh start.sh > iris.log 2>&1 &
```
查看日志：
```$xslt
tail -f iris.log
```
查看账户列表：
```$xslt
iriscli keys list
```

### Running in docker
这个命令会启动一个基于最新develop分支构建的镜像的docker容器，启动前需要确认`26656`和`26657`端口没有被占用
```
sh start-docker.sh
```
查看日志：
```$xslt
docker logs -f irishub-sandbox
```
查看账户列表：
```$xslt
docker exec -it irishub-sandbox iriscli keys list
```