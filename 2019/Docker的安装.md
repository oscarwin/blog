# Docker 的安装

从 2017 年 3 月开始 docker 在原来的基础上分为两个分支版本: Docker CE 和 Docker EE。Docker CE 即社区免费版，Docker EE 即企业版，强调安全，但需付费使用。本文介绍 Docker CE 的安装使用。

# centos下安装docker

## 安装条件

**CentOS版本要求**

CentOS 7 (64-bit)，CentOS 6.5 (64-bit)或更高的版本

**内核要求**
- CentOS 仅发行版本中的内核支持 Docker。
- Docker 运行在CentOS 7 上，要求系统为64位、系统内核版本为 3.10 以上。
- Docker 运行在CentOS 6.5或更高的版本的 CentOS 上，要求系统为64位、系统内核版本为 2.6.32-431 或者更高版本

### 方式1：使用脚本安装

#### 1. 确保 yum 包更新到最新。

    $ sudo yum update

#### 2. 切换到root用户

    $ sudo su -

#### 3. 执行脚本安装

    $ curl -sSL https://get.docker.com/ | sh

### 方式2：使用yum安装

#### 1. 查看版本是否满足要求

使用`uname -r`命令查看系统内核版本，看是否为3.10以上
    
    [root@hostname ~]# uname -r
    3.10.0-514.26.2.el7.x86_64

#### 2. 移除旧版本docker(如果安装过)

    sudo yum remove docker \
                    docker-client \
                    docker-client-latest \
                    docker-common \
                    docker-latest \
                    docker-latest-logrotate \
                    docker-logrotate \
                    docker-selinux \
                    docker-engine-selinux \
                    docker-engine

#### 3. 安装一些必要的系统工具

    sudo yum install -y yum-utils device-mapper-persistent-data lvm2

#### 4. 添加软件源信息

    sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
    
#### 5. 更新 yum 缓存

    sudo yum makecache fast
    
#### 6. 安装 Docker-ce：

    sudo yum -y install docker-ce

#### 7. 启动 Docker 后台服务

    sudo systemctl start docker

#### 8. 测试运行 hello-world

    [root@hostname ~]# docker run hello-world
    
运行结果如下，则表示docker安装成功
![](https://user-gold-cdn.xitu.io/2019/2/21/16910b45fe4893c4?w=798&h=363&f=png&s=36671)


# 运行第一个容器

刚刚测试启动的hello-world就是一个容器，不过没有什么实际作用。接下来运行一个httpd的容器，来直观感受一下。

    $ docker run -d -p 80:80 httpd 

运行过程如下：

    [root@izwz9alpqga9jjum6tmmkyz ~]# docker run -d -p 80:80 httpd
    Unable to find image 'httpd:latest' locally
    latest: Pulling from library/httpd
    6ae821421a7d: Pull complete 
    0ceda4df88c8: Pull complete 
    24f08eb4db68: Pull complete 
    ddf4fc318081: Pull complete 
    fc5812428ac0: Pull complete 
    Digest: sha256:214019bfc77677ac1f0c86b3a96e2b91600de7e2224f195b446cb13572cebb6b
    Status: Downloaded newer image for httpd:latest
    b554654b25376a48d41b1e4df53703cd7891fdbdcbc497c8f4d449000c4b9913
    
查看一下，容器是否在运行状态。可以看到STATUS为Up 2 seconds，说明在2秒钟前启动的。

    $ docker ps -a
    
    CONTAINER ID        IMAGE                      COMMAND              CREATED             STATUS                      PORTS                  NAMES
    7f1c454fc41a        httpd                      "httpd-foreground"   3 seconds ago       Up 2 seconds                0.0.0.0:80->80/tcp   vigorous_driscoll

可以通过浏览器访问主机的80端口，可以看到Apache服务器已经安装好了


![](https://user-gold-cdn.xitu.io/2019/2/22/169138878a97e9cc?w=635&h=163&f=png&s=11408)

上面的docker命令实际上做了一下几部工作：

1. 在本地查找httpd的镜像
2. 本地没有镜像，从公共仓库下载httpd最新的镜像文件
3. 启动httpd容器，将本地的80端口映射到容器的80端口


# 镜像加速

我们启动容器都需要一个镜像，但是docker官方的公共镜像仓库(https://hub.docker.com/)服务器在国外，下载速度可能会比较慢。还好国内有很多厂商提供了镜像下载的服务，我们需要在docker的配置文件中配置国内镜像的地址。

centos系统可以直接执行下面的命令配置加速器。

    curl -sSL https://get.daocloud.io/daotools/set_mirror.sh | sh -s http://f1361db2.m.daocloud.io

这个是`https://www.daocloud.io/mirror`该网站提供的镜像加速服务。

![](https://user-gold-cdn.xitu.io/2019/3/5/1694df9e5ae28074?w=1097&h=866&f=png&s=81680)

除了这种方法以外，我们也可以自己手动添加任意的加速镜像地址：

在配置文件`/etc/docker/daemon.json`中添加如下数据，该文件不存在的话，则新建一个。其中中括号里就是加速镜像的地址。

    {
      "registry-mirrors": ["https://registry.docker-cn.com"]
    }