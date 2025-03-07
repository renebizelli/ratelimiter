package main

import (
	"fmt"
	"net/http"

	"github.com/renebizelli/ratelimiter/configs"
	"github.com/renebizelli/ratelimiter/internal/infra/middlewares"
	"github.com/renebizelli/ratelimiter/internal/infra/webserver"
)

func main() {

	configs := configs.LoadConfig("./")

	mux := http.NewServeMux()

	songWebServer := webserver.NewSongWebServer()
	loginWebserver := webserver.NewLoginWebServer(configs.JWTToken, configs.JWTExpires)

	mux.HandleFunc("POST /login", loginWebserver.LoginHandler)
	mux.Handle("GET /songs", middlewares.RateLimiter(http.HandlerFunc(songWebServer.GetSongsHandler)))

	fmt.Printf("Web server running on port %s\n", configs.WebServerPort)

	http.ListenAndServe(fmt.Sprintf(":%s", configs.WebServerPort), mux)
}
