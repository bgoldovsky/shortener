package generator

import (
	crypto "crypto/rand"
	"fmt"
	math "math/rand"
	"time"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	math.Seed(time.Now().Unix())
}

type generator struct{}

func NewGenerator() *generator {
	return &generator{}
}

// RandomString генерирует строку случайных символов
func (g *generator) RandomString(n int64) (string, error) {
	bytes := make([]byte, n)
	if _, err := crypto.Read(bytes); err != nil {
		return "", fmt.Errorf("random string generation error: %w", err)
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes), nil
}
