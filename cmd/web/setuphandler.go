package main

import (
	"context"
	"net/http"

	"github.com/hculpan/areweplaying/cmd/web/templates"
	"github.com/hculpan/areweplaying/pkg/data"
)

func setupHandler(w http.ResponseWriter, r *http.Request) {
	session, dataErr := data.ReadSessionData()
	if dataErr != nil {
		http.Error(w, dataErr.Error(), http.StatusInternalServerError)
		return
	}

	if !ValidateJWTFromCookie(r, appKey) {
		r.ParseForm()
		pword := r.FormValue("password")
		if len(pword) == 0 || pword != session.Password {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		GenerateJWTAndStoreInCookie(w, appKey)
	}

	comp := templates.Setup(session)
	comp.Render(context.Background(), w)
}
