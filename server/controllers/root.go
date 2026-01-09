package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/LuisanaMTDev/spaced_learning/server/frontend/views"
	"github.com/rs/zerolog/log"
)

// TODO: Add auth for the CLI.
func (sc *ServerConfig) Root(w http.ResponseWriter, r *http.Request) {
	client := strings.ToUpper(r.Header.Get("SL-Client-Type"))

	if client == "SL-CLI" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	accessToken, err := sc.DBQueries.GetAccessToken(r.Context())
	if err != nil && err.Error() == "sql: no rows in result set" {
		err = views.Index(
			false,
			sc.OAuthConfig.Scopes[0],
			sc.OAuthConfig.ClientID,
			0).Render(r.Context(), w)
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

	log.Debug().Str("access_token", accessToken.String).Msg("SUCCESS: db has access token")

	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		log.Error().AnErr("error", err).Msg("failed creating request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken.String))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().AnErr("error", err).Msg("failed retrieving user details")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	userInfo := UserInfoResponse{}

	if err := json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		log.Error().AnErr("error", err).Msg("failed decoding user details")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Debug().Int64("user_id", userInfo.ID).Str("user_name", userInfo.Login).Msg("User Information retrieved.")

	if strconv.FormatInt(userInfo.ID, 10) == os.Getenv("MY_GITHUB_USER_ID") {
		err = sc.DBQueries.AddUserID(r.Context(), sql.NullInt64{Valid: true, Int64: userInfo.ID})
		if err != nil {
			log.Error().AnErr("error", err).Msg("failed retrieving user details")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	err = views.Index(
		true,
		"",
		"",
		userInfo.ID,
	).Render(r.Context(), w)
	if err != nil {
		log.Fatal().AnErr("error", err).Msg("ERROR: while sending main page")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
