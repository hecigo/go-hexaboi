package kafka

import (
	"context"
	"errors"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
	Config Config
}

func NewProducer(connName string) *Producer {
	return &Producer{
		Writer: WriterByName(connName),
		Config: GetConfig(connName),
	}
}

func (kw *Producer) Produce(msgs ...kafka.Message) error {
	var err error
	for i := 0; i < kw.Config.MaxAttempts; i++ {
		// Cancel write if it spents too long time
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// attempt to create topic prior to publishing the message
		err = kw.Writer.WriteMessages(ctx, msgs...)
		if err == nil {
			return nil
		}

		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			cancel()
			time.Sleep(kw.Config.AttemptWait)
			continue
		}

		if err != nil {
			cancel()
			return err
		}
	}

	return err
}
