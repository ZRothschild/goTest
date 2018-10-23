package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main()  {
	db, _ = sql.Open("mysql",  "root:my06lucky12dog.@tcp(127.0.0.1:3306)/lv")

	//Tx,_ := db.Begin()

	//mysql,err := db.Conn(context.Background())
	//if err != nil {
	//	fmt.Printf("db.conn error %s",err.Error())
	//}
	//mysql.
	defer db.Close()
}