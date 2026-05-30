package utils

import(
	"github.com/golang-jwt/jwt/v5"

	"time"

	"recruitmentportal/config"
)

func GenerateToken(userID int, email string, role string) (string, error){
	claims := jwt.MapClaims{
		"user_id": userID,
		"email": email,
		"role": role,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(config.JWTSecret)
}