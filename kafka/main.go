package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

const (
	topic         = "test-topic"
	brokerAddress = "localhost:9092"
)

func main() {
	// Create a context that will be cancelled on SIGINT/SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown gracefully
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	go func() {
		<-signals
		cancel()
	}()

	var wg sync.WaitGroup

	// Start producer and consumer in separate goroutines
	//wg.Add(1)
	//go runProducer(ctx, &wg)
	wg.Add(1)
	go runConsumer(ctx, &wg)

	wg.Wait()
}

func runProducer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// Producer config
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// Create producer
	producer, err := sarama.NewSyncProducer([]string{brokerAddress}, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Send messages until context is cancelled
	count := 1
	for {
		select {
		case <-ctx.Done():
			return
		default:
			message := fmt.Sprintf("Message %d", count)
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(message),
			}

			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			} else {
				log.Printf(
					"Message sent: %s, Partition: %d, Offset: %d", message, partition, offset,
				)
			}

			count++
			time.Sleep(time.Second)
		}
	}
}

func runConsumer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// Consumer config
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}

	// Create consumer group
	group := "test-group"
	consumer, err := sarama.NewConsumerGroup([]string{brokerAddress}, group, config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// Consumer handler
	handler := &ConsumerGroupHandler{}

	// Consume messages until context is cancelled
	for {
		err := consumer.Consume(ctx, []string{topic}, handler)
		if err != nil {
			if err == sarama.ErrClosedConsumerGroup {
				return
			}
			log.Printf("Error from consumer: %v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

// ConsumerGroupHandler implements sarama.ConsumerGroupHandler interface
type ConsumerGroupHandler struct{}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim,
) error {
	for message := range claim.Messages() {
		log.Printf(
			"Message received: %s, Partition: %d, Offset: %d",
			string(message.Value), message.Partition, message.Offset,
		)
		session.MarkMessage(message, "")
	}
	return nil
}
