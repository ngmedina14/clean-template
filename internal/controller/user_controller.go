package controller

import "github.com/labstack/echo/v4"

type UserController interface {
	GetUser(c echo.Context) error
	UpdateUser(c echo.Context) error
	LoginUser(c echo.Context) error
	RefreshToken(c echo.Context) error
	RevokeToken(c echo.Context) error
}

type UserControllerImpl struct {
	UserController
}

func NewUserController(userController UserController) *UserControllerImpl {
	return &UserControllerImpl{
		UserController: userController,
	}
}
