package middlewares

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
	"golang.org/x/time/rate"
)

var limiterSet = cache.New(5*time.Minute, 10*time.Minute)

func NewRateLimiter(path string, newLimiter *rate.Limiter, duration time.Duration) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := c.RealIP() + path
			limiter, ok := limiterSet.Get(key)
			if !ok || limiter == nil {
				limiterSet.Set(key, newLimiter, duration)
				limiter = newLimiter
			}

			if !limiter.(*rate.Limiter).Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"status":  "fail",
					"message": "Too many requests.",
				})
			}

			return next(c)
		}
	}
}
