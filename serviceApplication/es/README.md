# 1. elasticsearch 运行

1. http://127.0.0.1:9200/ 查看是否运行中
2. /usr/local/soft/es/bin/elasticsearch
3. ./bin/elasticsearch -d 后台运行

> elasticsearch.yml
```yaml
node.name: node-1

network.host: 0.0.0.0

cluster.initial_master_nodes: ["node-1"]
```

# 2. kibana 

1.  /bin/kibana -d
> ps -ef|grep kibana
> ps -ef|grep 5601
> ps -ef|grep node 或 netstat -anltp|grep 5601
> nohup ../bin/kibana &

# 3. apm 安装
1. https://www.elastic.co/cn/downloads/apm 网站
2. 下载 https://artifacts.elastic.co/downloads/apm-server/apm-server-7.6.2-linux-x86_64.tar.gz 自己选版本 千万要对应上与kibana elasticsearch 不然就是运行不了
3. 配置一下然后就  nohup ./apm-server > nohup.log 2>&1 & 看好可执行文件路径与nohup.log是否创建
