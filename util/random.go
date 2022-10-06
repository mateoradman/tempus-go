package util

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
	username := RandomString(rand.Intn(100))
	domain := RandomString(rand.Intn(100))
	return fmt.Sprintf("%s@%s.com", username, domain)
}

func RandomGender() string {
	genders := []string{Male, Female, Other, Unknown}
	randomInt := RandomInt(0, int64(len(genders)-1))
	return genders[randomInt]
}
