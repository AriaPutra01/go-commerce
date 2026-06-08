package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type userClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey}
}

func (maker *JWTMaker) GenerateToken(id, email, role string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &userClaims{
		ID:    id,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(15 * time.Minute))),
		},
	})
	return token.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(tokenString string) (*userClaims, error) {
	claims := new(userClaims)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("Invalid token signing method")
		}
		return []byte(maker.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, err
}
