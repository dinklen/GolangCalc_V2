package middleware

import (
	"githib.com/dinklen/GolangCalc_V2/internal/jwt"
	"github.com/dinklen/GolangCalc_V2/internal/config"

	"github.com/redis/go-redis/v9"
	"github.com/labstack/echo/v4"
)

func CreateJWTMiddleware(cfg *config.Config, redisClient *redis.Client) error {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			tokenString := ctx.Request().Header.Get("Authorization")
			if tokenString == "" {
				return ctx.JSON(401, map[string]string{"error":"empty jwt"})
			}

			claims, err := jwt.ValidateAccessToken(tokenString, cfg, redisClient)
			if err != nil {
				return ctx.JSON(401, map[string]string{"error":"invalod token"})
			}

			ctx.Set("userID", claims.UserID)
			return next(ctx)
		}
	}
}
