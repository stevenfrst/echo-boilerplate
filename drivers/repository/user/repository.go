package user

import (
	"echo-boilerplate/usecase/user"
	"gorm.io/gorm"
	"log"
)

type repository struct {
	Mysql *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.IUserRepository {
	return &repository{Mysql: db}
}

func (r *repository) GetByEmail(email string) (user.Domain, error) {
	model := User{}
	err := r.Mysql.Find(&model, "email = ?", email)
	if err.Error != nil {
		return user.Domain{}, err.Error
	}
	return model.ToDomain(), nil
}

func (r *repository) Create(data user.Domain) error {
	model := FromDomain(data)
	err := r.Mysql.Create(&model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetUserByID(id uint) (user.Domain, error) {
	model := User{}
	err := r.Mysql.Where("id = ?", id).First(&model).Error
	if err != nil {
		return user.Domain{}, err
	}
	return model.ToDomain(), err
}

func (r *repository) UpdateUser(user user.Domain) error {
	model := FromDomain(user)
	log.Println(model)
	err := r.Mysql.Save(&model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) DeleteUser(id uint) error {
	model := User{}
	err := r.Mysql.Where("id = ?", id).Delete(&model).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetAllUsers(offset, limit int) ([]user.Domain, error) {
	var models []User
	err := r.Mysql.Offset(offset).Limit(limit).Find(&models).Error
	if err != nil {
		return []user.Domain{}, err
	}
	return UserToDomainList(models), nil
}

func (r *repository) GetAllUsersVerify(offset, limit int, verified bool) ([]user.Domain, error) {
	var models []User
	err := r.Mysql.Offset(offset).Limit(limit).Where("is_verified = ?", verified).Find(&models).Error
	if err != nil {
		return []user.Domain{}, err
	}
	return UserToDomainList(models), nil
}

func (r *repository) UpdateStatus(id uint, state bool) error {
	model := User{}
	err := r.Mysql.Model(&model).Where("id = ?", id).Updates(User{
		IsVerified: state,
	}).Error
	if err != nil {
		return err
	}
	return nil
}
