package utils

import (
	"time"

	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

func SignToken(payload *domain.JwtPayload, secret []byte, ttl time.Duration) (string, error) {
	iat := time.Now()
	exp := iat.Add(ttl)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   payload.EntityID,
		Audience:  jwt.ClaimStrings{payload.Role},
		IssuedAt:  jwt.NewNumericDate(iat),
		ExpiresAt: jwt.NewNumericDate(exp),
	})
	return token.SignedString(secret)
}

func VerifyToken(token string, secret []byte) (*domain.JwtPayload, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})
	if err != nil || !parsed.Valid {
		return nil, jwt.ErrTokenNotValidYet
	}

	aud, err := parsed.Claims.GetAudience()
	if err != nil {
		return nil, jwt.ErrTokenInvalidClaims
	}
	if len(aud) == 0 {
		return nil, jwt.ErrTokenInvalidClaims
	}

	id, err := parsed.Claims.GetSubject()
	if err != nil {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return &domain.JwtPayload{EntityID: id, Role: aud[0]}, nil
}
