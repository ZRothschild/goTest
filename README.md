1. mongod -f /etc/mongodb.conf --fork
2. 视频转换 https://online-video-cutter.com/cn/


## redis

### redis 使用场景

1. 用户排名 积分，金额
2. 计数器相关问题 短信发送次数，用户试错次数，请求次数
3. 队列
4. 热点数据的缓存
5. 限时业务的运用，过期无效
6. 分布式锁

### redis 缓存使用

> 如果交互字段不是主键索引 需要使用转换 主键索引最好是数值类型
> 需要查询的字段 维护与主键的对应关系 eg. nickname phone
> redis 可以存储单条，然后使用 es 如果是大数据简单
> 获取数据 Multi Get 获取或根据主键查数据，有可能缓存有遗漏
> 可以运行一个额外程序去维护，mysql 与 redis

1. 缓存键名 前缀+表名+缓存类型+时间戳
2. 数据分两种种类型
   (1) 普通类型使用hash数据结构存储 缓存键名 + 主键

```shell

HSET key field value [field value ...]

hset csdn id 1 name wyk company CSDN

hset csdn gender male

```

(2) 有分页需求的数据Sorted Set 分数为排序字段值数值类型 member为主键，查询字段映射主键，往普通数据查，或用主键覆盖索引

```shell

ZADD key score1 member1 [score2 member2] 

ZRANGE w3ckey 0 10 WITHSCORES

ZINCRBY key increment member // 有序集合中对指定成员的分数加上增量 increment

ZRANGE w3ckey 0 10 WITHSCORES

ZCOUNT key min max 计算在有序集合中指定区间分数的成员数

```

### redis 秒杀场景使用

### 缓存更新策略 被动与主动更新

#### 提供主动接口

#### 被动删除

1. 延迟双删策略，为了把读然后缓存的不一致数据删除，因为读数据库是旧数据，写缓存出现延迟，把旧数据写入缓存，所以要延迟删除。可能第二次删除失败，需要从试机制
   这么做就可以在一定程度上缓解上述两个问题，第一次删除缓存相当于检测下缓存服务是否可用，网络是否有问题，第二次延迟一定时间，再次删除缓存，是因为要保证读的请求在写的请求之前完成。
2. 解耦删除；

### docker run --name grpc-micro -d


## 循环提交问题
1. 用户多次并发提交