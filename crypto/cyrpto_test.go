package crypto

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	plainText = "Hello GCM GCM GCM"
)

func TestEncryptGCM(t *testing.T) {
	aesKey, _ := base64.RawStdEncoding.DecodeString("4Iie50E1fnMNL9ikbFI1Fy0fIT4FUHzLcNqzQGJ8pUs")
	iv, _ := base64.RawStdEncoding.DecodeString("lL0rgPl6quWApwdq")

	cipher, err := EncryptGCM([]byte(plainText), aesKey, iv)
	require.NoError(t, err)
	require.Equal(t, "eb0937f6ee417088c785ee65e3678b2ae0f54d00b0b54c4c9763822c3821af795c", hex.EncodeToString(cipher))
}

func TestDecryptGCM(t *testing.T) {
	aesKey, _ := base64.RawStdEncoding.DecodeString("4Iie50E1fnMNL9ikbFI1Fy0fIT4FUHzLcNqzQGJ8pUs")
	iv, _ := base64.RawStdEncoding.DecodeString("lL0rgPl6quWApwdq")

	cipherBytes, _ := hex.DecodeString("eb0937f6ee417088c785ee65e3678b2ae0f54d00b0b54c4c9763822c3821af795c")
	plain, err := DecryptGCM(cipherBytes, aesKey, iv)
	require.NoError(t, err)
	require.Equal(t, plainText, string(plain))
}
