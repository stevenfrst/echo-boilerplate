package migrations

import (
	driver "echo-boilerplate/drivers/mysql"
	"echo-boilerplate/drivers/repository/user"
)

func AutoMigrate() {
	err := driver.Mysql.AutoMigrate(
			&user.User{},
	)
	if err != nil {
		return
	}
}