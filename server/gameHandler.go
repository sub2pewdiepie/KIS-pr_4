package server

import (
	"pr_4/logic"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// func (s *Server) handleGuess(conn *websocket.Conn, msg *Message) {
// 	sessionID := msg.SessionID
// 	guess := msg.Guess
// 	playerName := msg.Name

// 	// Найти сессию по sessionID
// 	session, exists := s.GameManager.Sessions[sessionID]
// 	if !exists {
// 		conn.WriteJSON(map[string]string{"error": "сессия не найдена"})
// 		return
// 	}

// 	if !session.Started || session.GameOver {
// 		conn.WriteJSON(map[string]string{"error": "игра ещё не началась или уже завершена"})
// 		return
// 	}

// 	// Проверка попытки
// 	black, white := logic.CheckGuess(session.SecretCode, guess)

// 	// Отправка результата игроку
// 	conn.WriteJSON(map[string]interface{}{
// 		"status": "guess_result",
// 		"guess":  guess,
// 		"black":  black,
// 		"white":  white,
// 	})

// 	// Проверка на победу
// 	if black == 4 {
// 		session.GameOver = true
// 		session.Winner = playerName
// 		session.EndTime = time.Now()

// 		// Уведомление всех игроков о завершении игры
// 		for c := range session.Players {
// 			c.WriteJSON(map[string]string{
// 				"status":  "game_over",
// 				"winner":  playerName,
// 				"message": playerName + " угадал код!",
// 			})
// 		}
// 	}
// }

func (s *Server) handleGuess(conn *websocket.Conn, msg *Message) {
	sessionID := msg.SessionID
	guess := msg.Guess

	// Найти сессию по sessionID
	session, exists := s.GameManager.Sessions[sessionID]
	if !exists {
		conn.WriteJSON(map[string]string{"error": "сессия не найдена"})
		return
	}

	if !session.Started || session.GameOver {
		conn.WriteJSON(map[string]string{"error": "игра ещё не началась или уже завершена"})
		return
	}

	// Получить игрока по соединению
	player, ok := session.Players[conn]
	if !ok {
		conn.WriteJSON(map[string]string{"error": "игрок не найден в сессии"})
		return
	}
	playerName := player.Name

	// Проверка попытки
	black, white := logic.CheckGuess(session.SecretCode, guess)

	// Отправка результата игроку
	player.Moves++
	conn.WriteJSON(map[string]interface{}{
		"status": "guess_result",
		"guess":  guess,
		"black":  black,
		"white":  white,
		"moves":  player.Moves,
	})
	if player.Moves > 30 {
		s.SendMSG2ALL(session, playerName, "ихрасходовал все попытки")
	} else if black == 4 {
		session.GameOver = true
		session.Winner = playerName
		session.EndTime = time.Now()

		// Уведомление всех игроков о завершении игры
		for c := range session.Players {
			c.WriteJSON(map[string]string{
				"status":  "game_over",
				"winner":  playerName,
				"message": playerName + " угадал код за ",
				"moves":   "за " + strconv.Itoa(player.Moves),
			})
		}
	}
}

func (s *Server) SendMSG2ALL(session *logic.GameSession, playerName string, msg string) {
	for c := range session.Players {
		c.WriteJSON(map[string]string{
			"message": msg,
			"player":  playerName,
		})
	}
}
