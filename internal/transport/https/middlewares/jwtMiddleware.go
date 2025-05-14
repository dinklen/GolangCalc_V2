package middlewares

import (
	"github.com/dinklen/GolangCalc_V2/internal/config"
	"github.com/dinklen/GolangCalc_V2/internal/jwt"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func CreateJWTMiddleware(cfg *config.Config, redisClient *redis.Client, logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			tokenString := ctx.Request().Header.Get("Authorization")
			if tokenString == "" {
				logger.Error("failed to get token")
				return ctx.JSON(401, map[string]string{"error": "empty jwt"})
			}

			claims, err := jwt.ValidateAccessToken(tokenString, cfg, redisClient, logger)
			if err != nil {
				logger.Error(err.Error())
				return ctx.JSON(401, map[string]string{"error": "invalod token"})
			}

			ctx.Set("userID", claims.UserID)

			logger.Info("JWT authorization success")
			return next(ctx)
		}
	}
}
