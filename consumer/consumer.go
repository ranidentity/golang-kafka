package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	topic := "my-topic-segmentio"     // Replace with the topic you want to consume from
	brokerAddress := "localhost:9092" // Replace with your Kafka broker address
	groupID := "my-segmentio-group"   // Replace with your desired consumer group ID

	// Create a new Kafka reader (consumer)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		// StartOffset: kafka.FirstOffset, // Uncomment to start from the beginning of the topic
		ReadLagInterval: -1, // Disable lag tracking for simplicity, or set a duration like 5 * time.Second
		MaxAttempts:     3,  // Retry up to 3 times on temporary errors
	})
	defer func() {
		if err := reader.Close(); err != nil {
			log.Fatalf("failed to close reader: %v", err)
		}
	}()

	fmt.Println("Consumer starting...")

	// Set up a channel to handle graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigchan
		fmt.Println("\nCaught signal, shutting down consumer...")
		cancel() // Signal context cancellation
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context cancelled, consumer exiting.")
			return
		default:
			// ReadMessage blocks until a message is received or the context is cancelled.
			// You can also use ReadMessage with a context.WithTimeout if you prefer a polling mechanism with a hard timeout.
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				// If the context was cancelled, we expect this error
				if err == context.Canceled {
					return
				}
				log.Printf("failed to read message: %v", err)
				time.Sleep(time.Second) // Small delay before retrying on other errors
				continue
			}

			fmt.Printf("Consumed message from topic %s, partition %d, offset %d: Key=%s, Value=%s\n",
				m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

			// Manually commit offsets if you set `CommitInterval: 0` in ReaderConfig
			// For automatic commits (default), you don't need this.
			// err = reader.CommitMessages(ctx, m)
			// if err != nil {
			// 	log.Printf("failed to commit offset: %v", err)
			// }
		}
	}
}
