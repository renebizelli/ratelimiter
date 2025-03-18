package middlewares_ratelimiter

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	pkg_utils "github.com/renebizelli/ratelimiter/pkg/utils"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
)

var jwt = pkg_utils.NewJwt("SECRET", 500)

func TestLimiterOff(t *testing.T) {

	b := &BasedOnToken{
		on: false,
	}

	ch := make(chan Response)

	b.Limiter(&http.Request{}, ch)

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

func TestLimiter(t *testing.T) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	core := NewCoreRedis(rdb)

	basedOnToken := &BasedOnToken{
		core:           core,
		headerByStuffs: &HeaderByStuffs{},
		on:             false,
	}

	type TestParam struct {
		TimeSpleep int
		StatusCode int
	}

	testParams := []TestParam{
		TestParam{TimeSpleep: 0, StatusCode: http.StatusOK},
		TestParam{TimeSpleep: 200, StatusCode: http.StatusOK},
		TestParam{TimeSpleep: 200, StatusCode: http.StatusOK},
		TestParam{TimeSpleep: 200, StatusCode: http.StatusTooManyRequests},
	}

	claims := map[string]interface{}{
		"key":                "rene.oliveira",
		"rl-max-requests":    3,
		"rl-seconds-blocked": 5,
		"exp":                time.Now().Add(time.Minute * time.Duration(60)).Unix(),
	}

	tokenString, _ := jwt.Generate(claims)

	request := &http.Request{
		Header: http.Header{},
	}

	basedOnToken.headerByStuffs.SetAPIKey(request, tokenString)

	tt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	ttt := tt.Add(time.Duration(1) * time.Second)
	log.Println(ttt)

	time.Sleep(time.Until(ttt))
	log.Println(time.Now())

	for i, p := range testParams {
		log.Println(i, p, time.Now())
		time.Sleep(time.Duration(p.TimeSpleep) * time.Millisecond)
		ch := make(chan Response)
		go basedOnToken.Limiter(request, ch)
		response := <-ch
		assert.Equal(t, p.StatusCode, response.HttpStatus, i)
	}
}
