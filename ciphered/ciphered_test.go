package ciphered

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	_ "embed"
	"testing"
)

func TestEncode(t *testing.T) {
	input := "hello world"
	expected := "H4sIAAAAAAAA/8pIzcnJVyjPL8pJAQQAAP//hRFKDQsAAAA="
	result, err := Encode(input)
	if err != nil {
		t.Errorf("Encode failed with error: %v", err)
	}

	if result != expected {
		t.Errorf("Encode result doesn't match expected")
	}
}

func TestDecode(t *testing.T) {
	input := "H4sIAAAAAAAA/8pIzcnJVyjPL8pJAQQAAP//hRFKDQsAAAA="
	expected := "hello world"

	result, err := Decode(input)
	if err != nil {
		t.Errorf("Decode failed with error: %v", err)
	}

	if result != expected {
		t.Errorf("Decode result doesn't match expected")
	}
}

func TestEncodeDecode(t *testing.T) {
	credentials := "H4sIAAAAAAAA/8pIzcnJVyjPL8pJAQQAAP//hRFKDQsAAAA="
	result, err := Encode(credentials)
	if err != nil {
		t.Errorf("Encode failed with error: %v", err)
	}

	s, err := Decode(result)
	if err != nil {
		t.Errorf("Decode failed with error: %v", err)
	}

	if s != credentials {
		t.Errorf("Decode result doesn't match expected")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key32 := make([]byte, 32) // 32-byte key for AES-256
	_, err := rand.Read(key32)
	if err != nil {
		t.Errorf("Error generating key: %v", err)
	}

	const originalString = "This is a secret message."
	ciphertext, err := Encrypt(originalString, key32)
	if err != nil {
		t.Errorf("Error during encryption: %v", err)
	}

	decryptedString, err := Decrypt(ciphertext, key32)
	if err != nil {
		t.Errorf("Error during decryption: %v", err)
	}

	if originalString != decryptedString {
		t.Errorf("EncryptDecrypt: failed")
	}
}

func TestEncryptionDecryption(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{input: ""},                               // Empty string
		{input: "a"},                              // Single character
		{input: "abcdefgh"},                       // 8 characters
		{input: "abcdefghijklmno"},                // 15 characters
		{input: "abcdefghijklmnop"},               // 16 characters
		{input: "abcdefghijklmnopq"},              // 17 characters
		{input: "abcdefghijklmnopqrstuvwxy"},      // 31 characters
		{input: "abcdefghijklmnopqrstuvwxyz0123"}, // 32 characters
	}

	key := []byte("753c69838adc4b64b0c8ae694a9a273c")
	for _, tc := range testCases {
		encryptedText, err := EncryptFixed(tc.input, false, key)
		if err != nil {
			t.Errorf("Error encrypting '%s': %v", tc.input, err)
			continue
		}

		decryptedText, err := DecryptFixed(encryptedText, false, key)
		if err != nil {
			t.Errorf("Error decrypting '%s': %v", encryptedText, err)
			continue
		}

		if decryptedText != tc.input {
			t.Errorf("Decrypted text '%s' does not match original '%s'", decryptedText, tc.input)
		}
	}
}

func TestPadding(t *testing.T) {
	testCases := []struct {
		input string
	}{
		{input: ""},                               // Empty string
		{input: "a"},                              // Single character
		{input: "abcdefgh"},                       // 8 characters
		{input: "abcdefghijklmno"},                // 15 characters
		{input: "abcdefghijklmnop"},               // 16 characters
		{input: "abcdefghijklmnopq"},              // 17 characters
		{input: "abcdefghijklmnopqrstuvwxy"},      // 31 characters
		{input: "abcdefghijklmnopqrstuvwxyz0123"}, // 32 characters
	}

	for _, tc := range testCases {
		paddedPlaintext := PadPKCS7([]byte(tc.input), aes.BlockSize)
		unpaddedPlaintext := UnpadPKCS7(paddedPlaintext)

		if len(paddedPlaintext)%aes.BlockSize != 0 {
			t.Errorf("Padded plaintext length is not a multiple of block size for input '%s'", tc.input)
		}

		if !bytes.Equal(unpaddedPlaintext, []byte(tc.input)) {
			t.Errorf("Unpadded plaintext does not match original for input '%s'", tc.input)
		}
	}
}
