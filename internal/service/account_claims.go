package service

import (
	"github.com/dgrijalva/jwt-go"
)

type AccountClaims struct {
	jwt.StandardClaims
	Email string
	Role  string
}
