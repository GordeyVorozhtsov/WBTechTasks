package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	kafka "wb/kafka"
	db "wb/postgresql"

	"github.com/gin-gonic/gin"
)

const (
	broker            = "kafka:9092"
	topicName         = "orders"
	consumerGroup     = "order-consumer-group"
	numPartitions     = 1
	replicationFactor = 1
	ginRout           = ":8081"
	cacheTTL          = 10 * time.Minute
	dbTimeout         = 5 * time.Second
)

// кеш с TTL
type (
	CacheItem struct {
		Value      *db.FullOrder
		Expiration int64
	}
	Cache struct {
		items map[string]CacheItem
		mu    sync.RWMutex
		ttl   time.Duration
	}
)

func NewCache(ttl time.Duration) *Cache {
	return &Cache{items: make(map[string]CacheItem), ttl: ttl}
}

// загрузка последних n заказов в кеш при старте программы(такое усовие задачи есть)
func preloadCache(ctx context.Context, cache *Cache, limit int) error {
	// Получаем 10 последних order_uid из orders по дате создания
	rows, err := db.Pool.Query(ctx, `
		SELECT order_uid FROM orders
		ORDER BY date_created DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return fmt.Errorf("preloadCache: query order_uids: %w", err)
	}
	defer rows.Close()

	var orderUIDs []string
	for rows.Next() {
		var orderUID string
		if err := rows.Scan(&orderUID); err != nil {
			return fmt.Errorf("preloadCache: scan order_uid: %w", err)
		}
		orderUIDs = append(orderUIDs, orderUID)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("preloadCache: rows error: %w", err)
	}

	// Для каждого order_uid получаем полный заказ и кладем в кеш
	for _, uid := range orderUIDs {
		fullOrder, err := db.GetFullOrder(ctx, uid)
		if err != nil {
			log.Printf("preloadCache: failed to get full order %s: %v", uid, err)
			continue
		}
		cache.Set(uid, fullOrder)
		log.Printf("preloadCache: cached order %s", uid)
	}
	return nil
}

func (c *Cache) Get(key string) (*db.FullOrder, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()
	if !found || time.Now().UnixNano() > item.Expiration {
		if found {
			c.mu.Lock()
			delete(c.items, key)
			c.mu.Unlock()
		}
		return nil, false
	}
	return item.Value, true
}
func (c *Cache) Set(key string, value *db.FullOrder) {
	c.mu.Lock()
	c.items[key] = CacheItem{Value: value, Expiration: time.Now().Add(c.ttl).UnixNano()}
	c.mu.Unlock()
}

// кафка десериализует и вставляет в бд
func handleMessage(data []byte) (*db.FullOrder, error) {
	var order db.FullOrder
	return &order, json.Unmarshal(data, &order)
}
func insertOrderToDB(ctx context.Context, order *db.FullOrder) error {
	// такая конструкция с таймаутом чтобы не зависать на багах с бд
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()
	return db.InsertFullOrder(ctx, order)
}

// gin http
// тест запросы curl localhost:8081/order/?
func startHTTPServer(cache *Cache) {
	router := gin.Default()
	router.GET("/order/:order_uid", func(c *gin.Context) {
		orderUID := c.Param("order_uid")
		if cached, found := cache.Get(orderUID); found {
			c.JSON(http.StatusOK, mapFullOrderToResponse(cached))
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
		defer cancel()
		fullOrder, err := db.GetFullOrder(ctx, orderUID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		cache.Set(orderUID, fullOrder)
		c.JSON(http.StatusOK, mapFullOrderToResponse(fullOrder))
	})
	router.Static("/static", "./web")
	log.Printf("Server running on http://localhost%s\n", ginRout)
	log.Fatal(router.Run(ginRout))
}

func mapFullOrderToResponse(o *db.FullOrder) gin.H {
	getString := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}
	items := make([]gin.H, 0, len(o.Items))
	for _, i := range o.Items {
		items = append(items, gin.H{
			"order_uid":    i.OrderUID,
			"chrt_id":      i.ChrtID,
			"track_number": i.TrackNumber,
			"price":        i.Price,
			"rid":          i.Rid,
			"name":         i.Name,
			"sale":         i.Sale,
			"size":         i.Size,
			"total_price":  i.TotalPrice,
			"nm_id":        i.NmID,
			"brand":        i.Brand,
			"status":       i.Status,
		})
	}
	return gin.H{
		"order_uid":          o.Orders.OrderUID,
		"track_number":       o.Orders.TrackNumber,
		"entry":              o.Orders.Entry,
		"locale":             o.Orders.Locale,
		"internal_signature": getString(o.Orders.InternalSignature),
		"customer_id":        o.Orders.CustomerID,
		"delivery_service":   o.Orders.DeliveryService,
		"shardkey":           o.Orders.Shardkey,
		"sm_id":              o.Orders.SmID,
		"date_created":       o.Orders.DateCreated.Format(time.RFC3339),
		"oof_shard":          o.Orders.OofShard,
		"delivery": gin.H{
			"order_uid": o.Delivery.OrderUID, "name": o.Delivery.Name, "phone": o.Delivery.Phone,
			"zip": o.Delivery.Zip, "city": o.Delivery.City, "address": o.Delivery.Address,
			"region": o.Delivery.Region, "email": o.Delivery.Email,
		},
		"payment": gin.H{
			"order_uid": o.Payment.OrderUID, "transaction": o.Payment.Transaction,
			"request_id": getString(o.Payment.RequestID), "currency": o.Payment.Currency,
			"provider": o.Payment.Provider, "amount": o.Payment.Amount, "payment_dt": o.Payment.PaymentDT,
			"bank": o.Payment.Bank, "delivery_cost": o.Payment.DeliveryCost,
			"goods_total": o.Payment.GoodsTotal, "custom_fee": o.Payment.CustomFee,
		},
		"items": items,
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" { // если не запарсилась env через докер
		connStr = "postgres://fucku:pass@postgres:5432/wb?sslmode=disable"
	}
	db.SetConnectionString(connStr)
	if err := db.Connect(ctx); err != nil {
		log.Fatalf("DB connect failed: %v", err)
	}
	defer db.Close()
	if err := db.CreateTables(ctx); err != nil {
		log.Fatalf("Create tables failed: %v", err)
	}

	cache := NewCache(cacheTTL)
	// Предзагрузка последних 10 заказов в кеш
	if err := preloadCache(ctx, cache, 10); err != nil {
		log.Printf("Warning: preload cache failed: %v", err)
	}

	err := kafka.CreateTopic(broker, topicName, numPartitions, replicationFactor)
	if err != nil {
		log.Fatalf("Failed to create Kafka topic: %v", err)
	}
	log.Printf("Kafka topic %q is ready", topicName)

	msgChan, err := kafka.RunKafkaConsumer(ctx, broker, topicName, consumerGroup)
	if err != nil {
		log.Fatalf("Kafka consumer failed: %v", err)
	}

	// читаем кафку, json строка в байтах передается в канал
	go func() {
		for msg := range msgChan {
			log.Printf("Received message: %s", string(msg))
			order, err := handleMessage(msg)
			if err != nil {
				log.Printf("JSON unmarshal error: %v", err)
				continue
			}
			if err := insertOrderToDB(ctx, order); err != nil {
				log.Printf("DB insert error: %v", err)
				continue
			}
			log.Printf("Order %s inserted successfully", order.Orders.OrderUID)
		}
	}()

	startHTTPServer(cache)
}
