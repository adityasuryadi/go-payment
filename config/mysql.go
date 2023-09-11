package config

import (
	"payment/exception"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlDB(configuration Config) *gorm.DB {
	host := configuration.Get("MYSQL_HOST")
	user := configuration.Get("MYSQL_USER")
	password := configuration.Get("MYSQL_PASSWORD")
	port := configuration.Get("MYSQL_PORT")
	db_name := configuration.Get("MYSQL_DB")

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + db_name + "?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	exception.PanicIfNeeded(err)

	return db
}
