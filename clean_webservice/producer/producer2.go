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
	internalSig := "sig456"
	requestID2 := "req789"

	fullOrder := FullOrder{
		Orders: Orders{
			OrderUID:          "6",
			TrackNumber:       "trackBack33",
			Entry:             "entry2",
			Locale:            "ru",
			InternalSignature: &internalSig,
			CustomerID:        "cust123",
			DeliveryService:   "FedEx",
			Shardkey:          "shard2",
			SmID:              43,
			DateCreated:       time.Now().Add(-time.Hour).Format(time.RFC3339), // час назад
			OofShard:          "oof2",
		},
		Delivery: Delivery{
			OrderUID: "6",
			Name:     "Anna Ivanova",
			Phone:    "+7987654321",
			Zip:      "654321",
			City:     "Saint Petersburg",
			Address:  "Nevsky pr. 10",
			Region:   "Leningrad Region",
			Email:    "anna@example.com",
		},
		Payment: Payment{
			OrderUID:     "6",
			Transaction:  "txn456",
			RequestID:    &requestID2,
			Currency:     "RUB",
			Provider:     "Mastercard",
			Amount:       2000,
			PaymentDT:    time.Now().Add(-time.Minute * 30).Unix(), // 30 минут назад
			Bank:         "Sberbank",
			DeliveryCost: 300,
			GoodsTotal:   1700,
			CustomFee:    0,
		},
		Items: []Item{
			{
				OrderUID:    "6",
				ChrtID:      3333,
				TrackNumber: "trackBack33",
				Price:       1200,
				Rid:         "rid3",
				Name:        "Watch",
				Sale:        5,
				Size:        "L",
				TotalPrice:  1140,
				NmID:        20001,
				Brand:       "Casio",
				Status:      1,
			},
			{
				OrderUID:    "6",
				ChrtID:      4444,
				TrackNumber: "trackBack33",
				Price:       800,
				Rid:         "rid4",
				Name:        "Backpack",
				Sale:        0,
				Size:        "M",
				TotalPrice:  800,
				NmID:        20002,
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
