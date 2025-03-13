package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/renebizelli/ratelimiter/configs"
	middlewares_ratelimiter "github.com/renebizelli/ratelimiter/internal/infra/middlewares/ratelimiter"
	"github.com/renebizelli/ratelimiter/internal/infra/webserver"

	"github.com/redis/go-redis/v9"
)

func main() {

	ctx := context.Background()

	rdb := createRedis(ctx)

	cfgs := configs.LoadConfig("./")

	jwt := createJWTAuth(cfgs)

	mux := http.NewServeMux()

	ratelimiter := createRateLimiterMiddleware(jwt, cfgs, rdb)
	songWebServer := webserver.NewSongWebServer(mux, ratelimiter)
	loginWebserver := webserver.NewLoginWebServer(mux, jwt, cfgs.JWTExpires)

	webservers := []webserver.RegisterRoutesInterface{
		songWebServer,
		loginWebserver,
	}

	for _, ws := range webservers {
		ws.RegisterRoutes()
	}

	fmt.Printf("Web server running on port %s\n", cfgs.WebServerPort)

	http.ListenAndServe(fmt.Sprintf(":%s", cfgs.WebServerPort), mux)
}

func createRedis(ctx context.Context) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong := rdb.Ping(ctx)
	if pong.Val() != "PONG" {
		fmt.Println("Error connecting to redis")
		return nil
	}

	fmt.Println("Redis connected")
	return rdb
}

func createJWTAuth(cfgs *configs.Config) *jwtauth.JWTAuth {
	return jwtauth.New("HS256", []byte(cfgs.JWTSecret), nil)
}

func createRateLimiterMiddleware(jwt *jwtauth.JWTAuth, cfgs *configs.Config, rdb *redis.Client) *middlewares_ratelimiter.RateLimiterMiddleware {

	basedOnToken := middlewares_ratelimiter.NewBasedOnToken(
		jwt,
		cfgs.RATELIMITER_TOKEN_DEFAULT_MAX_REQUESTS,
		cfgs.RATELIMITER_TOKEN_DEFAULT_BLOCKED_SECONDS)

	basedOnIP := middlewares_ratelimiter.NewBasedOnIP(
		cfgs.RATELIMITER_IP_MAX_REQUESTS,
		cfgs.RATELIMITER_IP_BLOCKED_SECONDS)

	core := middlewares_ratelimiter.NewCoreRedis(rdb)

	ratelimiter := middlewares_ratelimiter.NewRateLimiterMiddleware(
		cfgs.RATELIMITER_IP_ON,
		cfgs.RATELIMITER_TOKEN_ON,
		basedOnToken,
		basedOnIP,
		core,
	)

	return ratelimiter

}
