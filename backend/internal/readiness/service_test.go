package readiness

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCheckReportsDependencyFailure(t *testing.T) {
	service := New(time.Second, map[string]Probe{
		"mysql": func(context.Context) error { return nil },
		"redis": func(context.Context) error { return errors.New("unavailable") },
	})

	report := service.Check(context.Background())

	if report.Status != StatusNotReady {
		t.Fatalf("status = %q", report.Status)
	}
	if report.Dependencies["mysql"].Status != StatusOK ||
		report.Dependencies["redis"].Status != StatusFailed {
		t.Fatalf("unexpected dependencies: %+v", report.Dependencies)
	}
}

func TestCheckBoundsProbeDuration(t *testing.T) {
	service := New(10*time.Millisecond, map[string]Probe{
		"slow": func(ctx context.Context) error {
			<-ctx.Done()
			return ctx.Err()
		},
	})

	start := time.Now()
	report := service.Check(context.Background())

	if elapsed := time.Since(start); elapsed > 250*time.Millisecond {
		t.Fatalf("readiness check took %s", elapsed)
	}
	if report.Status != StatusNotReady {
		t.Fatalf("status = %q", report.Status)
	}
}
