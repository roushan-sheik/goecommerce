package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenExpiry  = 15 * time.Minute
	refreshTokenExpiry = 7 * 24 * time.Hour

	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
	Issuer           = "goecommerce"
)

type JWTClaims struct {
	UserId    uint   `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateAccessToken(userId uint, name string, email string) (string, error)
	GenerateRefreshToken(userId uint, name string, email string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}

type jwtService struct {
	secretKey string
}

func NewJWTService(secretkey string) JWTService {
	if secretkey == "" {
		panic("JWT secret is not defined")
	}

	return &jwtService{
		secretKey: secretkey,
	}
}

func (s *jwtService) GenerateAccessToken(userId uint, name string, email string) (string, error) {
	claims := &JWTClaims{
		UserId:    userId,
		Name:      name,
		Email:     email,
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) GenerateRefreshToken(userId uint, name string, email string) (string, error) {
	claims := &JWTClaims{
		UserId:    userId,
		Name:      name,
		Email:     email,
		TokenType: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
