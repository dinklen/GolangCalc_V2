package jwt

import (
	"context"
	"time"

	"github.com/dinklen/GolangCalc_V2/internal/config"
	"github.com/dinklen/GolangCalc_V2/internal/service/errors/jwterr"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type AccessClaims struct {
	UserID string `json:"sub"`
	Jti    string `json:"jti"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	UserID string `json:"sub"`
	jwt.StandardClaims
}

func GenerateTokens(userID string, cfg *config.Config, redisClient *redis.Client, logger *zap.Logger) (string, string, error) {
	accessJti := uuid.New().String()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims{
		UserID: userID,
		Jti:    accessJti,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(cfg.JWT.AccessExpiry).Unix(),
		},
	})

	access, err := accessToken.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", "", jwterr.ErrAccessTokenGeneratingFailed
	}
	logger.Info("access token creating success")

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(cfg.JWT.RefreshExpiry).Unix(),
		},
	})

	refresh, err := refreshToken.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", "", jwterr.ErrRefreshTokenGeneratingFailed
	}
	logger.Info("refresh token creating success")

	ctx := context.Background()
	redisClient.Set(ctx, "access:"+accessJti, userID, cfg.JWT.AccessExpiry)
	redisClient.Set(ctx, "refresh:"+refresh, userID, cfg.JWT.RefreshExpiry)

	return access, refresh, nil
}

func ValidateAccessToken(tokenString string, cfg *config.Config, redisClient *redis.Client, logger *zap.Logger) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwterr.ErrUnexpectedSigningMethod
			}
			return []byte(cfg.JWT.Secret), nil
		})

	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		if redisClient.Get(context.Background(), "access:"+claims.Jti).Err() == redis.Nil {
			return nil, jwt.NewValidationError("token revoked", jwt.ValidationErrorClaimsInvalid)
		}
		return claims, nil
	}

	if err != nil {
		return nil, jwterr.ErrAccessTokenValidationFailed
	}
	logger.Info("access token validation success")
	return nil, nil
}

func ValidateRefreshToken(tokenString string, cfg *config.Config, redisClient *redis.Client, logger *zap.Logger) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwterr.ErrUnexpectedSigningMethod
			}
			return []byte(cfg.JWT.Secret), nil
		})

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		if redisClient.Get(context.Background(), "refresh:"+tokenString).Err() == redis.Nil {
			return nil, jwt.NewValidationError("token revoked", jwt.ValidationErrorClaimsInvalid)
		}
		return claims, nil
	}

	if err != nil {
		return nil, jwterr.ErrRefreshTokenValidationFailed
	}
	logger.Info("refresh token validation success")
	return nil, nil
}
