package retrier

import (
	"log"
	"time"

	"github.com/Ashwin7867/go-kafka-retrier/config"
	"github.com/IBM/sarama"
)

type messageHandler struct {
	config *config.Config
	handle func(msg *sarama.ConsumerMessage) error
}

func (h *messageHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *messageHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *messageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		retryCount := 0
		for {
			err := h.handle(msg)
			if err == nil {
				sess.MarkMessage(msg, "")
				break
			}

			if retryCount >= h.config.MaxRetries {
				log.Printf("Max retries reached for message: %s", string(msg.Value))
				producer, _ := sarama.NewSyncProducer([]string{h.config.KafkaBrokers}, nil)
				dlqMsg := &sarama.ProducerMessage{
					Topic: h.config.DLQTopic,
					Value: sarama.StringEncoder(msg.Value),
				}
				producer.SendMessage(dlqMsg)
				break
			}

			log.Printf("Error processing message, retrying: %v", err)
			time.Sleep(h.config.RetryDelay)
			retryCount++
		}
	}
	return nil
}
