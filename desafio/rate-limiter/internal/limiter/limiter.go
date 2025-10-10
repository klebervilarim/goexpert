package limiter

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	RedisClient       *redis.Client
	IPLimit           int
	IPBlockSeconds    int
	TokenLimit        int
	TokenBlockSeconds int
	WindowSeconds     int // janela para limitar por segundo
}

func NewLimiter() *Limiter {
	ipLimit, _ := strconv.Atoi(os.Getenv("IP_LIMIT"))
	ipBlock, _ := strconv.Atoi(os.Getenv("IP_BLOCK_SECONDS"))
	tokenLimit, _ := strconv.Atoi(os.Getenv("TOKEN_LIMIT"))
	tokenBlock, _ := strconv.Atoi(os.Getenv("TOKEN_BLOCK_SECONDS"))
	window := 1 // limite por segundo

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &Limiter{
		RedisClient:       rdb,
		IPLimit:           ipLimit,
		IPBlockSeconds:    ipBlock,
		TokenLimit:        tokenLimit,
		TokenBlockSeconds: tokenBlock,
		WindowSeconds:     window,
	}
}

// AllowIP verifica limite por IP
func (l *Limiter) AllowIP(ip string) (bool, error) {
	return l.allow("ip:"+ip, l.IPLimit, l.IPBlockSeconds)
}

// AllowToken verifica limite por Token
func (l *Limiter) AllowToken(token string) (bool, error) {
	return l.allow("token:"+token, l.TokenLimit, l.TokenBlockSeconds)
}

// lógica interna
func (l *Limiter) allow(key string, limit int, blockSeconds int) (bool, error) {
	ctx := context.Background()

	// verifica se está bloqueado
	blockKey := "block:" + key
	blocked, _ := l.RedisClient.Get(ctx, blockKey).Result()
	if blocked != "" {
		return false, nil
	}

	// chave por segundo para janela de rate limit
	windowKey := fmt.Sprintf("%s:%d", key, time.Now().Unix())
	count, err := l.RedisClient.Incr(ctx, windowKey).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		l.RedisClient.Expire(ctx, windowKey, time.Duration(l.WindowSeconds)*time.Second)
	}

	if int(count) > limit {
		// seta bloqueio por N segundos
		l.RedisClient.Set(ctx, blockKey, 1, time.Duration(blockSeconds)*time.Second)
		return false, nil
	}

	return true, nil
}
