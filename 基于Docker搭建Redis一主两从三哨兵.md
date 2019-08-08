# 基于Docker搭建Redis一主两从三哨兵

这段时间正在学习Redis和容器相关的内容，因此想通过docker搭建一套redis主从系统来加深理解。看这篇文章可能你需要一定的docker基础，以及对redis主从和哨兵机制有所了解。

这次实验准备了三台云主机，系统为Debian，ip分别为：35.236.172.131 ，35.201.200.251，34.80.172.42。

首先分别在这三台主机上安装docker，然后每台主机上启动一个redis容器，运行redis-server服务，其中35.236.172.131作为master，另外两台机器作为slave，最后在三台主机上再分别启动一个redis容器，运行redis-sentinel。为什么还是redis容器呢？因为sentinel实际上还是一个redis-server，只不过是以sentinel模式执行，只能处理sentinel需要的一些命令。

## 安装docker

docker的安装有很多种方法，这里就不介绍了。本次使用脚本安装docker，Debian系统脚本安装如下，其他系统可以参考Docker官网的安装方法：[https://docs.docker.com/install/linux/docker-ce/debian/](https://docs.docker.com/install/linux/docker-ce/debian/)。不过下面的命令在官网命令的基础上修改镜像源为阿里云，因为国内镜像往往会快一些。

### 脚本安装docker
在物理主机或者云虚拟主机上运行下面的命令就可以完成docker安装了，当然我是在Debian系统上，其他系统相应参考官网上的方法。
```
$ curl -fsSL https://get.docker.com -o get-docker.sh
$ sudo sh get-docker.sh --mirror Aliyun
```

### 启动docker CE
docker是以客户端和服务器模型运行的，因此需要先运行docker的服务器，服务器以daemon的形式运行。docker CE是docker的社区版本。
```
$ sudo systemctl enable docker
$ sudo systemctl start docker
```

### 验证docker是否安装成功

下面的这条命令是从docker的官方仓库拉取一个名为hello-world的镜像，并通过这个镜像启动一个容器。
```
$ docker run hello-world
```
如果运行结果如下，出现了`Hello from Docker!`，说明docker安装成功了
```
$ docker run hello-world
Unable to find image 'hello-world:latest' locally
latest: Pulling from library/hello-world
1b930d010525: Pull complete 
Digest: sha256:6540fc08ee6e6b7b63468dc3317e3303aae178cb8a45ed3123180328bcc1d20f
Status: Downloaded newer image for hello-world:latest

Hello from Docker!
This message shows that your installation appears to be working correctly.

To generate this message, Docker took the following steps:
 1. The Docker client contacted the Docker daemon.
 2. The Docker daemon pulled the "hello-world" image from the Docker Hub.
    (amd64)
 3. The Docker daemon created a new container from that image which runs the
    executable that produces the output you are currently reading.
 4. The Docker daemon streamed that output to the Docker client, which sent it
    to your terminal.

To try something more ambitious, you can run an Ubuntu container with:
 $ docker run -it ubuntu bash

Share images, automate workflows, and more with a free Docker ID:
 https://hub.docker.com/

For more examples and ideas, visit:
 https://docs.docker.com/get-started/
```

## 启动容器搭建主从

docker安装成功后，可以开始部署redis服务了。先从docker官方公共仓库拉取redis镜像，然后修改redis服务的配置文件，最后启动容器，启动redis服务器。在多台机器上运行redis服务器，并建立主从关系。

redis的主从是实现redis集群和redis哨兵高可用的基础，redis的主从结构使从可以复制主上的数据，如果从与主之间网络断开，从会自动重连到主上。

![](https://user-gold-cdn.xitu.io/2019/7/12/16be65279459a7e9?w=403&h=71&f=png&s=3519)

### 获取Redis镜像
下面的命令会拉取最新的官方版本的redis镜像
```
$ docker pull redis
```
查看镜像
```
$ docker image ls
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
redis               latest              bb0ab8a99fe6        7 days ago          95MB
hello-world         latest              fce289e99eb9        6 months ago        1.84kB
```

### 获取并修改redis配置文件

redis官方提供了一个配置文件样例，通过wget工具下载下来。我用的root用户，就直接下载到/root目录里了。

```
$ wget http://download.redis.io/redis-stable/redis.conf
```

打开下载下来的文件后，可以看到配置有很多。我只是搭建服务进行试验所以只修改必要的几项。如果要运用到线上，那必须所有的配置都按需求进行修改。

其中redis服务器的master和slave角色使用的配置文件还会有些不同，下面分别进行说明。

对于master而言，配置文件修改以下几项
```
# 注释这一行，表示Redis可以接受任意ip的连接
# bind 127.0.0.1 

# 关闭保护模式
protected-mode no 

# 让redis服务后台运行
daemonize yes 

# 设定密码(可选，如果这里开启了密码要求，slave的配置里就要加这个密码. 只是练习配置，就不使用密码认证了)
# requirepass masterpassword 

# 配置日志路径，为了便于排查问题，指定redis的日志文件目录
logfile "/var/log/redis/redis.log"
```

对于slave而言，配置文件修改以下几项：
```
# 注释这一行，表示Redis可以接受任意ip的连接
# bind 127.0.0.1 

# 关闭保护模式
protected-mode no 

# 让redis服务后台运行
daemonize yes 

# 设定密码(可选，如果这里开启了密码要求，slave的配置里就要加这个密码)
requirepass masterpassword 

# 设定主库的密码，用于认证，如果主库开启了requirepass选项这里就必须填相应的密码
masterauth <master-password>

# 设定master的IP和端口号，redis配置文件中的默认端口号是6379
# 低版本的redis这里会是slaveof，意思是一样的，因为slave是比较敏感的词汇，所以在redis后面的版本中不在使用slave的概念，取而代之的是replica
# 将35.236.172.131做为主，其余两台机器做从。ip和端口号按照机器和配置做相应修改。
replicaof 35.236.172.131 6379

# 配置日志路径，为了便于排查问题，指定redis的日志文件目录
logfile "/var/log/redis/redis.log"
```

### 启动容器

分别在主机和从机上按照上面的方法建立好配置文件，检查无误后就可以开始启动容器了。

我们在三台机器上分别将容器别名指定为`redis-1, redis-2, redis-3`，这样便于区分与说明，docker通过`--name`参数来指定容器的别名。redis-1是master上容器的别名，redis-2和redis-3是两个slave上的别名。

下面以运行redis-3容器为例说明容器的启动过程。另外两台机器上的容器redis-1和redis-2操作是相同的，只是要注意master的配置文件和slave不同。不过首先要启动主服务器，也就是redis-1容器。然后再启动redis-2和redis-3。
```
# 首先以后台模式运行容器
$ docker run -it --name redis-3 -v /root/redis.conf:/usr/local/etc/redis/redis.conf -d -p 6379:6379 redis /bin/bash
# 容器成功启动后，会打印一个长串的容器ID
a3952342094dfd5a56838cb6becb5faa7a34f1dbafb7e8c506e9bd7bb1c2951b
# 通过ps命令查看容器的状态，可以看到redis-3已经启动
$ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
a3952342094d        redis               "docker-entrypoint.s…"   8 minutes ago       Up 8 minutes        0.0.0.0:6379->6379/tcp   redis-3
```

上面已经启动了容器，接下来进入容器里启动redis服务器。
```
# 以交互模式进入容器redis-3
$ docker exec -it redis-3 bash

# 创建日志文件目录
$ mkdir /var/log/redis/
$ touch /var/log/redis/redis.log

# 启动redis服务器，如果没有任何输出，就说明成功了
$ redis-server /usr/local/etc/redis/redis.conf

# 在容器里启动一个redis客户端
$ redis-cli 

# 执行info命令，查看服务器状态
127.0.0.1:6379> info
...
# 如果是主，这里的role的值会是master，如果是从，这里的role的值会是slave
role:slave
# 对于slave，还要查看master_link_status这个属性值。slave上这个属性值为up就说明主从复制是OK的，否者就有问题。如果从机状态不为up，首先排查主机的端口是否被限，然后查看redis日志排查原因
master_link_status:up
...

# 最后退出容器
$ exit
```

### 验证主从复制

主从搭建成功后，可以通过在master上写入一个key-value值，查看是否会同步到slave上，来验证主从同步是否能成功。
```
# 以交互模式进入容器redis-1中
$ docker exec -it redis-1 bash
```
运行一个redis-cli，向test_key写入一个值
```
$ redis-cli
127.0.0.1:6379> set test_key hello-world
OK
```

在任意slave机器上进入容器，也运行一个redis-cli，查询这个key的值。如果能查询到这个值，且与主机上的值相同，说明主从同步成功。经测试，主动同步成功。
```
127.0.0.1:6379> get test_key 
"hello-world"
```

## 添加哨兵

主从结构搭建成功了，系统的可用性变高了，但是如果主发生故障，需要人工手动切换从机为主机。这种切换工作不仅浪费人力资源，更大的影响是主从切换期间这段时间redis是无法对外提供服务的。因此，哨兵系统被开发出来了，哨兵可以在主发生故障后，自动进行故障转移，从从机里选出一台升级为主机，并持续监听着原来的主机，当原来的主机恢复后，会将其作为新主的从机。

哨兵先监听主，通过对主发送info命令，获取到从的信息，然后也会监听到从。另外哨兵都会像主订阅__sentinel__:hello频道，当有新的哨兵加入时，会向这个频道发送一条信息，这条信息包含了该哨兵的IP和端口等信息，那么其他已经订阅了该频道的哨兵就会收到这条信息，就知道有一个新的哨兵加入。这些哨兵会与新加入和哨兵建立连接，选主是需要通过这个连接来进行投票。这个关系可以用下面这个图来描述


![](https://user-gold-cdn.xitu.io/2019/7/30/16c403f124eedd98?w=880&h=586&f=png&s=53121)

### 获取并修改sentinel配置文件
通过wget命令获取sentinel的配置文件
```
wget http://download.redis.io/redis-stable/sentinel.conf
```

修改配置文件以下几项
```
# 让sentinel服务后台运行
daemonize yes 

# 修改日志文件的路径
logfile "/var/log/redis/sentinel.log"

# 修改监控的主redis服务器
# 最后一个2表示，两台机器判定主被动下线后，就进行failover(故障转移)
sentinel monitor mymaster 35.236.172.131 6379 2
```

### 启动容器

与启动redis容器类似，启动一个别名为sentinel的容器
```
$ docker run -it --name sentinel -p 26379:26379 -v /root/sentinel.conf:/usr/local/etc/redis/sentinel.conf -d redis /bin/bash
```

### 运行哨兵
```
# 进入容器
$ docker exec -it sentinel bash

# 创建日志目录和文件
$ mkdir /var/log/redis
$ touch /var/log/redis/sentinel.log

# 启动哨兵
redis-sentinel /usr/local/etc/redis/sentinel.conf 

# 查看日志，哨兵成功监听到一主和两从的机器
18:X 11 Jul 2019 13:25:55.416 # +monitor master mymaster 35.236.172.131 6379 quorum 2
18:X 11 Jul 2019 13:25:55.418 * +slave slave 35.201.200.251:6379 35.201.200.251 6379 @ mymaster 35.236.172.131 6379
18:X 11 Jul 2019 13:25:55.421 * +slave slave 34.80.172.42:6379 34.80.172.42 6379 @ mymaster 35.236.172.131 6379
```

在另外两台机器上按照同样的方法在一个容器中运行sentinel，sentinel都使用相同的配置文件。

### 验证failover(故障转移)

为了验证哨兵机制下的自动主从切换，我们将主上的redis进程kill掉。

稍等几秒钟后，就有另外一台从升级为主机，实验时是第三台机器，也就是redis-3升级为了主，用info命令查询可以看到redis-3服务器的角色变成的master。说明自动主从切换成功。
```
127.0.0.1:6379> info
...
# Replication
role:master
...
```

然后重新启动之前被kill掉的master服务器，启动后用info命令查看，可以发现其变成了redis-3的从服务器。

下面这段日志，描述了35.236.172.131作为主启动，执行故障转移的master sentinel选举，执行故障转移，建立新的主从关系。

```
root@4355ca3260c5:/var/log/redis# cat sentinel.log 
17:X 11 Jul 2019 13:25:55.395 # oO0OoO0OoO0Oo Redis is starting oO0OoO0OoO0Oo
17:X 11 Jul 2019 13:25:55.395 # Redis version=5.0.5, bits=64, commit=00000000, modified=0, pid=17, just started
17:X 11 Jul 2019 13:25:55.395 # Configuration loaded
18:X 11 Jul 2019 13:25:55.398 * Running mode=sentinel, port=26379.
18:X 11 Jul 2019 13:25:55.398 # WARNING: The TCP backlog setting of 511 cannot be enforced because /proc/sys/net/core/somaxconn is set to the lower value of 128.
18:X 11 Jul 2019 13:25:55.416 # Sentinel ID is 7d9a7877d4cffb6fec5877f605b975e00e7953c1
18:X 11 Jul 2019 13:25:55.416 # +monitor master mymaster 35.236.172.131 6379 quorum 2
18:X 11 Jul 2019 13:25:55.418 * +slave slave 35.201.200.251:6379 35.201.200.251 6379 @ mymaster 35.236.172.131 6379
18:X 11 Jul 2019 13:25:55.421 * +slave slave 34.80.172.42:6379 34.80.172.42 6379 @ mymaster 35.236.172.131 6379
18:X 11 Jul 2019 13:26:25.460 # +sdown slave 35.201.200.251:6379 35.201.200.251 6379 @ mymaster 35.236.172.131 6379
18:X 11 Jul 2019 14:04:23.390 * +sentinel sentinel 09aa7d2098ad2dc52e6e07d7bc6670f00f5ff3e3 172.17.0.3 26379 @ mymaster 35.236.172.131 6379
18:X 11 Jul 2019 14:04:25.418 * +sentinel-invalid-addr sentinel 09aa7d2098ad2dc52e6e07d7bc6670f00f5ff3e3 172.17.0.3 26379 @ mymaster 35.236.172.131 6379
18:X 11 Jul 2019 14:04:25.418 * +sentinel sentinel 7d9a7877d4cffb6fec5877f605b975e00e7953c1 172.17.0.3 26379 @ mymaster 35.236.172.131 6379
18:X 11 Jul 2019 14:04:25.456 * +sentinel-address-switch master mymaster 35.236.172.131 6379 ip 172.17.0.3 port 26379 for 09aa7d2098ad2dc52e6e07d7bc6670f00f5ff3e3
18:X 11 Jul 2019 14:08:34.338 * +sentinel-invalid-addr sentinel 09aa7d2098ad2dc52e6e07d7bc6670f00f5ff3e3 172.17.0.3 26379 @ mymaster 35.236.172.131 6379
18:X 11 Jul 2019 14:08:34.338 * +sentinel sentinel 28d3c0e636fa29ac9fb5c3cc2be00432c1b0ead9 172.17.0.3 26379 @ mymaster 35.236.172.131 6379
18:X 11 Jul 2019 14:08:36.236 * +sentinel-address-switch master mymaster 35.236.172.131 6379 ip 172.17.0.3 port 26379 for 09aa7d2098ad2dc52e6e07d7bc6670f00f5ff3e3
18:X 11 Jul 2019 14:11:12.151 # +sdown master mymaster 35.236.172.131 6379
18:X 11 Jul 2019 14:11:12.214 # +odown master mymaster 35.236.172.131 6379 #quorum 4/2
18:X 11 Jul 2019 14:11:12.214 # +new-epoch 1
18:X 11 Jul 2019 14:11:12.214 # +try-failover master mymaster 35.236.172.131 6379
18:X 11 Jul 2019 14:11:12.235 # +vote-for-leader 7d9a7877d4cffb6fec5877f605b975e00e7953c1 1
18:X 11 Jul 2019 14:11:12.235 # 7d9a7877d4cffb6fec5877f605b975e00e7953c1 voted for 7d9a7877d4cffb6fec5877f605b975e00e7953c1 1
18:X 11 Jul 2019 14:11:12.235 # 28d3c0e636fa29ac9fb5c3cc2be00432c1b0ead9 voted for 7d9a7877d4cffb6fec5877f605b975e00e7953c1 1
18:X 11 Jul 2019 14:11:12.235 # 09aa7d2098ad2dc52e6e07d7bc6670f00f5ff3e3 voted for 7d9a7877d4cffb6fec5877f605b975e00e7953c1 1
18:X 11 Jul 2019 14:11:12.294 # +elected-leader master mymaster 35.236.172.131 6379
18:X 11 Jul 2019 14:11:12.294 # +failover-state-select-slave master mymaster 35.236.172.131 6379
18:X 11 Jul 2019 14:11:12.394 # -failover-abort-no-good-slave master mymaster 35.236.172.131 6379
18:X 11 Jul 2019 14:11:12.453 # Next failover delay: I will not start a failover before Thu Jul 11 14:17:12 2019
18:X 11 Jul 2019 14:11:13.050 # +config-update-from sentinel 28d3c0e636fa29ac9fb5c3cc2be00432c1b0ead9 172.17.0.3 26379 @ mymaster 35.236.172.131 6379
18:X 11 Jul 2019 14:11:13.050 # +switch-master mymaster 35.236.172.131 6379 34.80.172.42 6379
18:X 11 Jul 2019 14:11:13.050 * +slave slave 35.201.200.251:6379 35.201.200.251 6379 @ mymaster 34.80.172.42 6379
18:X 11 Jul 2019 14:11:13.050 * +slave slave 35.236.172.131:6379 35.236.172.131 6379 @ mymaster 34.80.172.42 6379
18:X 11 Jul 2019 14:11:43.077 # +sdown slave 35.236.172.131:6379 35.236.172.131 6379 @ mymaster 34.80.172.42 6379
18:X 11 Jul 2019 14:11:43.077 # +sdown slave 35.201.200.251:6379 35.201.200.251 6379 @ mymaster 34.80.172.42 6379
18:X 12 Jul 2019 01:54:05.142 # -sdown slave 35.236.172.131:6379 35.236.172.131 6379 @ mymaster 34.80.172.42 6379
18:X 12 Jul 2019 01:54:15.087 * +convert-to-slave slave 35.236.172.131:6379 35.236.172.131 6379 @ mymaster 34.80.172.42 6379
```

## 总结

redis通过主从复制来实现高可用，但是发生故障时需要人工进行主从切换，效率低下。哨兵机制实现了redis主从的自动切换，提高了redis集群的可用性，提高了redis集群的故障转移效率。