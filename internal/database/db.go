package database

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	DB   *sqlx.DB
	once sync.Once
)

func InitDB(dbUrl string) {
	once.Do(func() {
		var err error
		DB, err = sqlx.Connect("postgres", dbUrl)
		if err != nil {
			fmt.Println(err.Error())
			panic(err)
		}
	})
}

func CloseDB() error {
	return DB.Close()
}

func BeginTransaction() (*sqlx.Tx, error) {
	return DB.Beginx()
}

func CommitTransaction(tx *sqlx.Tx) error {
	return tx.Commit()
}

func RollbackTransaction(tx *sqlx.Tx) error {
	return tx.Rollback()
}

func CheckHealth() error {
	_, err := DB.Exec("SELECT 1")
	return err
}

func LoadInitialData() error {
	query := `
        INSERT INTO users (id, name, created_at, updated_at) VALUES
        (1, 'User 1', '2024-01-01 00:00:00', '2024-01-01 00:00:00'),
        (2, 'User 2', '2024-01-02 00:00:00', '2024-01-02 00:00:00'),
        (3, 'User 3', '2024-01-03 00:00:00', '2024-01-03 00:00:00'),
        (4, 'User 4', '2024-01-04 00:00:00', '2024-01-04 00:00:00'),
        (5, 'User 5', '2024-01-05 00:00:00', '2024-01-05 00:00:00');
    `
	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
