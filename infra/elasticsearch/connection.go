package elasticsearch

import (
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v7"
	"hoangphuc.tech/go-hexaboi/infra/core"
)

type Config struct {
	ConnectionName    string
	Addresses         []string
	BasicAuth         []string
	MaxRetries        int
	SearchTimeout     time.Duration
	BatchIndexSize    int
	EnableDebugLogger bool
}

var clients map[string]*elasticsearch.Client = make(map[string]*elasticsearch.Client)
var configs map[string]Config = make(map[string]Config)

// Get the default Elasticsearch client
func Client() *elasticsearch.Client {
	return ClientByName("default")
}

// Get the default Elasticsearch client by name
func ClientByName(name string) *elasticsearch.Client {
	if len(clients) == 0 {
		panic("No client found")
	}
	return clients[name]
}

func GetConfig() Config {
	return GetConfigByName("default")
}

func GetConfigByName(name string) Config {
	if len(configs) == 0 {
		panic("No config found")
	}
	return configs[name]
}

// Initialize the default Elasticsearch client
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
	searchTimeout := core.GetDurationEnv(fmt.Sprintf("ELASTICSEARCH%s_SEARCH_TIMEOUT", _connName), 5*time.Second)
	batchIndexSize := core.GetIntEnv(fmt.Sprintf("ELASTICSEARCH%s_BATCH_INDEX_SIZE", _connName), 100)
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
		SearchTimeout:     searchTimeout,
		BatchIndexSize:    batchIndexSize,
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
			Transport:         &Transport{},
		}

		if cfg.EnableDebugLogger {
			esCfg.Logger = &elastictransport.ColorLogger{
				Output:             os.Stdout,
				EnableRequestBody:  true,
				EnableResponseBody: true,
			}
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
		configs[cfg.ConnectionName] = cfg
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
	fmt.Printf("| ELASTICSEARCH%s_BATCH_INDEX_SIZE: %d\r\n", _connName, cfg.BatchIndexSize)
	fmt.Printf("| ELASTICSEARCH%s_DEBUG: %v\r\n", _connName, cfg.EnableDebugLogger)
	fmt.Println("└──────────────────────────────────────────────")

}
