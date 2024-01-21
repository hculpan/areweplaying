package data

const (
	DATE_FORMAT string = "Monday Jan 2, 2006 3:04pm"
)

type Player struct {
	Name      string `json:"name"`
	Attending string `json:"attending"`
	ToggleUrl string `json:"-"`
	Email     string `json:"email"`
	Gm        bool   `json:"gm"`
}

type Session struct {
	Date          string   `json:"date"`
	Players       []Player `json:"players"`
	IncrementDays int      `json:"inc_days"`
	Status        string   `json:"status"`
	Password      string   `json:"password"`
}

func (s *Session) GetPlayerByEmail(addr string) *Player {
	for i, p := range s.Players {
		if addr == p.Email {
			return &s.Players[i]
		}
	}

	return nil
}
