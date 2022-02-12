package response

import "echo-boilerplate/usecase/user"

type LoginResponse struct {
	Token      string `json:"token"`
	IsVerified bool   `json:"is_verified"`
}

type UserResponse struct {
	LinkImage  string
	Email      string
	IsVerified bool
	Username   string
}

func FromDomainUser(u user.Domain) UserResponse {
	return UserResponse{
		LinkImage:  u.LinkImage,
		Email:      u.Email,
		IsVerified: u.IsVerified,
		Username:   u.Username,
	}
}

func FromDomainLogin(u user.Domain) LoginResponse {
	return LoginResponse{
		Token:      u.Token,
		IsVerified: u.IsVerified,
	}
}

func ListFromDomainUser(u []user.Domain) (list []UserResponse) {
	for _, x := range u {
		list = append(list, FromDomainUser(x))
	}
	return list
}
