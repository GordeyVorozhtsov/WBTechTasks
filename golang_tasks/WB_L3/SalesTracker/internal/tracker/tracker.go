package tracker

import (
	"context"
	"fmt"
	"time"

	"github.com/wb-go/wbf/dbpg"
)

type Item struct {
	ID         int       `json:"id"`
	Type       string    `json:"type"`
	Amount     float64   `json:"amount"`
	Category   string    `json:"category"`
	Note       string    `json:"note"`
	OccurredAt time.Time `json:"occurred_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type AnalyticsResult struct {
	Total   float64 `json:"total"`
	Average float64 `json:"average"`
	Count   int     `json:"count"`
	Median  float64 `json:"median"`
	P90     float64 `json:"p90"`
}

type Service struct {
	db *dbpg.DB
}

func NewService() *Service {
	connStr := "postgres://user:password@localhost:5433/mydb?sslmode=disable"
	db, err := dbpg.New(connStr, nil, nil)
	if err != nil {
		panic(fmt.Errorf("failed to connect to DB: %w", err))
	}
	return &Service{db: db}
}

// cписок всех записей
func (s *Service) List(ctx context.Context) ([]Item, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, type, amount, category, note, occurred_at, created_at
		FROM items
		ORDER BY occurred_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.ID, &i.Type, &i.Amount, &i.Category, &i.Note, &i.OccurredAt, &i.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, i)
	}
	return res, nil
}

// добавление новой записи
func (s *Service) Create(ctx context.Context, i *Item) error {
	query := `
		INSERT INTO items (type, amount, category, note, occurred_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := s.db.ExecContext(ctx, query, i.Type, i.Amount, i.Category, i.Note, i.OccurredAt)
	return err
}

// удаление по id
func (s *Service) Delete(ctx context.Context, id int) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM items WHERE id = $1`, id)
	return err
}

// обновление записи
func (s *Service) Update(ctx context.Context, i *Item) error {
	query := `
		UPDATE items
		SET type = $1, amount = $2, category = $3, note = $4, occurred_at = $5
		WHERE id = $6
	`
	_, err := s.db.ExecContext(ctx, query, i.Type, i.Amount, i.Category, i.Note, i.OccurredAt, i.ID)
	return err
}

// агрегированные метрики по диапазону дат
func (s *Service) Analytics(ctx context.Context, from, to time.Time) (*AnalyticsResult, error) {
	query := `
		WITH filtered AS (
			SELECT amount
			FROM items
			WHERE occurred_at BETWEEN $1 AND $2
		)
		SELECT
			COALESCE(SUM(amount), 0) AS total,
			COALESCE(AVG(amount), 0) AS average,
			COUNT(*) AS count,
			COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount), 0) AS median,
			COALESCE(PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount), 0) AS p90
		FROM filtered;
	`

	row := s.db.QueryRowContext(ctx, query, from, to)
	var res AnalyticsResult
	err := row.Scan(&res.Total, &res.Average, &res.Count, &res.Median, &res.P90)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *Service) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	if s.db.Master != nil {
		return s.db.Master.Close()
	}
	return nil
}
