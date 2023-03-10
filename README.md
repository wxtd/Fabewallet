# Fabewallet

支持版本Hyperledger Fabric v2.3.3

**并发冲突改进实验**

进入执行目录

```shell
cd fabewallet	
```

启动网络

```shell
./startFabric.sh
```

设置实验数据（存储在redis中）

```shell
cd gxtest/redis_test
go run inputdata.go
```

开启两个终端执行并查集程序并查看日志

```shell
cd ..
go run *.go
tail -f rollfunc.log
```

主窗口执行脚本查看结果

```shell
cd ../go/
# 执行原版实验
./run_fabric.sh
# 或 执行改进版本实验
./run_improve.sh
```