package data

import (
	"encoding/json"
	"os"
)

func ReadSessionData() (Session, error) {
	var session Session
	data, err := os.ReadFile("players.json")
	if err != nil {
		return session, err
	}
	err = json.Unmarshal(data, &session)

	if err == nil {
		for i, p := range session.Players {
			session.Players[i].ToggleUrl = "/sessioninfo?playerName=" + p.Name + "&playerAttending=toggle"
		}
	}

	return session, err
}

func PersistSession(session Session) error {
	// Convert the object to JSON
	jsonData, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}

	// Write JSON data to file
	err = os.WriteFile("players.json", jsonData, 0644)
	return err
}
