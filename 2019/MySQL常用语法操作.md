数据库的操作可以分为两大类，表数据变更(DML)和表结构变更(DDL)。表数据变更就是增、删、改、查。表结构变更就是创建数据表，修改表结构等。接下来主要从这两方面来总结MySQL常用的操作语法。MySQL里对语法的关键词不敏感，因此本文中存在大小写混用的情况。

### 数据定义语言 DDL

#### 定义数据库
**创建数据库**
  `create database <数据库名称>`

**删除数据库**
  `drop database <数据库名称>`

**列出所有管理的数据库**
`show databases`

**切换到指定的数据库**
 `use <数据库名称>`

**列出某数据库中的所有表**
`show tables`

#### 定义数据表

**创建数据表**

`create table table_name (column_name column_type)`
```
CREATE TABLE IF NOT EXISTS `runoob_tbl`(
   `runoob_id` INT UNSIGNED AUTO_INCREMENT,
   `runoob_title` VARCHAR(100) NOT NULL DEFAULT '',
   `runoob_author` VARCHAR(40) NOT NULL DEFAULT '',
   `submission_date` DATE,
   PRIMARY KEY ( `runoob_id` )，
   KEY ID_TITLE (`runoob_title`, `runoob_author`) 
   )ENGINE=InnoDB CHARSET=utf8;

1. 如果你不想字段为 NULL 可以设置字段的属性为 NOT NULL， 在操作数据库时如果输入该字段的数据为NULL ，就会报错；
2. AUTO_INCREMENT定义列为自增的属性，一般用于主键，数值会自动加1；
3. PRIMARY KEY关键字用于定义列为主键。 您可以使用多列来定义主键，列间以逗号分隔；
4. ENGINE 设置存储引擎；
5. CHARSET 设置编码, 常用的编码格式有utf8、utf8mb4、gbk；
6. UNIQUE 唯一索引
```

**删除数据表**

 `drop table <数据表名称>`

**查询建表操作**

`SHOW CREATE TABLE TABLE_NAME`

创建一个表t1，然后在其基础上进行表结构的变更，表t1如下：`CREATE TABLE t1 (a INTEGER, b CHAR(10));`
**修改表结构**
	
	# 添加字段c,d和e(ADD)
	ALTER TABLE t1 ADD c bigint(20) COMMENT 'xxx',
	ADD d bigint(20) NOT NULL COMMENT 'xxx',
	ADD e bigint(20) NOT NULL COMMENT 'xxx';

	# 将表名重命名为t2(RENAME) 
	ALTER TABLE t1 RENAME t2
	
	# 删除字段b(DROP)
	ALTER TABLE t2 DROP COLUMN b;

	# 修改字段b类型为VARCHAR(MODIFY)
	ALTER TABLE t2 MODIFY COLUMN b VARCHAR(255) NOT NULL DEFAULT '';
	
	# 修改字段名a为ach(CHANGE)
	ALTER TABLE t2 CHANGE a ach BIGINT(20) NOT NULL;

	# 添加一个名为key_a的索引
	ALTER TABLE t2 ADD INDEX key_a(a);

    # 删除名为key_a的索引
    ALTER TABLE t2 DROP INDEX key_a
	
	
### 数据操纵语言 DML

**查询数据**
```
SELECT column_name,column_name
FROM table_name
[WHERE Clause]
[LIMIT N][ OFFSET M]

1. 查询语句中你可以使用一个或者多个表，表之间使用逗号(,)分割，并使用WHERE语句来设定查询条件。
2. SELECT 命令可以读取一条或者多条记录。
3. 你可以使用星号（*）来代替其他字段，SELECT语句会返回表的所有字段数据
4. 你可以使用 WHERE 语句来包含任何条件。
5. 你可以使用 LIMIT 属性来设定返回的记录数。
6. 你可以通过OFFSET指定SELECT语句开始查询的数据偏移量。默认情况下偏移量为0。
```

```
SELECT * from t where id = 1 LIMIT 10;
```

**LIKE子句**

	SELECT field1, field2,...fieldN  FROM table_name
	WHERE field1 LIKE condition1 [AND [OR]] filed2 = 'somevalue'

1. LIKE子句相当于模糊搜索，用`%`表示任意字符

**插入数据**
`INSERT INTO t1 ( field1, field2,...fieldN )  VALUES  ( value1, value2,...valueN )`

插入的数据主键或者UNIQUE键冲突时，会执行后面的UPDATE命令，否者执行INSERT命令
`INSERT INTO t1 (a,b,c) VALUES (1,2,3) ON DUPLICATE KEY UPDATE c=c+1;`

**UPDATE修改数据**

`UPDATE table_name SET field1=new-value1, field2=new-value2 [WHERE Clause]`

**DELETE删除数据**

`DELETE FROM table_name [WHERE Clause] [LIMIT 20]`

1. 如果没有where子句，会删除整个表；
2. where子句用来限定删除指定的表。


## 其他

#### 查询表信息
`show columns from <table name>`
`show index from <table name>`	

#### 事务控制
**开启事务**
开启事务有两种方式
1. begin 常用的开启事务的方式，实际上BEGIN后事务没有立即开启，而是等到第一条语句执行后才开启；
2. start transaction与begin相同；
3. start transaction with consistent snapshot则是运行该语句后就立刻开启了事务，不用等到第一条语句的执行；

**提交事务**
commit

**回滚事务**
rollback
