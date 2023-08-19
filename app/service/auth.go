package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Brigant/TestTask/app/config"
	"github.com/Brigant/TestTask/app/model"
	"github.com/golang-jwt/jwt"
)

type UserStorager interface {
	FindUsernameWithPasword(c context.Context, userID, password string) (string, error)
}

type AuthService struct {
	Storage UserStorager
}

// AuthService implements AuthService interface.
func (s AuthService) GetToken(ctx context.Context, identity model.IdentityData, cfg config.AuthConfig) (model.Token, error) {
	stringPassword := strconv.Itoa(identity.Password)

	password := model.SHA256(stringPassword, cfg.Salt)

	userID, err := s.Storage.FindUsernameWithPasword(ctx, identity.Username, password)
	if err != nil {
		return model.Token{}, fmt.Errorf("service GetToken error: %w", err)
	}

	accessExpire := time.Now().Add(cfg.AccessTokenTTL)

	accessToken, err := generateJWT(accessExpire, cfg.SigningKey, userID)
	if err != nil {
		return model.Token{}, fmt.Errorf("service GetToken error: %w", err)
	}

	return model.Token{Token: accessToken}, nil
}

// generateJWT generates a new JWT with claims and signs it with the signing key.
func ParseToken(tokenString, signingKey string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, model.ErrInvalidSigningMethod
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("accessToken throws an error during parsing: %w", err)
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok {
		return "", model.ErrWrongTokenClaimType
	}

	return claims.UserID, nil
}

func generateJWT(expire time.Time, signingKey, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		model.Claims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expire.Unix(),
				IssuedAt:  time.Now().Unix(),
			},
			UserID: userID,
		})

	tokenValue, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", fmt.Errorf("cannot get SignetString token: %w", err)
	}

	return tokenValue, nil
}
