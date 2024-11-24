package repo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/SergeyBogomolovv/restaurant/common/config"
	er "github.com/SergeyBogomolovv/restaurant/common/errors"

	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type tokensRepo struct {
	db         *redis.Client
	refreshTTL time.Duration
	accessTTL  time.Duration
	jwtSecret  []byte
}

func NewTokensRepo(db *redis.Client, jwtConfig config.JwtConfig) *tokensRepo {
	return &tokensRepo{
		db:         db,
		refreshTTL: jwtConfig.RefreshTTL,
		accessTTL:  jwtConfig.AccessTTL,
		jwtSecret:  []byte(jwtConfig.Secret),
	}
}

func (r *tokensRepo) GenerateRefreshToken(ctx context.Context, userID string, role string) (string, error) {
	tokenID := uuid.NewString()
	iat := time.Now()
	exp := iat.Add(r.refreshTTL)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   tokenID,
		Audience:  jwt.ClaimStrings{role},
		IssuedAt:  jwt.NewNumericDate(iat),
		ExpiresAt: jwt.NewNumericDate(exp),
	}).SignedString(r.jwtSecret)
	if err != nil {
		return "", err
	}

	payload, err := json.Marshal(&domain.RefreshTokenPayload{
		UserID:    userID,
		ExpiresAt: exp,
		Role:      role,
	})
	if err != nil {
		return "", err
	}
	if err := r.db.Set(ctx, tokenKey(tokenID), payload, r.refreshTTL).Err(); err != nil {
		return "", err
	}

	return token, nil
}

func (r *tokensRepo) VerifyRefreshToken(ctx context.Context, token string) (string, error) {
	token, err := verifyToken(token, r.jwtSecret)
	if err != nil {
		if errors.Is(err, er.ErrInvalidToken) {
			return "", er.ErrUnauthorized
		}
		return "", err
	}

	res, err := r.db.Get(ctx, tokenKey(token)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", er.ErrUnauthorized
		}
		return "", err
	}

	var payload domain.RefreshTokenPayload
	if err := json.Unmarshal(res, &payload); err != nil {
		return "", err
	}

	if payload.ExpiresAt.Before(time.Now()) {
		return "", er.ErrUnauthorized
	}
	return payload.UserID, nil
}

func (r *tokensRepo) RevokeRefreshToken(ctx context.Context, token string) error {
	token, err := verifyToken(token, r.jwtSecret)
	if err != nil {
		if errors.Is(err, er.ErrInvalidToken) {
			return er.ErrUnauthorized
		}
		return err
	}
	return r.db.Del(ctx, tokenKey(token)).Err()
}

func verifyToken(token string, secret []byte) (string, error) {
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !parsed.Valid {
		return "", er.ErrInvalidToken
	}
	return parsed.Claims.GetSubject()
}

func tokenKey(tokenID string) string {
	return fmt.Sprintf("refresh_token:%s", tokenID)
}
