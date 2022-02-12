package driver

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var Mysql *gorm.DB

func SetupDatabase() {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("databases.mysql.user"),
		viper.GetString("databases.mysql.password"),
		viper.GetString("databases.mysql.host"),
		viper.GetString("databases.mysql.port"),
		viper.GetString("databases.mysql.dbname"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	Mysql = db
}
