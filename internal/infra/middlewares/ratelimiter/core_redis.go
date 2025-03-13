package middlewares_ratelimiter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type CoreRedis struct {
	db *redis.Client
}

func NewCoreRedis(db *redis.Client) *CoreRedis {
	return &CoreRedis{db: db}
}

func (l *CoreRedis) Limiter(ctx context.Context, key Key, parameters *Parameters) HttpStatus {

	counterKey := fmt.Sprintf("%s:%d", key, time.Now().Second())
	blockedKey := fmt.Sprintf("%s:blocked", key)

	blocked := l.db.Exists(ctx, blockedKey).Val()

	if blocked == 1 {
		return http.StatusTooManyRequests
	}

	count, _ := l.db.Get(ctx, counterKey).Int()

	count++

	if count > parameters.MaxRequests {
		l.db.Set(ctx, blockedKey, nil, time.Duration(parameters.BlockedSeconds)*time.Second)
		return http.StatusTooManyRequests
	}

	l.db.Set(ctx, counterKey, count, time.Second)

	return http.StatusOK
}
