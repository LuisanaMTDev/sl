package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/LuisanaMTDev/spaced_learning/server/database/gosql_queries"
	"github.com/LuisanaMTDev/spaced_learning/server/frontend/views"
	"github.com/rs/zerolog/log"
)

func (sc *ServerConfig) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := sc.OAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Error().AnErr("error", err).Msg("While exchanging code for token with oauth2.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Debug().Time("expiry", token.Expiry).Int64("expires_in", token.ExpiresIn).Msg("completed oauth exchange")

	// Store the access token and refresh token in in-memory session storage.
	err = sc.DBQueries.AddUser(r.Context(), gosql_queries.AddUserParams{
		AccessToken: sql.NullString{Valid: true, String: token.AccessToken},
		Showed:      0,
	})
	if err != nil {
		log.Error().AnErr("error", err).Msg("While saving tokens to db.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// TODO: Add logout.
func (sc *ServerConfig) Login(w http.ResponseWriter, r *http.Request) {
	client := strings.ToUpper(r.Header.Get("SL-Client-Type"))

	switch client {
	case "SL-CLI":
		_, err := sc.DBQueries.GetAccessToken(r.Context())
		if err != nil && err.Error() == "sql: no rows in result set" {
			if err := json.NewEncoder(w).Encode(struct {
				HasAccessToken bool   `json:"has_access_token"`
				LoginIn        string `json:"login_in"`
			}{
				HasAccessToken: false,
				LoginIn:        "/",
			}); err != nil {
				log.Fatal().AnErr("error", err).Msg("ERROR: while encoding response body")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		} else if err != nil {
			log.Fatal().AnErr("error", err).Msg("ERROR: while getting access token from the db")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		apiKey, err := generateApiKey()
		if err != nil {
			log.Fatal().AnErr("error", err).Msg("ERROR: while generating API key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		hashedAPIKey, err := hashApiKey(apiKey)
		if err != nil {
			log.Fatal().AnErr("error", err).Msg("ERROR: while hashing API key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = sc.DBQueries.UpdateAPIKEY(r.Context(), sql.NullString{Valid: true, String: hashedAPIKey})
		if err != nil {
			log.Fatal().AnErr("error", err).Msg("ERROR: while saving hashed API key")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(struct {
			APIKey string `json:"api_key"`
		}{APIKey: apiKey}); err != nil {
			log.Fatal().AnErr("error", err).Msg("ERROR: while encoding response body")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	case "SL-WEB-APP":
		r.ParseForm()

		if r.FormValue("password") == os.Getenv("LOGIN_PASSWORD") {
			if err := views.AuthButton().Render(r.Context(), w); err != nil {
				log.Fatal().AnErr("error", err).Msg("ERROR: while sending main page")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	default:
		w.WriteHeader(http.StatusForbidden)
		return
	}
}
