package jwt

import (
	"time"

	"github.com/Mosteben/hotel-booking-system/configs"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`

	Role string `json:"role"`

	jwt.RegisteredClaims
}

func GenerateToken(userID, role string) (string, error) {

	expire, err := time.ParseDuration(
		configs.GetEnv("JWT_EXPIRE"),
	)

	if err != nil {
		expire = 24 * time.Hour
	}

	claims := Claims{

		UserID: userID,

		Role: role,

		RegisteredClaims: jwt.RegisteredClaims{

			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(expire),
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

	return token.SignedString(
		[]byte(configs.GetEnv("JWT_SECRET")),
	)
}
func ParseToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(

		tokenString,

		&Claims{},

		func(token *jwt.Token) (interface{}, error) {

			return []byte(
				configs.GetEnv("JWT_SECRET"),
			), nil
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