package ratelimit

import (
	"context"
	"fmt"
	"time"

	"arvanch/config"

	gubernator "github.com/gubernator-io/gubernator/v2"
	"github.com/sirupsen/logrus"
)

const (
	RahyabBulkEvaluateKey = "rahyab_bulk"
	FakeBulkEvaluateKey   = "fake_bulk"
)

type OverLimitError struct {
	WaitTime time.Duration
}

func (o *OverLimitError) Error() string {
	return fmt.Sprintf("request overlimit. try again after %v", o.WaitTime)
}

type Evaluator interface {
	EvaluateWithWaitTime(hits int64) (bool, time.Duration, error)
}

type GubernatorEvaluator struct {
	*GubernatorLimiter
	key string
	cfg *config.RateLimitRule
}

type GubernatorLimiter struct {
	Client  gubernator.V1Client
	Timeout time.Duration
}

// func NewGubernatorLimiter(cfg mineCfg.Gubernator) (*GubernatorLimiter, error) {
// 	client, err := gubernator.DialV1Server(cfg.GRPCAddress, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to gubernator cluster: %w", err)
// 	}

// 	limiter := &GubernatorLimiter{
// 		Client:  client,
// 		Timeout: cfg.Timeout,
// 	}

// 	return limiter, nil
// }

func NewGubernatorEvaluator(limiter *GubernatorLimiter,
	rule *config.RateLimitRule, key string) (*GubernatorEvaluator, error) {
	if key == "" {
		return nil, fmt.Errorf("failed to create evaluator: empty key")
	}

	return &GubernatorEvaluator{
		GubernatorLimiter: limiter,
		key:               key,
		cfg:               rule,
	}, nil
}

// EvaluateWithWaitTime is used to evaluate request for specific rate limit.
// hit specifies how much this call costs (e.g. 10 from overall 340 limits).
// it also returns the wait time until the caller can try again.
func (l *GubernatorEvaluator) EvaluateWithWaitTime(hits int64) (bool, time.Duration, error) {
	resp, err := l.call(hits)

	logrus.Debugf("gubernator evaluate rate limit: hits: %d, gubernator resp: %v", hits, resp)

	if err != nil {
		return false, 0, fmt.Errorf("evaluate failed: %w", err)
	}

	if resp.GetStatus() == gubernator.Status_UNDER_LIMIT {
		return true, 0, nil
	} else {
		resetTime := time.Until(time.Unix(0, resp.GetResetTime()*int64(time.Millisecond)))
		return false, resetTime, nil
	}
}

func (l *GubernatorEvaluator) call(hits int64) (_ *gubernator.RateLimitResp, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(finalErr, startTime, fmt.Sprintf("%s:%s", l.cfg.Name, l.key)) }()

	rateLimitReq := gubernator.RateLimitReq{
		Name:      l.cfg.Name,
		UniqueKey: l.key,
		Hits:      hits,
		Limit:     l.cfg.Limit,
		Duration:  l.cfg.Duration.Milliseconds(),
		Algorithm: gubernator.Algorithm(l.cfg.Algorithm),
		Behavior:  gubernator.Behavior(l.cfg.Behaviour),
	}

	ctx, cancel := context.WithTimeout(context.Background(), l.Timeout)

	defer cancel()

	resp, err := l.Client.GetRateLimits(ctx, &gubernator.GetRateLimitsReq{
		Requests: []*gubernator.RateLimitReq{&rateLimitReq},
	})

	if err != nil {
		return nil, fmt.Errorf("could not get ratelimit from gubernator: %w", err)
	}

	if len(resp.GetResponses()) == 0 {
		return nil, fmt.Errorf("empty response from server")
	}

	return resp.GetResponses()[0], nil
}
