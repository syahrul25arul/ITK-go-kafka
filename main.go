package main

import (
	"fmt"
	"go-kafka-consumer/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "consumer1-copy",
		"auto.offset.reset": "smallest",
	})

	if err != nil {
		logger.Fatal(fmt.Sprintf("application terminate because error : %s", err.Error()), nil)
	}

	topic := "verification"
	err = c.Subscribe(topic, nil)
	if err != nil {
		logger.Fatal(fmt.Sprintf("failed subscribe topic %s because : %s", topic, err.Error()), nil)
	}

	for {
		ev := c.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			_, er := c.CommitMessage(e)
			if err != nil {
				logger.Error(fmt.Sprintf("Error commit message : %s", er.Error()), nil)
				return
			}

			fmt.Printf("consume from the topic : %s and partion : %v\n", string(e.Value), e.TopicPartition.Partition)
		case *kafka.Error:
			fmt.Printf("%v\n", e)
		}
	}
}
