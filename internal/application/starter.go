package application

import (
	"fmt"

	"github.com/dinklen/GolangCalc_V2/internal/config"
	"go.uber.org/zap"
)

type Application struct {
	Config *config.Config
	Logger *zap.Logger
}

func NewApplication() *Application {
	cfg, err := config.Load()
	if err != nil {
		// zap.Error(err.What())
		fmt.Printf("Failed to load config: %v\nSetting default values...\n", err)
		// set default with viper
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return &Application{
		Config: cfg,
		Logger: logger,
	}
}

func (app *Application) Run() error {
	return nil
}
