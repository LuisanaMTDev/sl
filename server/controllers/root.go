package controllers

import (
	"net/http"

	"github.com/LuisanaMTDev/spaced_learning/server/frontend/views"
	"github.com/rs/zerolog/log"
)

func (sc *ServerConfig) Root(w http.ResponseWriter, r *http.Request) {

	_, err := sc.DBQueries.GetAccessToken(r.Context())
	if err != nil && err.Error() == "sql: no rows in result set" {
		if err := views.Index(
			false,
			sc.OAuthConfig.Scopes[0],
			sc.OAuthConfig.ClientID,
		).Render(r.Context(), w); err != nil {
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

	if err := views.Index(
		true,
		"",
		"",
	).Render(r.Context(), w); err != nil {
		log.Fatal().AnErr("error", err).Msg("ERROR: while sending main page")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
