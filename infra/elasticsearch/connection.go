package elasticsearch

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/elastic/go-elasticsearch/v7"
	"hoangphuc.tech/go-hexaboi/infra/core"
)

type Config struct {
	ConnectionName    string
	Addresses         []string
	BasicAuth         []string
	MaxRetries        int
	EnableDebugLogger bool
}

var clients map[string]*elasticsearch.Client = make(map[string]*elasticsearch.Client)

// Get the default Elasticsearch client
func Client() *elasticsearch.Client {
	if len(clients) == 0 {
		panic("No client found")
	}
	return clients["default"]
}

// Get the default Elasticsearch client by name
func ClientByName(name string) *elasticsearch.Client {
	if len(clients) == 0 {
		panic("No client found")
	}
	return clients[name]
}

// Open the default connection to main database.
func OpenDefaultConnection() error {
	err := OpenConnectionByName("")
	if err == nil {
		return nil
	}

	log.Error(err)
	log.Fatal("Force stop application, cause of the DB connection has an error.")
	return err
}

func OpenConnectionByName(connName string) error {
	_connName := ""     // Emtpy is default connection
	if connName != "" { // Add _ as a prefix to the connection name
		_connName = "_" + connName
	}

	addresses := core.Getenv(fmt.Sprintf("ELASTICSEARCH%s_URL", _connName), "")
	basicAuth := core.Getenv(fmt.Sprintf("ELASTICSEARCH%s_BASIC_AUTH", _connName), "")
	maxRetries := core.GetIntEnv(fmt.Sprintf("ELASTICSEARCH%s_MAX_RETRIES", _connName), 3)
	enableDebugLogger := core.GetBoolEnv(fmt.Sprintf("ELASTICSEARCH%s_DEBUG", _connName), false)

	// Generate the default name as a key for the DB map
	if connName == "" {
		connName = "default"
	}

	err := OpenConnection(Config{
		ConnectionName:    connName,
		Addresses:         strings.Split(addresses, ";"),
		BasicAuth:         strings.Split(basicAuth, ":"),
		MaxRetries:        maxRetries,
		EnableDebugLogger: enableDebugLogger,
	})

	return err
}

func OpenConnection(config ...Config) error {
	for _, cfg := range config {
		esCfg := elasticsearch.Config{
			Addresses:         cfg.Addresses,
			MaxRetries:        cfg.MaxRetries,
			EnableDebugLogger: cfg.EnableDebugLogger,
		}

		if cfg.BasicAuth != nil && len(cfg.BasicAuth) == 2 {
			esCfg.Username = cfg.BasicAuth[0]
			esCfg.Password = cfg.BasicAuth[1]
		}

		client, err := elasticsearch.NewClient(esCfg)
		if err != nil {
			return err
		}

		// Clear existed connection to renew
		if clients[cfg.ConnectionName] != nil {
			clients[cfg.ConnectionName] = nil
		}

		clients[cfg.ConnectionName] = client
		Print(cfg)
	}

	return nil
}

func Print(cfg Config) {
	_connName := cfg.ConnectionName
	if cfg.ConnectionName == "default" {
		_connName = ""
	}

	fmt.Printf("\r\n┌─────── Elasticsearch/%s: Ready ─────────\r\n", cfg.ConnectionName)
	fmt.Printf("| ELASTICSEARCH%s_URL: %s\r\n", _connName, cfg.Addresses)
	if cfg.BasicAuth != nil && len(cfg.BasicAuth) == 2 {
		fmt.Printf("| ELASTICSEARCH%s_BASIC_AUTH: %s\r\n", _connName, cfg.BasicAuth)
	}
	fmt.Printf("| ELASTICSEARCH%s_MAX_RETRIES: %d\r\n", _connName, cfg.MaxRetries)
	fmt.Printf("| ELASTICSEARCH%s_DEBUG: %v\r\n", _connName, cfg.EnableDebugLogger)
	fmt.Println("└──────────────────────────────────────────────")

}
