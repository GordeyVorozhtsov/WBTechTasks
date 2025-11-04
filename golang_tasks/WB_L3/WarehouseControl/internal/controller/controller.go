package warehouse

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/wb-go/wbf/dbpg"
)

type Item struct {
	ID        int       `json:"id"`
	SKU       string    `json:"sku"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	Location  string    `json:"location"`
	Note      string    `json:"note"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuditRecord struct {
	ID        int            `json:"id"`
	ItemID    sql.NullInt64  `json:"item_id"`
	Operation string         `json:"operation"`
	ChangedBy sql.NullString `json:"changed_by"`
	ChangedAt time.Time      `json:"changed_at"`
	OldData   sql.NullString `json:"old_data"`
	NewData   sql.NullString `json:"new_data"`
}

type Service struct {
	db *dbpg.DB
}

func NewService() *Service {
	dsn := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(dsn) == "" {
		dsn = "postgres://user:password@localhost:5433/mydb?sslmode=disable"
	}

	db, err := dbpg.New(dsn, nil, &dbpg.Options{
		MaxOpenConns:    20,
		MaxIdleConns:    10,
		ConnMaxLifetime: 30 * time.Minute,
	})
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	return &Service{db: db}
}

func (s *Service) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	if s.db.Master != nil {
		_ = s.db.Master.Close()
	}
	for _, sl := range s.db.Slaves {
		if sl != nil {
			_ = sl.Close()
		}
	}
	return nil
}

func (s *Service) setCurrentActor(ctx context.Context, actor string) error {
	_, err := s.db.Master.ExecContext(ctx, `SELECT set_config('myapp.current_user', $1, true)`, actor)
	return err
}

// добавляет новый товар
func (s *Service) CreateItem(ctx context.Context, actor string, it *Item) error {
	if err := s.setCurrentActor(ctx, actor); err != nil {
		return err
	}

	q := `INSERT INTO items (sku, name, quantity, location, note)
	      VALUES ($1, $2, $3, $4, $5)
	      RETURNING id, updated_at`
	row := s.db.Master.QueryRowContext(ctx, q, it.SKU, it.Name, it.Quantity, it.Location, it.Note)
	return row.Scan(&it.ID, &it.UpdatedAt)
}

// возвращает все товары
func (s *Service) ListItems(ctx context.Context) ([]Item, error) {
	q := `SELECT id, sku, name, quantity, location, note, updated_at FROM items ORDER BY id`
	rows, err := s.db.Master.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Item
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.SKU, &it.Name, &it.Quantity, &it.Location, &it.Note, &it.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

// возвращает товар по ID
func (s *Service) GetItem(ctx context.Context, id int) (*Item, error) {
	q := `SELECT id, sku, name, quantity, location, note, updated_at FROM items WHERE id=$1`
	row := s.db.Master.QueryRowContext(ctx, q, id)
	var it Item
	if err := row.Scan(&it.ID, &it.SKU, &it.Name, &it.Quantity, &it.Location, &it.Note, &it.UpdatedAt); err != nil {
		return nil, err
	}
	return &it, nil
}

// обновляет товар
func (s *Service) UpdateItem(ctx context.Context, actor string, it *Item) error {
	if err := s.setCurrentActor(ctx, actor); err != nil {
		return err
	}

	q := `UPDATE items
	      SET sku=$1, name=$2, quantity=$3, location=$4, note=$5, updated_at=NOW()
	      WHERE id=$6`
	res, err := s.db.Master.ExecContext(ctx, q, it.SKU, it.Name, it.Quantity, it.Location, it.Note, it.ID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}

// удаляет товар
func (s *Service) DeleteItem(ctx context.Context, actor string, id int) error {
	if err := s.setCurrentActor(ctx, actor); err != nil {
		return err
	}
	_, err := s.db.Master.ExecContext(ctx, `DELETE FROM items WHERE id=$1`, id)
	return err
}

// возвращает историю изменений по товару
func (s *Service) ListAuditForItem(ctx context.Context, itemID int) ([]AuditRecord, error) {
	q := `SELECT id, item_id, operation, changed_by, changed_at, old_data::text, new_data::text
	      FROM items_audit
	      WHERE item_id=$1
	      ORDER BY changed_at DESC`
	rows, err := s.db.Master.QueryContext(ctx, q, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []AuditRecord
	for rows.Next() {
		var r AuditRecord
		if err := rows.Scan(&r.ID, &r.ItemID, &r.Operation, &r.ChangedBy, &r.ChangedAt, &r.OldData, &r.NewData); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
