package database

import (
	"fmt"
	"database/sql"

	"github.com/dinklen/GolangCalc_V2/internal/models"

	"github.com/jackc/pgx/v5"
)

func GetAccount(db *sql.DB, account *models.AccountData) (*models.AccountData, error) {
	var user models.AccountData

	err := db.QueryRow(
		`
		SELECT id, login, password_hash, created_at
		FROM users
		WHERE login = $1
		ORDERED BY created_at
		LIMIT 1
		`,
		account.Login,
	).Scan(&user.ID, &user.Login, &user.Password, &user.CreatingTime)
	if err != nil {
		return nil, fmt.Errorf("Failed to get account: %v", err)
	}
	return &user, nil
}

func CreateAccount(db *sql.DB, account *models.AccountData) error {
	_, err := db.Exec(
		`
		INSERT INTO users(login, password)
		VALUES ($1, $2)
		`,
		accoint.Login,
		account.Password,
	)
	if err != nil {
		return fmt.Errorf("Failed to create account: %v", err)
	}
	return nil
}
