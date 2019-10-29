
```
#!/bin/sh

HOST="IP"
PORT="PORT"
USERNAME="username"
PASSWORD="password"

MYSQL_CMD="mysql -h${HOST} -P${PORT} -u${USERNAME} -p${PASSWORD}"

CREATE_SQL="
        use database_name;
        create table table_name(
                i_id bigint(20) not null AUTO_INCREMENT comment 'primary key id',
                field1 bigint(20) not null default 0 comment 'xxx',
                field2 bigint(20) not null default 0 comment 'xxx',
                field3 int(8) not null default 0 comment 'xxx',
                i_status tinyint(1) not null default 0 comment 'xxx',
                PRIMARY KEY(i_id)
        )ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;"

${MYSQL_CMD} -e "${CREATE_SQL}"

echo "create table table_name success!"
``` 