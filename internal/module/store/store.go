package store

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func Init() {
	initSqlite()
}

func initSqlite() {
	dsn := viper.GetString("db")
	fmt.Sprintln("db dsn:", dsn)
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}
