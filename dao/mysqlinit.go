package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func MysqlInit() *gorm.DB {

	db, err := gorm.Open(mysql.Open("Sianao:Simple2002@tcp(localhost)/sql_test"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Println(err)
	}
	sqlset, err := db.DB()
	sqlset.SetMaxIdleConns(10)
	DB = db
	return DB

}
