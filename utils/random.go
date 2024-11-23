package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

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

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
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
	currencies := []string{"NGN", "USD", "EUR", "GBP"}
	k := len(currencies)
	return currencies[rand.Intn(k)]
}

// RandomEmail generates a random email address
func RandomEmail() string {
	tlds := []string{"gmail.com", "example.com", "yahoo.com", "danielokoronkwo.tech"}
	k := len(tlds)
	tld := tlds[rand.Intn(k)]
	s := RandomString(10)

	return fmt.Sprintf("%v@%v", s, tld)
}
