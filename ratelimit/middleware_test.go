package ratelimit_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"arvanch/ratelimit"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type RateLimiterMiddlewareSuite struct {
	suite.Suite
	rateLimiter *ShalghamEvaluator
}

type ShalghamEvaluator struct {
	*sync.Mutex

	duration   time.Duration
	remaining  int64
	shouldFail bool
}

func (g *ShalghamEvaluator) EvaluateWithWaitTime(hits int64) (bool, time.Duration, error) {
	if g.shouldFail {
		return false, 0, fmt.Errorf("limiter failed :(")
	}

	g.Lock()
	defer g.Unlock()

	if hits < g.remaining {
		return true, 0, nil
	}

	return false, g.duration, nil
}

func (suite *RateLimiterMiddlewareSuite) SetupSuite() {
	limiter := &ShalghamEvaluator{
		&sync.Mutex{}, 500 * time.Millisecond, 10, false,
	}

	suite.rateLimiter = limiter
}

//nolint:funlen,gomnd
func (suite *RateLimiterMiddlewareSuite) TestRateLimitMiddleware() {
	cases := []struct {
		name                string
		status              int
		hitValue            int
		waitDurationSeconds string
		evaluatorFail       bool
		middleware          *ratelimit.RateLimiterMiddleware
	}{
		{
			name:   "successfully pass rate limit",
			status: http.StatusOK,
			middleware: ratelimit.NewRateLimiterMiddleware(
				suite.rateLimiter, 1, false, false),
		},
		{
			name:     "successfully pass rate limit with custom hit",
			hitValue: 5,
			status:   http.StatusOK,
			middleware: ratelimit.NewRateLimiterMiddleware(
				suite.rateLimiter, 5, false, false),
		},
		{
			name:          "request pass with fail open",
			status:        http.StatusOK,
			evaluatorFail: true,
			middleware: ratelimit.NewRateLimiterMiddleware(
				suite.rateLimiter, 1, false, true),
		},
		{
			name:          "request fail with fail close",
			status:        http.StatusInternalServerError,
			evaluatorFail: true,
			middleware: ratelimit.NewRateLimiterMiddleware(
				suite.rateLimiter, 1, false, false),
		},
		{
			name:                "fail rate limit with high hit value",
			status:              http.StatusTooManyRequests,
			hitValue:            12,
			waitDurationSeconds: "1",
			middleware: ratelimit.NewRateLimiterMiddleware(
				suite.rateLimiter, 12, true, false),
		},
	}

	for i := range cases {
		tc := cases[i]
		suite.rateLimiter.shouldFail = tc.evaluatorFail
		middleware := tc.middleware.CheckLimit()(func(c echo.Context) error {
			return c.String(http.StatusOK, "test")
		})

		suite.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			res := httptest.NewRecorder()
			c := echo.New().NewContext(req, res)

			suite.NoError(middleware(c), tc.name)

			if tc.waitDurationSeconds != "" {
				retryAfter := res.Header().Get("Retry-After")
				suite.Equal(tc.waitDurationSeconds, retryAfter)
			}

			suite.Equal(tc.status, res.Code, tc.name)
		})
	}
}

func TestRateLimiterMiddlewareSuite(t *testing.T) {
	suite.Run(t, &RateLimiterMiddlewareSuite{})
}
