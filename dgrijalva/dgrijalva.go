package dgrijalva

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

func hash(param string) string {
	var t string
	h := sha256.New()
	h.Write([]byte(param))
	t = hex.EncodeToString(h.Sum(nil))

	hasher := sha256.New()
	hasher.Write([]byte(t))
	t = hex.EncodeToString(hasher.Sum(nil))
	return t
}

func secureRandomBytes(length int) []byte {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil
	}
	return randomBytes
}

func GenerateHashToken(salt string) string {
	rs := ""
	rb := secureRandomBytes(4)
	if rb != nil {
		rs = string(rb)
	}

	return hash(salt + GenerateUUID() + rs)
}

func GenerateUUID() string {
	id := uuid.New()
	return id.String()
}

func GenerateToken(secretKey string, claims map[string]interface{}, expireDuration time.Duration) (accessToken string, expiration int64, err error) {
	to := jwt.New(jwt.SigningMethodHS256)
	expiration = time.Now().Add(expireDuration).Unix()

	cls := make(jwt.MapClaims)

	for key, value := range claims {
		cls[key] = value
	}

	cls["exp"] = expiration

	to.Claims = cls
	tokenString, err := to.SignedString([]byte(secretKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiration, nil
}

func ValidateToken(secretKey string, token string) ([]byte, error) {
	tkn, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !tkn.Valid {
		return nil, errors.New("invalid token")
	}

	return json.Marshal(tkn.Claims)
}

func ValidateWithClaim(secretKey string, token string, claim jwt.Claims) error {
	tkn, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !tkn.Valid {
		return errors.New("invalid token")
	}

	return nil
}
