package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	min := int64(1)
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