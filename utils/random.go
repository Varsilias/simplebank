package utils

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const numbers = "1234567890"
const uppercaseAlphabets = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const specialCharacters = "!@#$%^&*()"

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generate a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generate a random string of length n
func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for range n {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomPassword(n int) string {
	var sb strings.Builder
	minLength := 8

	if n < minLength {
		n = minLength
	}

	l := len(alphabet)
	u := len(uppercaseAlphabets)
	s := len(specialCharacters)
	nu := len(numbers)

	for range int(math.Round(float64(n / 2))) {
		sb.WriteByte(uppercaseAlphabets[rand.Intn(u)])
		sb.WriteByte(specialCharacters[rand.Intn(s)])
		sb.WriteByte(numbers[rand.Intn(nu)])
	}

	for range int(math.Round(float64(n / 2))) {
		sb.WriteByte(alphabet[rand.Intn(l)])
	}

	return sb.String()
}

// RandomName generates a random name
func RandomName() string {
	return RandomString(8)
}

// RandomAmount generates a random amount of money
func RandomAmount() int64 {
	return RandomInt(100, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{NGN, GBP, EUR, USD}
	k := len(currencies)
	return currencies[rand.Intn(k)]
}

// RandomEmail generates a random email address
func RandomEmail() string {
	tlds := []string{"gmail.com", "example.com", "yahoo.com", "danielokoronkwo.com"}
	k := len(tlds)
	tld := tlds[rand.Intn(k)]
	s := RandomString(10)

	return fmt.Sprintf("%v@%v", s, tld)
}
