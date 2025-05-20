package logic

import (
	"errors"
	"time"

	"github.com/gorilla/websocket"
)

type GameSessionInterface interface {
	SetOwner(owner *Player) error
}

// GameSession представляет собой одну игровую сессию
type GameSession struct {
	ID string
	// Players     map[*websocket.Conn]string // Список игроков (сокет -> имя)
	Players      map[*websocket.Conn]*Player
	PlayerCount  int       // Текущее количество игроков
	MaxPlayers   int       // Максимум игроков
	SecretCode   string    // Секретный код
	Started      bool      // Флаг начала игры
	GameOver     bool      // Флаг завершения игры
	Winner       string    // Имя победителя
	StartTime    time.Time // Время начала игры
	EndTime      time.Time // Время завершения игры
	Owner        *Player
	LastActivity time.Time //Время последнего действия
}

// Новый конструктор GameSession
func NewGameSession(sessionID string, maxPlayers int) *GameSession {
	// Генерация случайного секретного кода
	secretCode := GenerateSecretCode()

	// Создание новой сессии
	session := &GameSession{
		ID:           sessionID,
		Players:      make(map[*websocket.Conn]*Player),
		PlayerCount:  0,
		MaxPlayers:   maxPlayers,
		SecretCode:   secretCode,
		Started:      false,
		GameOver:     false,
		StartTime:    time.Time{},
		EndTime:      time.Time{},
		Owner:        nil,
		LastActivity: time.Now(),
	}

	return session
}

func (session *GameSession) SetOwner(owner *Player) error {
	if session.Owner == nil {
		session.Owner = owner
		return nil
	} else {
		return errors.New("кто-то хочет сменить владельца")
	}
}

func (session *GameSession) UpdateActivity() {
	session.LastActivity = time.Now()
}
