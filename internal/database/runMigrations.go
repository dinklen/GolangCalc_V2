package database

import (
	"fmt"

	"github.com/dinklen/GolangCalc_V2/internal/config"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/dberr"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func RunMigrations(cfg *config.Config, logger *zap.Logger) error {
	m, err := migrate.New(
		"file://../../migrations",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		),
	)

	if err != nil {
		return dberr.ErrMigrationsReadingFailed
	}
	logger.Info("migrations creating success")

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return dberr.ErrMigrationsUpFailed
	}
	logger.Info("migrations up success")

	return nil
}
