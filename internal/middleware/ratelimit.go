package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// IPRateLimit IP 级别速率限制
func IPRateLimit(rdb *redis.Client, limit int) gin.HandlerFunc {
	return rateLimit(rdb, "ratelimit:ip:", func(c *gin.Context) string {
		return c.ClientIP()
	}, limit, 1*time.Minute)
}

// AccountRateLimit 账号级别速率限制
func AccountRateLimit(rdb *redis.Client, limit int) gin.HandlerFunc {
	return rateLimit(rdb, "ratelimit:acct:", func(c *gin.Context) string {
		return c.GetString("username")
	}, limit, 1*time.Minute)
}

func rateLimit(rdb *redis.Client, prefix string, keyFn func(*gin.Context) string, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := prefix + keyFn(c)
		if key == prefix {
			c.Next()
			return
		}

		ctx := context.Background()
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}
		if count == 1 {
			rdb.Expire(ctx, key, window)
		}

		if count > int64(limit) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
			})
			return
		}
		c.Next()
	}
}
