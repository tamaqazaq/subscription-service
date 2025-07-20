package domain

import (
	"github.com/google/uuid"
)

// swagger:model Subscription
type Subscription struct {
	ID          uuid.UUID `json:"id" db:"id" format:"uuid"`
	ServiceName string    `json:"service_name" db:"service_name"`
	Price       int       `json:"price" db:"price"`
	UserID      uuid.UUID `json:"user_id" db:"user_id" format:"uuid"`
	StartDate   DateOnly  `json:"start_date" db:"start_date"`
	EndDate     *DateOnly `json:"end_date,omitempty" db:"end_date"`
}
