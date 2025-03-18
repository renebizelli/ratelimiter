package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/renebizelli/ratelimiter/configs"
	middlewares_ratelimiter "github.com/renebizelli/ratelimiter/internal/infra/middlewares/ratelimiter"
	"github.com/renebizelli/ratelimiter/internal/infra/webserver"
	pkg_utils "github.com/renebizelli/ratelimiter/pkg/utils"

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

	registerRoutes(songWebServer, loginWebserver)

	fmt.Printf("Web server running on port %s\n", cfgs.WebServerPort)

	http.ListenAndServe(fmt.Sprintf(":%s", cfgs.WebServerPort), mux)
}

func registerRoutes(webservers ...webserver.RegisterRoutesInterface) {
	for _, ws := range webservers {
		ws.RegisterRoutes()
	}
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

func createJWTAuth(cfgs *configs.Config) *pkg_utils.Jwt {
	return pkg_utils.NewJwt(cfgs.JWTSecret, cfgs.JWTExpires)
}

func createRateLimiterMiddleware(jwt *pkg_utils.Jwt, cfgs *configs.Config, rdb *redis.Client) *middlewares_ratelimiter.RateLimiterMiddleware {

	core := middlewares_ratelimiter.NewCoreRedis(rdb)

	basedOnToken := middlewares_ratelimiter.NewBasedOnToken(
		core,
		jwt,
		&middlewares_ratelimiter.HeaderByStuffs{},
		cfgs.RATELIMITER_TOKEN_ON,
		cfgs.RATELIMITER_TOKEN_DEFAULT_MAX_REQUESTS,
		cfgs.RATELIMITER_TOKEN_DEFAULT_BLOCKED_SECONDS)

	basedOnIP := middlewares_ratelimiter.NewBasedOnIP(
		core,
		&middlewares_ratelimiter.HeaderByStuffs{},
		cfgs.RATELIMITER_IP_ON,
		cfgs.RATELIMITER_IP_MAX_REQUESTS,
		cfgs.RATELIMITER_IP_BLOCKED_SECONDS)

	ratelimiter := middlewares_ratelimiter.NewRateLimiterMiddleware(
		[]middlewares_ratelimiter.BasedonInterface{
			basedOnToken,
			basedOnIP,
		})

	return ratelimiter

}
