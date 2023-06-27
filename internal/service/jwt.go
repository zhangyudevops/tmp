package service

import (
	"github.com/dgrijalva/jwt-go"
	"pack/internal/model"
	"time"
)

type sJwt struct{}

func Jwt() *sJwt {
	return &sJwt{}
}

// GenerateToken generate jwt token
func (s *sJwt) GenerateToken(username string) (string, error) {
	claims := model.MyClaims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("user_key"))
}

// ParseToken parse jwt token
func (s *sJwt) ParseToken(tokenString string) (*model.MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("user_key"), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*model.MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// RefreshToken refresh jwt token
func (s *sJwt) RefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("user_key"), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*model.MyClaims); ok && token.Valid {
		return s.GenerateToken(claims.Username)
	}

	return "", err
}
