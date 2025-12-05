package common

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJTWToken(claims jwt.MapClaims, jwtSecret string, jwtExpiresIn time.Duration) (string, error) {
	// Compute token expiration
	exp := time.Now().Add(jwtExpiresIn)
	claims["exp"] = jwt.NewNumericDate(exp)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	singedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed generating JWT token: %w", err)
	}

	return singedToken, nil
}

func ParseJWTToken(tokenString string, jwtSecret string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
