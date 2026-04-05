package health

import "context"

// Status is the health check response payload.
type Status struct {
	Status string `json:"status"`
}

// Service defines the health check business logic.
type Service interface {
	Check(ctx context.Context) Status
}

type service struct{}

// NewService returns a Service implementation.
func NewService() Service {
	return &service{}
}

func (s *service) Check(_ context.Context) Status {
	return Status{Status: "OK"}
}
