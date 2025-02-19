package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"movieexample.com/rating/pkg/model"

	"github.com/IBM/sarama"
)

// Ingester defines a Kafka ingester.
type Ingester struct {
	consumer sarama.Consumer
	topic    string
}

// NewIngester creates a new Kafka ingester.
func NewIngester(brokers []string, topic string) (*Ingester, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}
	return &Ingester{consumer, topic}, nil
}

// Ingest starts ingestion from Kafka and returns a channel containing rating events
// representing the data consumed from the topic.
func (i *Ingester) Ingest(ctx context.Context) (chan model.RatingEvent, error) {
	fmt.Println("Starting Kafka ingester")

	consumer, err := i.consumer.ConsumePartition(i.topic, 0, sarama.OffsetOldest)
	if err != nil {
		return nil, err
	}
	fmt.Println("Consumer started")

	ch := make(chan model.RatingEvent, 1)
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Printf("Error with consumer: %v", err)
			case msg := <-consumer.Messages():
				fmt.Printf("Received order: Topic(%s) | Message(%s)\n", msg.Topic, string(msg.Value))
				var event model.RatingEvent
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					fmt.Printf("Unmarshal error: %v\n", err)
					continue
				}
				ch <- event
			case <-ctx.Done():
				close(ch)
				i.consumer.Close()
				return
			}
		}
	}()
	return ch, nil
}
