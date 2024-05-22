package util

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SecretKey = os.Getenv("SECRET")

func GenerateJwt(issuer string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	return claims.SignedString([]byte(SecretKey))

}

// validate the token
func Parsejwt(signedToken string) (string, error) {
	token, err := jwt.ParseWithClaims(signedToken, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil || !token.Valid {
		log.Println("error validating the token")
		return "", err
	}

	//asserts that the type of token.Claims is *SignedDetails.If successful, it assigns the claims to the claims variable, and ok is set to true
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		msg := fmt.Errorf("the token is invalid")
		log.Println(err)
		return "", msg

	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg := fmt.Errorf("the token is expired")
		return "", msg
	}
	return claims.Issuer, nil

}
