package base

type KafkaPayload interface {
	ToPayload() map[string]interface{}
	ClassName() string
}
