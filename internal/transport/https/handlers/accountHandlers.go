package handlers

import (
	"database/sql"
	"net/http"

	"github.com/dinklen/GolangCalc_V2/internal/config"
	"github.com/dinklen/GolangCalc_V2/internal/database"
	"github.com/dinklen/GolangCalc_V2/internal/jwt"
	"github.com/dinklen/GolangCalc_V2/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func CreateRegisterHandler(db *sql.DB, logger *zap.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		accountData := new(models.AccountData)

		if err := ctx.Bind(accountData); err != nil {
			logger.Error("binding request error")
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(accountData.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("failed to generate password hash")
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate password hash"})
		}
		accountData.PasswordHash = string(hashedPassword)

		if err := database.CreateAccount(db, accountData, logger); err != nil {
			logger.Error(err.Error())
			return ctx.JSON(http.StatusConflict, map[string]string{"error": "user is already exists"})
		}

		logger.Info("user registering success")
		return ctx.JSON(http.StatusCreated, map[string]string{"message": "success to create account"})
	}
}

func CreateLoginHandler(db *sql.DB, cfg *config.Config, redisClient *redis.Client, logger *zap.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		credentials := new(models.AccountData)
		if err := ctx.Bind(credentials); err != nil {
			logger.Error("binding failed")
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
		}

		user, err := database.GetAccount(db, credentials, logger)
		if err != nil {
			logger.Error(err.Error())
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.PasswordHash)); err != nil {
			logger.Error("compare password hash failed")
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		access, refresh, err := jwt.GenerateTokens(user.ID.String(), cfg, redisClient, logger)
		if err != nil {
			logger.Error(err.Error())
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		logger.Info("login success")
		return ctx.JSON(http.StatusOK, map[string]string{
			"access_token":  access,
			"refresh_token": refresh,
		})
	}
}

func CreateRefreshHandler(redisClient *redis.Client, cfg *config.Config, logger *zap.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		refreshToken := new(models.RefreshToken)
		if err := ctx.Bind(refreshToken); err != nil {
			logger.Error("refresh token bind failed")
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "failed to bind request"})
		}

		claims, err := jwt.ValidateRefreshToken(refreshToken.Token, cfg, redisClient, logger)
		if err != nil {
			logger.Error("refresh token is incorrect")
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}

		redisClient.Del(ctx.Request().Context(), refreshToken.Token)

		access, refresh, err := jwt.GenerateTokens(claims.UserID, cfg, redisClient, logger)
		if err != nil {
			logger.Error(err.Error())
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		logger.Info("refresh success")
		return ctx.JSON(http.StatusOK, map[string]string{
			"access_token":  access,
			"refresh_token": refresh,
		})
	}
}

func CreateLogOutHandler(redisClient *redis.Client, cfg *config.Config, logger *zap.Logger) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		accessToken := new(models.AccessToken)
		if err := ctx.Bind(accessToken); err != nil {
			log.Error("access token bind failed")
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "failed to bind request"})
		}

		_, err := jwt.ValidateAccessToken(accessToken.Token, cfg, redisClient, logger)
		if err != nil {
			logger.Error(err.Error())
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}

		keys, _ := redisClient.Keys(ctx.Request().Context(), "access:*").Result()
		redisClient.Del(ctx.Request().Context(), keys...)

		logger.Info("log out success")
		return ctx.JSON(http.StatusOK, map[string]string{"message": "logged out"})
	}
}
