package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hculpan/areweplaying/cmd/web/templates"
	"github.com/hculpan/areweplaying/pkg/data"
)

func sendReminderHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidateJWTFromCookie(r, appKey) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	session, err := data.ReadSessionData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	to := buildToList(session)
	subject := fmt.Sprintf("Game Reminder for %s", session.Date)
	text := fmt.Sprintf(
		`This is a game reminder for %s. Please respond with a "yes" or "no" to indicate
if you plan to attend.

If you wish to check on the status of the game, you may do so at https://awp.culpanvtt.org.

Thank you,
Your friendly "Are We Playing?" automated system`, session.Date)

	now := time.Now()
	sessionDate, err := time.Parse(data.DATE_FORMAT, session.Date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := ""
	if now.After(sessionDate) {
		msg = "Session date is in the past."
	}

	comp := templates.SendReminder(to, subject, text, msg)
	comp.Render(context.Background(), w)
}

func buildToList(session data.Session) string {
	emails := []string{}
	for _, player := range session.Players {
		if len(player.Email) > 0 {
			emails = append(emails, player.Email)
		}
	}

	return strings.Join(emails, ", ")
}
