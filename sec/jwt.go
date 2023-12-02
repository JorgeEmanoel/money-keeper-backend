package sec

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getPrivateKeyFromFile() (*rsa.PrivateKey, error) {
	privateKey, err := os.ReadFile(os.Getenv("JWT_KEY_FILE"))
	if err != nil {
		log.Printf("Failed to load JWT_KEY_FILE: %s", os.Getenv("JWT_KEY_FILE"))
		return nil, err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		log.Printf("error parsing RSA private key: %v\n", err)
		return nil, err
	}

	return key, nil
}

func getHMACSecret() []byte {
	return []byte(os.Getenv("JWT_HMAC_SECRET"))
}

func JWTFromPayload(payload map[string]any, ttl time.Duration) (string, error) {
	privateKey := getHMACSecret()

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["dat"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Printf("Fail while create JWT string: %v", err)
		return "", err
	}

	return tokenString, nil
}

func JWTToPayload(tokenString string) (map[string]any, error) {
	privateKey := getHMACSecret()

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		log.Printf("Method: %v", token.Method)
		log.Printf("privateKey: %v", privateKey)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("failed to get token claims")
	}

	payload := claims["dat"]

	return payload.(map[string]any), nil
}
