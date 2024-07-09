package retrier

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Ashwin7867/go-kafka-retrier/config"
	"github.com/IBM/sarama"
)

type Retrier struct {
	config   *config.Config
	consumer sarama.ConsumerGroup
}

func NewRetrier(config *config.Config) (*Retrier, error) {
	consumer, err := sarama.NewConsumerGroup([]string{config.KafkaBrokers}, config.ConsumerGroup, nil)
	if err != nil {
		return nil, err
	}

	return &Retrier{
		config:   config,
		consumer: consumer,
	}, nil
}

func (r *Retrier) ProcessMessages(ctx context.Context) error {
	handler := &messageHandler{
		config: r.config,
		handle: r.HandleMessage,
	}

	for {
		err := r.consumer.Consume(ctx, []string{r.config.Topic}, handler)
		if err != nil {
			log.Printf("Error from consumer: %v", err)
			time.Sleep(1 * time.Second)
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (r *Retrier) HandleMessage(msg *sarama.ConsumerMessage) error {
	// Implement your message handling logic here
	return fmt.Errorf("dummy error")
}
