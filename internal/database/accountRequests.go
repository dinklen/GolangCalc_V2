package database

import (
	"database/sql"

	"github.com/dinklen/GolangCalc_V2/internal/models"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/dberr"

	"go.uber.org/zap"
)

func GetAccount(db *sql.DB, account *models.AccountData, logger *zap.Logger) (*models.AccountData, error) {
	var user models.AccountData

	err := db.QueryRow(
		`
		SELECT id, login, password, created_at
		FROM users
		WHERE login = $1
		LIMIT 1
		`,
		account.Login,
	).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatingTime)

	if err != nil {
		return nil, dberr.ErrAccountGettingFailed
	}

	logger.Info("account get success")
	return &user, nil
}

func CreateAccount(db *sql.DB, account *models.AccountData, logger *zap.Logger) error {
	_, err := db.Exec(
		`
		INSERT INTO users(login, password)
		VALUES ($1, $2)
		`,
		account.Login,
		account.PasswordHash,
	)
	if err != nil {
		return dberr.ErrAccountCreatingFailed
	}

	logger.Info("account create success")
	return nil
}
