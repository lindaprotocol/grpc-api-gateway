package middleware

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lindaprotocol/grpc-api-gateway/internal/services/auth"
	"github.com/lindaprotocol/grpc-api-gateway/pkg/utils"
)

type RateLimiter struct {
	redisClient *redis.Client
	config      RateLimitConfig
}

type RateLimitConfig struct {
	Enabled     bool
	DefaultQPS  int
	DefaultBurst int
	Strategy    string // token_bucket, sliding_window, leaky_bucket
	Store       string // redis, memory
}

func RateLimit(redisClient *redis.Client, cfg RateLimitConfig) func(http.Handler) http.Handler {
	limiter := &RateLimiter{
		redisClient: redisClient,
		config:      cfg,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !cfg.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Get user from context (set by auth middleware)
			user, ok := r.Context().Value(AuthUserKey).(*auth.User)
			if !ok {
				// Anonymous user
				user = &auth.User{
					IsAnonymous: true,
					RateLimit:   cfg.DefaultQPS,
				}
			}

			// Get client IP as fallback
			clientIP := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				clientIP = forwarded
			}

			// Determine key for rate limiting
			var key string
			if user.ID != "" {
				key = "rate:user:" + user.ID
			} else {
				key = "rate:ip:" + clientIP
			}

			// Apply rate limiting strategy
			var allowed bool
			var remaining int64
			var err error

			switch cfg.Strategy {
			case "token_bucket":
				allowed, remaining, err = limiter.tokenBucket(key, user.RateLimit)
			case "sliding_window":
				allowed, remaining, err = limiter.slidingWindow(key, user.RateLimit)
			default:
				allowed, remaining, err = limiter.tokenBucket(key, user.RateLimit)
			}

			if err != nil {
				// If rate limiting fails, allow request but log error
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				// Check if this violation should trigger a block
				if user.ID != "" {
					limiter.handleViolation(user)
				}

				utils.RespondWithErrorHTTP(w, http.StatusTooManyRequests, "Rate limit exceeded. Please slow down.")
				return
			}

			// Set rate limit headers
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(user.RateLimit))
			w.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Second).Unix(), 10))

			next.ServeHTTP(w, r)
		})
	}
}

func (l *RateLimiter) tokenBucket(key string, rate int) (bool, int64, error) {
	ctx := context.Background()
	now := time.Now().UnixNano()

	// Lua script for token bucket algorithm
	script := `
		local key = KEYS[1]
		local rate = tonumber(ARGV[1])
		local now = tonumber(ARGV[2])
		local capacity = rate

		local last = redis.call('GET', key .. ':last')
		local tokens = redis.call('GET', key .. ':tokens')

		if last == false then
			last = now
			tokens = capacity
		else
			local elapsed = now - tonumber(last)
			local new_tokens = math.floor(elapsed * rate / 1000000000)
			tokens = math.min(capacity, tonumber(tokens) + new_tokens)
			last = now
		end

		if tokens >= 1 then
			tokens = tokens - 1
			redis.call('SET', key .. ':last', last)
			redis.call('SET', key .. ':tokens', tokens)
			redis.call('EXPIRE', key .. ':last', 3600)
			redis.call('EXPIRE', key .. ':tokens', 3600)
			return {1, tokens}
		else
			return {0, tokens}
		end
	`

	result, err := l.redisClient.Eval(ctx, script, []string{key}, rate, now).Result()
	if err != nil {
		return false, 0, err
	}

	vals := result.([]interface{})
	allowed := vals[0].(int64) == 1
	remaining := int64(vals[1].(int64))

	return allowed, remaining, nil
}

func (l *RateLimiter) slidingWindow(key string, rate int) (bool, int64, error) {
	ctx := context.Background()
	now := time.Now().Unix()

	// Use sorted set for sliding window
	windowKey := key + ":window"
	windowSize := int64(1) // 1 second window

	// Remove old entries
	l.redisClient.ZRemRangeByScore(ctx, windowKey, "0", strconv.FormatInt(now-windowSize, 10))

	// Count current window entries
	count, err := l.redisClient.ZCard(ctx, windowKey).Result()
	if err != nil {
		return false, 0, err
	}

	if count < int64(rate) {
		// Add current request
		member := strconv.FormatInt(now, 10) + ":" + strconv.FormatInt(rand.Int63(), 10)
		l.redisClient.ZAdd(ctx, windowKey, &redis.Z{
			Score:  float64(now),
			Member: member,
		})
		l.redisClient.Expire(ctx, windowKey, time.Duration(windowSize)*time.Second)
		return true, int64(rate) - count - 1, nil
	}

	return false, 0, nil
}

func (l *RateLimiter) handleViolation(user *auth.User) {
	ctx := context.Background()
	violationKey := "violation:user:" + user.ID

	// Increment violation count
	count, err := l.redisClient.Incr(ctx, violationKey).Result()
	if err != nil {
		return
	}
	l.redisClient.Expire(ctx, violationKey, 24*time.Hour)

	// Block user after multiple violations
	if count >= 3 {
		blockedUntil := time.Now().Add(30 * time.Second)
		l.redisClient.Set(ctx, "blocked:user:"+user.ID, blockedUntil.Unix(), 24*time.Hour)
	}
}