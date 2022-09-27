package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	Reader *kafka.Reader
	Config Config
}

func NewConsumer(connName string) *Consumer {
	return &Consumer{
		Reader: ReaderByName(connName),
		Config: GetConfig(connName),
	}
}

func (kw *Consumer) Consume(process func(KafkaMessage) error) error {
	ctx := context.Background()
	for {
		m, err := kw.Reader.ReadMessage(ctx)
		if err == nil {
			fmt.Printf("message at %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, m.Key, string(m.Value))
			msg, err := ScanMessage(m.Value)
			if err != nil {
				// Message structure is invalid
				log.Printf("[WARNING] Message is invalid: %s\n", string(m.Value))
				log.Println(err)
				continue
			}

			if !msg.IsValid() {
				// Message structure is invalid
				log.Printf("[WARNING] Message is invalid: %s\n", string(m.Value))
				continue
			}

			// Max attempts for processing
			for i := 0; i < kw.Config.MaxAttempts; i++ {
				// Ensure error is cleared
				err = nil

				err = process(*msg)
				if err == nil {
					break
				}

				log.Println(err)
				time.Sleep(kw.Config.AttemptWait)
			}
		} else {
			return err
		}
	}
}
