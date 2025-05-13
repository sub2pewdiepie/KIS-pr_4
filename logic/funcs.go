package logic

import (
	"log"
	"math/rand"
	"time"
)

const codeLength = 4

var symbols = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateSecretCode() string {
	rand.Seed(time.Now().UnixNano())
	code := make([]rune, codeLength)
	for i := range code {
		// code[i] = symbols[rand.Intn(len(symbols))]
		code[i] = symbols[i]
	}
	log.Println("Ошибка чтения JSON:", string(code))
	return string(code)
}

func CheckGuess(secret, guess string) (black, white int) {
	secretRunes := []rune(secret)
	guessRunes := []rune(guess)

	usedSecret := make([]bool, len(secretRunes)) // Помечаем использованные символы в secret
	usedGuess := make([]bool, len(guessRunes))   // Помечаем использованные символы в guess

	// Сначала ищем "чёрные" маркеры
	for i := range secretRunes {
		if secretRunes[i] == guessRunes[i] {
			black++
			usedSecret[i] = true
			usedGuess[i] = true
		}
	}

	// Теперь ищем "белые" маркеры
	for i := range guessRunes {
		if usedGuess[i] {
			continue
		}
		for j := range secretRunes {
			if !usedSecret[j] && guessRunes[i] == secretRunes[j] {
				white++
				usedSecret[j] = true
				break
			}
		}
	}

	return
}
