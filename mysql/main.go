package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *sql.DB

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
	ID   int64
	Name sql.NullString
	Age  int
}

type Result struct {
	Name string
	Age  int
}

func main() {
	db, _ := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")

	defer db.Close()

	var result Result
	db.Table("users").Select("name, age").Where("name = ?", "Antonio").Scan(&result)

	db, _ = sql.Open("mysql", "root:nm123456@tcp(127.0.0.1:3306)/test")
	defer db.Close()
	Stmt, err := db.Prepare("select name,age from users where id = ?")
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	user := new(User)
	err = Stmt.QueryRow(1).Scan(&user.Name, &user.Age)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Println(user.Name)
	fmt.Println(user.Age)
}
