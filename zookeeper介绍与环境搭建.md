# zookeeper介绍与环境搭建

# 简介

zookeeper是一个分布式服务框架，是Apache Hadoop 的一个子项目，它主要是用来解决分布式应用中经常遇到的一些数据管理问题，如：统一命名服务、状态同步服务、集群管理、分布式应用配置项的管理等。

在介绍zookeeper集群是，zookeeper的机器称为服务端，用zookeeper来管理的分布式系统的机器称为客户端。

首先zookeeper是用来服务于分布式系统的，而zookeeper集群本身也是一个分布式系统，zookeeper集群至少需要三台机器，因为zookeeper集群会选举出一个leader，而投票选举leader需要严格意义上的多数成员赞同(只要半数以上节点存活，ZooKeeper 就能正常服务)，所以zookeeper集群一般是奇数个，但也并不是强制要求。一般采用3台或5台的配置。zookeeper集群是典型的cp系统，保证一致性和分区容错性的条件下，牺牲了一定的可用性。3台配置可以容忍一台机器挂掉，5台的配置可以容忍第二台挂掉。4台机器相比3台机器并没有任何优势，因为4台机器也只能容忍一台机器挂掉。

机器越多，zookeeper的读性能越好，可以同时服务更多的客户端。但是，机器越多写性能会下降。因为，zookeeper为了保证严格的一致性，必须保证所有可用的服务都成功写入。因此，机器的数量也不可盲目增多。5台机器是比较合理的配置，可以同时容忍两台机器挂掉，而写入的成本也可以接受。

![](https://user-gold-cdn.xitu.io/2019/3/8/1695d14f9f432680?w=646&h=219&f=png&s=39012)

ZooKeeper还有一个重要的概念——节点。节点是zookeeper的一种文件存储模型，类似于linux的目录结构。每一个节点可以存储数据，也可以作为下一个节点的父节点。目录树的结构用来保证每个节点的唯一性。可以存储在节点中的数据的默认最大大小为 1 MB。因此，即使ZooKeeper的层次结构看起来与文件系统相似，也不应该将它用作一个通用的文件系统。相反，应该只将它用作少量数据的存储机制，以便为分布式应用程序提供可靠性、可用性和协调。

![](https://user-gold-cdn.xitu.io/2019/3/15/16981b0d7b560fa7?w=576&h=328&f=jpeg&s=11010)

## zookeeper可以用来做什么

**注册中心：**
在分布式系统中，常见的场景是通过RPC调用其他机器上的服务。而调用其他机器上的服务首先要去配置中心读取提供该服务主机的信息。服务与主机之间的映射关系是在服务注册的时候建立的。比如有两台机器提供了一个称为app1的服务，这两台提供服务的机器在启动的时候调用zookeeper的api在服务的根节点app1上进行注册，生成了子节点/app1/c1和/app1/c2。而调用方则从zookeeper订阅该服务，获得该提供该服务的主机信息。当有新的服务提供者注册到zookeeper的/app1/c3节点上时，注册中心会通知所有的订阅方发生的变更。

![](https://user-gold-cdn.xitu.io/2019/3/15/16981d516a096839?w=807&h=379&f=png&s=23201)

**配置管理：**后端的服务很多地方都需要使用配置文件，分布式系统中，一个配置文件的更改需要一台台的修改，这样的体力活有了zookeeper就不在需要了。可以将配置文件放在zookeeper的一个节点中，所有需要使用这个配置文件的去订阅这个节点，这样当配置文件发生变更时，zookeeper会将变更后的配置文件下发到所有订阅的机器上。

**分布式锁：**在分布式系统中可能会出现同时修改竞争资源的场景，使用redis做分布式锁的应用场景比较多，而zookeeper同样可以实现分布式锁。

**领导人选举：**在后端服务中，很多时候只需要单机去运行，但是单机运行会造成单点。为了避免单点，可以在多机上部署服务，然后通过zookeeper给订阅了该服务的客户端选出一个leader。我们可以在程序中进行逻辑判断，如果是leader才执行相应的业务逻辑，否者不执行。当leader的机器宕掉后，又会有新的leader被选出，这样就解决单点问题。

leader的选举原理是，客户端抢占去创建一个指定的节点，并将客户端的ip地址写入到节点的data中，如果节点不存在，那么就将成功创建了这个节点的客户端定义为leader，如果节点已经存在，则获取该节点的数据，将节点中存储的ip与本机的ip进行比较，如果相同说明就是leader，否者就不是。

zookeeper实现简单的分布式锁，也可以采用相同的原理，创建节点成功就表示抢锁成功。

# zookeeper集群搭建

## 1 使用yum安装JVM

zookeeper是运行在JVM环境下的，所以首先要安装JVM环境。

系统：centos

### 查看yum库中jdk的版本

    [root@localhost ~]# yum search java | grep jdk


### 选择1.8的版本安装

    yum install java-1.8.0-openjdk
    
### 修改环境变量

jdk的默认安装路径为`/usr/lib/jvm`

    [root@localhost ~]# cd /usr/lib/jvm
    [root@localhost ~]# ll
    total 4
    drwxr-xr-x 3 root root 4096 Mar  8 12:48 java-1.8.0-openjdk-1.8.0.201.b09-2.el7_6.x86_64
    lrwxrwxrwx 1 root root   21 Mar  8 12:48 jre -> /etc/alternatives/jre
    lrwxrwxrwx 1 root root   27 Mar  8 12:48 jre-1.8.0 -> /etc/alternatives/jre_1.8.0
    lrwxrwxrwx 1 root root   35 Mar  8 12:48 jre-1.8.0-openjdk -> /etc/alternatives/jre_1.8.0_openjdk
    lrwxrwxrwx 1 root root   51 Mar  8 12:48 jre-1.8.0-openjdk-1.8.0.201.b09-2.el7_6.x86_64 -> java-1.8.0-openjdk-1.8.0.201.b09-2.el7_6.x86_64/jre
    lrwxrwxrwx 1 root root   29 Mar  8 12:48 jre-openjdk -> /etc/alternatives/jre_openjdk

可以看到有一个软连接，将这个软连接的路径写到home目录中

    jre-1.8.0-openjdk-1.8.0.201.b09-2.el7_6.x86_64 -> java-1.8.0-openjdk-1.8.0.201.b09-2.el7_6.x86_64/jre
    
在`/etc/profile`文件中添加环境变量

    #set java environment
    JAVA_HOME=/usr/lib/jvm/jre-1.8.0-openjdk-1.8.0.201.b09-2.el7_6.x86_64
    JRE_HOME=$JAVA_HOME/jre
    CLASS_PATH=.:$JAVA_HOME/lib/dt.jar:$JAVA_HOME/lib/tools.jar:$JRE_HOME/lib
    PATH=$PATH:$JAVA_HOME/bin:$JRE_HOME/bin
    export JAVA_HOME JRE_HOME CLASS_PATH PATH
    
运行`source /etc/profile`让环境变量生效

验证安装是否成功

    [root@localhost ~]# java -version
    openjdk version "1.8.0_201"
    OpenJDK Runtime Environment (build 1.8.0_201-b09)
    OpenJDK 64-Bit Server VM (build 25.201-b09, mixed mode)

## 2 安装zookeeper

### 下载

使用清华大学开源软件镜像源下载

    wget https://mirrors.tuna.tsinghua.edu.cn/apache/zookeeper/zookeeper-3.4.13/zookeeper-3.4.13.tar.gz

解压文件，放到一个目录中，我的放在opt目录

    tar -xzvf zookeeper-3.4.13.tar.gz

新建一个目录用来保存zookeeper的本地数据

    mkdir /var/lib/zookeeper
    
修改配置文件，在目录`zookeeper-3.4.13/conf/`中拷贝配置文件`zoo_sample.cfg`，并将其重命名为`zoo.cfg`。修改后的配置文件如下所示：

    # The number of milliseconds of each tick
    tickTime=2000
    # The number of ticks that the initial 
    # synchronization phase can take
    initLimit=10
    # The number of ticks that can pass between 
    # sending a request and getting an acknowledgement
    syncLimit=5
    # the directory where the snapshot is stored.
    # do not use /tmp for storage, /tmp here is just 
    # example sakes.
    dataDir=/var/lib/zookeeper
    # the port at which the clients will connect
    clientPort=2181
    # the maximum number of client connections.
    # increase this if you need to handle more clients
    #maxClientCnxns=60
    #
    # Be sure to read the maintenance section of the 
    # administrator guide before turning on autopurge.
    #
    # http://zookeeper.apache.org/doc/current/zookeeperAdmin.html#sc_maintenance
    #
    # The number of snapshots to retain in dataDir
    #autopurge.snapRetainCount=3
    # Purge task interval in hours
    # Set to "0" to disable auto purge feature
    #autopurge.purgeInterval=1
    
    server.1=47.107.41.24:2888:3888
    server.2=34.73.24.64:2888:3888
    server.3=35.220.130.110:2888:3888
    

启动zookeeper

    [root@localhost]# cd zookeeper-3.4.13/bin
    [root@localhost]# ./zkServer.sh start
    



    