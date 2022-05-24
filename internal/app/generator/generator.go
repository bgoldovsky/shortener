package generator

import (
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type generator struct {
}

func NewGenerator() *generator {
	return &generator{}
}

// ID Генерирует десятизначный строковый ID
func (g *generator) ID() string {
	return generate(10)
}

func generate(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
