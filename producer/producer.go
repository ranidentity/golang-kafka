package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	topic := "my-topic-segmentio"     // Replace with your desired topic name
	brokerAddress := "localhost:9092" // Replace with your Kafka broker address

	// Create a new Kafka writer (producer)
	// Balancer: &kafka.LeastBytes{} distributes messages among partitions by sending to the partition with the fewest bytes in flight.
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddress),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	defer func() {
		if err := writer.Close(); err != nil {
			log.Fatalf("failed to close writer: %v", err)
		}
	}()

	fmt.Println("Producer starting...")

	for i := 0; i < 10; i++ {
		messageValue := fmt.Sprintf("Hello Kafka from segmentio/kafka-go! Message %d", i)
		messageKey := fmt.Sprintf("key-%d", i)

		msg := kafka.Message{
			Key:   []byte(messageKey),
			Value: []byte(messageValue),
			Headers: []kafka.Header{
				{Key: "source", Value: []byte("go-producer")},
			},
		}

		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			log.Printf("failed to write message: %v", err)
		} else {
			fmt.Printf("Produced message: Key=%s, Value=%s\n", messageKey, messageValue)
		}
		time.Sleep(100 * time.Millisecond) // Simulate some work
	}

	fmt.Println("Producer finished.")
}
