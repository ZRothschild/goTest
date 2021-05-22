# `go(golang)`语言使用`NoSql redis`数据库

## `send`方式

- `send`只是写在客户端缓存里面而已，并没有写入服务端数据库，需要`c.Flush()`往服务推送，`c.Receive()`是接收服务端操作返回是否成功，信号
```go
package main

import (
	"github.com/gomodule/redigo/redis"
	"fmt"
)

func main()  {
	//连接到redis数据库
	//c, err := redis.Dial("tcp", ":6379")
	//if err != nil {
	//	fmt.Printf("dial err %s\n",err)
	//}
	//defer c.Close()

	c, err := redis.DialURL("redis://localhost:6379")
	if err != nil {
		fmt.Printf("DialURL err %s\n",err)
	}
	defer c.Close()
	//选择1号数据库
	rep,err := c.Do("select",1)
	if err != nil {
		fmt.Printf("c.Do select err %s\n",err)
	}
	fmt.Printf("replay %v\n",rep)
	//set 一个key为go value值为iris的值
	err = c.Send("set","go","iris")
	if err != nil {
		fmt.Printf("set %v\n",rep)
	}
	c.Send("get","go")
	if err != nil {
		fmt.Printf("Send get %v\n",rep)
	}
	c.Flush()
	c.Receive()
	reply,err := redis.String(c.Receive())
	if err != nil {
		fmt.Printf("redis.Values %v\n",err)
	}
	fmt.Printf("get value %v\n",reply)
}
```

## redis 缓存穿透，无底洞，雪崩优化
1. 穿透：频繁查询一个不存在的数据，由于缓存不命中，每次都要查询持久层。从而失去缓存的意义。

> 解决办法：
>  ①用一个bitmap和n个hash函数做布隆过滤器过滤没有在缓存的键。
>  ②持久层查询不到就缓存空结果，有效时间为数分钟。

2. 雪崩：缓存大量失效的时候，引发大量查询数据库。
> 解决办法：
>  ①保证缓存曾服务高可用性。
>  ②依赖隔离组件为后端限流并降级。

3. 无底洞：为了满足业务大量加节点，但是性能没提升反而下降。

热点key:某个key访问非常频繁，当key失效的时候有大量线程来构建缓存，导致负载增加，系统崩溃。

> 解决办法：
>  ①使用锁，单机用synchronized,lock等，分布式用分布式锁。 
>  ②缓存过期时间不设置，而是设置在key对应的value里。如果检测到存的时间超过过期时间则异步更新缓存。
>  ③在value设置一个比过期时间t0小的过期时间值t1，当t1过期的时候，延长t1并做更新缓存操作。

4. 数据就持久化

5. 缓存和数据库双写一致性问题

6.缓存并发竞争

7. 单线程的redis为什么这么快
>  纯内存操作
>  单线程操作，避免了频繁的上下文切换
>  采用了非阻塞I/O多路复用机制

8. redis的数据类型，以及每种数据类型的使用场景

> String 这个其实没啥好说的，最常规的set/get操作，value可以是String也可以是数字。一般做一些复杂的计数功能的缓存。
 
> hash 这里value存放的是结构化的对象，比较方便的就是操作其中的某个字段。就是用这种数据结构存储用户信息，
以cookieId作为key，设置30分钟为缓存过期时间，能很好的模拟出类似session的效果,手机验证码。 

 赋值
 
 hset key  field value : 为指定的key设定field/value对
 hmset key field1 value1 field2 value2  field3 value3     为指定的key设定多个field/value对
 
 取值
 
 hget key field : 返回指定的key中的field的值
 
 hmget key field1 field2 field3 : 获取key中的多个field值
 
 hkeys key : 获取所有的key
 
 hvals key :获取所有的value
 
 hgetall key : 获取key中的所有field 中的所有field-value


> list 使用List的数据结构，可以做简单的消息队列的功能。另外还有一个就是，可以利用lrange命令，做基于redis的分页功能，性能极佳，用户体验好

> set 因为set堆放的是一堆不重复值的集合。所以可以做全局去重的功能。为什么不用JVM自带的Set进行去重？因为我们的系统一般都是集群部署，使用JVM自带的Set，比较麻烦，难道为了一个做一个全局去重，再起一个公共服务，太麻烦了


etcd  https://blog.csdn.net/bbwangj/article/details/82584988

redis https://blog.csdn.net/itcats_cn/article/details/82391719
rabbitMq http://rabbitmq.mr-ping.com/AMQP/AMQP_0-9-1_Model_Explained.html
mysql https://www.jianshu.com/p/19cb0f16dea4

elk https://blog.csdn.net/lvtula/article/details/88285650


mysql binlog日志有三种模式，分别为：

    ROW: 记录每一行数据被修改的情况，但是日志量太大
    STATEMENT: 记录每一条修改数据的SQL语句，减少了日志量，但是SQL语句使用函数或触发器时容易出现主从不一致
    MIXED: 结合了ROW和STATEMENT的优点，根据具体执行数据操作的SQL语句选择使用ROW或者STATEMENT记录日志

mysql es https://www.jianshu.com/p/c3faa26bc221