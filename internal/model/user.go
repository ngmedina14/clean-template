package model

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int          `db:"id" json:"id" form:"id"`
	Name      string       `db:"name" json:"name" form:"name"`
	Username  string       `db:"username" json:"username" form:"username"`
	Email     string       `db:"email" json:"email" form:"email"`
	Password  string       `db:"password" json:"password" form:"password"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
}

func (u *User) HashedPassword() ([]byte, error) {
	// Generate a bcrypt hash of the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Convert the hashed password to a string and return it
	return hashedPassword, nil
}

func (u *User) GenerateAccessToken() (string, error) {
	// Define the claims
	claims := &jwt.RegisteredClaims{
		ID:        strconv.Itoa(u.ID),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
		Issuer:    "access",
	}

	// Create the token using your secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := viper.GetString("JWT_SECRET")

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *User) GenerateRefreshToken() (string, time.Time, error) {
	// Define the claims

	expiryDate := time.Now().Add(720 * time.Hour)
	claims := &jwt.RegisteredClaims{
		ID:        strconv.Itoa(u.ID),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expiryDate), // Set the expiration time to 30 days from now
		Issuer:    "refresh",
	}

	// Create the token using your secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := viper.GetString("JWT_SECRET")

	// Sign and get the complete encoded token as a string using the secret
	refreshToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", time.Now(), err
	}

	return refreshToken, expiryDate, nil
}
