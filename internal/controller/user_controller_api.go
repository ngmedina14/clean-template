package controller

import (
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/ngmedina14/clean-template/internal/model"
	"github.com/ngmedina14/clean-template/internal/service"
)

type APIUserController struct {
	UserService service.UserService
}

func NewAPIUserController(userService service.UserService) *APIUserController {
	return &APIUserController{
		UserService: userService,
	}
}

func checkAndHandleToken(c echo.Context, uc APIUserController) (id int, fullAccess bool, err error) {
	// Get ID from the token
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	refreshToken := token.Raw
	if err := uc.UserService.CheckRevokedToken(refreshToken); err != nil {
		return 0, false, echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	userIDStr, ok := claims["jti"].(string)
	if !ok {
		return 0, false, echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": "invalid payload"})
	}

	issuer, ok := claims["iss"]
	if !ok {
		return 0, false, echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "invalid payload"})
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, false, echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": "string conversion failed"})
	}

	if issuer == "refresh" {
		return userID, false, nil
	}
	return userID, true, nil
}

// GetUser returns the user with the given userID
// @Summary return a specific user based on given ID
// @Description get the user by given ID
// @Tags USERS
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "Users ID"
// @Success 200 {object} model.User
// @Failure 400 {object} echo.HTTPError
// @Router /users/{id} [get]
func (uc APIUserController) GetUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := uc.UserService.GetUserByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

// UpdateUser returns Status Accepted
// @Summary return a specific user based on given ID
// @Description get the user by given ID
// @Tags USERS
// @Produce  json
// @Security BearerToken
// @Param name path string true "Users Name"
// @Success 202 {object} model.User
// @Failure 400 {object} echo.HTTPError
// @Router /users [patch]
func (uc APIUserController) UpdateUser(c echo.Context) error {
	userID, access, err := checkAndHandleToken(c, uc)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}
	if !access {
		// TODO: redirect to login page
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "you need to re-login"})
	}
	userPatch := new(model.User)
	if err := c.Bind(userPatch); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	userPatch.ID = userID

	if err = uc.UserService.UpdateProfile(userPatch); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusAccepted, echo.Map{"status": "ok"})
}

func (uc APIUserController) LoginUser(c echo.Context) error {
	type LoginRequest struct {
		ID       string `json:"id" form:"id"`
		Password string `json:"password" form:"password"`
	}
	req := new(LoginRequest)
	if err := c.Bind(req); err != nil {
		return err
	}

	user, err := uc.UserService.ValidateLoginUser(req.ID, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	accessToken, err := user.GenerateAccessToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	generatedRefreshToken, err := uc.UserService.GenerateSaveRefreshToken(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"id": user.ID, "email": user.Email, "token": accessToken, "refresh": generatedRefreshToken})
}

func (uc APIUserController) RefreshToken(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["iss"] != "refresh" {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
	}
	refreshToken := token.Raw

	if err := uc.UserService.CheckRevokedToken(refreshToken); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	if err := uc.UserService.ValidateRefreshToken(refreshToken); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, echo.Map{"error": err.Error()})
	}

	// Generate a new refresh token
	userIDStr, ok := claims["jti"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": "invalid token"})
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{"error": "string conversion failed"})
	}
	userRefreshToken := new(model.User)
	userRefreshToken.ID = userID
	generatedRefreshToken, err := uc.UserService.GenerateSaveRefreshToken(userRefreshToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	// Store the refresh token in an HttpOnly cookie
	// cookie := new(http.Cookie)
	// cookie.Name = "refresh_token"
	// cookie.Value = generatedRefreshToken
	// cookie.Path = "/"
	// cookie.HttpOnly = true
	// c.SetCookie(cookie)

	// return c.JSON(http.StatusOK, map[string]string{
	//     "message": "Refresh token updated",
	// })

	return c.JSON(http.StatusOK, map[string]string{
		"token": generatedRefreshToken,
	})
}

func (uc APIUserController) RevokeToken(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["iss"] != "refresh" {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
	}

	refreshToken := token.Raw
	if err := uc.UserService.ExecuteTokenRevoke(refreshToken); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func (uc APIUserController) SessionAutoTokenCleanup(c echo.Context) error {
	_, access, err := checkAndHandleToken(c, uc)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}
	if !access {
		// TODO: redirect to login page
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": "you need to re-login"})
	}

	err = uc.UserService.RemoveOldToken()
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
