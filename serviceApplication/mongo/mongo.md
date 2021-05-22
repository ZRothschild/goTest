# `go(golang)`语言使用`NoSql mongodb`数据库

## 链接`mongodb`数据库

> 链接格式 `mongodb://[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]`

### 类型有如下几种
> 不填写端口，则为默认端口

连接`mongodb`

```go
package main

import (
	"github.com/globalsign/mgo"
	"fmt"
)

func main()  {
	url := "mongodb://localhost"
	//默认链接类型
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Printf("err %s\n",err)
	}
	defer session.Close()
    //获取mongodb所有库
    dbs,err := session.DatabaseNames()
    if err != nil {
        fmt.Printf("DatabaseNames err %s\n",err)
    }
    for k,v := range dbs{
        fmt.Printf("第%d个库名为%s\n",k,v)
    }
}
```
连接`mongodb`指定用户名与密码

```go
package main

import (
	"github.com/globalsign/mgo"
	"fmt"
)

func main()  {
	url := "mongodb://user:password@localhost:27017"
	//默认链接类型
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Printf("err %s\n",err)
	}
	defer session.Close()
    //获取mongodb所有库
    dbs,err := session.DatabaseNames()
    if err != nil {
        fmt.Printf("DatabaseNames err %s\n",err)
    }
    for k,v := range dbs{
        fmt.Printf("第%d个库名为%s\n",k,v)
    }
}
```

连接`mongodb`指定用户名与密码,选择指定数据库

```go
package main

import (
	"github.com/globalsign/mgo"
	"fmt"
)

func main()  {
	url := "mongodb://user:password@localhost:27017/records"
	//默认链接类型
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Printf("err %s\n",err)
	}
	defer session.Close()
	//查看是否链接
    err = session.Ping()
    if err != nil {
        fmt.Printf("ping err %s\n",err)
    }
}
```
## `mongodb`建库

>  `mongodb`建库有点特殊，数据库下面必须有集合才会成功

### 1.建库与集合并插入数据

```go
package main

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"fmt"
)

func main()  {
	url := "localhost:27017"
	//默认链接类型
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Printf("err %s\n",err)
	}
	defer session.Close()
	//如果不存在就创建go database并在go database中创建test collection
	c := session.DB("go").C("test")
	//插入两条数据
	err = c.Insert(&Language{"go", "iris"},
		&Language{"php", "laravel"})
	if err != nil {
		fmt.Printf("insert err %s \n",err)
	}
	result := Language{}
	//查询数据lang为go的数据
	err = c.Find(bson.M{"lang": "go"}).One(&result)
	if err != nil {
		fmt.Printf("find err %s \n",err)
	}
	fmt.Println("Name:", result.Name)
	fmt.Println("Lang:", result.Lang)
}

type Language struct {
	Lang  string
	Name string
}
```
### 2.建库与集合不插入数据

```go
package main

import (
	"github.com/globalsign/mgo"
	"fmt"
	"github.com/globalsign/mgo/bson"
)

func main()  {
	url := "mongodb://localhost"
	//默认链接类型
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Printf("err %s\n",err)
	}
	defer session.Close()
	db := session.DB("go")
	var result map[string]interface{}
	//如果集合已存在，就会报错
	err = db.Run(bson.D{{"create", "test"}, {"size", 1024}},&result)
	if err != nil {
		fmt.Printf("Run err :%v", err)
		//Run err :a collection 'go.test' already existsresult => :map[ok:false code:true]
	}
	fmt.Printf("result => :%v", result)
	//查看是否建立数据库与集合
	dbs,err := session.DatabaseNames()
	if err != nil {
		fmt.Printf("DatabaseNames err %s\n",err)
	}
	for k,v := range dbs{
		fmt.Printf("第%d个库名为%s\n",k,v)
	}
}
```

## 创建用户

```go
package main

import (
	"github.com/globalsign/mgo"
	"fmt"
	"github.com/globalsign/mgo/bson"
)

func main()  {
	url := "mongodb://localhost"
	//默认链接类型
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Printf("err %s\n",err)
	}
	defer session.Close()
	db := session.DB("go")

	user := mgo.User{
		Username: "zhao",
		Password: "123456",
		Roles: []mgo.Role{mgo.RoleReadWrite},
	}
	var res map[string]interface{}
	//AddUser 以及被弃用
	err = db.UpsertUser(&user)
	if err != nil {
		fmt.Printf("UpsertUser err %v\n",err)
	}
	err = db.Run(bson.D{{"usersInfo", map[string]string{"user":"readWrite","db":"go"}}, {"showPrivileges", true}},&res)
	fmt.Printf("Run err %s\n",err)
	fmt.Printf("res err %v\n",res)
}
```
> 需要验证登录

```go
package main

import (
	"github.com/globalsign/mgo"
	"fmt"
	"github.com/globalsign/mgo/bson"
)

func main()  {
	url := "mongodb://zhao:123456@localhost:27017/go"
	//默认链接类型
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Printf("err %s\n",err)
	}
	defer session.Close()
	//err = session.DB("go").Login("zhao","123456")
	//fmt.Printf("loin %v\n",err)
	db := session.DB("go").C("test")
	var res map[string]interface{}
	err = db.Find(bson.M{"lang": "go"}).One(&res)
	fmt.Printf("res %v\n",res)
	fmt.Printf("err %v\n",err)
}
```