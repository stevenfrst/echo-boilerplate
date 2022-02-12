package main

import (
	routes "echo-boilerplate/app/cmd/http"
	"echo-boilerplate/drivers/mysql"
	"echo-boilerplate/migrations"
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.SetConfigFile(`app/config/config.yaml`)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	driver.SetupDatabase()
	migrations.AutoMigrate()
}

func main() {
	routes.InitHttp()
}
