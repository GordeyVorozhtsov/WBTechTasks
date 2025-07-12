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
	internalSig := "sig123"
	requestID := "req456"

	fullOrder := FullOrder{
		Orders: Orders{
			OrderUID:          "5",
			TrackNumber:       "trackBack22",
			Entry:             "entry1",
			Locale:            "en",
			InternalSignature: &internalSig,
			CustomerID:        "cust789",
			DeliveryService:   "DHL",
			Shardkey:          "shard1",
			SmID:              42,
			DateCreated:       time.Now().Format(time.RFC3339),
			OofShard:          "oof1",
		},
		Delivery: Delivery{
			OrderUID: "5",
			Name:     "John Doe",
			Phone:    "+1234567890",
			Zip:      "123456",
			City:     "Moscow",
			Address:  "Lenina 1",
			Region:   "Moscow Region",
			Email:    "john@example.com",
		},
		Payment: Payment{
			OrderUID:     "5",
			Transaction:  "txn789",
			RequestID:    &requestID,
			Currency:     "USD",
			Provider:     "Visa",
			Amount:       1000,
			PaymentDT:    time.Now().Unix(),
			Bank:         "BankName",
			DeliveryCost: 50,
			GoodsTotal:   950,
			CustomFee:    0,
		},
		Items: []Item{
			{
				OrderUID:    "5",
				ChrtID:      1111,
				TrackNumber: "trackBack22",
				Price:       500,
				Rid:         "rid1",
				Name:        "Glasses",
				Sale:        0,
				Size:        "M",
				TotalPrice:  500,
				NmID:        10001,
				Brand:       "RayBan",
				Status:      1,
			},
			{
				OrderUID:    "5",
				ChrtID:      2222,
				TrackNumber: "trackBack22",
				Price:       450,
				Rid:         "rid2",
				Name:        "Sneakers",
				Sale:        10,
				Size:        "L",
				TotalPrice:  450,
				NmID:        10002,
				Brand:       "Nike",
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
