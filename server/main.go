package main

import (
	"net/http"

	"github.com/LuisanaMTDev/spaced_learning/server/controllers"
	"github.com/joho/godotenv"

	// "github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "modernc.org/sqlite"
)

func main() {
	godotenv.Load()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	handler := http.NewServeMux()
	serverConfig := controllers.NewServerConfig()

	//End points
	handler.Handle("GET /app/", http.StripPrefix("/app/", http.FileServer(http.Dir("./frontend/assets/"))))

	handler.HandleFunc("GET /", serverConfig.Root)

	handler.HandleFunc("GET /oauth2/callback", serverConfig.OAuthCallback)
	handler.HandleFunc("POST /login", serverConfig.Login)

	handler.HandleFunc("POST /lesson/add", serverConfig.AddLesson)

	server := http.Server{Handler: handler, Addr: ":8090"}
	log.Info().Str("running_platfotm", serverConfig.Platform).Msg("Running...")
	err := server.ListenAndServe()
	log.Fatal().AnErr("server_error", err)
}
