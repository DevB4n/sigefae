package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Secret = []byte("CAMBIAR_ESTO_DESPUES")

type Claims struct {
	UserID uint `json:"user_id"`
	RolID  uint `json:"rol_id"`

	jwt.RegisteredClaims
}

func GenerateToken(userID uint, rolID uint) (string, error) {

	claims := Claims{
		UserID: userID,
		RolID:  rolID,

		RegisteredClaims: jwt.RegisteredClaims{

			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),

			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString(Secret)
}

func ValidateToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return Secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}