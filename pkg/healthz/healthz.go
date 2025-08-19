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

func (hm *Healthz) Server() *health.Server {
	return hm.server
}

func (hm *Healthz) RegisterLiveness(service string) {
	hm.liveness[service] = struct{}{}
	hm.server.SetServingStatus(service, grpc_health_v1.HealthCheckResponse_SERVING)
}

func (hm *Healthz) RegisterReadiness(service string, checks ...Checker) {
	hm.readiness[service] = checks
	hm.server.SetServingStatus(service, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
}

func (hm *Healthz) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	hm.cancelFunc = cancel

	go func() {
		ticker := time.NewTicker(hm.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				for svc := range hm.liveness {
					hm.server.SetServingStatus(svc, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
				}
				for svc := range hm.readiness {
					hm.server.SetServingStatus(svc, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
				}
				return
			case <-ticker.C:
				hm.evaluateReadiness(ctx)
			}
		}
	}()
}

func (hm *Healthz) evaluateReadiness(ctx context.Context) {
	for svc, checks := range hm.readiness {
		ready := true
		for _, check := range checks {
			if !check(ctx) {
				ready = false
				break
			}
		}
		if ready {
			hm.server.SetServingStatus(svc, grpc_health_v1.HealthCheckResponse_SERVING)
		} else {
			hm.server.SetServingStatus(svc, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
		}
	}
}

func (hm *Healthz) Stop() {
	if hm.cancelFunc != nil {
		hm.cancelFunc()
	}
	for svc := range hm.liveness {
		hm.server.SetServingStatus(svc, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}
	for svc := range hm.readiness {
		hm.server.SetServingStatus(svc, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	}
}
