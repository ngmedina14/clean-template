package main

import (
	"fmt"
	"log"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echomiddlware "github.com/labstack/echo/v4/middleware"
	"github.com/ngmedina14/clean-template/internal/controller"
	"github.com/ngmedina14/clean-template/internal/database"
	"github.com/ngmedina14/clean-template/internal/handler"
	"github.com/ngmedina14/clean-template/internal/repository"
	"github.com/ngmedina14/clean-template/internal/service"
	_ "github.com/ngmedina14/clean-template/internal/swagger"
	val "github.com/ngmedina14/clean-template/pkg/validation"
	"github.com/spf13/viper"
	echoswagger "github.com/swaggo/echo-swagger"
)

// @title Clean API Reference
// @version 3.0
// @description Clean Architecture

// @contact.name API Support
// @contact.email ngmedina14@gmail.com

// @SecurityDefinitions.apikey BearerToken
// @in header
// @name Authorization

// @host localhost:8080
func main() {

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	jwtSecret := viper.GetString("JWT_SECRET")
	dbUrl := viper.GetString("DB_URL")

	e := echo.New()

	// Use the middlewares
	e.Use(echomiddlware.LoggerWithConfig(echomiddlware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(echomiddlware.CORS())

	// Initialize database
	database.InitDB(dbUrl)

	// User initialization
	postgresRepo := repository.NewPostgresUserRepository(database.DB) // low-level module
	userRepo := repository.NewUserRepository(postgresRepo)            // interface
	primaryService := service.NewPrimaryUserService(userRepo)         // high-level module
	userService := service.NewUserService(primaryService)             // interface
	apiController := controller.NewAPIUserController(userService)     // handle http request and response
	userController := controller.NewUserController(apiController)     // interface

	// Handler initialization
	// handlerService := handler.NewHandlerService(database.DB)

	// Register the routes
	e.GET("/health", handler.HealthCheck)

	e.GET("/users/:id", userController.GetUser)
	e.POST("/login", userController.LoginUser)

	apiGroup := e.Group("/api")
	apiGroup.Use(echojwt.JWT([]byte(jwtSecret)))
	e.PATCH("/users", userController.UpdateUser)
	e.POST("/refresh", userController.RefreshToken)
	e.POST("/revoke", userController.RevokeToken)

	// API Docs
	e.GET("/swagger/*", echoswagger.WrapHandler)

	// Custom Interface Implementation
	e.Validator = val.NewCustomValidator()

	// Start the server
	log.Fatal(e.Start(":8080"))
}
