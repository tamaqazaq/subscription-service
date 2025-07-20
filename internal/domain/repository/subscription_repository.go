package repository

import (
	"github.com/google/uuid"
	"time"

	"github.com/tamaqazaq/subscription-service/internal/domain"
)

type SubscriptionRepository interface {
	Create(sub *domain.Subscription) error
	GetAll() ([]domain.Subscription, error)
	GetByID(id uuid.UUID) (*domain.Subscription, error)
	Update(id uuid.UUID, sub *domain.Subscription) error
	Delete(id uuid.UUID) error
	GetTotal(userID *uuid.UUID, serviceName *string, start, end time.Time) (int, error)
}
