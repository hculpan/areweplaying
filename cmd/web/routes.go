package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hculpan/areweplaying/cmd/web/templates"
	"github.com/hculpan/areweplaying/pkg/data"
)

func routes(r *chi.Mux) {
	r.Post("/save-setup", func(w http.ResponseWriter, r *http.Request) {
		saveSetupHandler(w, r)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		comp := templates.MainPage()
		comp.Render(context.Background(), w)
	})

	r.Post("/setup", func(w http.ResponseWriter, r *http.Request) {
		setupHandler(w, r)
	})

	r.Get("/setup", func(w http.ResponseWriter, r *http.Request) {
		setupHandler(w, r)
	})

	r.Post("/send-email", func(w http.ResponseWriter, r *http.Request) {
		sendEmailHandler(w, r)
	})

	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		if ValidateJWTFromCookie(r, appKey) {
			http.Redirect(w, r, "/setup", http.StatusSeeOther)
		}
		comp := templates.Login()
		comp.Render(context.Background(), w)
	})

	r.Get("/send-reminder", func(w http.ResponseWriter, r *http.Request) {
		sendReminderHandler(w, r)
	})

	r.Get("/sessioninfo", func(w http.ResponseWriter, r *http.Request) {
		session, dataErr := data.ReadSessionData()
		if dataErr != nil {
			http.Error(w, dataErr.Error(), http.StatusInternalServerError)
			return
		}

		playerName := r.URL.Query().Get("playerName")
		playerAttending := r.URL.Query().Get("playerAttending")

		if len(playerName) > 0 && len(playerAttending) > 0 {
			if err := updatePlayerAttending(playerName, playerAttending, session); err != nil {
				log.Default().Println(err.Error())
			}
		}

		comp := templates.PageBody(session)
		comp.Render(context.Background(), w)
	})

}

func ValidateJWTFromCookie(r *http.Request, key []byte) bool {
	// Retrieve the JWT token from the cookie
	cookie, err := r.Cookie("jwt_token")
	if err != nil {
		// Cookie not found
		return false
	}

	tokenString := cookie.Value

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm used for the token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return key, nil
	})

	if err != nil {
		return false
	}

	// Check if token is valid
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Here you can also add more checks on claims if needed
		// For example, check the issuer, subject, or expiration
		return true
	}

	return false
}

func GenerateJWTAndStoreInCookie(w http.ResponseWriter, key []byte) error {
	// Create a new token object, specifying signing method and the claims you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(90 * 24 * time.Hour).Unix(), // 90 days expiration
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(key)
	if err != nil {
		return err
	}

	// Create a cookie
	cookie := http.Cookie{
		Name:     "jwt_token",
		Value:    tokenString,
		Expires:  time.Now().Add(90 * 24 * time.Hour),
		HttpOnly: true, // HttpOnly if you want to prevent access to the cookie from JavaScript
	}

	// Set the cookie
	http.SetCookie(w, &cookie)

	return nil
}

func updatePlayerAttending(name, attending string, session data.Session) error {
	for i, p := range session.Players {
		if p.Name == name {
			switch strings.ToLower(attending) {
			case "yes", "no":
				session.Players[i].Attending = strings.ToLower(attending)
				return data.PersistSession(session)
			case "toggle":
				newAttending := session.Players[i].Attending
				switch newAttending {
				case "unknown", "no":
					newAttending = "yes"
				case "yes":
					newAttending = "no"
				}
				session.Players[i].Attending = newAttending
				return data.PersistSession(session)
			default:
				return fmt.Errorf("unrecognized attending indicator %q", attending)
			}
		}
	}

	return fmt.Errorf("player %q not found", name)
}
