package servicios

import (
	"github.com/Amari05fc/practica-7/database"
	"gorm.io/gorm"
	)

type UserService interface {
	CreateUser(user database.User) error
}

type userService struct{
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService{
	return &userService{db}
}

func (i *userService) CreateUser(user database.User) error{
	err:= i.db.Create(&user)

	return err.Error
}
