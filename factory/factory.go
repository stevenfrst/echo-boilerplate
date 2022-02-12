package factory

import (
	"echo-boilerplate/app/cmd/http/middlewares"
	d_user "echo-boilerplate/delivery/user"
	driver "echo-boilerplate/drivers/mysql"
	r_user "echo-boilerplate/drivers/repository/user"
	u_user "echo-boilerplate/usecase/user"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

type DeliveryHTTP struct {
	User      *d_user.UserDelivery
	ConfigJWT middleware.JWTConfig
}

func InitFactoryHTTP() DeliveryHTTP {
	configJWT := middlewares.ConfigJWT{
		SecretJWT:       viper.GetString(`jwt.secret`),
		ExpiresDuration: viper.GetInt64(`jwt.expired`),
	}

	userRepo := r_user.NewUserRepository(driver.Mysql)
	userUsecase := u_user.NewUsecase(userRepo, &configJWT)
	userDelivery := d_user.NewUserDelivery(userUsecase)

	return DeliveryHTTP{userDelivery, configJWT.Init()}

}
