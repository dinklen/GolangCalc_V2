package https

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dinklen/GolangCalc_V2/internal/config"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/neterr"
	"github.com/dinklen/GolangCalc_V2/internal/tm"
	"github.com/dinklen/GolangCalc_V2/internal/transport/https/handlers"
	"github.com/dinklen/GolangCalc_V2/internal/transport/https/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Start(cfg *config.Config, redisClient *redis.Client, db *sql.DB, taskManager *tm.TaskManager, logger *zap.Logger) error {
	router := echo.New()

	router.Use(middleware.Recover())
	router.Use(middleware.Secure())

	protectedEndpoints := router.Group("api/v1", middlewares.CreateJWTMiddleware(cfg, redisClient, logger))
	{
		protectedEndpoints.POST("/calculate", handlers.CreateCalculatorHandler(db, taskManager, logger))
		protectedEndpoints.GET("/expressions", handlers.CreateExpressionsHandler(db, logger))
		protectedEndpoints.GET("/expressions/:id", handlers.CreateCurrentExpressionHandler(db, logger))
	}

	accountEndpoints := router.Group("")
	{
		accountEndpoints.POST("/login", handlers.CreateLoginHandler(db, cfg, redisClient, logger))
		accountEndpoints.POST("/register", handlers.CreateRegisterHandler(db, logger))
		accountEndpoints.POST("/refresh", handlers.CreateRefreshHandler(redisClient, cfg, logger))
		accountEndpoints.POST("/logout", handlers.CreateLogOutHandler(redisClient, cfg, logger))
	}

	go func() {
		err := router.StartTLS(
			fmt.Sprintf(
				"%s:%s",
				cfg.Server.Host,
				cfg.Server.Port,
			),
			"../../certs/cert.pem",
			"../../certs/key.pem",
		)
		if err != nil && err != http.ErrServerClosed {
			logger.Error(neterr.ErrStartingFailed.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	cancel()

	if err := router.Shutdown(ctx); err != nil {
		logger.Error(neterr.ErrShuttingDownFailed.Error())
	}

	db.Close()

	fmt.Println()
	logger.Info("Ctrl+C is a superpower! :)")
	return nil
}
