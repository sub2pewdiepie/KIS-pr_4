package logic

import (
	"testing"
)

func TestGenerateSecretCode(t *testing.T) {
	code := GenerateSecretCode()

	if len(code) != codeLength {
		t.Errorf("Expected code length of %d, but got %d", codeLength, len(code))
	}

	// Проверим, что в коде только допустимые символы
	for _, char := range code {
		if !isValidSymbol(char) {
			t.Errorf("Invalid character %c in generated code", char)
		}
	}
}

// Проверяем, что символ является допустимым (A-Z, 0-9)
func isValidSymbol(c rune) bool {
	return (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func TestCheckGuess(t *testing.T) {
	tests := []struct {
		secret string
		guess  string
		black  int
		white  int
	}{
		{"ABCD", "ABCD", 4, 0}, // Все символы на правильных позициях
		{"ABCD", "DCBA", 0, 4}, // Все символы на неправильных позициях
		{"ABCD", "AACD", 3, 0}, // Один символ на правильной позиции
		{"ABCD", "AABC", 1, 2}, // Два символа на правильных позициях, один на неправильной
		{"ABCD", "ABDC", 2, 2}, // Один символ на неправильной позиции
	}

	for _, tt := range tests {
		t.Run(tt.secret+"-"+tt.guess, func(t *testing.T) {
			black, white := CheckGuess(tt.secret, tt.guess)

			if black != tt.black {
				t.Errorf("Expected %d black markers, got %d", tt.black, black)
			}

			if white != tt.white {
				t.Errorf("Expected %d white markers, got %d", tt.white, white)
			}
		})
	}
}
