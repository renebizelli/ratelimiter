package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

	fmt.Println(time.Now().UnixMilli())

	e := rdb.Set(ctx, "key", "1", time.Second).Err()
	if e != nil {
		fmt.Println(e)
	}

	v, e := rdb.Get(ctx, "key").Result()
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(v)

	time.Sleep(1 * time.Second)
	v = rdb.Get(ctx, "key").Val()
	fmt.Println(v)
	cfgs := configs.LoadConfig("./")

	mux := http.NewServeMux()

	songWebServer := webserver.NewSongWebServer()
	loginWebserver := webserver.NewLoginWebServer(cfgs.JWTToken, cfgs.JWTExpires)

	mux.HandleFunc("POST /login", loginWebserver.LoginHandler)
	mux.Handle("GET /songs", middlewares.RateLimiter(cfgs, http.HandlerFunc(songWebServer.GetSongsHandler)))

	fmt.Printf("Web server running on port %s\n", cfgs.WebServerPort)

	http.ListenAndServe(fmt.Sprintf(":%s", cfgs.WebServerPort), mux)
}
