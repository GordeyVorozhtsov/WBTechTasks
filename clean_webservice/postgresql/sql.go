package postgresql

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	connStr string
	Pool    *pgxpool.Pool
)

type (
	Orders struct {
		OrderUID          string    `json:"order_uid"`
		TrackNumber       string    `json:"track_number"`
		Entry             string    `json:"entry"`
		Locale            string    `json:"locale"`
		InternalSignature *string   `json:"internal_signature,omitempty"`
		CustomerID        string    `json:"customer_id"`
		DeliveryService   string    `json:"delivery_service"`
		Shardkey          string    `json:"shardkey"`
		SmID              int32     `json:"sm_id"`
		DateCreated       time.Time `json:"date_created"`
		OofShard          string    `json:"oof_shard"`
	}
	Delivery struct {
		OrderUID string `json:"order_uid"`
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Zip      string `json:"zip"`
		City     string `json:"city"`
		Address  string `json:"address"`
		Region   string `json:"region"`
		Email    string `json:"email"`
	}
	Payment struct {
		OrderUID     string  `json:"order_uid"`
		Transaction  string  `json:"transaction"`
		RequestID    *string `json:"request_id,omitempty"`
		Currency     string  `json:"currency"`
		Provider     string  `json:"provider"`
		Amount       int32   `json:"amount"`
		PaymentDT    int64   `json:"payment_dt"`
		Bank         string  `json:"bank"`
		DeliveryCost int32   `json:"delivery_cost"`
		GoodsTotal   int32   `json:"goods_total"`
		CustomFee    int32   `json:"custom_fee"`
	}
	Item struct {
		OrderUID    string `json:"order_uid"`
		ChrtID      int64  `json:"chrt_id"`
		TrackNumber string `json:"track_number"`
		Price       int32  `json:"price"`
		Rid         string `json:"rid"`
		Name        string `json:"name"`
		Sale        int32  `json:"sale"`
		Size        string `json:"size"`
		TotalPrice  int32  `json:"total_price"`
		NmID        int64  `json:"nm_id"`
		Brand       string `json:"brand"`
		Status      int32  `json:"status"`
	}
	FullOrder struct {
		Orders   Orders   `json:"orders"`
		Delivery Delivery `json:"delivery"`
		Payment  Payment  `json:"payment"`
		Items    []Item   `json:"items"`
	}
)

func SetConnectionString(connectionString string) {
	connStr = connectionString
}

func Connect(ctx context.Context) error {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	Pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("create connection pool: %w", err)
	}
	return nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}

func readCreateTables() string {
	data, err := os.ReadFile("postgresql/create_tables.sql")
	if err != nil {
		log.Fatalf("unavailable read file: %v", err)
	}
	return string(data)
}

func CreateTables(ctx context.Context) error {
	s := readCreateTables()
	for i, q := range strings.Split(s, ";") {
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}
		if _, err := Pool.Exec(ctx, q); err != nil {
			return fmt.Errorf("exec query #%d failed: %w\nquery: %s", i+1, err, q)
		}
	}
	log.Println("All create table queries executed successfully")
	return nil
}

// функция собирает выходную структуру, из мелких функций по таблицам
func GetFullOrder(ctx context.Context, orderUID string) (*FullOrder, error) {
	order, err := getOrder(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	delivery, err := getDelivery(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	payment, err := getPayment(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	items, err := getItems(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	return &FullOrder{
		Orders:   *order,
		Delivery: *delivery,
		Payment:  *payment,
		Items:    items,
	}, nil
}

func getOrder(ctx context.Context, orderUID string) (*Orders, error) {
	var o Orders
	err := Pool.QueryRow(ctx, `
		SELECT order_uid, track_number, entry, locale, internal_signature, customer_id,
		       delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders WHERE order_uid=$1`, orderUID).
		Scan(&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature,
			&o.CustomerID, &o.DeliveryService, &o.Shardkey, &o.SmID, &o.DateCreated, &o.OofShard)
	if err != nil {
		return nil, fmt.Errorf("getOrder: %w", err)
	}
	return &o, nil
}

func getDelivery(ctx context.Context, orderUID string) (*Delivery, error) {
	var d Delivery
	err := Pool.QueryRow(ctx, `
		SELECT order_uid, name, phone, zip, city, address, region, email
		FROM delivery WHERE order_uid=$1`, orderUID).
		Scan(&d.OrderUID, &d.Name, &d.Phone, &d.Zip, &d.City, &d.Address, &d.Region, &d.Email)
	if err != nil {
		return nil, fmt.Errorf("getDelivery: %w", err)
	}
	return &d, nil
}

func getPayment(ctx context.Context, orderUID string) (*Payment, error) {
	var p Payment
	err := Pool.QueryRow(ctx, `
		SELECT order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment WHERE order_uid=$1`, orderUID).
		Scan(&p.OrderUID, &p.Transaction, &p.RequestID, &p.Currency, &p.Provider,
			&p.Amount, &p.PaymentDT, &p.Bank, &p.DeliveryCost, &p.GoodsTotal, &p.CustomFee)
	if err != nil {
		return nil, fmt.Errorf("getPayment: %w", err)
	}
	return &p, nil
}

func getItems(ctx context.Context, orderUID string) ([]Item, error) {
	rows, err := Pool.Query(ctx, `
		SELECT order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items WHERE order_uid=$1`, orderUID)
	if err != nil {
		return nil, fmt.Errorf("getItems query: %w", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(
			&i.OrderUID, &i.ChrtID, &i.TrackNumber, &i.Price, &i.Rid,
			&i.Name, &i.Sale, &i.Size, &i.TotalPrice, &i.NmID, &i.Brand, &i.Status,
		); err != nil {
			return nil, fmt.Errorf("getItems scan: %w", err)
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getItems rows error: %w", err)
	}
	return items, nil
}

// удобная обертка без которой была какаято ошибка с on conflict
func insertOrUpdate(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) error {
	_, err := tx.Exec(ctx, query, args...)
	return err
}

func InsertFullOrder(ctx context.Context, order *FullOrder) (err error) {

	log.Printf("InsertFullOrder: start inserting order %s", order.Orders.OrderUID)
	tx, err := Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Printf("rollback error: %v (original error: %v)", rbErr, err)
				err = fmt.Errorf("rollback error: %v, original error: %w", rbErr, err)
			} else {
				log.Printf("transaction rolled back due to error: %v", err)
			}
		} else {
			if cmErr := tx.Commit(ctx); cmErr != nil {
				log.Printf("commit error: %v", cmErr)
				err = cmErr
			} else {
				log.Printf("transaction committed successfully for order %s", order.Orders.OrderUID)
			}
		}
	}()

	err = insertOrUpdate(ctx, tx, `
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id,
			delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (order_uid) DO UPDATE SET
			track_number = EXCLUDED.track_number,
			entry = EXCLUDED.entry,
			locale = EXCLUDED.locale,
			internal_signature = EXCLUDED.internal_signature,
			customer_id = EXCLUDED.customer_id,
			delivery_service = EXCLUDED.delivery_service,
			shardkey = EXCLUDED.shardkey,
			sm_id = EXCLUDED.sm_id,
			date_created = EXCLUDED.date_created,
			oof_shard = EXCLUDED.oof_shard
	`, order.Orders.OrderUID, order.Orders.TrackNumber, order.Orders.Entry, order.Orders.Locale, order.Orders.InternalSignature,
		order.Orders.CustomerID, order.Orders.DeliveryService, order.Orders.Shardkey, order.Orders.SmID, order.Orders.DateCreated, order.Orders.OofShard)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	err = insertOrUpdate(ctx, tx, `
		INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (order_uid) DO UPDATE SET
			name = EXCLUDED.name,
			phone = EXCLUDED.phone,
			zip = EXCLUDED.zip,
			city = EXCLUDED.city,
			address = EXCLUDED.address,
			region = EXCLUDED.region,
			email = EXCLUDED.email
	`, order.Delivery.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
		order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("insert delivery: %w", err)
	}

	err = insertOrUpdate(ctx, tx, `
		INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (order_uid) DO UPDATE SET
			transaction = EXCLUDED.transaction,
			request_id = EXCLUDED.request_id,
			currency = EXCLUDED.currency,
			provider = EXCLUDED.provider,
			amount = EXCLUDED.amount,
			payment_dt = EXCLUDED.payment_dt,
			bank = EXCLUDED.bank,
			delivery_cost = EXCLUDED.delivery_cost,
			goods_total = EXCLUDED.goods_total,
			custom_fee = EXCLUDED.custom_fee
	`, order.Payment.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}

	for _, item := range order.Items {
		err = insertOrUpdate(ctx, tx, `
			INSERT INTO items (
				order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
			) VALUES (
				$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12
			)
			ON CONFLICT (order_uid, chrt_id) DO UPDATE SET
				track_number = EXCLUDED.track_number,
				price = EXCLUDED.price,
				rid = EXCLUDED.rid,
				name = EXCLUDED.name,
				sale = EXCLUDED.sale,
				size = EXCLUDED.size,
				total_price = EXCLUDED.total_price,
				nm_id = EXCLUDED.nm_id,
				brand = EXCLUDED.brand,
				status = EXCLUDED.status
		`, item.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("insert item (order_uid=%s, chrt_id=%d): %w", item.OrderUID, item.ChrtID, err)
		}
	}

	log.Printf("InsertFullOrder: finished inserting order %s successfully", order.Orders.OrderUID)
	return nil
}
