package dgrijalva

import "testing"

func TestGenerateCode(t *testing.T) {
	code := GenerateCode(5)
	t.Log(code)

	code = GenerateCode(7)
	t.Logf(code)

	code = GenerateCode(12)
	t.Logf(code)
}

func TestGenerateForgetPasswordCode(t *testing.T) {
	code := GenerateForgetPasswordCode(12)
	t.Log(code)

	code = GenerateForgetPasswordCode(4)
	t.Log(code)
}
