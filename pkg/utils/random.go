package utils

import (
	"crypto/rand"
	"math/big"
)

func randomElement(list []string) (string, error) {
	maxIndex := big.NewInt(int64(len(list)))
	n, err := rand.Int(rand.Reader, maxIndex)
	if err != nil {
		return "", err
	}
	return list[n.Int64()], nil
}

func GenerateUsername() (string, error) {
	adjectives := []string{"Веселый", "Ленивый", "Сумасшедший", "Саркастичный", "Бодрый"}
	nouns := []string{"Кот", "Бобер", "Енот", "Хомяк", "Лемур"}

	adj, err := randomElement(adjectives)
	if err != nil {
		return "", err
	}

	noun, err := randomElement(nouns)
	if err != nil {
		return "", err
	}

	return adj + "_" + noun, nil
}
