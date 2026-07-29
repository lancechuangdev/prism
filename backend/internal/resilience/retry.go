package resilience

import (
	"context"
	"fmt"
	"time"
)

type Policy struct {
	Attempts       int
	AttemptTimeout time.Duration
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

func Do(ctx context.Context, policy Policy, operation func(context.Context) error) error {
	_, err := Value(ctx, policy, func(attemptCtx context.Context) (struct{}, error) {
		return struct{}{}, operation(attemptCtx)
	})
	return err
}

func Value[T any](ctx context.Context, policy Policy, operation func(context.Context) (T, error)) (T, error) {
	var zero T
	if policy.Attempts < 1 || policy.AttemptTimeout <= 0 || policy.InitialBackoff < 0 || policy.MaxBackoff < policy.InitialBackoff {
		return zero, fmt.Errorf("invalid retry policy")
	}

	backoff := policy.InitialBackoff
	var lastErr error
	for attempt := 1; attempt <= policy.Attempts; attempt++ {
		attemptCtx, cancel := context.WithTimeout(ctx, policy.AttemptTimeout)
		value, err := operation(attemptCtx)
		cancel()
		if err == nil {
			return value, nil
		}
		lastErr = err
		if attempt == policy.Attempts {
			break
		}

		timer := time.NewTimer(backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return zero, ctx.Err()
		case <-timer.C:
		}
		backoff *= 2
		if backoff > policy.MaxBackoff {
			backoff = policy.MaxBackoff
		}
	}
	return zero, lastErr
}
