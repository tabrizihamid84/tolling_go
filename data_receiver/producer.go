package main

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fulltimegodev/tolling/types"
)

type DataProducer interface {
	ProduceData(types.OBUData) error
}

type KafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer() (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
	})

	// Start another go routine to check if we have deviverd the data
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Devivered message to %v\n", ev.TopicPartition)
				}
			}

		}
	}()

	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: p,
	}, nil
}

func (p *KafkaProducer) ProduceData(data types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kafkaTopic,
			Partition: kafka.PartitionAny,
		},
		Value: b,
	}, nil)

}
