package user

import (
	"github.com/labstack/echo/v4"
)

type Domain struct {
	ID         uint
	LinkImage  string
	Email      string
	IsVerified bool
	Username   string
	Password   string
	Token      string
}

type IUserUsecase interface {
	Register(c echo.Context, user Domain) error
	Login(c echo.Context, email, password string) (Domain, error)
	//Login(email, password string) (Domain, error)
	ChangePassword(id int, oldPassword, newPassword string) error
	GetUserByID(id int) (Domain, error)
	Delete(id uint) error
	ListAllUsers(offset, limit int) ([]Domain, error)
	ListAllUsersVerified(offset, limit int, IsVerified bool) ([]Domain, error)
	VerifyUser(c echo.Context, emailBase64, encrypt string) error
}

type IUserRepository interface {
	GetByEmail(c echo.Context, email string) (Domain, error)
	Create(c echo.Context, data Domain) error
	GetUserByID(id uint) (Domain, error)
	UpdateUser(user Domain) error
	DeleteUser(id uint) error
	GetAllUsers(offset, limit int) ([]Domain, error)
	GetAllUsersVerify(offset, limit int, verified bool) ([]Domain, error)
	UpdateStatus(id uint, state bool) error
}
