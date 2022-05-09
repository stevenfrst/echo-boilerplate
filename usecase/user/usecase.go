package user

import (
	"echo-boilerplate/app/cmd/http/middlewares"
	"echo-boilerplate/helpers/encoder"
	smtpEmail "echo-boilerplate/helpers/smtp"
	"echo-boilerplate/utils/baseErrors"
	"echo-boilerplate/utils/hash"
	"github.com/labstack/echo/v4"
	"go.elastic.co/apm"
	"log"
)

type Usecase struct {
	repo    IUserRepository
	jwtAuth *middlewares.ConfigJWT
}

func NewUsecase(userRepo IUserRepository, jwt *middlewares.ConfigJWT) IUserUsecase {
	return &Usecase{
		repo:    userRepo,
		jwtAuth: jwt,
	}
}

func (uc *Usecase) Delete(id uint) error {
	err := uc.repo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

func (uc *Usecase) Register(c echo.Context, user Domain) error {
	span, _ := apm.StartSpan(c.Request().Context(), "UseCase Layer -> Register", "request")
	defer span.End()

	u, err := uc.repo.GetByEmail(c, user.Email)
	if err != nil {
		return err
	}
	if u.ID != 0 {
		return baseErrors.ErrUserEmailUsed
	}
	user.Password, err = hash.HashPassword(user.Password)
	if err != nil {
		return err
	}
	err = uc.repo.Create(c, user)
	if err != nil {
		return err
	}
	url := encoder.EncodeUrlEmailVerify(user.Email)
	bodyEmail := `
		<h2>Hello ` + user.Username + `!</h2>
		Please verify your email with click this link : ` + url +
		`<br><br>Regards,<br>Boilerplate Admin`
	err = smtpEmail.SendMail(c, []string{user.Email}, "Boilerplate: Email Registration Confirm", bodyEmail)
	if err != nil {
		return err
	}
	return nil
}

func (uc *Usecase) GetUserByID(id int) (Domain, error) {
	resp, err := uc.repo.GetUserByID(uint(id))
	if err != nil {
		return Domain{}, err
	}
	return resp, nil
}

func (uc *Usecase) Login(c echo.Context, email, password string) (Domain, error) {
	//func (uc *Usecase) Login(email, password string) (Domain, error) {

	user, err := uc.repo.GetByEmail(c, email)
	if !hash.CheckPassword(password, user.Password) {
		return Domain{}, baseErrors.ErrPasswordNotMatch
	} else if user.IsVerified == false {
		return Domain{}, baseErrors.ErrUserInactive
	} else if err != nil {
		return Domain{}, err
	}

	user.Token = uc.jwtAuth.GenerateToken(int(user.ID), user.IsVerified)
	if user.Token == "" {
		return Domain{}, err
	}
	return user, nil
}

func (uc *Usecase) ChangePassword(id int, oldPassword, newPassword string) error {
	resp, err := uc.repo.GetUserByID(uint(id))
	if err != nil {
		return err
	}

	if !hash.CheckPassword(oldPassword, resp.Password) {
		return baseErrors.ErrOldPasswordNotMatch
	}

	resp.Password, err = hash.HashPassword(newPassword)
	err = uc.repo.UpdateUser(resp)
	if err != nil {
		return err
	}
	return nil
}

func (uc *Usecase) ListAllUsers(offset, limit int) ([]Domain, error) {
	resp, err := uc.repo.GetAllUsers(offset, limit)
	if err != nil {
		return []Domain{}, err
	}
	log.Println(len(resp))
	return resp, nil
}

func (uc *Usecase) ListAllUsersVerified(offset, limit int, IsVerified bool) ([]Domain, error) {
	resp, err := uc.repo.GetAllUsersVerify(offset, limit, IsVerified)
	if err != nil {
		return []Domain{}, err
	}
	log.Println(len(resp))
	return resp, nil
}

func (uc *Usecase) VerifyUser(c echo.Context, emailBase64, encrypt string) error {
	email, _ := encoder.DecodeEmailVerify(emailBase64, encrypt)
	if email == "" {
		return baseErrors.ErrInvalidPayload
	}
	u, err := uc.repo.GetByEmail(c, email)
	if err != nil {
		return err
	}
	if u.ID == 0 {
		return baseErrors.ErrUserEmailNotFound
	}
	err = uc.repo.UpdateStatus(u.ID, true)
	if err != nil {
		return err
	}
	bodyEmail := `
		<h2>Hello ` + u.Username + `!</h2>
		Your account has been <font color="green"><b>actived</b></font> :)<br><br>Regards,<br>Boilerplate Admin`
	err = smtpEmail.SendMail(c, []string{u.Email}, "Boilerplate: Email Verified!", bodyEmail)

	return nil
}
