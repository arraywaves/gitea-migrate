package core

import (
	"context"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter *rate.Limiter
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		limiter: rate.NewLimiter(r, b),
	}
}

func (rl *RateLimiter) Wait(context context.Context) error {
	return rl.limiter.Wait(context)
}

func (rl *RateLimiter) Allow() bool {
	return rl.limiter.Allow()
}
