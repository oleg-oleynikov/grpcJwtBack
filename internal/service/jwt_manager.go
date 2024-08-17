package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTManager struct {
	secretKey       string
	tokenExpiration time.Duration
}

func NewJWTManager(secretKey string, tokenExpiration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:       secretKey,
		tokenExpiration: tokenExpiration,
	}
}

func (manager *JWTManager) GenerateAccessToken(account *Account) (string, error) {
	claims := AccountClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenExpiration).Unix(),
			Subject:   account.Id,
		},
		Email: account.Email,
		Role:  account.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

func (manager *JWTManager) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (manager *JWTManager) Verify(accessToken string) (*AccountClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&AccountClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(manager.secretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	claims, ok := token.Claims.(*AccountClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims: %v", err)
	}

	return claims, nil
}
