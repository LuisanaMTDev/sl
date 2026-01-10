package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/LuisanaMTDev/spaced_learning/server/database/gosql_queries"
	"github.com/rs/zerolog/log"
)

func (sc *ServerConfig) AddLesson(w http.ResponseWriter, r *http.Request) {
	client := r.Header.Get("SL-Client-Type")

	switch client {
	case "SL-CLI":
		body := AddLessonRequest{}
		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			log.Error().AnErr("error", err).Msg("While decoding json request.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Debug().Any("request_body", body).Msg("Decoded request body.")

		repetitionsDatesAsJson, err := json.Marshal(body.RepetitionsDates)

		if err != nil {
			log.Error().AnErr("error", err).Msg("While marshaling repetitions dates on decoded request body.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		repetitionsDatesAsString := string(repetitionsDatesAsJson)

		err = sc.DBQueries.AddLesson(r.Context(), gosql_queries.AddLessonParams{
			Topic:         body.Topic,
			StartedDate:   body.RepetitionsDates[0],
			Json:          repetitionsDatesAsString,
			AmountOfCards: body.AmountOfCards,
		})

		if err != nil {
			log.Error().AnErr("error", err).Msg("While adding lesson to the database.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		currentLessons, err := sc.DBQueries.GetAllLessons(r.Context())

		if err != nil {
			log.Error().AnErr("error", err).Msg("While getting lessons from the database.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Debug().Any("current_lessons", currentLessons).Msg("Current lessons in the database.")

		w.Write([]byte("Request has success."))
	case "SL-WEB-APP":
		err := r.ParseForm()
		// TODO: Implement.
		if err != nil {
			log.Error().AnErr("error", err).Msg("While parsing request form.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusForbidden)
		return
	}
}
