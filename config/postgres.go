package config

import (
	"payment/exception"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(configuration Config) *gorm.DB {
	host := configuration.Get("POSTGRES_HOST")
	user := configuration.Get("POSTGRES_USER")
	password := configuration.Get("POSTGRES_PASSWORD")
	port := configuration.Get("POSTGRES_PORT")
	db_name := configuration.Get("POSTGRES_DB")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + db_name + " port=" + port + " sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	exception.PanicIfNeeded(err)

	return db
}

func NewTestPostgresDB(configuration Config) *gorm.DB {
	host := "localhost"
	user := configuration.Get("POSTGRES_USER")
	password := configuration.Get("POSTGRES_PASSWORD")
	port := "5433"
	db_name := configuration.Get("POSTGRES_DB")

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + db_name + " port=" + port + " sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	exception.PanicIfNeeded(err)

	return db
}
