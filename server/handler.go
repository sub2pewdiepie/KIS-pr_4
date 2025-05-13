package server

import (
	"log"
	"net/http"
	"pr_4/logic"

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
	return &Server{
		GameManager: logic.NewGameSessionManager(maxSessions),
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
			conn.WriteJSON(map[string]string{
				"status":    "joined",
				"sessionID": session.ID,
			})

		case "start":
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
			s.handleGuess(conn, &msg)

		default:
			conn.WriteJSON(map[string]string{"error": "неизвестная команда"})
		}
	}

}
