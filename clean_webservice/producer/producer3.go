package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Определяем структуры с json тегами для корректной сериализации

type Orders struct {
	OrderUID          string  `json:"order_uid"`
	TrackNumber       string  `json:"track_number"`
	Entry             string  `json:"entry"`
	Locale            string  `json:"locale"`
	InternalSignature *string `json:"internal_signature,omitempty"`
	CustomerID        string  `json:"customer_id"`
	DeliveryService   string  `json:"delivery_service"`
	Shardkey          string  `json:"shardkey"`
	SmID              int32   `json:"sm_id"`
	DateCreated       string  `json:"date_created"` // В формате RFC3339
	OofShard          string  `json:"oof_shard"`
}

type Delivery struct {
	OrderUID string `json:"order_uid"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Zip      string `json:"zip"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Region   string `json:"region"`
	Email    string `json:"email"`
}

type Payment struct {
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

type Item struct {
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

type FullOrder struct {
	Orders   Orders   `json:"orders"`
	Delivery Delivery `json:"delivery"`
	Payment  Payment  `json:"payment"`
	Items    []Item   `json:"items"`
}

func main() {
	// Создаём тестовые данные
	internalSig := "sig789"
	requestID3 := "req101112"

	fullOrder := FullOrder{
		Orders: Orders{
			OrderUID:          "7",
			TrackNumber:       "trackBack44",
			Entry:             "entry3",
			Locale:            "en",
			InternalSignature: &internalSig,
			CustomerID:        "cust456",
			DeliveryService:   "UPS",
			Shardkey:          "shard3",
			SmID:              44,
			DateCreated:       time.Now().Format(time.RFC3339),
			OofShard:          "oof3",
		},
		Delivery: Delivery{
			OrderUID: "7",
			Name:     "Michael Smith",
			Phone:    "+15551234567",
			Zip:      "789012",
			City:     "New York",
			Address:  "5th Avenue 101",
			Region:   "NY",
			Email:    "michael@example.com",
		},
		Payment: Payment{
			OrderUID:     "7",
			Transaction:  "txn101112",
			RequestID:    &requestID3,
			Currency:     "USD",
			Provider:     "PayPal",
			Amount:       3500,
			PaymentDT:    time.Now().Unix(),
			Bank:         "Chase",
			DeliveryCost: 250,
			GoodsTotal:   3250,
			CustomFee:    0,
		},
		Items: []Item{
			{
				OrderUID:    "7",
				ChrtID:      5001,
				TrackNumber: "trackBack44",
				Price:       500,
				Rid:         "rid5001",
				Name:        "Laptop",
				Sale:        5,
				Size:        "15 inch",
				TotalPrice:  475,
				NmID:        30001,
				Brand:       "Dell",
				Status:      1,
			},
			{
				OrderUID:    "7",
				ChrtID:      5002,
				TrackNumber: "trackBack44",
				Price:       300,
				Rid:         "rid5002",
				Name:        "Wireless Mouse",
				Sale:        0,
				Size:        "Standard",
				TotalPrice:  300,
				NmID:        30002,
				Brand:       "Logitech",
				Status:      1,
			},
			{
				OrderUID:    "7",
				ChrtID:      5003,
				TrackNumber: "trackBack44",
				Price:       400,
				Rid:         "rid5003",
				Name:        "Mechanical Keyboard",
				Sale:        10,
				Size:        "Full-size",
				TotalPrice:  360,
				NmID:        30003,
				Brand:       "Corsair",
				Status:      1,
			},
			{
				OrderUID:    "7",
				ChrtID:      5004,
				TrackNumber: "trackBack44",
				Price:       150,
				Rid:         "rid5004",
				Name:        "USB-C Hub",
				Sale:        0,
				Size:        "Compact",
				TotalPrice:  150,
				NmID:        30004,
				Brand:       "Anker",
				Status:      1,
			},
			{
				OrderUID:    "7",
				ChrtID:      5005,
				TrackNumber: "trackBack44",
				Price:       200,
				Rid:         "rid5005",
				Name:        "External SSD",
				Sale:        15,
				Size:        "1TB",
				TotalPrice:  170,
				NmID:        30005,
				Brand:       "Samsung",
				Status:      1,
			},
			{
				OrderUID:    "7",
				ChrtID:      5006,
				TrackNumber: "trackBack44",
				Price:       100,
				Rid:         "rid5006",
				Name:        "Webcam",
				Sale:        5,
				Size:        "1080p",
				TotalPrice:  95,
				NmID:        30006,
				Brand:       "Logitech",
				Status:      1,
			},
			{
				OrderUID:    "7",
				ChrtID:      5007,
				TrackNumber: "trackBack44",
				Price:       400,
				Rid:         "rid5007",
				Name:        "Gaming Chair",
				Sale:        20,
				Size:        "Standard",
				TotalPrice:  320,
				NmID:        30007,
				Brand:       "DXRacer",
				Status:      1,
			},
		},
	}

	// Сериализуем в JSON
	jsonData, err := json.Marshal(fullOrder)
	if err != nil {
		log.Fatalf("Ошибка сериализации: %v", err)
	}
	fmt.Println(jsonData)

	// Создаём Kafka продюсера
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9092"})
	if err != nil {
		log.Fatalf("Ошибка создания продюсера: %v", err)
	}
	defer producer.Close()

	topic := "orders"

	// Отправляем сообщение
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          jsonData,
	}, nil)
	if err != nil {
		log.Fatalf("Ошибка отправки сообщения: %v", err)
	}

	// Ждём подтверждения доставки
	producer.Flush(10000)

	log.Println("Тестовое сообщение успешно отправлено в Kafka")
}
