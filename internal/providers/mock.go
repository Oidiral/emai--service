package providers

import (
	"context"

	"github.com/Oidiral/emai--service/internal/models"
)

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) Slug() string {
	return models.ProviderMock
}

func (p *MockProvider) Send(ctx context.Context, email *models.ProviderEmail) error {
	return nil
}
