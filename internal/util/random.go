package util

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/mateoradman/tempus/internal/types"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@%s.com", uuid.New().String(), uuid.New().String())
}

func RandomGender() string {
	genders := []string{types.Male, types.Female, types.Other, types.Unknown}
	randomInt := RandomInt(0, int64(len(genders)-1))
	return genders[randomInt]
}

func RandomLanguage() string {
	langs := []string{types.English, types.German, types.Croatian}
	randomInt := RandomInt(0, int64(len(langs)-1))
	return langs[randomInt]
}
