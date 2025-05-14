package database

import (
	"database/sql"
	"testing"
	"time"

	"github.com/dinklen/GolangCalc_V2/internal/models"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/dberr"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateAccount_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := zap.NewNop()

	mock.ExpectExec("INSERT INTO users").
		WithArgs("testuser", "hashedpassword").
		WillReturnResult(sqlmock.NewResult(1, 1))

	account := &models.AccountData{
		Login:        "testuser",
		PasswordHash: "hashedpassword",
	}

	err = CreateAccount(db, account, logger)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateAccount_Duplicate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := zap.NewNop()

	mock.ExpectExec("INSERT INTO users").
		WithArgs("existinguser", "hash").
		WillReturnError(sql.ErrNoRows)

	account := &models.AccountData{
		Login:        "existinguser",
		PasswordHash: "hash",
	}

	err = CreateAccount(db, account, logger)

	assert.ErrorIs(t, err, dberr.ErrAccountCreatingFailed)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAccount_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := zap.NewNop()

	testUUID := uuid.New()
	testTime := time.Now()

	mock.ExpectQuery("SELECT id, login, password, created_at FROM users").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "created_at"}).
			AddRow(testUUID, "testuser", "hash", testTime))

	account := &models.AccountData{Login: "testuser"}
	result, err := GetAccount(db, account, logger)

	assert.NoError(t, err)
	assert.Equal(t, "testuser", result.Login)
	assert.Equal(t, testUUID, result.ID)
	assert.Equal(t, testTime, result.CreatingTime)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAccount_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := zap.NewNop()

	mock.ExpectQuery("SELECT id, login, password, created_at FROM users").
		WithArgs("unknown").
		WillReturnError(sql.ErrNoRows)

	account := &models.AccountData{Login: "unknown"}
	_, err = GetAccount(db, account, logger)

	assert.ErrorIs(t, err, dberr.ErrAccountGettingFailed)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAccount_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	logger := zap.NewNop()

	mock.ExpectQuery("SELECT id, login, password, created_at FROM users").
		WithArgs("baduser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "login"}). // Missing columns
										AddRow(uuid.New(), "baduser"))

	account := &models.AccountData{Login: "baduser"}
	_, err = GetAccount(db, account, logger)

	assert.ErrorIs(t, err, dberr.ErrAccountGettingFailed)
	assert.NoError(t, mock.ExpectationsWereMet())
}
