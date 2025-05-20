package logic

import (
	"encoding/xml"
	"os"
	"time"
)

type GameResult struct {
	XMLName    xml.Name     `xml:"GameResult"`
	StartTime  string       `xml:"StartTime"`
	EndTime    string       `xml:"EndTime"`
	SecretCode string       `xml:"SecretCode"`
	Winner     string       `xml:"Winner"`
	Players    []PlayerInfo `xml:"Players>Player"`
}

type PlayerInfo struct {
	Name  string `xml:"Name"`
	Moves int    `xml:"Moves"`
}

func SaveGameResultToXML(session *GameSession, filename string) error {
	result := GameResult{
		StartTime:  session.StartTime.Format(time.RFC3339),
		EndTime:    session.EndTime.Format(time.RFC3339),
		SecretCode: session.SecretCode,
		Winner:     session.Winner,
	}

	for _, player := range session.Players {
		result.Players = append(result.Players, PlayerInfo{
			Name:  player.Name,
			Moves: player.Moves,
		})
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	return encoder.Encode(result)
}
