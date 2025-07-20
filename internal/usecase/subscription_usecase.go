package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/tamaqazaq/subscription-service/internal/domain"
	"github.com/tamaqazaq/subscription-service/internal/domain/application"
	"github.com/tamaqazaq/subscription-service/internal/domain/repository"
)

type subscriptionUsecase struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionUsecase(r repository.SubscriptionRepository) application.SubscriptionUsecase {
	return &subscriptionUsecase{repo: r}
}

func (u *subscriptionUsecase) Create(sub *domain.Subscription) error {
	return u.repo.Create(sub)
}

func (u *subscriptionUsecase) GetAll() ([]domain.Subscription, error) {
	return u.repo.GetAll()
}

func (u *subscriptionUsecase) GetByID(id uuid.UUID) (*domain.Subscription, error) {
	return u.repo.GetByID(id)
}

func (u *subscriptionUsecase) Update(id uuid.UUID, sub *domain.Subscription) error {
	return u.repo.Update(id, sub)
}

func (u *subscriptionUsecase) Delete(id uuid.UUID) error {
	return u.repo.Delete(id)
}

func (u *subscriptionUsecase) GetTotal(userID *uuid.UUID, serviceName *string, start, end time.Time) (int, error) {
	return u.repo.GetTotal(userID, serviceName, start, end)
}
