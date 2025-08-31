package models

import "context"

const (
	ProviderMock  = "mock"
	EmailProvider = "email"
)

type ProviderInterface interface {
	Slug() string
	Send(ctx context.Context, email *ProviderEmail) error
}
