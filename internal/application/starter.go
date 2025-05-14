package application

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "github.com/dinklen/GolangCalc_V2/api/proto/generated"
	"github.com/dinklen/GolangCalc_V2/internal/config"
	"github.com/dinklen/GolangCalc_V2/internal/database"
	redisc "github.com/dinklen/GolangCalc_V2/internal/redis"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/dberr"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/tmerr"
	"github.com/dinklen/GolangCalc_V2/internal/tm"
	grpcs "github.com/dinklen/GolangCalc_V2/internal/transport/grpc"
	"github.com/dinklen/GolangCalc_V2/internal/transport/https"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Application struct {
	Config      *config.Config
	Logger      *zap.Logger
	DB          *sql.DB
	TaskManager *tm.TaskManager
	RedisClient *redis.Client
}

func initLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("\x1b[38;5;87m" + t.Format("02.01.2006 15:04:05.000") + "\x1b[0m")
	}

	config.EncoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		color := ""
		text := l.CapitalString()

		switch l {
		case zapcore.DebugLevel:
			color = "\x1b[1;35m"
			text = "BYE"
		case zapcore.InfoLevel:
			color = "\x1b[1;34m"
		case zapcore.WarnLevel:
			color = "\x1b[1;33m"
		case zapcore.ErrorLevel:
			color = "\x1b[1;31m"
		case zapcore.FatalLevel:
			color = "\x1b[1;38;5;88m"
		}

		enc.AppendString(color + "[" + text + "]\x1b[0m")
	}

	config.DisableCaller = true
	config.DisableStacktrace = true

	logger, _ := config.Build()

	return logger
}

func initConfig(logger *zap.Logger) *config.Config {
	cfg, err := config.Load(logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return cfg
}

func initDatabase(cfg *config.Config, logger *zap.Logger) *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	))
	if err != nil {
		logger.Fatal(dberr.ErrDatabaseCreatedFailed.Error())
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal(dberr.ErrDatabasePingFailed.Error())
	}

	return db
}

func initTaskManager(cfg *config.Config, db *sql.DB, logger *zap.Logger) *tm.TaskManager {
	creds, err := credentials.NewClientTLSFromFile("../../certs/cert.pem", "")
	if err != nil {
		logger.Fatal(tmerr.ErrClientCreatingFailed.Error())
	}

	connection, err := grpc.Dial(
		fmt.Sprintf(
			"%s:%s",
			cfg.Microservice.Host,
			cfg.Microservice.Port,
		),
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		logger.Fatal(tmerr.ErrClientConnectionRefused.Error())
	}

	client := pb.NewCalculatorClient(connection)

	taskManager, err := tm.NewTaskManager(client, cfg, db, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return taskManager
}

func initRedis(cfg *config.Config, logger *zap.Logger) *redis.Client {
	redisClient, err := redisc.NewClient(cfg, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return redisClient
}

func NewApplication() *Application {
	logger := initLogger()
	config := initConfig(logger)

	return &Application{
		Config: config,
		Logger: logger,
	}
}

func (app *Application) Run() error {
	err := database.RunMigrations(app.Config, app.Logger)
	if err != nil {
		app.Logger.Fatal(err.Error())
	}

	_, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := grpcs.Start(app.Config, app.Logger)
		if err != nil {
			app.Logger.Fatal(err.Error())
		}
	}()

	app.DB = initDatabase(app.Config, app.Logger)
	app.TaskManager = initTaskManager(app.Config, app.DB, app.Logger)
	app.RedisClient = initRedis(app.Config, app.Logger)

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := https.Start(app.Config, app.RedisClient, app.DB, app.TaskManager, app.Logger)
		if err != nil {
			app.Logger.Fatal(err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	cancel()

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer timeoutCancel()

	go func() {
		wg.Wait()
		timeoutCancel()
	}()

	<-timeoutCtx.Done()

	app.Logger.Debug("Thank you for a testing, goodbye!")
	return nil
}
