package dgrijalva

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateCode
// return an integer code with given len
func GenerateCode(len int) string {
	code := ""
	for i := 0; i < len; i++ {
		x := rand.Intn(10)
		code = fmt.Sprintf("%s%d", code, x)
	}
	return code
}

const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateForgetPasswordCode
// generate a code contains digit and texts
func GenerateForgetPasswordCode(length int) string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
