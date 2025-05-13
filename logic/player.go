package logic

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Player struct {
	ID    string
	Conn  *websocket.Conn
	Name  string
	Moves int
}

func NewPlayer(Conn *websocket.Conn, Name string) *Player {

	userID := uuid.New().String()

	return &Player{
		Conn:  Conn,
		Name:  Name,
		ID:    userID,
		Moves: 0,
	}
}
