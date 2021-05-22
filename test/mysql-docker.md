Docker 启动mysql
docker run -p 3307:3306  --restart=always  --privileged=true --name mysql -v D:\dockerData\mysql\data:/var/lib/mysql -v D:\dockerData\mysql\my.cnf/my.cnf:/etc/mysql/my.cnf -e MYSQL_ROOT_PASSWORD="123456" -d mysql:5.7

命令说明：
-p 3306:3306：将容器的3306端口映射到主机的3306端口
-e MYSQL_ROOT_PASSWORD=123456：初始化root用户的密码
-d: 后台运行容器，并返回容器ID

下面这两行配置是挂载数据

-v /opt/mysql/data:/var/lib/mysql

挂载mysql配置文件
-v /opt/mysql/my.cnf:/etc/mysql/my.cnf

my.cnf配置文件添加以下内容：

[mysqld]
skip-name-resolve
user=root
character-set-server=utf8
default_authentication_plugin=mysql_native_password
sql_mode=ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION
lower_case_table_names=1 #忽略表名大小写
skip-grant-tables
[client]
default-character-set=utf8
[mysql]
default-character-set=utf8

mysql更改密码：

use mysql;

下面我们就要修改密码了

以前的版本我们用的是以下修改

update user set password=password('root') where user='root';

但是在5.7版本中不存在password字段，所有我们要用以下修改进行重置密码

update mysql.user set authentication_string=password('root') where user='root';
update user set authentication_string=password('root') where user='root';
重新授权

flush privileges;

grant all privileges on *.* to root@'%' identified by  'root';
退出mysql
quit;