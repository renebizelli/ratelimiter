package middlewares_ratelimiter

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	pkg_utils "github.com/renebizelli/ratelimiter/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var jwt = pkg_utils.NewJwt("SECRET", 500)

func TestLimiterTokenOff(t *testing.T) {

	b := &BasedOnToken{
		on: false,
	}

	ch := make(chan Response)

	go b.Limiter(&http.Request{}, ch)

	e := <-ch

	assert.Equal(t, e.HttpStatus, http.StatusOK)

}

func TestLimiterOnWithNoAPIKey(t *testing.T) {

	b := &BasedOnToken{
		on: true,
	}

	ch := make(chan Response)

	go b.Limiter(&http.Request{}, ch)

	e := <-ch

	assert.Equal(t, e.HttpStatus, http.StatusOK)
}

func TestLimiterToken(t *testing.T) {

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong := rdb.Ping(ctx)

	if pong.Val() != "PONG" {
		panic(pong.Err())
	}

	core := NewCoreRedis(ctx, rdb)

	basedOnToken := NewBasedOnToken(core, jwt, &HeaderByStuffs{}, true, 3, 5)

	scenarious := []Scenario{
		ScenarioGenerateToken("rene", 10, 0),
		ScenarioGenerateToken("rene_1", 10, 0),
		ScenarioGenerateToken("rene_2", 10, 0),
		ScenarioGenerateToken("rene_3", 10, 0),
		ScenarioGenerateToken("rene_4", 10, 0),
	}

	for _, scenario := range scenarious {

		claims := map[string]interface{}{
			"key":                scenario.Email,
			"rl-max-requests":    scenario.AllowedRequest,
			"rl-seconds-blocked": 5,
			"exp":                time.Now().Add(time.Minute * time.Duration(60)).Unix(),
		}

		tokenString, _ := jwt.Generate(claims)

		request := &http.Request{
			Header: http.Header{},
		}

		basedOnToken.headerByStuffs.SetAPIKey(request, tokenString)

		cha := make(chan Response, 1000)

		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		ticker := time.NewTicker(time.Duration(100) * time.Millisecond)

		defer ticker.Stop()

		counter := 0

		c := true

		for c {

			go basedOnToken.Limiter(request, cha)
			counter++

			select {
			case <-ticker.C:
			case <-ctx.Done():
				counterOk := 0
				counter429 := 0
				for i := 0; i < counter; i++ {
					x := <-cha

					if x.HttpStatus == http.StatusOK {
						counterOk += 1
					} else {
						counter429 += 1
					}

				}
				log.Println(counterOk, counter429)
				assert.Equal(t, scenario.AllowedRequest, counterOk)
				assert.Greater(t, counter429, 0)
				c = false
			}

		}

	}
}

func TestLimiterBlockedToken(t *testing.T) {

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong := rdb.Ping(ctx)

	if pong.Val() != "PONG" {
		panic(pong.Err())
	}

	core := NewCoreRedis(ctx, rdb)

	basedOnToken := NewBasedOnToken(core, jwt, &HeaderByStuffs{}, true, 3, 5)

	scenarious := []Scenario{
		ScenarioGenerateToken("rene", 1, 5),
		ScenarioGenerateToken("rene_1", 1, 15),
		ScenarioGenerateToken("rene_2", 1, 10),
	}

	for _, scenario := range scenarious {

		claims := map[string]interface{}{
			"key":                scenario.Email,
			"rl-max-requests":    scenario.AllowedRequest,
			"rl-seconds-blocked": scenario.BlockedSeconds,
			"exp":                time.Now().Add(time.Minute * time.Duration(60)).Unix(),
		}

		tokenString, _ := jwt.Generate(claims)

		request := &http.Request{
			Header: http.Header{},
		}

		basedOnToken.headerByStuffs.SetAPIKey(request, tokenString)

		cha := make(chan Response)

		go basedOnToken.Limiter(request, cha)
		<-cha
		time.Sleep(time.Duration(100) * time.Millisecond)
		go basedOnToken.Limiter(request, cha)

		c := <-cha

		assert.Equal(t, 429, c.HttpStatus)

	}
}

func ScenarioGenerateToken(email string, allowedRequest int, blockedSeconds int) Scenario {
	s := Scenario{Email: email, AllowedRequest: allowedRequest, BlockedSeconds: blockedSeconds}
	return s
}

type Scenario struct {
	Email          string
	AllowedRequest int
	BlockedSeconds int
}
