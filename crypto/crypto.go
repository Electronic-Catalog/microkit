package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privkey, &privkey.PublicKey, nil
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, nil
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte, password []byte) (*rsa.PrivateKey, error) {
	var b []byte
	block, _ := pem.Decode(priv)
	if block == nil {
		b = priv
	} else {
		enc := x509.IsEncryptedPEMBlock(block)
		b = block.Bytes
		var err error
		if enc {
			b, err = x509.DecryptPEMBlock(block, password)
			if err != nil {
				return nil, err
			}
		}
	}

	if ke, err := x509.ParsePKCS1PrivateKey(b); err == nil {
		return ke, nil
	}

	var key *rsa.PrivateKey
	ke, err := x509.ParsePKCS8PrivateKey(b)
	if err != nil {
		return nil, err
	}

	switch ke.(type) {
	case *rsa.PrivateKey:
		key = ke.(*rsa.PrivateKey)
	default:
		return nil, fmt.Errorf("unknown private key format !! %T", ke)
	}

	return key, nil
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	var b []byte
	block, _ := pem.Decode(pub)
	if block == nil {
		b = pub
	} else {
		enc := x509.IsEncryptedPEMBlock(block)
		b = block.Bytes
		var err error
		if enc {
			b, err = x509.DecryptPEMBlock(block, nil)
			if err != nil {
				return nil, err
			}
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not a public key")
	}

	return key, nil
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, msg)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// DecryptAES decrypts data with aes key and initial vector
func DecryptAES(ciphertext []byte, aesKey []byte, initialVector []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	if (len(ciphertext) % aes.BlockSize) != 0 {
		return nil, errors.New("blocksize must be multipe of decoded message length")
	}

	cbc := cipher.NewCBCDecrypter(block, initialVector)
	cbc.CryptBlocks(ciphertext, ciphertext)

	unpad, err := UnpadAESBlockSize(ciphertext)
	if err != nil {
		return nil, err
	}

	return unpad, nil
}

// EncryptAES encrypts data with aes key and initial vector
func EncryptAES(plaintext []byte, aesKey []byte, initialVector []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	msg := PadAESBlockSize(plaintext)

	cbc := cipher.NewCBCEncrypter(block, initialVector)
	cbc.CryptBlocks(msg, msg)
	return msg, nil
}

func PadAESBlockSize(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func UnpadAESBlockSize(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, errors.New("unpad error. src is empty")
	}

	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}

func EncryptGCM(plainBytes []byte, aesKey []byte, initialVector []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCMWithNonceSize(block, len(initialVector))
	if err != nil {
		return nil, err
	}

	return aesGCM.Seal(nil, initialVector, plainBytes, nil), nil
}

func DecryptGCM(cipherBytes []byte, aesKey []byte, initialVector []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCMWithNonceSize(block, len(initialVector))
	if err != nil {
		return nil, err
	}

	return aesGCM.Open(nil, initialVector, cipherBytes, nil)
}
