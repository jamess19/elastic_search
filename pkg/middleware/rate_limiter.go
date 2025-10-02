package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	redisClient *redis.Client
	windowSize  time.Duration
	maxRequests int
}

func NewRateLimiter(redisClient *redis.Client, windowSize time.Duration, maxRequests int) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		windowSize:  windowSize,
		maxRequests: maxRequests,
	}
}

func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		userID := c.Request.Header.Get("x-user-id")

		if userID == "" {
			c.Next()
			return
		}

		ctx := context.Background()
		now := time.Now()
		nowMs := now.UnixNano() / int64(time.Millisecond)
		windowStartMs := now.Add(-rl.windowSize).UnixNano() / int64(time.Millisecond)

		// Tạo key Redis cho user + IP
		redisKey := fmt.Sprintf("rate_limit:%s:%s", ip, userID)

		// Xóa các request quá thời gian sliding window
		rl.redisClient.ZRemRangeByScore(ctx, redisKey, "0", fmt.Sprintf("%d", windowStartMs))

		// Đếm số request còn lại trong sliding window
		reqCount, err := rl.redisClient.ZCard(ctx, redisKey).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if reqCount >= int64(rl.maxRequests) {
			fmt.Printf("Rate limit exceeded: IP=%s, UserID=%s, Count=%d\n", ip, userID, reqCount)
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		// Thêm request hiện tại vào ZSET
		member := fmt.Sprintf("%d:%s", nowMs, strings.ReplaceAll(c.Request.URL.Path, "/", "_"))
		_, err = rl.redisClient.ZAdd(ctx, redisKey, &redis.Z{
			Score:  float64(nowMs),
			Member: member,
		}).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		// Set TTL để tự xóa key sau khi không hoạt động
		rl.redisClient.Expire(ctx, redisKey, rl.windowSize*2)

		c.Next()
	}
}
