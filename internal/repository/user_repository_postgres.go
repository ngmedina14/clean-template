package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ngmedina14/clean-template/common"
	"github.com/ngmedina14/clean-template/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type PostgresUserRepository struct {
	DB *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		DB: db,
	}
}

func (ur *PostgresUserRepository) SaveUser(user *model.User) error {
	userjson, _ := json.Marshal(user)
	var userMap map[string]interface{}
	_ = json.Unmarshal(userjson, &userMap)

	_, err := ur.DB.NamedExec(`INSERT INTO users (name, created_at, updated_at) VALUES (:name,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)`, userMap)
	if err != nil {
		return err
	}
	return nil
}

func (ur *PostgresUserRepository) GetUserByID(id int) (*model.User, error) {
	var user model.User
	err := ur.DB.Get(&user, "SELECT * FROM users WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *PostgresUserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := ur.DB.Get(&user, "SELECT * FROM users WHERE email=$1 AND deleted_at IS NULL", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *PostgresUserRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := ur.DB.Get(&user, "SELECT * FROM users WHERE username=$1 AND deleted_at IS NULL", username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// users, err := ur.FilterUsers("username", true, 10, 0, "username = $1", "John")
func (ur *PostgresUserRepository) FilterUsers(field string, asc bool, limit int, offset int, condition string, args ...interface{}) ([]*model.User, error) {
	users := []*model.User{}
	order := "ASC"
	if !asc {
		order = "DESC"
	}
	// Sanitize field parameter
	field = strings.ReplaceAll(field, "'", `"`)

	query := fmt.Sprintf(`SELECT * FROM users WHERE %s ORDER BY %s %s LIMIT $1 OFFSET $2`, condition, field, order)
	err := ur.DB.Select(&users, query, append(args, limit, offset)...)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *PostgresUserRepository) PutUser(user *model.User) error {
	userjson, _ := json.Marshal(user)
	var userMap map[string]interface{}
	_ = json.Unmarshal(userjson, &userMap)

	_, err := ur.DB.NamedExec(`UPDATE users SET name=:name, updated_at=:updated_at WHERE id=:id`, userMap)
	if err != nil {
		return err
	}
	return nil
}

func (ur *PostgresUserRepository) PatchUser(user *model.User) error {
	userjson, _ := json.Marshal(user)
	var userMap map[string]interface{}
	_ = json.Unmarshal(userjson, &userMap)

	// Remove fields with zero values
	for k, v := range userMap {
		switch t := v.(type) {
		case string:
			if common.IsEmptyOrWhitespace(t) {
				userMap[k] = nil
			}
		case int:
			if t == 0 {
				userMap[k] = nil
			}
		}
	}

	_, err := ur.DB.NamedExec(`UPDATE users SET name=COALESCE(:name, name), username=COALESCE(:username, username), email=COALESCE(:email, email), password=COALESCE(:password, password), updated_at=CURRENT_TIMESTAMP WHERE id=:id`, userMap)
	if err != nil {
		return err
	}
	return nil
}

func (ur *PostgresUserRepository) SoftDeleteUser(id int) error {
	_, err := ur.DB.Exec(`UPDATE users SET deleted_at=CURRENT_TIMESTAMP WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (ur *PostgresUserRepository) DeleteUser(id int) error {
	_, err := ur.DB.Exec(`DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (ur *PostgresUserRepository) CheckPassword(user *model.User, password string) error {
	hashedpassword, err := user.HashedPassword()
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hashedpassword, []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func (ur *PostgresUserRepository) SaveUserRefreshToken(id int, token string, expiryDate time.Time) error {
	_, err := ur.DB.NamedExec(`INSERT INTO refresh_tokens (token, user_id, expiry_date) VALUES (:token,:user_id,:expiry_date)`, map[string]interface{}{"token": token, "user_id": id, "expiry_date": expiryDate})
	if err != nil {
		return err
	}
	return nil
}

func (ur *PostgresUserRepository) CheckRefreshToken(token string) error {
	var count int
	err := ur.DB.Get(&count, "SELECT Count(*) FROM refresh_tokens WHERE token=$1 AND expiry_date>=CURRENT_TIMESTAMP", token)
	if err != nil {
		return err
	}
	if count != 0 {
		return errors.New("token is still valid")
	}
	return nil
}

func (ur *PostgresUserRepository) SaveRevokedToken(token string) error {
	_, err := ur.DB.NamedExec(`INSERT INTO revoked_tokens (token) VALUES (:token)`, map[string]interface{}{"token": token})
	if err != nil {
		return err
	}
	return nil
}

func (ur *PostgresUserRepository) CheckRevokedToken(token string) error {
	var count int
	err := ur.DB.Get(&count, "SELECT Count(*) FROM revoked_tokens WHERE token=$1 ", token)
	if err != nil {
		return err
	}
	if count != 0 {
		return errors.New("token is revoked")
	}
	return nil
}

func (ur *PostgresUserRepository) DeleteOldToken(cutoffDate time.Time) error {
	_, err := ur.DB.Exec(`DELETE FROM refresh_tokens WHERE expiry_date < $1`, cutoffDate)
	if err != nil {
		return err
	}

	_, err = ur.DB.Exec(`DELETE FROM revoke_tokens WHERE expiry_date < $1`, cutoffDate)
	if err != nil {
		return err
	}

	return nil
}
