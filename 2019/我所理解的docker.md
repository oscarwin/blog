# 我所理解的docker

## docker是什么

docker是容器领域应用最为广泛的开源软件。docker的中文意思中有集装箱的意思，正如其名，容器就像是一个集装箱，在这个集装箱里你可以放任何东西，集装箱与集装箱之间互不干涉。只不过这个集装箱里装的是运行的软件而已。这个虚拟机的特点有很多相似的地方。

但是docker与虚拟机有很大的区别。

我们知道操作系统分为用户空间和内核空间。
虚拟机是在宿主主机之上虚拟出了操作系统的内核空间和用户空间，每个虚拟机都有独立的内核空间。而docker则公用了宿主主机的内核空间，只在其基础上虚拟出用户空间。docker之所可以这么干，一是因为不同的Linux系统都是在同一套内核级基础上进行开发的；二是用户运行的软件，大多数时间运行在用户空间，而且数据也存储在用户空间。内核空间只是负责调度、内存管理等底层的工作。


![](https://user-gold-cdn.xitu.io/2019/3/5/1694e10db0e24c09?w=916&h=466&f=png&s=130252)


## 为什么要用docker

用过python的人都知道，python分为python2和python3，但是这两哥们并不像我们之前的认知一样，一般软件或语言的升级都是像下兼容的，他们两区别太大，不能兼容，更烦人的是他们依赖的库也是分版本的。因此我们开发程序时，为了避免两个环境混合，需要设置虚拟环境来实现隔离，但是配置这种东西一是麻烦，二是对新手不友好，而且还很容易出问题。

docker的出现就是为解决这种事情而生的。现在我们可以分别在两个容器中运行python2和python3，两者相互不影响。除此之外，容器为服务的部署也提供了极大的便利性。假设有个服务是安装在centos上，假设现在要在另一台ubuntu机器上同样的安装，无需做任何改变，直接可以运行，因为可以在ubuntu的宿主主机上可以虚拟出一个centos。docker的出现为微服务架构的设计提供了极大的便利。

当然，同样的事情通过虚拟机也能做到，但是一个虚拟机需要占用的资源太多，一个ubuntu系统可能就要好几G，而一个docker容器中安装一个ubuntu只要上百兆，而且docker的管理比虚拟机要灵活方便太多。

## docker的组成

docker分为以下几个部分：

1. 客户端(client)
2. 服务端(daemon)
3. 镜像(image)
4. 仓库(registry)
5. 容器(container)

docker架构图：
![](https://user-gold-cdn.xitu.io/2019/3/5/1694e2830d0c0eeb?w=884&h=463&f=png&s=300521)

### 客户端

docker采用的是c/s模式，客户端是用来与用户进行交互的，客户端将命令发送给服务端进行处理。

![](https://user-gold-cdn.xitu.io/2019/3/5/1694e2f4c0683339?w=1059&h=772&f=png&s=84061)

### 服务端

docker服务端是一个后台daemon，只有启动了后台daemon，docker才能执行来至于客户端的命令，否则会出现如下错误提示。

    Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?

管理docker daemon的命令：

启动：

    systemctl start docker.service
    或
    service docker start

停止：

    systemctl stop docker.service
    或
    service docker stop
    
重启：

    systemctl restart docker.service
    或
    service docker restart
    
### 镜像

如果把容器看做是一个软件，那么镜像就可以看做是这个容器的安装包。可以通过一个镜像来启动多个容器。

镜像与传统的软件安装包有个很大的区别在于，镜像是可以堆叠的，在一个镜像的基础上


### 容器

容器就是镜像生成的运行实例。用户可以通过客户端管理docker容器。也可以进入到docker容器内，docker容器内就像是一个虚拟的操作系统环境，里面一样有/bin /home /var等目录

### 仓库

仓库就是存放docker镜像的地方，仓库分为公有和私有。

Docker Hub（https://hub.docker.com/） 是默认的 Registry，由 Docker 公司维护，上面有数以万计的镜像，用户可以自由下载和使用。

出于对速度或安全的考虑，用户也可以创建自己的私有 Registry。