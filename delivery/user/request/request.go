package request

import "echo-boilerplate/usecase/user"

type UserRegister struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserChangePassword struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

func (u *UserLogin) LoginToDomain() user.Domain {
	return user.Domain{
		ID:         0,
		LinkImage:  "",
		Email:      u.Email,
		IsVerified: false,
		Username:   "",
		Password:   u.Password,
		Token: "",
	}
}


func (u *UserRegister) RegisterToDomain() user.Domain {
	return user.Domain{
		ID:         0,
		LinkImage:  "",
		Email:      u.Email,
		IsVerified: false,
		Username:   u.Username,
		Password:   u.Password,
		Token: "",
	}
}
