package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	UserID   string `json:"sub"`
	jwt.StandardClaims
}

type JWTService struct {
	jwtKey []byte
}

func NewJWTService(jwtKey string) *JWTService {
	return &JWTService{
		jwtKey: []byte(jwtKey),
	}
}

func (js *JWTService) GenerateJWT(email, username, userID string) (string, error) {

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Email:    email,
		Username: username,
		UserID:   userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(js.jwtKey)
}

func (js *JWTService) ValidateToken(signedToken string) error {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return js.jwtKey, nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed to parse and validate token '%s': %v", signedToken, err)
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return fmt.Errorf("couldn't parse claims from token '%s'", signedToken)
	}

	if time.Now().Unix() > claims.ExpiresAt {
		return fmt.Errorf("token '%s' expired", signedToken)
	}

	return nil
}

func (js *JWTService) ExtractUserIDFromToken(signedToken string) (string, error) {

	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return js.jwtKey, nil
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to extract user ID: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return "", errors.New("couldn't parse claims")
	}

	return claims.UserID, nil
}
