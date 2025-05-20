package server

import (
	"log"
	"net/http"
	"pr_4/logic"
	"time"

	"github.com/gorilla/websocket"
)

type Server struct {
	GameManager *logic.GameSessionManager
	Upgrader    websocket.Upgrader
}

type Message struct {
	Command    string `json:"command"`
	Name       string `json:"name,omitempty"`
	SessionID  string `json:"sessionID,omitempty"`
	Guess      string `json:"guess,omitempty"`
	MaxPlayers int    `json:"maxPlayers,omitempty"`
}

func NewServer(maxSessions int) *Server {
	gameManager := logic.NewGameSessionManager(maxSessions)
	gameManager.StartCleanupRoutine(time.Minute, 30*time.Minute)
	return &Server{
		GameManager: gameManager,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *Server) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка апгрейда WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("Ошибка чтения JSON:", err)
			break
		}

		switch msg.Command {
		case "create":
			session, err := s.GameManager.CreateSession(conn, msg.Name, msg.MaxPlayers)
			if err != nil {
				conn.WriteJSON(map[string]string{"error": err.Error()})
				continue
			}
			session.UpdateActivity()
			conn.WriteJSON(map[string]string{
				"status":    "created",
				"sessionID": session.ID,
			})

		case "join":
			session, err := s.GameManager.JoinSession(msg.SessionID, conn, msg.Name)
			if err != nil {
				conn.WriteJSON(map[string]string{"error": err.Error()})
				continue
			}
			session.UpdateActivity()
			conn.WriteJSON(map[string]string{
				"status":    "joined",
				"sessionID": session.ID,
			})

		case "start":
			//Обновление времени в SartGame
			err := s.GameManager.StartGame(msg.SessionID)
			if err != nil {
				conn.WriteJSON(map[string]string{"error": err.Error()})
				continue
			}
			conn.WriteJSON(map[string]string{
				"status":    "started",
				"sessionID": msg.SessionID,
			})

		case "guess":
			//Обновление времени в handleGuess
			err = s.handleGuess(conn, &msg)
			if err != nil {
				conn.WriteJSON(map[string]string{"error": err.Error()})
				continue
			}

		default:
			conn.WriteJSON(map[string]string{"error": "неизвестная команда"})
		}
	}

}
