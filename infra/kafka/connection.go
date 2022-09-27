package kafka

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"hoangphuc.tech/go-hexaboi/infra/core"
)

type Config struct {
	ConnectionName         string
	Brokers                []string
	SASLMechanism          string
	DialTimeout            time.Duration
	IncomingTopic          string
	IncomingConsumerGroup  string
	OutgoingTopic          string
	Logger                 kafka.Logger
	MaxAttempts            int
	AttemptWait            time.Duration
	AllowAutoTopicCreation bool
}

var readers map[string]*kafka.Reader = make(map[string]*kafka.Reader)
var writers map[string]*kafka.Writer = make(map[string]*kafka.Writer)
var configs map[string]Config = make(map[string]Config)

// Get the default Kafka topic reader
func Reader() *kafka.Reader {
	if len(readers) == 0 {
		panic("No reader found")
	}
	return readers["default"]
}

func Writer() *kafka.Writer {
	if len(writers) == 0 {
		panic("No writer found")
	}
	return writers["default"]
}

// Get a reader by name
func ReaderByName(name string) *kafka.Reader {
	if len(readers) == 0 {
		panic("No reader found")
	}
	if name == "" {
		return readers["default"]
	}
	return readers[name]
}

func WriterByName(name string) *kafka.Writer {
	if len(writers) == 0 {
		panic("No writer found")
	}
	if name == "" {
		return writers["default"]
	}
	return writers[name]
}

func GetConfig(name string) Config {
	if len(configs) == 0 {
		panic("No config found")
	}
	if name == "" {
		return configs["default"]
	}
	return configs[name]
}

// Open the default connection to main Kafka topic.
func OpenDefaultConnection() error {
	err := OpenConnectionByName("")
	if err == nil {
		return nil
	}

	log.Println(fmt.Errorf("%w", err))
	log.Fatal("Force stop application, cause of the Kafka connection has an error.")
	return err
}

// Open a connection with name specified from ENV
func OpenConnectionByName(connName string) error {
	_connName := ""     // Emtpy is default connection
	if connName != "" { // Add _ as a prefix to the connection name
		_connName = "_" + connName
	}
	_connName = strings.ToUpper(_connName)

	// Generate the default name as a key for the Kafka map
	if connName == "" {
		connName = "default"
	}

	brokers := core.Getenv("KAFKA_BROKERS", "")
	sasl := core.Getenv("KAFKA_SASL", "")
	showLog := core.GetBoolEnv("KAFKA_SHOW_LOG", false)
	dialTimeout := core.GetDurationEnv("KAFKA_DIAL_TIMEOUT", 3*time.Minute)
	incomingGroupId := core.Getenv("KAFKA_INCOMING_CONSUMER_GROUP", "")
	incomingTopic := core.Getenv("KAFKA_INCOMING_TOPIC", "")
	outgoingTopic := core.Getenv("KAFKA_OUTGOING_TOPIC", "")
	retry := core.GetIntEnv("KAFKA_MAX_ATTEMPTS", 3)
	attemptWait := core.GetDurationEnv("KAFKA_ATTEMPT_WAIT", 100*time.Millisecond)
	allowAutoTopicCreation := core.GetBoolEnv("KAFKA_ALLOW_AUTO_TOPIC_CREATION", false)

	// Override if specialize a profile
	if _connName != "" {
		profileBrokers := core.Getenv(fmt.Sprintf("KAFKA%s_BROKERS", _connName), "")
		if profileBrokers != "" {
			brokers = profileBrokers
		}

		profileSasl := core.Getenv(fmt.Sprintf("KAFKA%s_SASL", _connName), "")
		if profileSasl != "" {
			sasl = profileSasl
		}

		profileShowLog := core.Getenv(fmt.Sprintf("KAFKA%s_SHOW_LOG", _connName), "")
		if profileShowLog != "" {
			val, e := strconv.ParseBool(profileShowLog)
			if e == nil {
				showLog = val
			}
		}

		profileDialTimeout := core.Getenv(fmt.Sprintf("KAFKA%s_DIAL_TIMEOUT", _connName), "")
		if profileDialTimeout != "" {
			val, e := time.ParseDuration(profileDialTimeout)
			if e == nil {
				dialTimeout = val
			}
		}

		profileIncomingGroupId := core.Getenv(fmt.Sprintf("KAFKA%s_INCOMING_CONSUMER_GROUP", _connName), "")
		if profileIncomingGroupId != "" {
			incomingGroupId = profileIncomingGroupId
		}

		profileIncomingTopic := core.Getenv(fmt.Sprintf("KAFKA%s_INCOMING_TOPIC", _connName), "")
		if profileIncomingTopic != "" {
			incomingTopic = profileIncomingTopic
		}

		profileOutcomingTopic := core.Getenv(fmt.Sprintf("KAFKA%s_OUTGOING_TOPIC", _connName), "")
		if profileOutcomingTopic != "" {
			outgoingTopic = profileOutcomingTopic
		}

		profileRetry := core.Getenv(fmt.Sprintf("KAFKA%s_MAX_ATTEMPTS", _connName), "")
		if profileRetry != "" {
			val, e := strconv.Atoi(profileRetry)
			if e == nil {
				retry = val
			}
		}

		profileAttempWait := core.Getenv(fmt.Sprintf("KAFKA%s_ATTEMPT_WAIT", _connName), "")
		if profileAttempWait != "" {
			val, e := time.ParseDuration(profileAttempWait)
			if e == nil {
				attemptWait = val
			}
		}

		profileAllowAutoTopicCreation := core.Getenv(fmt.Sprintf("KAFKA%s_ALLOW_AUTO_TOPIC_CREATION", _connName), "")
		if profileAllowAutoTopicCreation != "" {
			val, e := strconv.ParseBool(profileAllowAutoTopicCreation)
			if e == nil {
				allowAutoTopicCreation = val
			}
		}
	}

	// Always requires consumer group ID if incoming topic is specified
	if incomingTopic != "" && incomingGroupId == "" {
		return fmt.Errorf("consumer group ID of %s is required", connName)
	}

	// Init logger
	var logger kafka.Logger
	if showLog {
		logger = kafka.LoggerFunc(logf)
	}

	err := OpenConnection(Config{
		ConnectionName:         connName,
		Brokers:                strings.Split(brokers, ";"),
		SASLMechanism:          sasl,
		DialTimeout:            dialTimeout,
		IncomingTopic:          incomingTopic,
		IncomingConsumerGroup:  incomingGroupId,
		OutgoingTopic:          outgoingTopic,
		Logger:                 logger,
		MaxAttempts:            retry,
		AttemptWait:            attemptWait,
		AllowAutoTopicCreation: allowAutoTopicCreation,
	})

	return err
}

func OpenConnection(config ...Config) error {
	for _, cfg := range config {

		var mechanism sasl.Mechanism
		if cfg.SASLMechanism != "" {
			mecConfig := strings.Split(cfg.SASLMechanism, ":")
			if len(mecConfig) != 3 {
				return errors.New("SASLMechanism config is wrong")
			}
			if strings.HasPrefix(cfg.SASLMechanism, "PLAIN") {
				mechanism = plain.Mechanism{
					Username: strings.TrimLeft(mecConfig[1], "/"),
					Password: mecConfig[2],
				}
			} else if strings.HasPrefix(cfg.SASLMechanism, "SCRAM") {
				m, err := scram.Mechanism(scram.SHA512, strings.TrimLeft(mecConfig[1], "/"), mecConfig[2])
				if err != nil {
					return err
				}
				mechanism = m
			}
		}

		// Init dialer
		dialer := &kafka.Dialer{
			Timeout:       cfg.DialTimeout,
			DualStack:     true,
			SASLMechanism: mechanism,
		}

		if cfg.IncomingTopic != "" {
			reader := kafka.NewReader(kafka.ReaderConfig{
				Dialer:         dialer,
				Brokers:        cfg.Brokers,
				Topic:          cfg.IncomingTopic,
				StartOffset:    kafka.LastOffset,
				GroupID:        cfg.IncomingConsumerGroup,
				CommitInterval: time.Second,
				MaxAttempts:    cfg.MaxAttempts,
				Logger:         cfg.Logger,
				ErrorLogger:    cfg.Logger,
			})

			if readers == nil {
				readers = make(map[string]*kafka.Reader)
			}

			if readers[cfg.ConnectionName] != nil {
				CloseReader(cfg.ConnectionName)
			}
			readers[cfg.ConnectionName] = reader
		}

		if cfg.OutgoingTopic != "" {
			writer := kafka.NewWriter(kafka.WriterConfig{
				Dialer:      dialer,
				Brokers:     cfg.Brokers,
				Topic:       cfg.OutgoingTopic,
				Balancer:    &kafka.Hash{},
				MaxAttempts: cfg.MaxAttempts,
				Logger:      cfg.Logger,
				ErrorLogger: cfg.Logger,
			})
			writer.AllowAutoTopicCreation = cfg.AllowAutoTopicCreation

			if writers == nil {
				writers = make(map[string]*kafka.Writer)
			}

			if writers[cfg.ConnectionName] != nil {
				CloseWriter(cfg.ConnectionName)
			}
			writers[cfg.ConnectionName] = writer
		}

		if configs == nil {
			configs = make(map[string]Config)
		}
		configs[cfg.ConnectionName] = cfg
		Print(cfg)
	}
	return nil
}

func logf(msg string, a ...interface{}) {
	log.Printf(msg, a...)
}

func CloseReader(connName string) error {
	reader := readers[connName]

	if reader != nil {
		if err := reader.Close(); err != nil {
			return err
		}
		readers[connName] = nil
		log.Printf("Closed Kafka/%s reader connection\n", connName)
	}

	return nil
}

func CloseWriter(connName string) error {
	writer := writers[connName]

	if writer != nil {
		if err := writer.Close(); err != nil {
			return err
		}
		writers[connName] = nil
		log.Printf("Closed Kafka/%s writer connection\n", connName)
	}

	return nil
}

func CloseAll() error {
	for connName := range readers {
		if err := CloseReader(connName); err != nil {
			return err
		}
	}

	// Althrough number of readers is greater than or equal to writers,
	// so we still scan writers for sure.
	for connName := range writers {
		if err := CloseWriter(connName); err != nil {
			return err
		}
	}

	return nil
}

func Print(cfg Config) {
	_connName := strings.ToUpper(cfg.ConnectionName)
	if cfg.ConnectionName == "DEFAULT" {
		_connName = ""
	} else if cfg.ConnectionName != "" {
		_connName = "_" + _connName
	}

	fmt.Printf("\r\n┌─────── Kafka/%s ─────────────────────\r\n", cfg.ConnectionName)
	fmt.Printf("│ %s: %s\r\n", fmt.Sprintf("KAFKA%s_BROKERS", _connName), cfg.Brokers)
	fmt.Printf("│ %s: %s\r\n", fmt.Sprintf("KAFKA%s_SASL", _connName), cfg.SASLMechanism)
	fmt.Printf("│ %s: %v\r\n", fmt.Sprintf("KAFKA%s_DIAL_TIMEOUT", _connName), cfg.DialTimeout)
	fmt.Printf("│ %s: %s\r\n", fmt.Sprintf("KAFKA%s_INCOMING_CONSUMER_GROUP", _connName), cfg.IncomingConsumerGroup)
	fmt.Printf("│ %s: %s\r\n", fmt.Sprintf("KAFKA%s_INCOMING_TOPIC", _connName), cfg.IncomingTopic)
	fmt.Printf("│ %s: %s\r\n", fmt.Sprintf("KAFKA%s_OUTGOING_TOPIC", _connName), cfg.OutgoingTopic)
	fmt.Printf("│ %s: %d\r\n", fmt.Sprintf("KAFKA%s_MAX_ATTEMPTS", _connName), cfg.MaxAttempts)
	fmt.Println("└───────────────────────────────────────────")

}
