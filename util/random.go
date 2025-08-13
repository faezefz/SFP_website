package util

import (
	"math/rand"
	"os"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
const domain = "@gmail.com"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for range n {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomEmail() string {
	return RandomString(10) + domain
}

func RandomName() string {
	firstName := RandomString(5)
	lastName := RandomString(7)
	return firstName + " " + lastName
}

func RandomPassword() string {
	var sb strings.Builder
	k := len(alphabet)

	for range 6 {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	for range 2 {
		c := strings.ToUpper(string(alphabet[rand.Intn(k)]))
		sb.WriteByte(c[0])
	}

	for range 2 {
		sb.WriteByte(byte(rand.Intn(10) + '0'))
	}

	return sb.String()
}

func ReadCSV() ([]byte, error) {
	// خواندن محتوای فایل CSV به صورت بایت
	fileContent, err := os.ReadFile("../../util/testfile.csv") // مسیر فایل CSV
	if err != nil {
		return nil, err // در صورت بروز خطا، آن را برمی‌گرداند
	}

	// بازگشت محتویات فایل به صورت بایت
	return fileContent, nil
}
