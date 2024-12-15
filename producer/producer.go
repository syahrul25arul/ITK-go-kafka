package producer

import (
	"fmt"
	"go-kafka-consumer/logger"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type OrderPlacer struct {
	producer   *kafka.Producer
	topic      string
	deliverych chan kafka.Event
}

func NewProducer(p *kafka.Producer, topic string) *OrderPlacer {
	return &OrderPlacer{
		producer:   p,
		topic:      topic,
		deliverych: make(chan kafka.Event, 10000),
	}
}

func (op *OrderPlacer) placeOrder(orderType string, size int) {
	var (
		format  = fmt.Sprintf("%s - %d", orderType, size)
		payload = []byte(format)
	)
	err := op.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &op.topic,
			Partition: kafka.PartitionAny,
		},
		Value: payload,
	}, op.deliverych)

	if err != nil {
		logger.Fatal(fmt.Sprintf("error produce data to topic : %s", err.Error()), nil)
	}
}

func main() {
	kafkaInstance, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "producer-2",
		"acks":              "all",
	})

	if err != nil {
		logger.Fatal(fmt.Sprintf("application terminate because error : %s", err.Error()), nil)
	}

	producer := NewProducer(kafkaInstance, "testoh")
	defer producer.producer.Flush(100)
	defer producer.producer.Close()

	for i := 0; i < 1000; i++ {
		producer.placeOrder("producer-2", i+1)
		e := <-producer.deliverych
		logger.Info(fmt.Sprintf("produce data to kafka : %s", e.String()), nil)

		time.Sleep(2 * time.Second)
	}
}
