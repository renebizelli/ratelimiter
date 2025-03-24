package middlewares_ratelimiter

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestLimiterIpOff(t *testing.T) {

	b := &BasedOnIP{
		on: false,
	}

	ch := make(chan Response)

	go b.Limiter(&http.Request{}, ch)

	e := <-ch

	assert.Equal(t, e.HttpStatus, http.StatusOK)

}

func TestLimiterIp(t *testing.T) {

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

	allowedRequest := 10

	basedOnIP := NewBasedOnIP(core, &HeaderByStuffs{}, true, allowedRequest, 5)

	ips := []string{
		"111.111.111.111",
		"222.222.222.222",
		"333.333.333.333",
	}

	for _, ip := range ips {

		request := &http.Request{
			RemoteAddr: ip,
		}

		cha := make(chan Response, 1000)

		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		ticker := time.NewTicker(time.Duration(100) * time.Millisecond)

		defer ticker.Stop()

		counter := 0

		c := true

		for c {

			go basedOnIP.Limiter(request, cha)
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
				assert.Equal(t, allowedRequest, counterOk)
				assert.Greater(t, counter429, 0)
				c = false
			}

		}

	}
}

func TestLimiterBlockedIP(t *testing.T) {

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

	allowedRequest := 1
	blockedSeconds := 15

	basedOnIP := NewBasedOnIP(core, &HeaderByStuffs{}, true, allowedRequest, blockedSeconds)

	ip := time.Now().Format("20060102150405.999999999")

	cha := make(chan Response)

	request := &http.Request{
		RemoteAddr: string(ip),
	}

	go basedOnIP.Limiter(request, cha)
	<-cha

	time.Sleep(time.Duration(200) * time.Millisecond)
	go basedOnIP.Limiter(request, cha)

	c := <-cha

	assert.Equal(t, 429, c.HttpStatus)
}
