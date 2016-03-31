package core

import "math/rand"

func RandString(n int) string {
	runes := []rune(" \n\tabcdef\n\t123456\n\t")
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Int63()%int64(len(runes))]
	}
	return string(b)
}

func SyntaxHighlighting() bool {
	if Ed != nil {
		return Ed.Config().SyntaxHighlighting
	}
	return false
}
