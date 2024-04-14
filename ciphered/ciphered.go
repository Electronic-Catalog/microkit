// Package ciphered provides functions for encoding and decoding strings using AES, gzip compression and base64 encoding.
package ciphered

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

// Encode the input string using gzip compression followed by base64 encoding,
// and returns the resulting encoded string.
func Encode(input string) (result string, err error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err = gw.Write([]byte(input))
	if err != nil {
		return
	}
	err = gw.Close()
	if err != nil {
		return
	}

	result = base64.StdEncoding.EncodeToString(buf.Bytes())
	return
}

// Decode the input string using base64 decoding followed by gzip decompression,
// and returns the resulting string. If any error occurs during decoding or decompression,
// it returns an empty string and the corresponding error.
func Decode(input string) (result string, err error) {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return
	}

	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return
	}
	defer zr.Close()

	var sb strings.Builder
	_, err = io.Copy(&sb, zr)
	result = sb.String()
	return
}

// Encrypt a variable-length string using AES-GCM
func Encrypt(plaintext string, key32 []byte) (ciphertext string, err error) {
	block, err := aes.NewCipher(key32)
	if err != nil {
		return
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return
	}

	b := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)
	b = append(nonce, b...)
	ciphertext = base64.StdEncoding.EncodeToString(b)
	return
}

// Decrypt a variable-length string using AES-GCM
func Decrypt(ciphertext string, key32 []byte) (plaintext string, err error) {
	s, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return
	}
	block, err := aes.NewCipher(key32)
	if err != nil {
		return
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	if len(ciphertext) < aesGCM.NonceSize() {
		err = errors.New("too short")
		return
	}

	nonce := s[:aesGCM.NonceSize()]
	s = s[aesGCM.NonceSize():]

	b, err := aesGCM.Open(nil, nonce, s, nil)
	if err != nil {
		return
	}
	plaintext = string(b)
	return
}

// EncryptFixed encrypts a fixed-length string using AES encryption.
func EncryptFixed(plainText string, random bool, fixKey []byte) (encodedCipherText string, err error) {
	block, err := aes.NewCipher(fixKey)
	if err != nil {
		return
	}

	iv := make([]byte, aes.BlockSize)
	if random {
		_, err = io.ReadFull(rand.Reader, iv)
		if err != nil {
			return
		}
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	paddedPlaintext := PadPKCS7([]byte(plainText), aes.BlockSize)
	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	if random {
		ciphertext = append(iv, ciphertext...)
	}
	encodedCipherText = base64.StdEncoding.EncodeToString(ciphertext)
	return
}

// DecryptFixed decrypts an encrypted string using AES encryption.
func DecryptFixed(encodedCipherText string, random bool, fixKey []byte) (plaintext string, err error) {
	combined, err := base64.StdEncoding.DecodeString(encodedCipherText)
	if err != nil {
		return
	}

	block, err := aes.NewCipher(fixKey)
	if err != nil {
		return
	}

	var iv, ciphertext []byte
	if random {
		iv = combined[:aes.BlockSize]
		ciphertext = combined[aes.BlockSize:]
	} else {
		iv = make([]byte, aes.BlockSize)
		ciphertext = combined
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(ciphertext, ciphertext)

	plaintextBytes := UnpadPKCS7(ciphertext)
	plaintext = string(plaintextBytes)
	return
}

func PadPKCS7(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	pad := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, pad...)
}

func UnpadPKCS7(data []byte) []byte {
	length := len(data)
	padding := int(data[length-1])
	return data[:length-padding]
}
