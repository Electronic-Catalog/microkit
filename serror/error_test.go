package serror

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestError_Error(t *testing.T) {
	e := New().WithMessage("test.test.error").WithTrace("test1")

	require.True(t, Equals(e, e))
	et, ok := Is(e)
	require.True(t, ok)
	t.Log(et)
}
