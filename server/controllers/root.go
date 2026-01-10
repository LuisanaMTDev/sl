package controllers

import (
	"net/http"
	"strings"

	"github.com/LuisanaMTDev/spaced_learning/server/frontend/views"
	"github.com/rs/zerolog/log"
)

// TODO: Add auth for the CLI.
func (sc *ServerConfig) Root(w http.ResponseWriter, r *http.Request) {
	client := strings.ToUpper(r.Header.Get("SL-Client-Type"))

	// TODO: Generate API key, send it to the user, and hash it and save it in the db.
	if client == "SL-CLI" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	_, err := sc.DBQueries.GetAccessToken(r.Context())
	if err != nil && err.Error() == "sql: no rows in result set" {
		err = views.Index(
			false,
			sc.OAuthConfig.Scopes[0],
			sc.OAuthConfig.ClientID,
		).Render(r.Context(), w)
		if err != nil {
			log.Fatal().AnErr("error", err).Msg("ERROR: while sending main page")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	} else if err != nil {
		log.Fatal().AnErr("error", err).Msg("ERROR: while getting access token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Debug().Msg("SUCCESS: db has access token")

	err = views.Index(
		true,
		"",
		"",
	).Render(r.Context(), w)
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("ERROR: while sending main page")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
