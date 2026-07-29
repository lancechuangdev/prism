package resilience

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestValueRetriesWithBackoff(t *testing.T) {
	attempts := 0
	value, err := Value(context.Background(), Policy{
		Attempts: 3, AttemptTimeout: time.Second, InitialBackoff: time.Millisecond, MaxBackoff: 2 * time.Millisecond,
	}, func(context.Context) (string, error) {
		attempts++
		if attempts < 3 {
			return "", errors.New("temporary")
		}
		return "ok", nil
	})
	if err != nil || value != "ok" || attempts != 3 {
		t.Fatalf("value = %q, attempts = %d, err = %v", value, attempts, err)
	}
}

func TestValueHonorsParentCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	attempts := 0
	_, err := Value(ctx, Policy{
		Attempts: 3, AttemptTimeout: time.Second, InitialBackoff: time.Second, MaxBackoff: time.Second,
	}, func(context.Context) (string, error) {
		attempts++
		return "", errors.New("temporary")
	})
	if !errors.Is(err, context.Canceled) || attempts != 1 {
		t.Fatalf("attempts = %d, err = %v", attempts, err)
	}
}
