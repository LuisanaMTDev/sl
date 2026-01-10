package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (sc *ServerConfig) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().Str("pattern_matched", r.Pattern).Str("http_method", r.Method).Str("url", r.URL.String()).Any("Headers", r.Header).Msg("Incoming request")

		next.ServeHTTP(w, r)
	})
}

func (sc *ServerConfig) HasAPIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
		log.Debug().Any("sended_api_key", apiKey).Msg("API Key sended by the user")
		if apiKey == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		hashedAPIKey, err := sc.DBQueries.GetAPIKEY(r.Context())
		if err != nil {
			log.Fatal().AnErr("error", err).Msg("ERROR: while getting API key from the db")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Debug().Any("hashed_api_key", hashedAPIKey).Msg("Hashed API Key obtained from the db")
		if hashedAPIKey.Valid {
			sameAPIKey := checkHashedApiKey(hashedAPIKey.String, apiKey)
			if sameAPIKey {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Header().Add("WWW-Authenticate", "Basic realm=Send API Key generated with sl login command.")
				return
			}
		} else {
			if err := json.NewEncoder(w).Encode(struct {
				HasAPIKeyInDB bool   `json:"has_api_key_in_db"`
				LoginIn       string `json:"login_in"`
			}{
				HasAPIKeyInDB: false,
				LoginIn:       "Run sl login to get you API Key.",
			}); err != nil {
				log.Fatal().AnErr("error", err).Msg("ERROR: while encoding response body")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	})
}
