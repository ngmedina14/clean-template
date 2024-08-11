package service

import (
	"github.com/ngmedina14/clean-template/internal/model"
)

type UserService interface {
	RegisterUser(id, password string) error
	ValidateLoginUser(input, password string) (*model.User, error)
	GetUserByID(id int) (*model.User, error)
	UpdateProfile(user *model.User) error
	ValidateRefreshToken(token string) error
	GenerateSaveRefreshToken(user *model.User) (string, error)
	CheckRevokedToken(token string) error
	ExecuteTokenRevoke(token string) error
	RemoveOldToken() error
}

type UserServiceImpl struct {
	UserService
}

func NewUserService(userService UserService) *UserServiceImpl {
	return &UserServiceImpl{
		UserService: userService,
	}
}
