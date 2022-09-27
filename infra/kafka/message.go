package kafka

import (
	log "github.com/sirupsen/logrus"

	"github.com/goccy/go-json"

	"hoangphuc.tech/go-hexaboi/domain/base"
	"hoangphuc.tech/go-hexaboi/infra/core"

	"github.com/segmentio/kafka-go"
)

type KafkaMessage struct {
	UUID      string                 `json:"uuid"`
	ClassName string                 `json:"class_name"`
	Event     string                 `json:"event"`
	CreatedBy string                 `json:"created_by"`
	Timestamp string                 `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
}

func NewMessage(uuid string, event string, payload base.KafkaPayload) *KafkaMessage {
	return &KafkaMessage{
		UUID:      uuid,
		ClassName: payload.ClassName(),
		Event:     event,
		CreatedBy: core.AppName(),
		Timestamp: core.Utils.NowStr(),
		Payload:   payload.ToPayload(),
	}
}

func ScanMessage(data []byte) (*KafkaMessage, error) {
	var msg *KafkaMessage
	if err := core.UnmarshalNoPanic(data, &msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (me *KafkaMessage) ToBytes() []byte {
	bytes, err := json.Marshal(me)
	if err != nil {
		log.Println(err)
	}
	return bytes
}

func (me *KafkaMessage) ToKafkaGo() kafka.Message {
	return kafka.Message{
		Key:   []byte(me.UUID),
		Value: me.ToBytes(),
	}
}

// Detects if the message is a valid KafkaMessage
func (me *KafkaMessage) IsValid() bool {
	return me.UUID != "" && me.ClassName != "" && me.Event != "" &&
		me.CreatedBy != "" && me.Timestamp != "" && me.Payload != nil
}
