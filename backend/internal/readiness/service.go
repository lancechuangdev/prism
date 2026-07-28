package readiness

import (
	"context"
	"time"
)

const (
	StatusReady    = "ready"
	StatusNotReady = "not_ready"
	StatusOK       = "ok"
	StatusFailed   = "failed"
)

type Probe func(context.Context) error

type Dependency struct {
	Status string `json:"status"`
}

type Report struct {
	Status       string                `json:"status"`
	Dependencies map[string]Dependency `json:"dependencies"`
}

type Checker interface {
	Check(context.Context) Report
}

type Service struct {
	timeout time.Duration
	probes  map[string]Probe
}

func New(timeout time.Duration, probes map[string]Probe) *Service {
	return &Service{timeout: timeout, probes: probes}
}

func (s *Service) Check(ctx context.Context) Report {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	type result struct {
		name string
		err  error
	}
	results := make(chan result, len(s.probes))
	pending := make(map[string]struct{}, len(s.probes))
	for name, probe := range s.probes {
		pending[name] = struct{}{}
		go func() {
			results <- result{name: name, err: probe(ctx)}
		}()
	}

	report := Report{
		Status:       StatusReady,
		Dependencies: make(map[string]Dependency, len(s.probes)),
	}
	for len(pending) > 0 {
		select {
		case result := <-results:
			delete(pending, result.name)
			status := StatusOK
			if result.err != nil {
				status = StatusFailed
				report.Status = StatusNotReady
			}
			report.Dependencies[result.name] = Dependency{Status: status}
		case <-ctx.Done():
			report.Status = StatusNotReady
			for name := range pending {
				report.Dependencies[name] = Dependency{Status: StatusFailed}
			}
			return report
		}
	}
	return report
}
