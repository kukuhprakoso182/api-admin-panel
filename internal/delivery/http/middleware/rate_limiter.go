package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// bucket menyimpan state rate limit per key (IP / user).
type bucket struct {
	tokens    int
	lastRefil time.Time
	mu        sync.Mutex
}

// limiterStore menyimpan semua bucket dalam memory.
type limiterStore struct {
	mu      sync.RWMutex
	buckets map[string]*bucket
	max     int
	period  time.Duration
}

func newLimiterStore(max int, period time.Duration) *limiterStore {
	s := &limiterStore{
		buckets: make(map[string]*bucket),
		max:     max,
		period:  period,
	}
	// Goroutine untuk bersihkan bucket yang sudah kadaluarsa
	go s.cleanup()
	return s
}

func (s *limiterStore) allow(key string) bool {
	s.mu.Lock()
	b, ok := s.buckets[key]
	if !ok {
		b = &bucket{tokens: s.max, lastRefil: time.Now()}
		s.buckets[key] = b
	}
	s.mu.Unlock()

	b.mu.Lock()
	defer b.mu.Unlock()

	// Refill token jika sudah melewati period
	if time.Since(b.lastRefil) >= s.period {
		b.tokens = s.max
		b.lastRefil = time.Now()
	}

	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

func (s *limiterStore) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		s.mu.Lock()
		for key, b := range s.buckets {
			b.mu.Lock()
			expired := time.Since(b.lastRefil) > s.period*2
			b.mu.Unlock()
			if expired {
				delete(s.buckets, key)
			}
		}
		s.mu.Unlock()
	}
}

// RateLimit membatasi jumlah request per IP dalam periode tertentu.
//
// Contoh penggunaan:
//
//	// Maksimal 5 request per menit per IP
//	auth.POST("/login", middleware.RateLimit(5, time.Minute), handler.Login)
func RateLimit(max int, period time.Duration) gin.HandlerFunc {
	store := newLimiterStore(max, period)

	return func(c *gin.Context) {
		key := c.ClientIP()

		if !store.allow(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimitByUser membatasi request berdasarkan user_id (setelah autentikasi).
func RateLimitByUser(max int, period time.Duration) gin.HandlerFunc {
	store := newLimiterStore(max, period)

	return func(c *gin.Context) {
		key, exists := c.Get(ContextKeyUserID)
		if !exists {
			key = c.ClientIP()
		}

		if !store.allow(key.(string)) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
