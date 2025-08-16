package ratelimit_test

import (
	"testing"
	"time"

	"arvanch/ratelimit"

	"github.com/stretchr/testify/suite"

	"arvanch/config"
)

type GubernatorLimiterSuite struct {
	G *ratelimit.GubernatorLimiter
	suite.Suite
}

func (suite *GubernatorLimiterSuite) SetupSuite() {
	// cfg := config.Init()

	// gubernator, err := ratelimit.NewGubernatorLimiter(cfg.Gubernator)
	// suite.NoError(err)

	// suite.G = gubernator
}

// nolint:funlen, gomnd, gocognit
func (suite *GubernatorLimiterSuite) TestEvaluate() {
	rateLimiterBucketTime := 200 * time.Millisecond
	cases := []struct {
		name               string
		rule               config.RateLimitRule
		firstRequestBurst  int
		secondRequestBurst int
		firstTryResult     bool
		secondTry          bool
		secondTryResult    bool
	}{
		{
			name:              "10 request pass",
			firstRequestBurst: 10,
			rule: config.RateLimitRule{
				Name:      "test_rule1",
				Duration:  rateLimiterBucketTime,
				Limit:     10,
				Algorithm: 0,
				Behaviour: 0,
			},
			firstTryResult: true,
		},
		{
			name:              "10 request over limit",
			firstRequestBurst: 12,
			rule: config.RateLimitRule{
				Name:      "test_rule2",
				Duration:  rateLimiterBucketTime,
				Limit:     10,
				Algorithm: 0,
				Behaviour: 0,
			},
			firstTryResult: false,
		},
		{
			name:               "10 request pass for 2 periods(bucket)",
			firstRequestBurst:  10,
			secondRequestBurst: 10,
			rule: config.RateLimitRule{
				Name:      "test_rule4",
				Duration:  rateLimiterBucketTime,
				Limit:     10,
				Algorithm: 0,
				Behaviour: 0,
			},
			firstTryResult:  true,
			secondTry:       true,
			secondTryResult: true,
		},
		{
			name:               "10 request over limit then pass",
			firstRequestBurst:  12,
			secondRequestBurst: 10,
			rule: config.RateLimitRule{
				Name:      "test_rule5",
				Duration:  rateLimiterBucketTime,
				Limit:     10,
				Algorithm: 0,
				Behaviour: 0,
			},
			firstTryResult:  false,
			secondTry:       true,
			secondTryResult: true,
		},
		{
			name:               "10 request pass then over limit",
			firstRequestBurst:  10,
			secondRequestBurst: 12,
			rule: config.RateLimitRule{
				Name:      "test_rule6",
				Duration:  rateLimiterBucketTime,
				Limit:     10,
				Algorithm: 0,
				Behaviour: 0,
			},
			firstTryResult:  true,
			secondTry:       true,
			secondTryResult: false,
		},
	}

	for i := range cases {
		tc := cases[i]
		suite.Run(tc.name, func() {
			evaluateKey := "test"

			var (
				result bool
				err    error
			)

			underLimited := int64(0)
			overLimited := int64(0)

			evaluator, err := ratelimit.NewGubernatorEvaluator(suite.G, &tc.rule, evaluateKey)
			suite.NoError(err)

			for range tc.firstRequestBurst {
				result, _, err = evaluator.EvaluateWithWaitTime(1)

				if result {
					underLimited++
				} else {
					overLimited++
				}

				suite.NoError(err)
			}

			suite.Equal(tc.rule.Limit, underLimited)

			if tc.firstTryResult {
				suite.Zero(overLimited)
			} else {
				suite.Equal(int64(tc.firstRequestBurst)-tc.rule.Limit, overLimited)
			}

			if tc.secondTry {
				underLimited = 0
				overLimited = 0

				time.Sleep(rateLimiterBucketTime)

				for range tc.secondRequestBurst {
					result, _, err = evaluator.EvaluateWithWaitTime(1)
					suite.NoError(err)

					if result {
						underLimited++
					} else {
						overLimited++
					}
				}

				suite.Equal(tc.rule.Limit, underLimited)

				if tc.secondTryResult {
					suite.Zero(overLimited)
				} else {
					suite.Equal(int64(tc.secondRequestBurst)-tc.rule.Limit, overLimited)
				}
			}
		})
	}
}

// nolint:funlen, gomnd, gocognit
func (suite *GubernatorLimiterSuite) TestEvaluateWithWaitTime() {
	cases := []struct {
		name         string
		rule         config.RateLimitRule
		requestCount int
		result       bool
		waitTime     time.Duration
	}{
		{
			name:         "10 request no wait",
			requestCount: 10,
			result:       true,
			rule: config.RateLimitRule{
				Name:      "test_wait_rule1",
				Duration:  100 * time.Millisecond,
				Limit:     10,
				Algorithm: 0,
				Behaviour: 0,
			},
			waitTime: 0,
		},
		{
			name:         "12 request token bucket wait 1 duration",
			requestCount: 12,
			result:       false,
			rule: config.RateLimitRule{
				Name:      "test_wait_rule2",
				Duration:  200 * time.Millisecond,
				Limit:     10,
				Algorithm: 0,
				Behaviour: 0,
			},
			waitTime: 200 * time.Millisecond,
		},
		{
			name:         "12 request leaky bucket wait",
			requestCount: 12,
			result:       false,
			rule: config.RateLimitRule{
				Name:      "test_wait_rule3",
				Duration:  1000 * time.Millisecond,
				Limit:     10,
				Algorithm: 1,
				Behaviour: 0,
			},
			// in leaky bucket algorithm reset time is now + duration/limit
			waitTime: 100 * time.Millisecond,
		},
	}
	for i := range cases {
		tc := cases[i]
		suite.Run(tc.name, func() {
			evaluateKey := "test"

			underLimited := int64(0)
			overLimited := int64(0)

			evaluator, err := ratelimit.NewGubernatorEvaluator(suite.G, &tc.rule, evaluateKey)
			suite.NoError(err)

			for range tc.requestCount {
				result, waitTime, err := evaluator.EvaluateWithWaitTime(1)

				if result {
					underLimited++

					suite.Zero(waitTime)
				} else {
					suite.InDelta(waitTime.Nanoseconds(), tc.waitTime.Nanoseconds(), float64(50*time.Millisecond))

					overLimited++
				}

				suite.NoError(err)
			}

			suite.Equal(tc.rule.Limit, underLimited)

			if tc.result {
				suite.Zero(overLimited)
			} else {
				suite.Equal(int64(tc.requestCount)-tc.rule.Limit, overLimited)
			}
		})
	}
}

func TestGubernatorLimiterSuite(t *testing.T) {
	suite.Run(t, new(GubernatorLimiterSuite))
}
