package repository

import (
	"context"
	"errors"

	"github.com/Oidiral/emai--service/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type EmailRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewEmailRepository(db *gorm.DB, redis *redis.Client) *EmailRepository {
	return &EmailRepository{
		db:    db,
		redis: redis,
	}
}

func (r *EmailRepository) Create(ctx context.Context, s *models.Email) (*models.Email, error) {
	result := r.db.WithContext(ctx).Create(s)
	return s, result.Error
}

func (r *EmailRepository) Update(ctx context.Context, s *models.Email) error {
	if s.Uuid == "" {
		return errors.New("uuid is required")
	}

	err := r.db.WithContext(ctx).Model(&models.Email{}).Where("uuid = ?", s.Uuid).Updates(s).Error

	return err
}
