package handlers

import (
	"net/http"
	"time"

	"github.com/dinklen/GolangCalc_V2/internal/models"
	"github.com/dinklen/GolangCalc_V2/internal/database"
	"github.com/dinklen/GolangCalc_V2/internal/jwt"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"github.com/labstack/echo/v4"
)

func CreateRegisterHandler(db *sql.DB) error {
	return func(ctx echo.Context) error {
		accountData := new(models.AccountData)

		if err := ctx.Bind(accountData); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error":"invalid request"})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(accountData.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		accountData.Password = string(hashedPassword)

		if err := database.CreateAccount(db, accountData); err != nil {
			return ctx.JSON(http.StatusConflict, map[string]string{"error":"user already exists"})
		}
		return ctx.JSON(http.StatusCreated, map[string]string{"message":"success to create account"})
	}
}

func CreateLoginHandler(db *sql.DB, cfg *config.Config, redisClient *redis.Client) error {
	return func(ctx echo.Context) error {
		credentials := new(models.AccountData)
		if err := ctx.Bind(credentials); err != nil {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error":"invalid request"})
		}

		user := new(models.AccountData)
		if user, err = database.GetAccount(db, credentials); user.ID == "" {
			if err != nil {
				return err
			}
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error":"account not found"})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error":"invalid password"})
		}

		access, refresh, err := jwt.GenerateTokens(string(user.ID), cfg, redisClient)
		if err != nil {
			return err
		}

		ctx.JSON(http.StatusOK, map[string]string{
			"access_token": access,
			"refresh_token": refresh,
		})
	}
}

func CreateRefreshHandler(redisClient *redis.Client) error {
	return func(ctx echo.Context) error {
		refreshToken := ctx.Request().Header.Get("Authorization")
		if refreshToken == "" {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "refresh token required"})
		}

		claims, err := jwt.ValidateRefreshToken(refreshToken)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}

		redisClient.Del(ctx.Request().Context(), refreshToken)

		access, refresh, err := jwt.GenerateTokens(claims.UserID)
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, map[string]string{
			"access_token":  access,
			"refresh_token": refresh,
		})
	}
}

func CreateLogOutHandler(redisClient *redis.Client) error {
	return func(ctx echo.Context) error {
		accessToken := ctx.Request().Header.Get("Authorization")
		claims, err := ValidateAccessToken(accessToken)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}

		keys, _ := redisClient.Keys(ctx.Request().Context(), "access:*").Result()
		redisClient.Del(c.Request().Context(), keys...)

		return ctx.JSON(http.StatusOK, map[string]string{"status": "logged out"})
	}
}
