package config

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func MySqlDb() *gorm.DB {
	MySqlDb, err := gorm.Open("mysql", "root:Nm123456.@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		FailOnError(err, "gorm.Open")
	}
	return MySqlDb.LogMode(true) //打印mysql 日子
}
