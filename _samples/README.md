# 说明

```shell
kill -SIGTSTP pid    关闭
kill -SIGUSR1 pid    重启
```

## tcp 测试

* 启动server和client服务
* （如修改代码，重新编译出server）
* kill -SIGUSR1 server服务pid进行重启，此时有新、旧两个进程在跑，旧的server停止对外服务。
* 启动新的client服务，连接上新的server服务
* 关闭旧的client服务，旧的server自动关闭
