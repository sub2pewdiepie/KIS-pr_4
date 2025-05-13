package main

import (
	"fmt"
	"log"
	"net/http"
	"pr_4/server"
)

func main() {
	s := server.NewServer(100)

	http.HandleFunc("/ws", s.WebSocketHandler)

	fmt.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
