package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hculpan/areweplaying/cmd/web/templates"
	"github.com/hculpan/areweplaying/pkg/data"
)

func saveSetupHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nextSession := r.FormValue("nextSession")
	incDays := r.FormValue("incDays")
	sessionStatus := r.FormValue("sessionStatus")
	log.Default().Printf("sessionStatus = %s", sessionStatus)

	session, dataErr := data.ReadSessionData()
	if dataErr != nil {
		http.Error(w, dataErr.Error(), http.StatusInternalServerError)
		return
	}

	msg := "Save successful!"
	err := updateSession(session, incDays, nextSession, sessionStatus)
	if err != nil {
		msg = fmt.Sprintf("Error in update: %s", err.Error())
	}

	comp := templates.SaveSetup(msg)
	comp.Render(context.Background(), w)
}

func updateSession(session data.Session, incDays, nextSession, sessionStatus string) error {
	if val, err := strconv.Atoi(incDays); err == nil {
		session.IncrementDays = val
	} else {
		return err
	}

	if nextSession == "nextSession" {
		if d, err := time.Parse(data.DATE_FORMAT, session.Date); err == nil {
			d = d.AddDate(0, 0, session.IncrementDays)
			session.Date = d.Format(data.DATE_FORMAT)
		} else {
			return err
		}

		for i, p := range session.Players {
			if p.Gm {
				session.Players[i].Attending = "yes"
			} else {
				session.Players[i].Attending = "unknown"
			}
		}
	}

	if len(sessionStatus) > 0 {
		session.Status = sessionStatus
	}

	return data.PersistSession(session)
}
