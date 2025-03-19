package middlewares_ratelimiter

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

type CoreRedis struct {
	db *redis.Client
}

func NewCoreRedis(db *redis.Client) *CoreRedis {
	return &CoreRedis{db: db}
}

var count int32

func (l *CoreRedis) Limiter(ctx context.Context, key Key, parameters *Parameters) int {

	//counterKey := fmt.Sprintf("%s:%d", key, time.Now().Second())
	counterKey := string(key)
	blockedKey := fmt.Sprintf("%s:blocked", key)

	blocked := l.db.Exists(ctx, blockedKey).Val()

	if blocked == 1 {
		return http.StatusTooManyRequests
	}

	countInt, _ := l.db.Get(ctx, counterKey).Int()

	count = int32(countInt)

	fmt.Fprintln(os.Stdout, counterKey, count, time.Now())

	atomic.AddInt32(&count, 1)

	if count > int32(parameters.MaxRequests) {
		l.db.Set(ctx, blockedKey, nil, time.Duration(parameters.BlockedSeconds)*time.Second)
		return http.StatusTooManyRequests
	}

	l.db.Set(ctx, counterKey, count, time.Second)

	return http.StatusOK
}
