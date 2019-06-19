package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
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

func main() {
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
