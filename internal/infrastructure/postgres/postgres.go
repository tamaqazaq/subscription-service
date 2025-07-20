package postgres

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/tamaqazaq/subscription-service/internal/domain"
)

type SubscriptionRepository struct {
	DB *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{DB: db}
}

func (r *SubscriptionRepository) Create(sub *domain.Subscription) error {
	sub.ID = uuid.New()

	query := `
		INSERT INTO subscriptions (
			id, service_name, price, user_id, start_date, end_date
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	start := sub.StartDate.ToTime()

	var end *time.Time
	if sub.EndDate != nil {
		t := sub.EndDate.ToTime()
		end = &t
	}

	_, err := r.DB.Exec(query, sub.ID, sub.ServiceName, sub.Price, sub.UserID, start, end)
	return err
}

func (r *SubscriptionRepository) GetAll() ([]domain.Subscription, error) {
	rows, err := r.DB.Query(`SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		var start time.Time
		var end sql.NullTime

		err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &start, &end)
		if err != nil {
			return nil, err
		}

		sub.StartDate = domain.DateOnly(start)
		if end.Valid {
			d := domain.DateOnly(end.Time)
			sub.EndDate = &d
		}

		subs = append(subs, sub)
	}

	return subs, nil
}

func (r *SubscriptionRepository) GetByID(id uuid.UUID) (*domain.Subscription, error) {
	var sub domain.Subscription
	var start time.Time
	var end sql.NullTime

	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`

	err := r.DB.QueryRow(query, id).Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &start, &end)
	if err != nil {
		return nil, err
	}

	sub.StartDate = domain.DateOnly(start)
	if end.Valid {
		d := domain.DateOnly(end.Time)
		sub.EndDate = &d
	}

	return &sub, nil
}

func (r *SubscriptionRepository) Update(id uuid.UUID, sub *domain.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
		WHERE id = $6
	`

	start := sub.StartDate.ToTime()

	var end *time.Time
	if sub.EndDate != nil {
		t := sub.EndDate.ToTime()
		end = &t
	}

	_, err := r.DB.Exec(query, sub.ServiceName, sub.Price, sub.UserID, start, end, id)
	return err
}
func (r *SubscriptionRepository) Delete(id uuid.UUID) error {
	_, err := r.DB.Exec(`DELETE FROM subscriptions WHERE id = $1`, id)
	return err
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

func (r *SubscriptionRepository) GetTotal(userID *uuid.UUID, serviceName *string, start, end time.Time) (int, error) {
	query := `
		SELECT SUM(price)
		FROM subscriptions
		WHERE 
			start_date <= $2
			AND (end_date IS NULL OR end_date >= $1)
	`

	args := []interface{}{start, end}
	argIdx := 3

	if userID != nil {
		query += ` AND user_id = $` + itoa(argIdx)
		args = append(args, *userID)
		argIdx++
	}

	if serviceName != nil {
		query += ` AND service_name = $` + itoa(argIdx)
		args = append(args, *serviceName)
		argIdx++
	}

	var total sql.NullInt64
	err := r.DB.QueryRow(query, args...).Scan(&total)
	if err != nil {
		return 0, err
	}

	if total.Valid {
		return int(total.Int64), nil
	}
	return 0, nil
}
