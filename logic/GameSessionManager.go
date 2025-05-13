package logic

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// GameSessionManager управляет игровыми сессиями
type GameSessionManager struct {
	Sessions    map[string]*GameSession
	MaxSessions int
	mu          sync.Mutex
}

// NewGameSessionManager создаёт менеджер с заданным лимитом сессий
func NewGameSessionManager(maxSessions int) *GameSessionManager {
	return &GameSessionManager{
		Sessions:    make(map[string]*GameSession),
		MaxSessions: maxSessions,
	}
}

// CreateSession создаёт и регистрирует новую игровую сессию с авто-генерацией ID
func (manager *GameSessionManager) CreateSession(conn *websocket.Conn, playerName string, maxPlayers int) (*GameSession, error) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	if len(manager.Sessions) >= manager.MaxSessions {
		return nil, errors.New("достигнут лимит количества сессий")
	}

	// Генерация уникального ID
	var sessionID string
	for {
		sessionID = uuid.New().String()
		if _, exists := manager.Sessions[sessionID]; !exists {
			break
		}
	}
	owner := NewPlayer(conn, playerName)
	// Создание сессии
	session := NewGameSession(sessionID, maxPlayers)
	session.SetOwner(owner)
	session.Players[conn] = owner
	session.PlayerCount = 1

	// Добавление в менеджер
	manager.Sessions[sessionID] = session

	return session, nil
}

func (manager *GameSessionManager) JoinSession(sessionID string, conn *websocket.Conn, playerName string) (*GameSession, error) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	session, exists := manager.Sessions[sessionID]
	if !exists {
		return nil, errors.New("сессия не найдена")
	}

	if session.GameOver {
		return nil, errors.New("игра завершена")
	}

	if session.Started {
		return nil, errors.New("игра уже началась")
	}

	if session.PlayerCount >= session.MaxPlayers {
		return nil, errors.New("сессия заполнена")
	}

	session.Players[conn] = NewPlayer(conn, playerName)
	session.PlayerCount++

	return session, nil
}

func (manager *GameSessionManager) StartGame(sessionID string) error {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	session, exists := manager.Sessions[sessionID]
	if !exists {
		return errors.New("сессия не найдена")
	}

	if session.Started {
		return errors.New("игра уже началась")
	}

	if session.PlayerCount < 2 {
		return errors.New("недостаточно игроков для начала игры")
	}

	session.Started = true
	session.StartTime = time.Now()

	return nil
}

func (m *GameSessionManager) GetSession(sessionID string) (*GameSession, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.Sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("сессия с ID '%s' не найдена", sessionID)
	}

	return session, nil
}
