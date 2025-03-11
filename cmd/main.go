package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/renebizelli/ratelimiter/configs"
	"github.com/renebizelli/ratelimiter/internal/infra/middlewares"
	"github.com/renebizelli/ratelimiter/internal/infra/webserver"

	"github.com/redis/go-redis/v9"
)

func main() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	x := rdb.Ping(ctx)
	fmt.Println(x)

	cfgs := configs.LoadConfig("./")

	mux := http.NewServeMux()

	songWebServer := webserver.NewSongWebServer()
	loginWebserver := webserver.NewLoginWebServer(cfgs.JWTToken, cfgs.JWTExpires)

	ratelimiter := middlewares.NewRateLimiter(rdb, cfgs)

	mux.HandleFunc("POST /login", loginWebserver.LoginHandler)
	mux.Handle("GET /songs", ratelimiter.Limiter(http.HandlerFunc(songWebServer.GetSongsHandler)))

	fmt.Printf("Web server running on port %s\n", cfgs.WebServerPort)

	http.ListenAndServe(fmt.Sprintf(":%s", cfgs.WebServerPort), mux)
}
