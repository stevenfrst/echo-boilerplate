package user

import (
	"echo-boilerplate/usecase/user"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID         uint `gorm:"primarykey"`
	LinkImage  string
	IsVerified bool   `gorm:"not null"`
	Username   string `gorm:"unique;not null"`
	Password   string `gorm:"not null"`
	Email      string `gorm:"unique"`
	CreatedAt  *time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (u *User) ToDomain() user.Domain {
	return user.Domain{
		ID:         u.ID,
		LinkImage:  u.LinkImage,
		Email:      u.Email,
		IsVerified: u.IsVerified,
		Username:   u.Username,
		Password:   u.Password,
	}
}

func FromDomain(u user.Domain) User {
	return User{
		ID:         u.ID,
		LinkImage:  u.LinkImage,
		IsVerified: u.IsVerified,
		Username:   u.Username,
		Password:   u.Password,
		Email:      u.Email,
	}
}

func UserToDomainList(user []User) (list []user.Domain) {
	for _, v := range user {
		list = append(list, v.ToDomain())
	}
	return list
}
