package middlewares_ratelimiter

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type CoreRedis struct {
	db  *redis.Client
	ctx context.Context
}

func NewCoreRedis(ctx context.Context, db *redis.Client) *CoreRedis {
	return &CoreRedis{ctx: ctx, db: db}
}

func (l *CoreRedis) IsBlocked(blockedKey string) bool {
	return l.db.Exists(l.ctx, blockedKey).Val() > 0
}

func (l *CoreRedis) Limiter(key Key, parameters *Parameters) int {

	blockedKey := fmt.Sprintf("%s:blocked", key)

	if l.IsBlocked(blockedKey) {
		return http.StatusTooManyRequests
	}

	counterKey := string(key)
	count, _ := l.db.Get(l.ctx, counterKey).Int()

	count++

	if count > parameters.MaxRequests {
		l.db.Set(l.ctx, blockedKey, nil, time.Duration(parameters.BlockedSeconds)*time.Second)
		return http.StatusTooManyRequests
	}

	l.db.Set(l.ctx, counterKey, count, time.Second)

	return http.StatusOK
}
