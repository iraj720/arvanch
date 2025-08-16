package ratelimit

import (
	"math"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	BulkEvaluateKey      = "arvanch_bulk"
	DeliveryEvaluateKey  = "arvanch_delivery"
	DefaultMiddlewareHit = 1
)

type RateLimiterMiddleware struct {
	evaluator       Evaluator
	hitValue        int64
	enableRetryHint bool
	failOpen        bool
}

func NewRateLimiterMiddleware(
	evaluator Evaluator, hitValue int64, enableRetryHint, failOpen bool,
) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		evaluator:       evaluator,
		hitValue:        hitValue,
		enableRetryHint: enableRetryHint,
		failOpen:        failOpen,
	}
}

func (rl *RateLimiterMiddleware) CheckLimit() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			allowed, resetTime, err := rl.evaluator.EvaluateWithWaitTime(rl.hitValue)
			if err != nil {
				logrus.Errorf("rate limit middleware failed with error: %s", err.Error())

				if rl.failOpen {
					return next(c)
				}

				return c.String(http.StatusInternalServerError, "rate-limit middleware failed")
			}

			if rl.enableRetryHint && !allowed {
				c.Response().Header().Set("Retry-After", strconv.Itoa(int(math.Ceil(resetTime.Seconds()))))

				return c.NoContent(http.StatusTooManyRequests)
			}

			return next(c)
		}
	}
}
