package dgrijalva

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	testValidSecret   = "some-secret-value"
	testInvalidSecret = "some-different-secret-value"
)

func TestDgrijalva_GenerateAccessToken(t *testing.T) {
	claim := map[string]interface{}{
		"key-1":        1,
		"key-2-string": "some-data",
	}

	t.Log("-- original claim --", claim)

	token, expiration, err := GenerateToken(testValidSecret, claim, 2*time.Second)
	require.NoError(t, err, "generate token")
	t.Log(token)
	t.Log(expiration)

	cl, err := ValidateToken(testValidSecret, token)
	require.NoError(t, err, "validate token")
	t.Log("-- verified token claim --", string(cl))

	cl, err = ValidateToken(testInvalidSecret, token)
	require.Error(t, err, "validated with invalid secret !")
	t.Log(err)
	t.Log("-- unverified token claim --", string(cl))
}
