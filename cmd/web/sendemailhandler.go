package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/hculpan/areweplaying/cmd/web/templates"
	"github.com/hculpan/areweplaying/pkg/email"
)

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	if !ValidateJWTFromCookie(r, appKey) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	to := r.FormValue("to")
	from := "culpanvtt@gmail.com"
	subject := r.FormValue("subject")
	body := r.FormValue("text")
	toArray := strings.Split(to, ",")
	for i := range toArray {
		toArray[i] = strings.Trim(toArray[i], " \t\n\r")
	}
	e := email.NewEmail(toArray, from, subject, body)

	// log.Default().Printf("got email: To: %s, From: %s, Subject: %q, Body: %q", to, from, subject, body)
	err := email.SendEmail(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	comp := templates.SendEmail("Email sent!", "")
	comp.Render(context.Background(), w)
}
