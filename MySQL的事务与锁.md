# MySQL的事务与锁

提到事务首先想到的当然是事务的四个特性：原子性、一致性、隔离性、持久性

## 事务的四个特性

**原子性：** 事务的所有操作在数据库中要么全部正确的反映出来，要么完全不反映。

**一致性：** 事务执行前后数据库的数据保持一致。

**隔离性：** 多个事务并发执行时，对于任何一对事务Ti和Tj，在Ti看来，Tj 要么在 Ti 之前已经完成执行，或者在Ti完成之后开始执行。因此，每个事务都感觉不到系统中有其他事务在并发执行。

**持久性：** 一个事务成功完成后，它对数据库的改变必须是永久的，即使事务刚提交机器就宕机了数据也不能丢。

事务的原子性和持久性比较好理解，但是一致性会更加抽象一些。对于一致性经常有个转账的例子，A 给 B 转账，转账前后 A 和 B 的账户总和不变就是一致的。这个例子咋一看好像很清楚，但转念一想原子性是不是也能达到这个目的呢？答案是：不能，原子性可以保证 A 账户扣减和 B 账户增加同时成功或者同时失败，但是并不能保证 A 扣减的数量等于 B 增加的数量。实际上是为了达到一致性所以要同时满足其他三个条件。

还有一个事务的隔离性比较复杂，因为 MySQL 的事务可以有多种隔离级别，接下里一起看看。

## 事务的隔离级别

当多个事务并发执行时可能存在脏读(dirty read)，不可重复读(non-repeatable read)和幻读(phantom read)，为了解决这些问题因此引入了不同的隔离级别。

**脏读：** 事务 A 和事务 B 并发执行时，事务 B 可以读到事务 A 未提交的数据，就发生了脏读。脏读的本质在于事务 B 读了事务 A 未提交的数据，如果事务 A 发生了回滚，那么事务 B 读到的数据实际上是无效的。如下面案例所示：事务 B 查询到 value 的结果为100，但是因为事务 A 发生了回滚，因此 value 的值不一定是 100。

| 事务 A | 事务 B |
| :--- | :--- |
| begin | begin |
| update t set value = 100 |  |
|       | select value from t |
| rollback |  |
| commit | commit |

**不可重复读：** 在一个事务中，多次查询同一个数据会得到不同的结果，就叫不可重复读。如下面案例所示：事务 B 两次查询 value 的结果不一致。

| 事务 A | 事务 B |
| :---   | :---  |
| begin  | begin |
| update t set value = 100 |  |
|        | select value from t ( value = 100 ) |
| update t set value = 200 | |
|        | select value from t ( value = 200 ) |
| commit | commit |

**幻读：** 在一个事务中进行范围查询，查询到了一定条数的数据，但是这个时候又有新的数据插入就导致数据库中数据多了一行，这就是幻读。如下面案例所示：事务 B 两次查询到的数据行数不一样。

| 事务 A | 事务 B |
| :---   | :---  |
| begin  | begin |
|        | select * from t |
| insert into t ... | |
| commit | |
|        | select * from t |
|  | commit |

MySQL 的事务隔离级别包括：读未提交（read uncommitted）、读提交（read committed）可重复读（repeatable read）和串行化（serializable）。

**未提交读：** 一个事务还未提交，其造成的更新就可以被其他事务看到。这就造成了脏读。

**读提交：** 一个事务提交后，其更改才能被其他事务所看到。读提交解决了脏读的问题。

**可重复读：** 在一个事务中，多次读取同一个数据得到的结果总是相同的，即使有其他事务更新了这个数据并提交成功了。可重复读解决了不可重复读的问题。但是还是会出现幻读。InnoDB 引擎通过多版本并发控制（Multiversion concurrency control，MVCC）解决了幻读的问题。

**串行化：** 串行话是最严格的隔离级别，在事务中对读操作加读锁，对写操作加写锁，所以可能会出现大量锁争用的场景。

从上到下，隔离级别越来越高，效率相应也会随之降低，对于不同的隔离级别需要根据业务场景进行合理选择。

### 查询和修改事务的隔离级别

下面的命令可以查询 InnoDB 引擎全局的隔离级别和当前会话的隔离级别

```
mysql> select @@global.tx_isolation,@@tx_isolation;
+-----------------------+-----------------+
| @@global.tx_isolation | @@tx_isolation  |
+-----------------------+-----------------+
| REPEATABLE-READ       | REPEATABLE-READ |
+-----------------------+-----------------+
```

设置innodb的事务级别方法是：

```SQL
set 作用域 transaction isolation level 事务隔离级别

SET [SESSION | GLOBAL] TRANSACTION ISOLATION LEVEL {READ UNCOMMITTED | READ COMMITTED | REPEATABLE READ | SERIALIZABLE}

mysql> set global transaction isolation level read committed; // 设定全局的隔离级别为读提交

mysql> set session transaction isolation level read committed; // 设定当前会话的隔离级别为读提交
```

### 举例说明不同隔离级别的影响

接下来我们用一个案例来看不同隔离级别下会有怎样不同的结果。

```SQL
create table t (k int) ENGINE=InnoDB;
insert into t values (1);
```
| 事务 A | 事务 B |
| :---   | :---  |
| begin  |  |
| 1: select k from t | |
| | begin; update t set k = k + 1 |
| 2: select k from t | |
| | commit |
| 3: select k from t | |
| commit | |
| 4: select k from t | |

隔离级别为未提交读时：对于事务 A，第1条查询语句的结果是1，第2条查询语句的结果是2，第3条和第4条查询语句的结果也都是2。

隔离级别为读提交时：对于事务 A，第1条查询语句的结果是1，第2条查询语句的结果是1，第3条查询语句的结果是2，第4条查询语句的结果也是2。

隔离级别为可重复读时：对于事务 A，第1条、第2条和第3条查询语句的结果都是1，第4条查询语句的结果是2。

隔离级别为串行化时：对于事务 A，第1条查询语句的结果是1。这时事务 B 执行更新语句时会被阻塞，因为事务 A 在这条数据上加上了读锁，事务 B 要更新这个数据就必须加写锁，由于读锁和写锁冲突，因此事务 B 只能等到事务 A 提交后释放读锁才能进行更新。因此，事务 A 的第2条和第3条查询语句的结果也是1，第4条查询语句的结果是2。

## 事务隔离性的实现

## 参考

[1] 数据库系统概念（第6版）

[2] MySQL实战45讲，林晓斌

[3] 高性能MySQL（第3版）

[4] [事务的隔离级别和mysql事务隔离级别修改](https://www.cnblogs.com/549294286/p/5433318.html)

[5] [MySQL 加锁处理分析, 何登成](https://github.com/hedengcheng/tech/blob/master/database/MySQL/MySQL%20%E5%8A%A0%E9%94%81%E5%A4%84%E7%90%86%E5%88%86%E6%9E%90.pdf)