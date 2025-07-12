package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// RunKafkaConsumer запускает Kafka consumer и возвращает канал, в который будет поступать json
// канал закрывается, когда consumer остановится (например, при отмене контекста)
func RunKafkaConsumer(ctx context.Context, brokers, topic, groupID string) (<-chan []byte, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest", //  какойто дефолт на оффсет
	})
	if err != nil {
		return nil, err
	}

	if err := consumer.Subscribe(topic, nil); err != nil {
		consumer.Close()
		return nil, err
	}

	msgChan := make(chan []byte)

	go func() {
		defer func() {
			// при выходе из горутины закрываем(например контекст) consumer и канал сообщений
			consumer.Close()
			close(msgChan)
			log.Println("Kafka consumer stopped")
		}()

		// канал для получения системных сигналов для graceful shutdown
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

		run := true
		for run {
			select {
			case sig := <-sigchan:
				log.Printf("Caught signal %v: terminating consumer", sig)
				run = false // конструция для закрытия цикла
			case <-ctx.Done():
				log.Println("Context cancelled: terminating consumer")
				run = false // конструция для закрытия цикла
			default:
				// получаем сообщение из кафки с таймаутом 100 сек
				ev := consumer.Poll(100)
				if ev == nil {
					continue // если ничего непришло, ждем дальше
				}

				switch e := ev.(type) {
				case *kafka.Message:
					select {
					case msgChan <- e.Value:
						if _, err := consumer.CommitMessage(e); err != nil {
							log.Printf("Failed to commit message: %v", err)
						}
					case <-ctx.Done():
						run = false // конструция для закрытия цикла
					}
				case kafka.Error:
					log.Printf("Kafka error: %v", e)
					// если ошибка фатальная — завершаем работу consumer
					if e.IsFatal() {
						run = false // конструция для закрытия цикла
					}
				}
			}
		}
	}()

	return msgChan, nil
}

func CreateTopic(brokers, topicName string, numPartitions, replicationFactor int) error {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return err
	}
	defer adminClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	topicSpec := kafka.TopicSpecification{
		Topic:             topicName,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	results, err := adminClient.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{topicSpec},
		kafka.SetAdminOperationTimeout(30*time.Second),
	)
	if err != nil {
		return err
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError && result.Error.Code() != kafka.ErrTopicAlreadyExists {
			return result.Error
		}
		if result.Error.Code() == kafka.ErrTopicAlreadyExists {
			log.Printf("Topic %s already exists", topicName)
		} else {
			log.Printf("Topic %s created successfully", topicName)
		}
	}

	return nil
}
