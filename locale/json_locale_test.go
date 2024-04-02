package locale

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewFromDir(t *testing.T) {
	l, err := NewFromDir("testdata")
	require.NoError(t, err)
	require.NotNil(t, l)
	require.Equal(t, "Selam", l.Message(Turkish, "sample_message"))
	require.Equal(t, "Selam Ali", l.Message(Turkish, "formatted_message", "Ali"))
}
