package controllers

import (
	"database/sql"
	"net/http"
	"os"

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

func (sc *ServerConfig) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.FormValue("password") == os.Getenv("LOGIN_PASSWORD") {
		views.AuthButton().Render(r.Context(), w)
	}
}
