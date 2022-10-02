package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	min := int64(2)
	max := int64(2)
	result := RandomInt(min, max)
	require.Equal(t, min, result)
}

func TestRandomString(t *testing.T) {
	for i := 1; i < 10; i++ {
		str := RandomString(i)
		require.Equal(t, i, len(str))
		require.IsType(t, "", str)
	}
}

func TestRandomEmail(t *testing.T) {
	email := RandomEmail()
	require.Contains(t, email, "@")
	require.Contains(t, email, ".com")
	require.Less(t, len(email), 203)
}

func TestRandomGender(t *testing.T) {
	gender := RandomGender()
	require.Contains(t, [4]string{"male", "female", "other", "unknown"}, gender)
}
