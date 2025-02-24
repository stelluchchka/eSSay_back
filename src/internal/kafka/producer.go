package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"essay/src/internal/models"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ProduceEssay sends essay in Kafka queue
func ProduceEssay(broker, topic string, essay models.Essay) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"acks":              "all",
	})
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}
	defer p.Close()

	essayJSON, err := json.Marshal(essay)
	if err != nil {
		return fmt.Errorf("failed to marshal essay: %w", err)
	}

	deliveryChan := make(chan kafka.Event)

	log.Printf("Producing message to topic %s: %s", topic, string(essayJSON))
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          essayJSON,
	}, deliveryChan)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
	}

	log.Printf("Message delivered to topic %s [%d] at offset %d\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	return nil
}
