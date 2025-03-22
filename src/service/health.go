package service

import (
	"context"

	"github.com/highonsemicolon/aura/src/dal"
)

type HealthService interface {
	Liveness(context.Context) error
	Readiness(context.Context) readiness
}

type healthService struct {
	db dal.Database
}

type readiness struct {
	Database bool `json:"database"`
}

func NewHealthService(db dal.Database) *healthService {
	return &healthService{
		db: db,
	}
}

func (s *healthService) Liveness(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *healthService) Readiness(ctx context.Context) readiness {
	dbReady := s.db.PingContext(ctx) == nil
	return readiness{Database: dbReady}
}
