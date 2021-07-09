# 1. elasticsearch 运行

1. http://127.0.0.1:9200/ 查看是否运行中
2. /usr/local/soft/es/bin/elasticsearch
3. ./bin/elasticsearch -d 后台运行

> elasticsearch.yml

```yaml
node.name: node-1

network.host: 0.0.0.0

cluster.initial_master_nodes: [ "node-1" ]
```

# 2. kibana

1. /bin/kibana -d

> ps -ef|grep kibana
> ps -ef|grep 5601
> ps -ef|grep node 或 netstat -anltp|grep 5601
> nohup ../bin/kibana &

### 3. apm 安装

1. https://www.elastic.co/cn/downloads/apm 网站
2. 下载 https://artifacts.elastic.co/downloads/apm-server/apm-server-7.6.2-linux-x86_64.tar.gz 自己选版本 千万要对应上与kibana
   elasticsearch 不然就是运行不了
3. 配置一下然后就 nohup ./apm-server > nohup.log 2>&1 & 看好可执行文件路径与nohup.log是否创建

### 4. 安装elasticsearch

1. 安装地址  https://www.elastic.co/cn/downloads/elasticsearch
2. 安装命令

```shell
[zr@master-00 soft]$ sudo wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-7.13.2-linux-x86_64.tar.gz
[zr@master-00 soft]$ sudo tar -zxvf elasticsearch-7.13.2-linux-x86_64.tar.gz

```

### es 不能使用root用户运行 所有需要把es解压文件设置成非root 用户权限

> /usr/local/soft/elasticsearch-7.13.2/config/elasticsearch.yml

```shell
[zr@node-00 bin]$ sudo chown -R zr:zr  elasticsearch-7.13.2

```

### 创建 /usr/lib/systemd/system/es.service

```shell
[Unit]
Description=elasticsearch
[Service]
User=zr
LimitNOFILE=100000
LimitNPROC=100000
ExecStart=/usr/local/soft/elasticsearch-7.13.2/bin/elasticsearch
Restart=on-failure
[Install]
WantedBy=multi-user.target

```

### 设置服务启动 开机启动

```shell
[zr@node-00 system]$ sudo systemctl daemon-reload
[zr@node-00 system]$ sudo systemctl start es
[zr@node-00 system]$ sudo systemctl enable es
```

### 查看是否已经运行

```shell
[zr@node-00 ~]$ curl http://localhost:9200/
```

#### 下载地址  https://www.elastic.co/cn/start

```shell
[zr@node-00 bin]$ sudo wget https://artifacts.elastic.co/downloads/kibana/kibana-7.13.2-linux-x86_64.tar.gz
[zr@node-00 bin]$ sudo tar -zxvf kibana-7.13.2-linux-x86_64.tar.gz
[zr@node-00 bin]$ chmod zr:zr kibana-7.13.2-linux-x86_64
```

#### 创建 /usr/lib/systemd/system/kibana.service

> 配置路径 /usr/local/soft/kibana-7.13.2-linux-x86_64/config/kibana.yml

```shell
[Unit]
Description=kibana
[Service]
User=zr
LimitNOFILE=100000
LimitNPROC=100000
ExecStart=/usr/local/soft/kibana-7.13.2-linux-x86_64/bin/kibana
Restart=on-failure
[Install]
WantedBy=multi-user.target
```

```shell
[zr@node-00 system]$ sudo systemctl daemon-reload
[zr@node-00 system]$ sudo systemctl start kibana
[zr@node-00 system]$ sudo systemctl enable kibana
```

### 查看 http://localhost:5601

1. 如果要设置其他主机可访问需要添加 elasticsearch-7.13.2 如果要设置成0.0.0.0 xpack.security.enabled: true

```yaml
network.host: 0.0.0.0

xpack.security.enabled: true
# 节点设置为单个
discovery.type: single-node

```

2. 为默认用户设置密码 `elastic` 超级用户 `kibana_system` `logstash_system` `beats_system` 密码我都设置为 `nm123456`

```shell

[zr@node-00 elasticsearch-7.13.2]$ /usr/local/soft/elasticsearch-7.13.2/bin/elasticsearch-setup-passwords interactive

```

3. 如果使用postman请求就需要在`Auth`导航栏`TYPE`选择`Basic Auth`然后填入`Username: elastic` `Password: nm123456`

4. `elastic`设置了密码之后需要修改`kibana`配置，开启用户与用户密码

```yaml
elasticsearch.username: "kibana_system"
elasticsearch.password: "nm123456"
```

### [logstash安装](https://www.elastic.co/cn/downloads/logstash)

```shell

[zr@node-00 soft]$ sudo wget https://artifacts.elastic.co/downloads/logstash/logstash-7.13.3-linux-x86_64.tar.gz

[zr@node-00 soft]$ sudo tar -zxvf logstash-7.13.3-linux-x86_64.tar.gz 

```

1. [logstash使用文档](https://www.elastic.co/guide/en/logstash/current/installing-logstash.html)

简单使用

```shell
[zr@node-00 soft]$ /usr/local/soft/logstash-7.13.3/bin/logstash -e 'input { stdin { } } output { stdout {} }'
```

执行之后，就可以在交互窗口输入数据，然后logstash捕获，输入 hello world

`logstash`对应`filebeat`的配置

```shell

input {
    beats {
        port => "5044"
    }
}
# The filter part of this file is commented out to indicate that it is
# optional.
# filter {
#
# }
output {
    stdout { codec => rubydebug }
}

```

运行命令

```shell
bin/logstash -f first-pipeline.conf --config.test_and_exit # 测试配置是否可用

bin/logstash -f first-pipeline.conf --config.reload.automatic # 运行配置文件

```

### 下载更为轻量级 `filebeat`  [文档](https://www.elastic.co/guide/en/logstash/current/advanced-pipeline.html)

> 注意filebeat需要用root用户执行
> 删除

```shell

[zr@node-00 soft]$ sudo curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-7.13.3-linux-x86_64.tar.gz

[zr@node-00 soft]$ sudo tar -zxvf filebeat-7.13.3-linux-x86_64.tar.gz 
```

filebeat 修改配置文件

```yaml

filebeat.inputs:
  - type: log
    paths:
      - /path/to/file/logstash-tutorial.log
output.logstash:
  hosts: [ "localhost:5044" ]

```

下载测试数据地址 [https://download.elastic.co/demos/logstash/gettingstarted/logstash-tutorial.log.gz](https://download.elastic.co/demos/logstash/gettingstarted/logstash-tutorial.log.gz)
解压

```shell

[zr@node-00 soft]$ gzip -d  logstash-tutorial.log.gz

```

配置文件运行

```shell

sudo rm data/registry # 删除之前以读取到的位置
sudo ./filebeat -e -c filebeat.yml -d "publish"

```
