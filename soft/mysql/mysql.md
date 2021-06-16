## docker 安装mysql
```shell
Docker 启动mysql
# 主
docker run -p 3307:3306  --restart=always  --privileged=true --name mysql -v D:\dockerData\mysql\data:/var/lib/mysql -v D:\dockerData\mysql\my.cnf/my.cnf:/etc/mysql/my.cnf -e MYSQL_ROOT_PASSWORD="123456" -d mysql:5.7
# 从
docker run -p 3308:3306  --restart=always  --privileged=true --name mysql2 -v D:\dockerData\mysql2\data:/var/lib/mysql -v D:\dockerData\mysql2\conf\my.cnf:/etc/mysql/my.cnf -e MYSQL_ROOT_PASSWORD="123456" -d mysql:5.7

````
命令说明：
-p 3306:3306：将容器的3306端口映射到主机的3306端口
-e MYSQL_ROOT_PASSWORD=123456：初始化root用户的密码
-d: 后台运行容器，并返回容器ID

下面这两行配置是挂载数据

-v /opt/mysql/data:/var/lib/mysql

挂载mysql配置文件
-v /opt/mysql/my.cnf:/etc/mysql/my.cnf

my.cnf配置文件添加以下内容：

```shell
[mysqld]
server-id = 1        # 节点ID，确保唯一
skip-name-resolve
user=root
character-set-server=utf8
default_authentication_plugin=mysql_native_password
sql_mode=ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
lower_case_table_names=1 #忽略表名大小写
# skip-grant-tables

# log config
log-bin = mysql-bin     #开启mysql的binlog日志功能
sync_binlog = 1         #控制数据库的binlog刷到磁盘上去 , 0 不控制，性能最好，1每次事物提交都会刷到日志文件中，性能最差，最安全
binlog_format = mixed   #binlog日志格式，mysql默认采用statement，建议使用mixed
expire_logs_days = 7                           #binlog过期清理时间
max_binlog_size = 100m                    #binlog每个日志文件大小
binlog_cache_size = 4m                    #binlog缓存大小
max_binlog_cache_size= 512m              #最大binlog缓存大
binlog-ignore-db=mysql              #不生成日志文件的数据库，多个忽略数据库可以用逗号拼接，或者 复制这句话，写多行

auto-increment-offset = 1     # 自增值的偏移量
auto-increment-increment = 1  # 自增值的自增量
slave-skip-errors = all #跳过从库错误

[client]
default-character-set=utf8
[mysql]
default-character-set=utf8
```

mysql更改密码：

```shell

use mysql;

```

下面我们就要修改密码了,以前的版本我们用的是以下修改

```shell

update user set password=password('root') where user='root';

```

但是在5.7版本中不存在password字段，所有我们要用以下修改进行重置密码

```shell
update mysql.user set authentication_string=password('root') where user='root';
update user set authentication_string=password('root') where user='root';
```

重新授权

```shell

mysql> flush privileges;


mysql> grant all privileges on *.* to root@'%' identified by  'root';

# 退出mysql

mysql> quit;
```


## mysql 主从

###  主数据库执行

```shell
# 创建主库同步数据的用户，用于从库登录

mysql> grant replication slave on *.* to slave@'%' identified by '123456';
# 刷新权限

mysql> flush privileges;

# 关闭已有的主从

mysql> stop slave;

# 开启主从

mysql> start slave;

# 查看主从状态

mysql> show master status;

+------------------+----------+--------------+------------------+-------------------+
| File             | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
+------------------+----------+--------------+------------------+-------------------+
| mysql-bin.000001 |      590 |              | mysql            |                   |
+------------------+----------+--------------+------------------+-------------------+
1 row in set (0.00 sec)

```

#### 从数据库执行

> 配置需要同步的主数据库地址，登录从库执行

```shell

mysql> CHANGE MASTER TO MASTER_HOST = '192.168.244.106',  
MASTER_USER = 'slave',
MASTER_PASSWORD = '123456',
MASTER_PORT = 3307,
MASTER_LOG_FILE='mysql-bin.000001',
MASTER_LOG_POS=590,
MASTER_RETRY_COUNT = 60,
MASTER_HEARTBEAT_PERIOD = 10000;

```

#### 从库配置文件

```shell

[mysqld]
server-id = 2        # 节点ID，确保唯一
skip-name-resolve
user=root
character-set-server=utf8
default_authentication_plugin=mysql_native_password
sql_mode=ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
lower_case_table_names=1 #忽略表名大小写
# skip-grant-tables


log-bin=mysql-bin
relay-log=relay-log-bin   #添加
relay-log-index=slave-relay-bin.index #添加
replicate-wild-ignore-table=mysql.%
replicate-wild-ignore-table=test.%
replicate-wild-ignore-table=information_schema.%

[client]
default-character-set=utf8
[mysql]
default-character-set=utf8

```

