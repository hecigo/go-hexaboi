package redis

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis/v8"

	"hoangphuc.tech/dora/infra/core"
)

type Config struct {
	ConnectionName string
	Addresses      []string
	BasicAuth      []string
	DB             int
	MaxRetries     int
	DialTimeout    time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

var (
	clients map[string]redis.UniversalClient = make(map[string]redis.UniversalClient)
)

// Get the default Elasticsearch client
func Client() redis.UniversalClient {
	if len(clients) == 0 {
		panic("No client found")
	}
	return clients["default"]
}

// Get the default Elasticsearch client by name
func ClientByName(name string) redis.UniversalClient {
	if len(clients) == 0 {
		panic("No client found")
	}
	if clients[name] == nil {
		log.Errorf("Not found client: %s", name)
	}
	return clients[name]
}

// Initialize the default Redis client
func OpenDefaultConnection() error {
	err := OpenConnectionByName("")
	if err == nil {
		return nil
	}

	log.Fatalln("Force stop application, cause of the redis default connection has an error.", err)
	return err
}

func OpenConnectionByName(connName string) error {
	_connName := ""     // Emtpy is default connection
	if connName != "" { // Add _ as a prefix to the connection name
		_connName = "_" + connName
	}

	addresses := core.Getenv(fmt.Sprintf("REDIS%s_URL", _connName), "")
	basicAuth := core.Getenv(fmt.Sprintf("REDIS%s_BASIC_AUTH", _connName), "")
	db := core.GetIntEnv(fmt.Sprintf("REDIS%s_DB", _connName), 0)
	maxRetries := core.GetIntEnv(fmt.Sprintf("REDIS%s_MAX_RETRIES", _connName), 3)
	dialTimeout := core.GetDurationEnv(fmt.Sprintf("REDIS%s_DIAL_TIMEOUT", _connName), time.Second)
	readTimeout := core.GetDurationEnv(fmt.Sprintf("REDIS%s_READ_TIMEOUT", _connName), 3*time.Second)
	writeTimeout := core.GetDurationEnv(fmt.Sprintf("REDIS%s_WRITE_TIMEOUT", _connName), 3*time.Second)

	// Generate the default name as a key for the DB map
	if connName == "" {
		connName = "default"
	}

	err := OpenConnection(Config{
		ConnectionName: connName,
		Addresses:      strings.Split(addresses, ";"),
		BasicAuth:      strings.Split(basicAuth, ":"),
		DB:             db,
		MaxRetries:     maxRetries,
		DialTimeout:    dialTimeout,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
	})

	return err
}

func OpenConnection(config ...Config) error {
	for _, cfg := range config {
		rdsCfg := &redis.UniversalOptions{
			Addrs:        cfg.Addresses,
			MaxRetries:   cfg.MaxRetries,
			DialTimeout:  cfg.DialTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			MasterName:   "primary",
		}

		if cfg.BasicAuth != nil && len(cfg.BasicAuth) == 2 {
			rdsCfg.Username = cfg.BasicAuth[0]
			rdsCfg.Password = cfg.BasicAuth[1]
		}

		client := redis.NewUniversalClient(rdsCfg)

		// Clear existed connection to renew
		if clients[cfg.ConnectionName] != nil {
			Close(clients[cfg.ConnectionName])
			clients[cfg.ConnectionName] = nil
		}
		clients[cfg.ConnectionName] = client

		Print(cfg)
	}

	return nil
}

func Close(client redis.UniversalClient) error {
	return client.Close()
}

func CloseAll() error {
	for _, client := range clients {
		// By pass error to continue the next connection
		Close(client)
	}
	log.Info("Closed Redis clients...")
	return nil
}

func Print(cfg Config) {
	_connName := cfg.ConnectionName
	if cfg.ConnectionName == "default" {
		_connName = ""
	}

	fmt.Printf("\r\n┌─────── REDIS/%s: Ready ─────────\r\n", cfg.ConnectionName)
	fmt.Printf("| REDIS%s_URL: %s\r\n", _connName, cfg.Addresses)
	if cfg.BasicAuth != nil && len(cfg.BasicAuth) == 2 {
		fmt.Printf("| REDIS%s_BASIC_AUTH: %s\r\n", _connName, cfg.BasicAuth)
	}
	fmt.Printf("| REDIS%s_DB: %d\r\n", _connName, cfg.DB)
	fmt.Printf("| REDIS%s_MAX_RETRIES: %d\r\n", _connName, cfg.MaxRetries)
	fmt.Printf("| REDIS%s_DIAL_TIMEOUT: %v\r\n", _connName, cfg.DialTimeout)
	fmt.Printf("| REDIS%s_READ_TIMEOUT: %v\r\n", _connName, cfg.ReadTimeout)
	fmt.Printf("| REDIS%s_WRITE_TIMEOUT: %v\r\n", _connName, cfg.WriteTimeout)
	fmt.Println("└──────────────────────────────────────")

}
