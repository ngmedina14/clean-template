package handler

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/ngmedina14/clean-template/internal/database"
)

type HandlerService struct {
	DB *sqlx.DB
}

func NewHandlerService(db *sqlx.DB) *HandlerService {
	return &HandlerService{
		DB: db,
	}
}

func HealthCheck(c echo.Context) error {
	if err := database.CheckHealth(); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "OK")
}

func (hs *HandlerService) HanlderThatUsesDatabase(c echo.Context) error {

	_, err := hs.DB.Exec(`SELECT * FROM users WHERE id=$1`, 1)
	if err != nil {
		return err
	}
	return nil
}
