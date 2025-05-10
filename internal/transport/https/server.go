package https

import (
	"database/sql"

	"github.com/dinklen/GolangCalc_V2/internal/config"
	"github.com/dinklen/GolangCalc_V2/internal/transport/https/handlers"
	"github.com/dinklen/GolangCalc_V2/internal/transport/https/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func Start(cfg *config.Config, redisClient *redis.Client, db *sql.DB) error {
	router := echo.New()

	router.Use(middleware.Recover())
	router.Use(middleware.Secure())

	protectedEndpoints := router.Group("api/v1")
	{
		protectedEndpoints.POST("/calculate", handlers.CalculatorHandler)
		protectedEndpoints.GET("/expressions", handlers.ExpressionsHandler)
		protectedEndpoints.GET("/expressions/:id", handlers.CurrentExpressionHandler)
	}
	protectedEndpoints.Use(middlewares.CreateJWTMiddleware(cfg, redisClient))

	accountEndpoints := router.Group("")
	{
		accountEndpoints.POST("/login", handlers.CreateLoginHandler(db, cfg, redisClient))
		accountEndpoints.POST("/register", handlers.CreateRegisterHandler(db))
		accountEndpoints.POST("/refresh", handlers.CreateRefreshHandler(redisClient))
		accountEndpoints.POST("/logout", handlers.CreateLogOutHandler(redisClient))
	}

	return router.StartTLS(
		cfg.Server.Port,
		"../../../certs/cert.pem",
		"../../../certs/key.pem",
	)
}
