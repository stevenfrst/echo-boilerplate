package delivery

import (
	"echo-boilerplate/app/cmd/http/middlewares"
	"echo-boilerplate/delivery"
	"echo-boilerplate/delivery/user/request"
	"echo-boilerplate/delivery/user/response"
	"echo-boilerplate/usecase/user"
	"echo-boilerplate/utils/baseErrors"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type UserDelivery struct {
	usecase user.IUserUsecase
}

func NewUserDelivery(u user.IUserUsecase) *UserDelivery {
	return &UserDelivery{
		usecase: u,
	}
}

func (d *UserDelivery) Verify(c echo.Context) error {
	email := c.QueryParam("u")
	verify := c.QueryParam("v")

	err := d.usecase.VerifyUser(email, verify)
	if err != nil {
		return delivery.ErrorResponse(c, http.StatusInternalServerError, "", err)
	}
	return delivery.SuccessResponse(c, "verified")
}

func (d *UserDelivery) ListAllUsers(c echo.Context) error {
	var (
		users []user.Domain
		err   error
	)
	verifyParams := c.QueryParam("verified")

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 20
	}
	offset = (offset - 1) * limit

	if verifyParams != "" {
		verified, _ := strconv.Atoi(verifyParams)
		var isVerified bool
		if verified > 1 {
			return delivery.ErrorResponse(c, http.StatusBadRequest, "", baseErrors.ErrWrongQueryParams)
		} else if verified != 0 {
			isVerified = true
		} else {
			isVerified = false
		}
		users, err = d.usecase.ListAllUsersVerified(offset, limit, isVerified)
	} else {
		users, err = d.usecase.ListAllUsers(offset, limit)
	}

	if err != nil {
		return delivery.ErrorResponse(c, http.StatusInternalServerError, "", err)
	}

	return delivery.SuccessResponse(c, response.ListFromDomainUser(users))
}
func (d *UserDelivery) Delete(c echo.Context) error {
	jwtMeta := middlewares.GetUser(c)
	err := d.usecase.Delete(uint(jwtMeta.ID))
	if err != nil {
		return delivery.ErrorResponse(c, http.StatusInternalServerError, fmt.Sprintf("%v", err), baseErrors.ErrDeleteAccount)
	}
	return delivery.SuccessResponse(c, "success")
}

func (d *UserDelivery) GetDetail(c echo.Context) error {
	jwtMeta := middlewares.GetUser(c)
	domain, err := d.usecase.GetUserByID(jwtMeta.ID)
	if err != nil {
		return delivery.ErrorResponse(c, http.StatusInternalServerError, "", err)
	}
	return delivery.SuccessResponse(c, response.FromDomainUser(domain))
}
func (d *UserDelivery) ChangePassword(c echo.Context) error {
	req := request.UserChangePassword{}
	if err := c.Bind(&req); err != nil {
		return delivery.ErrorResponse(c, http.StatusBadRequest, "", err)
	}
	if err := c.Validate(&req); err != nil {
		return delivery.ErrorResponse(c, http.StatusBadRequest, "", err)
	}

	jwtMeta := middlewares.GetUser(c)

	err := d.usecase.ChangePassword(jwtMeta.ID, req.OldPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, baseErrors.ErrOldPasswordNotMatch) {
			return delivery.ErrorResponse(c, http.StatusBadRequest, "", err)
		}
		return delivery.ErrorResponse(c, http.StatusInternalServerError, "", err)
	}
	return delivery.SuccessResponse(c, "success")
}

func (d *UserDelivery) Register(c echo.Context) error {
	req := request.UserRegister{}
	if err := c.Bind(&req); err != nil {
		return delivery.ErrorResponse(c, http.StatusBadRequest, "", err)
	}
	if err := c.Validate(&req); err != nil {
		return delivery.ErrorResponse(c, http.StatusBadRequest, "", err)
	}

	err := d.usecase.Register(req.RegisterToDomain())
	if err != nil {
		return delivery.ErrorResponse(c, http.StatusInternalServerError, "", err)
	}
	return delivery.SuccessResponse(c, "success")
}

func (d *UserDelivery) Login(c echo.Context) error {
	req := request.UserLogin{}
	if err := c.Bind(&req); err != nil {
		return delivery.ErrorResponse(c, http.StatusBadRequest, "", err)
	}
	if err := c.Validate(&req); err != nil {
		return delivery.ErrorResponse(c, http.StatusBadRequest, "", err)
	}

	resp, err := d.usecase.Login(req.Email, req.Password)
	if err != nil {
		if errors.Is(err, baseErrors.ErrUserInactive) || errors.Is(err, baseErrors.ErrPasswordNotMatch) {
			return delivery.ErrorResponse(c, http.StatusForbidden, "", err)
		}
		return delivery.ErrorResponse(c, http.StatusInternalServerError, "", err)
	}
	return delivery.SuccessResponse(c, response.FromDomainLogin(resp))
}
