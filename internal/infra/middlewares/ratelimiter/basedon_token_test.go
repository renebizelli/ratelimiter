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

	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong := rdb.Ping(ctx)

	assert.Equal(t, "PONG", pong.Val())

	core := NewCoreRedis(rdb)

	basedOnToken := NewBasedOnToken(core, jwt, &HeaderByStuffs{}, true, 3, 5)

	scenarious := []Scenario{

		ScenarioGenerate("rene", 3, true),
		ScenarioGenerate("rene_1", 5, true),
		ScenarioGenerate("rene_2", 5, false),
		ScenarioGenerate("rene_3", 50, true),
	}

	for _, scenario := range scenarious {

		claims := map[string]interface{}{
			"key":                scenario.Email,
			"rl-max-requests":    scenario.MaxAllowedRequests,
			"rl-seconds-blocked": 5,
			"exp":                time.Now().Add(time.Minute * time.Duration(60)).Unix(),
		}

		tokenString, _ := jwt.Generate(claims)

		request := &http.Request{
			Header: http.Header{},
		}

		basedOnToken.headerByStuffs.SetAPIKey(request, tokenString)

		cha := make(chan Response, len(scenario.Params))

		for _, p := range scenario.Params {
			time.Sleep(time.Duration(p.TimeSpleep) * time.Millisecond)
			go basedOnToken.Limiter(request, cha)
		}

		for i := 0; i < len(scenario.Params); i++ {
			x := <-cha
			log.Println(i, scenario.Params[i].StatusCode, x.HttpStatus)
			assert.Equal(t, scenario.Params[i].StatusCode, x.HttpStatus)
		}

	}
}

func ScenarioGenerate(email string, quantity200 int, isLast429 bool) Scenario {
	s := Scenario{Email: email}
	s.AddParams(quantity200, isLast429)
	return s
}

type Scenario struct {
	Email              string
	Params             []TestParam
	MaxAllowedRequests int
}

func (s *Scenario) AddParams(quantity200 int, isLast429 bool) {

	millisecond := 1000 / quantity200

	s.MaxAllowedRequests = quantity200

	for i := 0; i < quantity200; i++ {
		s.Params = append(s.Params, TestParam{StatusCode: http.StatusOK, TimeSpleep: millisecond})
	}

	if isLast429 {
		s.Params[len(s.Params)-1].StatusCode = http.StatusTooManyRequests
		s.MaxAllowedRequests = quantity200 - 1
	}
}

type TestParam struct {
	TimeSpleep int
	StatusCode int
}
