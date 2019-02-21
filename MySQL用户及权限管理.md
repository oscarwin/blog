

在MySQL中用户数据放在mysql库的user表里，在该表中保存用用户的用户名，密码和权限信息。

## 创建用户

命令：`CREATE USER 'username'@'host' IDENTIFIED BY 'password'`;

说明：
- username：登录的用户名
- password：是登录的密码
- host：指定可以登录的主机，其中localhost表示本机，%表示所有主机

举例：

    CREATE USER 'testuser'@'%' IDENTIFIED BY '123';
    
创建一个用户后，user表中会插入一行新的数据，但是该用户是没有任何权限的

## 管理用户权限
**添加用户权限**

命令：`GRANT privileges ON databasename.tablename TO 'username'@'host'`
说明：

- privileges：用户的操作权限,如SELECT , INSERT , UPDATE，如果要授予所有权限用all privileges
- databasename：数据库名，指定要授予权限的库，如果要对所有库授予权限用`*`代替
- tablename：表名，指定要授予权限的表，如果要对所有表授予权限用`*`代替
- username：被授予权限的用户名
- host：被授予权限的主机，如果要对所有主机授予权限用%代替

举例：

    # 对用户testuser授予在所有主机上对test库的SELECT权限
    GRANT SELECT ON test.* TO 'testuser'@'%';
    
    # 对所有库都授予INSERT权限后，user表中的该用户Insert_priv的值被置为Y
    GRANT INSERT ON *.* TO 'testuser'@'%';
    
    # 给该用户赋予所有权限
    GRANT ALL privileges ON *.* TO 'testuser'@'%';
    
    # 执行完授予权限的命令后，必须执行以下命令使修改生效
    flush privileges;

**撤销用户权限**

命令：`REVOKE privileges ON databasename.tablename FROM 'username'@'host';`

例子：`REVOKE ALL privileges ON *.* FROM 'testuser'@'%';`

撤销权限的使用方法与授予权限类似，参照授予的方式进行处理即可。

## 修改用户密码

命令: `SET PASSWORD FOR 'username'@'host' = PASSWORD('newpassword');`

例子: `SET PASSWORD FOR 'testuser'@'%' = PASSWORD('abcdef');`

## 查看权限

查看本用户的权限：`SHOW GRANTS;`

查看指定用户的权限：`SHOW GRANTS FOR 'username'@'host';`

举例：

    # 查询用户testuser在所有主机上的权限
    SHOW GRANTS FOR 'testuser'@'%'

## 删除用户

删除用户：`DROP USER 'username'@'host';`

举例：`DROP USER 'testuser'@'%';`

用户被删除后user表中就不存在该用户的数据。但是有一点需要注意的是，user表中是通过用户名和主机名来唯一确定一个用户，如果运行`CREATE USER 'testuser'@'10.1.1.1' IDENTIFIED BY 'password';`和`CREATE USER 'testuser'@'10.1.1.2' IDENTIFIED BY 'password';`这两条语句，会在user表中插入两条信息。因此，删除用户的时候也只能分两次删除。
    

## 总结

本篇文章介绍了MySQL管理用户常用的一些命令，包括创建用户，添加用户权限，删除用户权限，修改用户密码，查询用户权限。另外，user表也是一种数据表，因此也可以用MySQL操作表的SQL语句去操作该表，比如删除用户可以使用`DELETE FORM USER WHERE User = testuser`。当然，这种方式是不被推荐的，存在更大的误操作风险。

时间：2019年01月08日