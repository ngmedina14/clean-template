package service

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/ngmedina14/clean-template/internal/model"
	"github.com/ngmedina14/clean-template/internal/repository"
)

type PrimaryUserService struct {
	IUserRepository repository.IUserRepository
}

func NewPrimaryUserService(userRepository repository.IUserRepository) *PrimaryUserService {
	return &PrimaryUserService{
		IUserRepository: userRepository,
	}
}

func (us *PrimaryUserService) GetUserByID(id int) (*model.User, error) {

	// check cookie validity
	// log action
	user, err := us.IUserRepository.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *PrimaryUserService) UpdateProfile(user *model.User) error {

	err := us.IUserRepository.PatchUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (us *PrimaryUserService) RegisterUser(email, password string) error {
	// Validate email and password
	// Encrypt password
	// Check if user already exists
	// If everything checks out, create user
	// return us.UserRepository.CreateUser(&User{Email: email, Password: encryptedPassword})
	return nil
}

func (us *PrimaryUserService) ValidateLoginUser(input, password string) (*model.User, error) {
	idType, err := identifyInput(input)
	if err != nil {
		return nil, err
	}

	var user *model.User
	switch idType {
	case "Email":
		user, err = us.IUserRepository.GetUserByEmail(input)
	case "Username":
		user, err = us.IUserRepository.GetUserByUsername(input)
	}
	if err != nil {
		return nil, err
	}

	if err := us.IUserRepository.CheckPassword(user, password); err != nil {
		return nil, err
	}

	return user, nil
}
func (us *PrimaryUserService) ValidateRefreshToken(token string) error {
	err := us.IUserRepository.CheckRefreshToken(token)
	if err != nil {
		return err
	}
	return nil
}

func (us *PrimaryUserService) GenerateSaveRefreshToken(user *model.User) (string, error) {
	refreshToken, expiryDate, err := user.GenerateRefreshToken()
	if err != nil {
		return "", err
	}
	if err := us.IUserRepository.SaveUserRefreshToken(user.ID, refreshToken, expiryDate); err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (us *PrimaryUserService) CheckRevokedToken(token string) error {
	err := us.IUserRepository.CheckRevokedToken(token)
	if err != nil {
		return err
	}
	return nil
}

func (us *PrimaryUserService) ExecuteTokenRevoke(token string) error {
	err := us.IUserRepository.SaveRevokedToken(token)
	if err != nil {
		return err
	}
	return nil
}

func (us *PrimaryUserService) RemoveOldToken() error {
	cutoffDate := time.Now().Add(-time.Hour * 24 * 365) // 1 year ago
	err := us.IUserRepository.DeleteOldToken(cutoffDate)
	if err != nil {
		return err
	}
	return nil
}

// Functions

func identifyInput(input string) (string, error) {
	emailPattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	emailRegex := regexp.MustCompile(emailPattern)
	if emailRegex.MatchString(input) {
		return "Email", nil
	}

	// If the input doesn't match the email pattern, assume it's a username
	// Here we're assuming usernames are alphanumeric strings of length 3-30
	// You may need to adjust this depending on your specific requirements
	if len(input) >= 8 && len(input) <= 30 && !strings.ContainsAny(input, "!@#$%^&*()_+={}[]|\\:;<>,.?/~`") {
		return "Username", nil
	}
	// If the input doesn't match either pattern, return an error
	return "", errors.New("invalid username or email input")
}
