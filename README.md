# gortcp [![Build Status](https://travis-ci.org/lovedboy/gortcp.svg?branch=master)](https://travis-ci.org/lovedboy/gortcp)

支持

* 内网穿透，灵感来源于[rtcp](https://github.com/knownsec/rtcp)。
* 文件上传
* 远程命令执行

# [Download](https://github.com/lovedboy/gortcp/releases/tag/v0.1)



# Usage

* cmd/server/server.go  转发服务器
* cmd/client/client.go  内网机器
* cmd/control/control.go 控制终端



## 转发服务器运行：

`./server -addr :33456 -auth 123456`

## 内网机器运行：

`./client -addr 10.68.102.49:33456`

## 控制终端

### 查看连接的机器

```
⇒  ./control -addr 10.68.102.49:33456 -action list
|ID        |ADDRESS
|1         |10.68.102.48:49725
```

### 远程执行命令
在远程1号机器上执行命令    
```
⇒  ./control -addr 10.68.102.49:33456 -action exec -id 1 -cmd "ping -c 4 baidu.com"
PING baidu.com (111.13.101.208): 56 data bytes
64 bytes from 111.13.101.208: icmp_seq=0 ttl=47 time=52.323 ms
64 bytes from 111.13.101.208: icmp_seq=1 ttl=47 time=48.281 ms
64 bytes from 111.13.101.208: icmp_seq=2 ttl=47 time=53.093 ms
64 bytes from 111.13.101.208: icmp_seq=3 ttl=47 time=50.535 ms

--- baidu.com ping statistics ---
4 packets transmitted, 4 packets received, 0.0% packet loss
round-trip min/avg/max/stddev = 48.281/51.058/53.093/1.852 ms
```

### 文件上传
上传本地的./control到远程1号机器的/tmp/test.bin   
```
⇒  ./control -addr 10.68.102.49:33456 -action upload -id 1 -src ./control -dst /tmp/test.bin
send: 3300KB | time: 0.06S | speed: 56970KB/S
CLIENT: receive complete, md5 verify passed
```

### 端口转发
转发本地33060端口的数据到远程机器的3306端口，当然也可以是远程机器能连接的其他主机    
```
⇒  ./control -addr 10.68.102.49:33456 -action forward -id 1 -laddr 127.0.0.1:33060 -raddr 127.0.0.1:3306
16-08-30 17:31:22.063 INFO @control.go:172 [tcp] listen on local 127.0.0.1:33060
```
连接本地的33060端口：   
```
⇒  mysql -uroot -h 127.0.0.1 -P 33060
Welcome to the MySQL monitor.  Commands end with ; or \g.
Your MySQL connection id is 882
Server version: 5.6.17 MySQL Community Server (GPL)

Copyright (c) 2000, 2014, Oracle and/or its affiliates. All rights reserved.

Oracle is a registered trademark of Oracle Corporation and/or its
affiliates. Other names may be trademarks of their respective
owners.

Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

mysql>
```

