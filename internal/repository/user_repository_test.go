package repository_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/ngmedina14/clean-template/internal/model"
	"github.com/ngmedina14/clean-template/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	// Initialize a new mock DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Initialize a new sqlx.DB with the mock DB
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	// Initialize a new PostgresUserRepository with the sqlx.DB
	repo := repository.NewPostgresUserRepository(sqlxDB)

	// Define the SQL query and arguments
	query := regexp.QuoteMeta("SELECT * FROM users WHERE id=$1")
	args := []int{1}

	// Define the expected result
	expectedUser := &model.User{
		ID:        1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      "User 1",
	}

	// Mock the query expectation
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "name"}).AddRow(expectedUser.ID, expectedUser.CreatedAt, expectedUser.UpdatedAt, expectedUser.Name)
	mock.ExpectQuery(query).WithArgs(args[0]).WillReturnRows(rows)

	// Call the method and check the result
	user, err := repo.GetUserByID(args[0])
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
}
