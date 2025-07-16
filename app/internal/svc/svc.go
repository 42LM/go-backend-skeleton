package svc

import (
	"context"
	"log/slog"

	"go-backend-skeleton/app/internal/db"
)

// Service provides business logic.
type Service interface {
	NoneService
}

type NoneService interface {
	FindNone(ctx context.Context) string
}

type service struct {
	noneRepo db.NoneRepository
	logger   *slog.Logger
}

var _ Service = (*service)(nil)

// ServiceConfig contains the configuration params of the service.
type ServiceConfig struct {
	NoneRepo db.NoneRepository
	Logger   *slog.Logger
}

// New returns a service with middleware wired in.
func New(config *ServiceConfig) Service {
	var svc Service
	svc = &service{
		noneRepo: config.NoneRepo,
		logger:   config.Logger,
	}
	svc = LoggingMiddleware(config.Logger)(svc)
	return svc
}
