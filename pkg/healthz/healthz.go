package healthz

import (
	"context"
	"time"

	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type Checker func(ctx context.Context) bool

type Healthz struct {
	server     *health.Server
	liveness   map[string]struct{}
	readiness  map[string][]Checker
	interval   time.Duration
	cancelFunc context.CancelFunc
}

func NewHealthz(interval time.Duration) *Healthz {
	return &Healthz{
		server:    health.NewServer(),
		liveness:  make(map[string]struct{}),
		readiness: make(map[string][]Checker),
		interval:  interval,
	}
}

func (h *Healthz) Server() *health.Server {
	return h.server
}

func (h *Healthz) RegisterLiveness(svc string) {
	h.liveness[svc] = struct{}{}
	h.server.SetServingStatus(svc, grpc_health_v1.HealthCheckResponse_SERVING)
}

func (h *Healthz) RegisterReadiness(svc string, checks ...Checker) {
	h.readiness[svc] = checks
	h.server.SetServingStatus(svc, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
}

func (h *Healthz) AddDynamicCheck(svc string, check Checker) {
	h.readiness[svc] = append(h.readiness[svc], check)
}

func (h *Healthz) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	h.cancelFunc = cancel

	go func() {
		ticker := time.NewTicker(h.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				h.SetAllNotServing()
				return
			case <-ticker.C:
				h.evaluateReadiness(ctx)
			}
		}
	}()
}

func (h *Healthz) evaluateReadiness(ctx context.Context) {
	for svc, checks := range h.readiness {
		ready := true
		for _, check := range checks {
			if !check(ctx) {
				ready = false
				break
			}
		}
		if ready {
			h.server.SetServingStatus(svc, grpc_health_v1.HealthCheckResponse_SERVING)
		} else {
			h.server.SetServingStatus(svc, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		}
	}
}

func (h *Healthz) Stop(ctx context.Context) error {
	if h.cancelFunc != nil {
		h.cancelFunc()
	}
	h.SetAllNotServing()
	return nil
}

func (h *Healthz) SetAllNotServing() {
	for svc := range h.liveness {
		h.server.SetServingStatus(svc, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}
	for svc := range h.readiness {
		h.server.SetServingStatus(svc, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}
}
