package jwt

import (
	"time"
	"fmt"
	// errors

	"github.com/dinklen/GolangCalc_V2/internal/config"
	
	"github.com/redis/go-redis/v9"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AccessClaims struct {
	UserID string   `json:"sub"`
	Jti    string   `json:"jti"`
	jwt.StandardClaims
}

type RefreshClaims struct {
	UserID string `json:"sub"`
	jwt.StandardClaims
}

func GenerateTokens(userID string, cfg *config.Config, redisClient *redis.Client) (string, string, error) {
	accessJti := uuid.New().String()
	
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims{
		UserID: userID,
		Jti:    accessJti,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	})

	access, err := accessToken.SignedString(cfg.JWT.Secret)
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	})

	refresh, err := refreshToken.SignedString(cfg.JWT.Secret)
	if err != nil {
		return "", "", err
	}

	redisClient.Set(redisClient.Context(), 
		"access:"+accessJti, 
		userID, 
		cfg.JWT.AccessExpiry,
	)

	redisClient.Set(redisClient.Context(), 
		"refresh:"+refresh, 
		userID, 
		cfg.JWT.RefreshExpiry,
	)

	return access, refresh, nil
}

func ValidateAccessToken(tokenString string, cfg *config.Config, redisClient *redis.Client) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, 
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return cfg.JWT.Secret, nil
		})
	
	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		if redisClient.Get(redisClient.Context(), "access:"+claims.Jti).Err() == redis.Nil {
			return nil, jwt.NewValidationError("token revoked", jwt.ValidationErrorClaimsInvalid)
		}
		return claims, nil
	}
	return nil, err
}

func ValidateRefreshToken(tokenString string, cfg *config.Config, redisClient *redis.Client) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, 
		func(t *jwt.Token) (interface{}, error) { 
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
		})
	
	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		if redisClient.Get(redisClient.Context(), "refresh:"+tokenString).Err() == redis.Nil {
			return nil, jwt.NewValidationError("token revoked", jwt.ValidationErrorClaimsInvalid)
		}
		return claims, nil
	}
	return nil, err
}
