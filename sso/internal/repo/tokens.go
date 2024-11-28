package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/SergeyBogomolovv/restaurant/common/config"

	"github.com/SergeyBogomolovv/restaurant/sso/internal/domain/entities"
	errs "github.com/SergeyBogomolovv/restaurant/sso/internal/domain/errors"
	"github.com/SergeyBogomolovv/restaurant/sso/pkg/payload"
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

func (r *tokensRepo) GenerateRefreshToken(ctx context.Context, entityID string, role string) (string, error) {
	tokenID := uuid.NewString()
	iat := time.Now()
	exp := iat.Add(r.refreshTTL)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   tokenID,
		IssuedAt:  jwt.NewNumericDate(iat),
		ExpiresAt: jwt.NewNumericDate(exp),
	}).SignedString(r.jwtSecret)
	if err != nil {
		return "", err
	}

	payload, err := json.Marshal(&entities.RefreshTokenEntity{
		EntityID:  entityID,
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

func (r *tokensRepo) VerifyRefreshToken(ctx context.Context, token string) (*payload.JwtPayload, error) {
	tokenID, err := r.verifyRefreshTokenID(token)
	if err != nil {
		return nil, errs.ErrInvalidJwtToken
	}

	res, err := r.db.Get(ctx, tokenKey(tokenID)).Bytes()
	if err != nil {
		return nil, errs.ErrInvalidJwtToken
	}

	var refreshToken entities.RefreshTokenEntity
	if err := json.Unmarshal(res, &refreshToken); err != nil {
		return nil, err
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, errs.ErrInvalidJwtToken
	}
	return &payload.JwtPayload{EntityID: refreshToken.EntityID, Role: refreshToken.Role}, nil
}

func (r *tokensRepo) RevokeRefreshToken(ctx context.Context, token string) error {
	tokenID, err := r.verifyRefreshTokenID(token)
	if err != nil {
		return errs.ErrInvalidJwtToken
	}
	return r.db.Del(ctx, tokenKey(tokenID)).Err()
}

func (r *tokensRepo) verifyRefreshTokenID(jwtToken string) (string, error) {
	parsed, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return r.jwtSecret, nil
	})
	if err != nil || !parsed.Valid {
		return "", errs.ErrInvalidJwtToken
	}
	return parsed.Claims.GetSubject()
}

func (r *tokensRepo) SignAccessToken(entityID string, role string) (string, error) {
	iat := time.Now()
	exp := iat.Add(r.accessTTL)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   entityID,
		Audience:  jwt.ClaimStrings{role},
		IssuedAt:  jwt.NewNumericDate(iat),
		ExpiresAt: jwt.NewNumericDate(exp),
	})
	return token.SignedString(r.jwtSecret)
}

func tokenKey(tokenID string) string {
	return fmt.Sprintf("refresh_token:%s", tokenID)
}
