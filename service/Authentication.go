package service

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService interface {
	GenerateToken(email string, isUser bool) string
	ValidateToken(token string) (*jwt.Token, error)
}

type authCustomClaims struct {
	Name string `json:"name"`
	User bool   `json:"user"`
	jwt.StandardClaims
}

type jwtServices struct {
	secretKey string
	issuer    string
}

func AuthService() JWTService {
	return &jwtServices{
		secretKey: getSecretKey(),
		issuer:    "Survey-api",
	}
}

func getSecretKey() string {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "Survey-secret-signature-key"
	}
	return secret
}

func (service *jwtServices) GenerateToken(email string, isUser bool) string {
	location, _ := time.LoadLocation("Asia/Bangkok")
	claims := &authCustomClaims{
		email,
		isUser,
		jwt.StandardClaims{
			ExpiresAt: time.Now().In(location).Add(time.Second * 20).Unix(),
			Issuer:    service.issuer,
			IssuedAt:  time.Now().In(location).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	endcodedToken, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return endcodedToken
}

func (service *jwtServices) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("invalid token: %v", token.Header["alg"])
		}

		return []byte(service.secretKey), nil
	})

}
