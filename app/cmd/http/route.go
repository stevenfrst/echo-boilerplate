package routes

import (
	"echo-boilerplate/factory"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func InitHttp() {
	f := factory.InitFactoryHTTP()
	e := echo.New()
	e.Validator = &CustomValidator{Validator: validator.New()}

	v1 := e.Group("api/v1")
	// Auth Endpoint
	v1.POST("/users", f.User.Register)
	v1.POST("/users/login", f.User.Login)
	v1.PATCH("/users/change-password", f.User.ChangePassword, middleware.JWTWithConfig(f.ConfigJWT))
	v1.GET("/users/me", f.User.GetDetail, middleware.JWTWithConfig(f.ConfigJWT))
	v1.DELETE("/users", f.User.Delete, middleware.JWTWithConfig(f.ConfigJWT))
	v1.GET("/users", f.User.ListAllUsers)

	v1.GET("/verify",f.User.Verify)
	err := e.Start(":1234")

	if err != nil {
		return
	}
}
