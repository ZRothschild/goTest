package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Db struct {
	DbName    string
	TableName string
}

type Curd interface {
	GetDbName() string
	TableName() string
	Find()
}

type User struct {
	Id   int64
	Name sql.NullString
	Age  int
	Sex  int64
}

type Result struct {
	Name string
	Age  int
}

// Go连接MySQL示例
var db *sql.DB // db是一个连接池对象，每次操作会从连接池中取一个db对象来服务

func main() {
	// db, _ := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	//
	// defer db.Close()
	// var result Result
	// db.Table("users").Select("name, age").Where("name = ?", "Antonio").Scan(&result)

	err := initDB()
	if err != nil {
		fmt.Printf("init DB failed, err:%v\n", err)
	}
	fmt.Println("连接数据库成功!")
	queryOne(1)
	// queryMore(2)
	// insert()
	// updateRow(9000, 2)
	// deleteRow(2)
	// prepareInsert()
}

func initDB() (err error) {
	// 数据库信息
	// 用户名:密码@tcp(ip:端口)/数据库的名字
	dsn := "root:@tcp(127.0.0.1:3306)/test"
	// 连接数据库
	db, err = sql.Open("mysql", dsn) // 不会校验用户名和密码是否正确
	if err != nil {                  // dsn格式不正确的时候会报错
		return
	}
	err = db.Ping() // 尝试连接数据库
	if err != nil {
		return
	}
	db.SetMaxOpenConns(10) // 设置数据库连接池的最大连接数（连接池中只有10个连接，如果全部被占用，再来请求就会阻塞住）
	db.SetMaxIdleConns(5)  // 设置最大空闲连接数
	return
}

// 查询单个记录
func queryOne(id int) {
	var u1 User // 用来接收查询结果
	// 1. 写查询单条记录的sql语句
	sqlStr := `select id, name, age,sex from users where sex = ?;` // ？占位 下面的id
	// 2. 执行并拿到结果
	// 必须对rowObj对象调用Scan方法,因为该方法会释放数据库链接 // 从连接池里拿一个连接出来去数据库查询单条记录
	db.QueryRow(sqlStr, 0).Scan(&u1.Id, &u1.Name, &u1.Age, &u1.Sex) // &u1.id, &u1.name, &u1.age初始化u1结构体对象（变量）
	// row一行

	// 打印结果
	fmt.Printf("u1:%#v\n", u1)
}

func queryMore(n int) {
	// 1. SQL语句
	sqlStr := `select id, name, age from users where id > ?;`
	// 2. 执行
	rows, err := db.Query(sqlStr, n)
	if err != nil {
		fmt.Printf("exec %s query failed, err:%v\n", sqlStr, err)
		return
	}
	// 3. 一定要关闭rows，才会释放连接（数据库的连接）
	defer rows.Close()
	// 4. 循环取值
	for rows.Next() {
		var u1 User
		err := rows.Scan(&u1.Id, &u1.Name, &u1.Age)
		if err != nil {
			fmt.Printf("scan failed,err:%v\n", err)
		}
		fmt.Printf("u1:%#v\n", u1)
	}
}

// GOPROXY=https://goproxy.cn,direct

// 插入数据
func insert() {
	// 1. 写SQL语句
	sqlStr := `insert into users(name, age) values("图朝阳", 28)`
	// 2. exec
	ret, err := db.Exec(sqlStr) // exec执行（Python中的exec就是执行字符串代码的，返回值是None，eval有返回值）
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	// 如果是插入数据的操作，能够拿到插入数据的id
	id, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get id failed,err:%v\n", err)
		return

	}
	fmt.Println("id:", id)
}

// 更新操作
func updateRow(newAge int, id int) {
	sqlStr := `update users set age=? where id > ?`
	ret, err := db.Exec(sqlStr, newAge, id)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // RowsAffected 受影响的行数
	if err != nil {
		fmt.Printf("get id failed,err:%v\n", err)
		return
	}
	fmt.Printf("更新了%d行数据\n", n)
}

// 删除
func deleteRow(id int) {
	sqlStr := `delete from users where id=?`
	ret, err := db.Exec(sqlStr, id)
	if err != nil {
		fmt.Printf("delete failed,err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get id failed,err:%v\n", err)
		return
	}
	fmt.Printf("删除了%d行数据\n", n)
}

// 预处理方式插入多条数据
func prepareInsert() {
	sqlStr := `insert into user(name, age) values(?,?)`
	stmt, err := db.Prepare(sqlStr) // 1.把SQL语句先发给MySQL预处理一下
	if err != nil {
		fmt.Printf("prepare failed ,err:%v\n", err)
		return
	}
	defer stmt.Close()
	// 后续只需要拿到stmt去执行一些操作
	var m = map[string]int{
		"六七强": 30,
		"王相机": 32,
		"天说":  72,
		"白慧姐": 40,
	}
	for k, v := range m {
		stmt.Exec(k, v) // 2.后续只需要传值
	}
}
