# MySQL 的锁，看这篇就够了

MySQL里有非常多锁的概念，经常可以听到的有：乐观锁、悲观锁、行锁、表锁、Gap锁（间隙锁）、MDL锁（元数据锁）、意向锁、读锁、写锁、共享锁、排它锁。这么锁一听就让人头大，于是去看一些博客，有些讲乐观锁、悲观锁，有些在讲读锁、写锁，于是乐观锁和悲观锁好像理解了，读锁写锁好像也理解了，但是我任然不知道怎么用，也不知道乐观锁与读锁写锁有没有什么关系？再看了很多文章后，逐渐弄懂了它们之间的关系，于是写下这篇文章来梳理思路。能力有限，难免有误，请酌情参考。

虽然上面列举了很多锁的名词，但是这些锁其实并不是在同一个维度上的，这就是我之所以含糊不清的原因。接下来从不同的维度来分析 MySQL 的锁。

## 读锁和写锁

首先读锁还有一个名称叫共享锁，写锁也相应的还有个名称叫排它锁，也就是说共享锁和读锁是同一个东西，排它锁和写锁是同一个东西。读锁、写锁是系统实现层面上的锁，也是最基础的锁。读锁和写锁还是锁的一种性质，比如行锁里，有行写锁和行读锁。MDL 锁里也有 MDL 写锁和 MDL 读锁。读锁和写锁的加锁关系如下，Y 表示可以共存，X 表示互斥。

| | 读锁 | 写锁 |
| :---: | :---: | :---: |
| 读锁 | Y | X |
| 写锁 | X | X |

从这个表格里可以知道读锁和写锁不能共存，请考虑这样一个场景，一个请求占用了读锁，这时又来了一个请求要求加写锁，但是资源已经被读锁占据，写锁阻塞。这样本没有问题，但是如果后续不断的有请求占用读锁，读锁一直没有释放，造成写锁一直等待。这样写锁就被饿死了，为了避免这种情况发生，数据库做了优化，当有写锁阻塞时，后面的读锁也会阻塞，这样就避免了饿死现象的发生。后面还会再次提到这个现象。

之前的文章已经介绍了 MySQL 的存储模型，对于 InnoDB 引擎而言，采用的是 B+ 树索引，假设需要将整个表锁住那么需要在整个 B+ 树的每个节点上都加上锁，显然这是个非常低效的做法。因此，MySQL 提出了意向锁的概念，意向锁就是如果要在一个节点上加锁就必须在其所有的祖先节点加上意向锁。关于意向锁还有更多复杂设计，如果想了解可以查看 《数据库系统概率》 一书。

## 表锁和行锁

表锁和行锁是从加锁的粒度来区分不同的锁。除了表锁和行锁以外还有更大粒度的锁——全局锁。需要注意的是 MyISAM 引擎并不支持行锁，InnoDB 引擎支持行锁。

**全局锁：** 全局锁会锁住整个数据库，MySQL 使用 flush tables with read lock 命令来加全局锁，使用 unlock tables 解锁。线程退出后锁也会自动释放。当加上全局锁以后，除了当前线程以外，其他线程的更新操作都会被阻塞，包括增删改数据表中的数据、建表、修改表结构等。全局锁的典型使用场景是全库的逻辑备份。

**表锁：** 表锁会锁住一张表，MySQL 使用 lock tables <table name> read/write 命令给表加上读锁或写锁，必须通过 unlock tables 命令释放表锁，如果没有主动释放，即使事务提交了表锁也不会被自动释放。通过 lock tables t read 给表 t 加上读锁后，当前线程只能访问表 t，不能访问数据库中的其他表，对表 t 也只有读权限，不能进行修改操作。通过 lock tables t write 给表 t 加上写锁后，当前线程只能访问表 t，不能访问数据库中的其他表，对表 t 有读写权限。

**行锁：** 行锁会锁锁住表中的某一行或者多行，MySQL 使用 lock in share mode 命令给行加读锁，用 for update 命令给行加写锁，行锁不需要显示释放，当事务被提交时，该事务中加的行锁就会被释放。通过 select k from t where k = 1 for update 命令可以锁住 k 为 1 的所有行。另外当使用 update 命令更新表数据时，会自动给命中的行加上行锁。另外 MySQL 加行锁时并不是一次性把所有的行都加上锁，执行一个 update 命令之后，server 层将命令发送给 InnoDB 引擎，InnoDB 引擎找到第一条满足条件的数据，并加锁后返回给 server 层，server 层更新这条数据然后传给 InnoDB 引擎。完成这条数据的更新后，server 层再取下一条数据。

我们用一个例子来验证这个过程，首先执行如下命令建表并插入几行数据

```SQL
mysql-> create table t(id int not null auto_increment, c int not null, primary key(id))ENGINE=InnoDB;
mysql-> insert into t(id, c) values (1, 1), (2, 2), (3, 3);
```

| 事务 A | 事务 B | 事务 C |
| :---: | :---: | :---: |
| begin; | 
| select * from t where id = 3 for update; |
| | update t set c = 0 where id = c; |
| | | set session transaction isolation level READ UNCOMMITTED; select *  from t; |
| commit |

事务 A 执行 select * from t where id = 3 for update 将 id 等于3的行锁住，事务 B 执行 update 命令的时候被阻塞。这时候再开启事务 C，并且将事务 C 的隔离级别修改为未提交读，得到的如下表所示，发现前两行已经被更新，最后 id 为 3 的行没有更新，说明事务 B 是阻塞在这里了。

```SQL
mysql> select *  from t;
+----+---+
| id | c |
+----+---+
|  1 | 0 |
|  2 | 0 |
|  3 | 3 |
+----+---+
```

**行锁与表锁的关系：**
行锁与表锁是独立存在的，可以分别进行加锁。如果把一张数据表比作一个大宅院，大宅院里有很多小房间，大宅院只有一个大门。那么表锁就好比大门的钥匙，而行锁是每个房间的钥匙。如果要进入大宅院，要么是有大门的钥匙，要么是大门压根没有关闭。如果要进入某个小房间，首先要有大门的钥匙或者大门没有关，而且有这个小房间的钥匙或者这个小房间的门也没有关。

| 事务 A | 事务 B |
| :---: | :---: |
| begin; | begin; |
| select * from t where id = 2 for update; |
| | lock table t write; |
| | select * from t where id = 1 for update; |
| update t set a = 0 where id = 2; <执行被阻塞> | |
| | select * from t where id = 2 for update; <执行被阻塞> |
| | unlock tables; |
| commit; | commit; |

事务 A 首先给 id 等于2的行加写锁，加锁成功。然后事务 B 给表 t 加表的写锁，加锁成功。然后事务 B 给 id 等于1的行加写锁，加锁成功。最后事务 A 对 id 等于2的行进行修改被阻塞，因为事务 A 没有拿到大门钥匙。拿到大门钥匙的事务 B 就可以给其他小房间加锁。但是事务 B 也不能给 id 为2的行加锁，会被阻塞。这时候事务 A 在等待事务 B 的表锁释放，事务 B 在等待事务 A 在 id 为2的行的行锁释放，从而发生死锁。发生死锁后，加锁会被因为超时而释放。

## 乐观锁和悲观锁

乐观锁和悲观锁与前面介绍的行锁与表锁又不是同一个层次的概念了，乐观锁与悲观锁我认为是一种加锁的思想，悲观锁是在对资源进行操作前就先锁住这个资源，然后进行操作，而乐观锁是通过对比条件来判断，条件达成就更新，条件没达成就不更新。悲观锁需要依赖于行锁或者表锁来实现。而乐观锁并不需要依赖锁资源。

### 乐观锁

乐观锁总是假设不会发生冲突，因此读取资源的时候不加锁，只有在更新的时候判断在整个事务期间是否有其他事务更新这个数据。如果没有其他事务更新这个数据那么本次更新成功，如果有其他事务更新本条数据，那么更新失败。

### 悲观锁

悲观锁总是假设会发生冲突，因此在读取数据时候就将数据加上锁，这样保证同时只有一个线程能更改数据。文章前面介绍的表锁、行锁等都是悲观锁。

乐观锁和悲观锁是两种不同的加锁策略。乐观锁假设的场景是冲突少，因此适合读多写少的场景。悲观锁则正好相反，合适写多读少的场景。乐观锁无需像悲观锁那样维护锁资源，做加锁阻塞等操作，因此更加轻量化。

### 乐观锁的实现

乐观锁的实现有两种方式：版本号和 CAS 算法

**版本号**

通过版本号来实现乐观锁主要有以下几个步骤：

1 给每条数据都加上一个 version 字段，表示版本号

2 开启事务后，先读取数据，并保存数据里的版本号 version1，然后做其他处理

3 最后更新的时候比较 version1 和数据库里当前的版本号是否相同。用 SQL 语句表示就是 update t set version = version + 1 where version = version1。
根据前面事务的文章我们知道，update 操作时会进行当前读，因此即使是在可重复读的隔离级别下，也会取到到最新的版本号。如果没有其他事务更新过这条数据，那么 version 等于 version1，于是更新成功。如果有其他事务更新过这条数据，那么 version 字段的值会被增加，那么 version 不等于 version1，于是更新没有生效。

**CAS 算法**

CAS 是 compare and swap 的缩写，翻译为中文就是先比较然后再交换。CAS 实现的伪代码：

```
<< atomic >>
bool cas(int* p, int old, int new)  
{
    if (*p != old)
    {
        return false
    }
    *p = new
    return true
}
```
其中，p 是要修改的变量指针，old 是修改前的旧值，new 是将要写入的新值。这段伪代码的意思就是，先比较 p 所指向的值与旧值是否相同，如果不同说明数据已经被其他线程修改过，返回 false。如果相同则将新值赋值给 p 所指向的对象，返回 true。这整个过程是通过硬件同步原语来实现，保证整个过程是原子的。

大多数语言都实现了 CAS 函数，比如 C 语言在 GCC 实现：
```
bool__sync_bool_compare_and_swap (type *ptr, type oldval type newval, ...)
type __sync_val_compare_and_swap (type *ptr, type oldval type newval, ...)
```
无锁编程实际上也是通过 CAS 来实现，比如无锁队列的实现。CAS 的引入也带来了 ABA 问题。关于 CAS 后面再开一篇专门的文章来总结无锁编程。

## MDL 锁和 Gap 锁

### MDL 锁

MDL 锁也是一种表级锁，MDL 锁不需要显示使用。MDL 锁是用来避免数据操作与表结构变更的冲突，试想当你执行一条查询语句时，这个时候另一个线程正在修改表结构——删除表中的一个字段，那么两者就发生冲突了，因此 MySQL 在5.5版本以后加上了 MDL 锁。当对一个表做增删查改等数据变更语句时会加 MDL 读锁，当对一个表做表结构变更时会加 MDL 写锁。读锁相互兼容，读锁与写锁不能兼容，写锁与写锁也不能兼容。

MDL 需要注意的就是避免 MDL 写锁阻塞 MDL 读锁。什么意思呢？下面用一个例子来说明。

| 事务 A | 事务 B | 事务 C | 事务 D |
| :---: | :---: | :---: | :---: |
| select * from t |
| | select * from t |
| | | alter table t add c int |
| | | | select * from t |

事务 A 执行 select 后给表 t 加 MDL 读锁。事务 B 执行 select 后给表再次加上 MDL 读锁，读锁和读锁可以兼容。事务 C 执行 alter 命令时会阻塞，需要对表 t 加 MDL 写锁。事务 C 被阻塞问题并不大，但是会导致后面所有的事务都被阻塞，比如事务 D。这是为了避免写锁饿死的情况发生，MySQL 对加锁所做的优化，当有写锁在等待的时候，新的读锁都需要等待。如果事务 C 长时间拿不到锁，或者事务 C 执行的时间很长都会导致数据库的操作被阻塞。

为了避免这种事情发生有以下几点优化思路：

1 避免长事务。事务 A 和事务 B 如果是长事务就可能导致事务 C 阻塞在 MDL 写锁的时间比较长。

2 对于大表，修改表结构的语句可以拆分成多个小的事务，这样每次修改表结构时占用 MDL 写锁的时间会缩短。

3 给 alter 命令加等待超时时间

### Gap 锁

Gap 锁就是间隙锁，是 InnoDB 引擎为了避免幻读而引入的。在 MySQL的事务一文中已经谈到，InnoDB 引擎在可重复读隔离级别下可以避免幻读。间隙锁就是锁住数据行之间的间隙，避免新的数据插入进来。只有在进行当前读的时候才会加 gap 锁。关于什么是当前读，文章后面会提及。

用下面一个例子来演示一下幻读的出现，仍然使用上面的表 t：

| 事务 A | 事务 B |
| :---: | :---: |
| begin; set session transaction isolation level read committed; | 
| select * from t where c in (1, 2, 3) for update; (1) | |
| | begin; set session transaction isolation level read committed; |
| | insert into t (c) values(1); |
| select * from t where c in (1, 2, 3) for update; (2) | |
| commit; | commit; |

事务 A 两次查询语句的结果分别是：
```
// 第一次查询的结果
+----+---+
| id | c |
+----+---+
|  1 | 1 |
|  2 | 2 |
|  3 | 3 |
+----+---+
```
```
// 第二次查询的结果
+----+---+
| id | c |
+----+---+
|  1 | 1 |
|  2 | 2 |
|  3 | 3 |
|  4 | 1 |
+----+---+
```

首先启动事务时，将当前会话的隔离级别设置为提交读，在提交读的隔离级别下 binlog 的日志格式要设置为 ROW，否者 MySQL 会报错，因为在提交读的隔离级别下 STATEMENT 格式的 binlog 日志可能会出现不一致的情况。在事务 A 中第一次执行 `select * from t where c in (1, 2, 3) for update;` 时会将所有 c 值为1-3的行都加上行锁，但是这并不能阻止事务 B 插入一个新的行且字段 c 的值也为1。这就导致，虽然事务 B 还未提交，但是事务 A 第二次查询比第一次查询得到的行数更多，这种现象就叫幻读。

## 加锁实践

上面已经讲解了 MySQL 里主要的锁概念，在工作中我们遇到更多的时要分析语句执行时会怎么加锁，一是要保证数据安全，二是要避免死锁，三是要考虑锁带来的性能瓶颈，下面一起来分析一下在 MySQL 执行一条语句时会怎么加锁。

再分析 MySQL 怎么加锁之前还要在啰嗦几个概念：MVCC、两阶段加锁、快照读和当前读。

### 几个概念

**MVCC**

MVCC(Multi Version Concurrency Control) 就是多版本并发控制，这个技术的作用实现了读不加锁，读写不冲突。就是说即使两个事物并发执行，A 事务对某行数据加了写锁，这时候 B 事务仍然可以读这行数据，不会被阻塞。这就大大提高了数据库的性能。

**快照读和当前读**

在 MVCC 里分为快照读和当前读，快照读不加锁，当前读会加锁。select * from t where ??? 就是快照读。当前读有以下这些情况：

```SQL
select * from t where ? lock in share mode;
select * from t where ? for update;
update t set ?;
insert into t ?;
delete from t where ?;
```

**两阶段加锁**

在 InnoDB 引擎中，行锁是需要的时候才加上，但是要等到事务结束的时候才会释放。加锁和解锁是在两个阶段分别进行的，因此称为两阶段加锁。

### 一条简单语句的加锁

```
select * from t where id = 0 for update; 
```
这条语句是一个常见的加锁语句，使用 for update 加写锁，也被称为 X 锁。要分析这条语句怎么加锁，首先要明确几个前提条件：

1. 使用的搜索引擎是什么？无特殊说明情况下，本文都是讨论 InnoDB 引擎。

2. 事务的隔离级别是什么？

3. id 是否为主键索引？

4. id 是否为唯一辅助索引？

5. id 是否为普通辅助索引？

6. id 无索引？

只有明确了隔离级别和索引的情况分析加锁才是有意义的。虽然要考虑的条件多，但是咋们不用虚，因为在实际环境下，事务的隔离级别和索引的类型是已经确定的，对应的加锁结果也是唯一的。这里我们为了全面分析，因此考虑事务隔离级别和索引的每一种组合情况。

在提交读隔离级别和可重复度隔离级别下的加锁差别较大，因此对这两种隔离级别下的加锁情况分别进行分析。假设表有两个字段 (id int, user varchar(10))。

### 提交读隔离级别下的加锁分析

隔离级别为读提交时，MySQL 不会加 Gap 锁，因此只有行锁加锁比较容易分析。分为四种情况：id 是主键索引，id 是唯一辅助索引，id 是普通辅助索引，id 上没有索引。加锁的情况如下图所示，数据为灰色背景的表示加上了 X 锁。

![读隔离级别下的加锁](./image/mysql_lock_rc.png)

（1）隔离级别为读提交，id 是主键索引。只在主键索引上满足 where 条件的所有行加 X 锁。

（2）隔离级别为读提交，id 是唯一辅助索引。在辅助索引上满足 where 条件的所有行加锁。并且这些行对应的主键索引上也要加 X 锁。

（3）隔离级别为读提交，id 是普通辅助索引。加锁的情况的与情况（2）相同，只是普通索引可能有多行数据满足等值查询，这些数据都会被加上锁。

（4）隔离级别为读提交，id 上没有索引。id 上没有索引，存储引擎只能走全表扫表，因此会在主键索引上将所有行都加上 X 锁。可想而知，加锁会产生极大的开销。


### 可重复读隔离级别下的加锁分析

在可重复读的隔离级别下，加锁更加复杂，因为在可重复的隔离级别下要防止幻读的产生，所以除了加 X 锁还会加间隙锁。而间隙锁的加锁规则还会因为等值查询和范围查询而有所不同。加锁的规则[2]我是直接参考丁奇老师的结论，这里我直接引用过来：

> 可重复读隔离级别下的加锁规则，包含了两个“原则”、两个“优化”和一个“bug”。
>
> 原则 1：加锁的基本单位是 next-key lock。希望你还记得，next-key lock 是前开后闭区间。
>
> 原则 2：查找过程中访问到的对象才会加锁。
>
> 优化 1：索引上的等值查询，给唯一索引加锁的时候，next-key lock 退化为行锁。
>
> 优化 2：索引上的等值查询，向右遍历时且最后一个值不满足等值条件的时候，next-key lock 退化为间隙锁。
>
> 一个 bug：唯一索引上的范围查询会访问到不满足条件的第一个值为止。

在这段规则中有一个 next-key lock 的概念，为了说清楚这个概念，我们先来分析一下普通索引上进行等值查询时会怎样加锁。

假设表结构为:

```
create table t (id int, user varchar(10), primary key(id), key key_u(user))engine=innodb;
insert into t (id, user) values (1, 'a'),(5, 'b'),(10, 'c'),(15, 'd'),(20, 'e');
```
事务的隔离级别为可重复读，id 是主键索引，字段 user 上是普通索引，那么对于语句 delete from t where user = 'c' 会怎么加锁？

1. 通过索引 key_u 找到第一条 user 字段为 c 的数据行，也就是 id 为10的行；

2. 在这行数据上加上 X 锁；

3. 为了避免幻读的出现，还要在 (5, 'b') 与 (10, 'c') 之间加上 gap 锁，这样在 id 为 (5, 10) 的范围内就不能插入数据了。这时候 (5, 10) 区间的间隙锁和 id = 10 的 X 锁合起来就是 next-key lock，我们用区间 (5, 10] 来表示；

4. 索引 key_u 不是唯一索引，因此还要继续往后找，找到了 id = 15 的行。根据原则1，会加上 (10, 15] 的 next-key lock，根据优化2，遍历时会退化为间隙锁，也就是 (10, 15)；

5. 因此最终加锁的范围就是 (5, 10], (10, 15)；

讨论完 next-key lock 的概念后，我们再根据可重复读隔离级别下的加锁规则分别探讨索引与查询的所有组合情况的加锁状态：

#### 1. 隔离级别为可重复读，id 是主键索引，等值查询；

加锁过程：通过主键索引进行等值查询找到 id 为10的行，加 (5, 10] 的 next-key lock，根据优化1会退化为行锁，因此最后只锁了 id 为10的行。加锁的结果如下图所示：

![](./image/mysql_lock_rr_1.png)

#### 2. 隔离级别为可重复读，id 是主键索引，范围查询；

加锁过程：
(1) 首先要找到第一个大于5的行，因此找到了10，加锁 (5, 10]；

(2) 然后向后遍历，因为范围是小于12，所以要找到第一个大于等于12的行，加锁 (10, 15]。这里要注意，这个地方的 next-key lock 并没有退化为间隙锁，因为这时不是等值查询，所以没有退化。

最终加锁的结果如下图所示：

![](./image/mysql_lock_rr_2.png)

ps：图中用灰色背景填充的行表示加 X 锁，用红色 V 字在两个行之间表示加间隙锁，文章后面若无单独说明，均表示此含义。

#### 3. 隔离级别为可重复读，id 是唯一索引，等值查询；

加锁过程：通过唯一辅助索引进行等值查询找到 id 为10的行，加 (5, 10] 的 next-key lock，根据优化1会退化为行锁，并且要将主键索引上对应的行加上 X 锁。

![](./image/mysql_lock_rr_3.png)

#### 4. 隔离级别为可重复读，id 是唯一索引，范围查询；

加锁过程：
(1) 首先要找到第一个大于5的行，因此找到了10，在辅助索引上加锁 (5, 10]，并且将主键索引上 id 为10的行加上行锁；

(2) 然后向后遍历，与情况2相似，在辅助索引上加锁 (10, 15]，并且将主键索引上的对应的 id 为15的行加上行锁；

最终加锁的结果如下图所示：

![](./image/mysql_lock_rr_4.png)

#### 5. 隔离级别为可重复读，id 是普通索引，等值查询；

加锁过程：
(1) 首先等值查询找第一个 id 为10的行，加 (5, 10] 的 next-key lock，并且将主键索引上 id 为10的行加上行锁；

(2) 因为不是唯一索引，所以还要向后找，加 (10, 15] 的 next-key lock，根据优化2这里退化为间隙锁，也就是加 (10, 15)；

最终加锁的结果如下图所示：

![](./image/mysql_lock_rr_5.png)

#### 6. 隔离级别为可重复读，id 是普通索引，范围查询；

这种情况的加锁状态与情况4相同，不在赘述，加锁结果如下：

![](./image/mysql_lock_rr_6.png)

#### 7. 隔离级别为可重复读，id 无索引，等值查询；

当 id 上无索引时就要在主键索引上，给所有的行加上行锁，给所有的间隙加上间隙锁。想象一下，如果你的数据表里100万条数据，这该是多么恐怖的场景。加锁的结果如下图所示：

![](./image/mysql_lock_rr_7.png)

#### 8. 隔离级别为可重复读，id 无索引，范围查询；

对于无索引的范围查询，加锁结果与等值查询一样，直接上图：

![](./image/mysql_lock_rr_8.png)

可重复读隔离级别下的加锁分析基本分析完了，那么再来看个更复杂的语句会怎么加锁？

```
select * from t where id >= 5 and id < 12 order by id desc for update;
```
加锁的范围是：(-∞, 1], (1, 5], (5, 10], (10, 15)

加锁过程：
(1) 因为是反向排序，要查询的范围小于12，因此就要先找打第一个大于等于12的位置，就找到了 id 为15的行，加锁 (10, 15]，这里是等值查询，根据优化2退化为间隙锁 (10, 15)。这里为什么是等值查询？等值查询就是指通过索引树去搜索的过程就是这里的等值查询，而遍历查询是指通过 B+ 树叶子节点之间的指针进行遍历的查询方式。

(2) 然后向 id 更小的方向进行遍历。遍历到 id 为10的行，加锁 (5, 10]。

(3) 继续遍历到 id 为5的行，满足where 条件，加锁 (1, 5]。

(4) 然而遍历并没有结束，根据一个 bug 的原则，还要继续向下找到第一个不满足 where 条件的行，那么就也是 id 为1的行，然后加锁 (-∞, 1]。

再举例的过程中我们用 select * from t for update 语句来举例的，但是这些加锁的分析思路对当前读的语句都是一样，哪些语句是当前读前面已经介绍过了，不记得的往前面翻一翻。注意单纯的 select * from t 是快照读，不会加锁。

## 参考

[1] 数据库系统概念（第6版）

[2] MySQL实战45讲，林晓斌

[3] 高性能MySQL（第3版）

[4] [事务的隔离级别和mysql事务隔离级别修改](https://www.cnblogs.com/549294286/p/5433318.html)

[5] [MySQL 加锁处理分析, 何登成](https://github.com/hedengcheng/tech/blob/master/database/MySQL/MySQL%20%E5%8A%A0%E9%94%81%E5%A4%84%E7%90%86%E5%88%86%E6%9E%90.pdf)

[6] [乐观锁、悲观锁，这一篇就够了！](https://segmentfault.com/a/1190000016611415)