package repository

import (
	"time"

	"github.com/ngmedina14/clean-template/internal/model"
)

type IUserRepository interface {
	SaveUser(user *model.User) error
	GetUserByID(id int) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	FilterUsers(field string, asc bool, limit int, offset int, condition string, args ...interface{}) ([]*model.User, error)
	PutUser(user *model.User) error
	PatchUser(user *model.User) error
	DeleteUser(id int) error
	SoftDeleteUser(id int) error

	CheckPassword(user *model.User, password string) error
	SaveUserRefreshToken(id int, token string, expiryDate time.Time) error
	CheckRefreshToken(token string) error
	SaveRevokedToken(token string) error
	CheckRevokedToken(token string) error
	DeleteOldToken(cutoffDate time.Time) error
}

type UserRepositoryImpl struct {
	IUserRepository
}

func NewUserRepository(userRepository IUserRepository) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		IUserRepository: userRepository,
	}
}
