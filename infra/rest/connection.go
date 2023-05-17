package rest

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-resty/resty/v2"
	"hecigo.com/go-hexaboi/infra/core"
)

type Config struct {
	ConnectionName string
	BaseURL        string
	EnableDebug    bool
	Timeout        time.Duration
	BasicAuth      []string
	BearerAuth     string
	RetryCount     int
}

var clients = make(map[string]*resty.Client)

// Get the default database
func Client() *resty.Client {
	if len(clients) == 0 {
		panic("No client found")
	}
	return clients["default"]
}

// Get a database by name
func ClientByName(name string) *resty.Client {
	if len(clients) == 0 {
		panic("No client found")
	}
	return clients[name]
}

func OpenDefaultConnection() error {
	err := OpenConnectionByName("")
	if err == nil {
		return nil
	}

	log.Errorln(err)
	log.Fatal("Force stop application, cause of the default connection of RESTful Client has an error.")
	return err
}

// Open a connection with name specified from ENV
func OpenConnectionByName(connName string) error {
	_connName := ""     // Emtpy is default connection
	if connName != "" { // Add _ as a prefix to the connection name
		_connName = "_" + connName
	}

	baseUrl := core.Getenv(fmt.Sprintf("REST%s_BASE_URL", _connName), "")
	enableDebug := core.GetBoolEnv(fmt.Sprintf("REST%s_ENABLE_DEBUG", _connName), false)
	timeout := core.GetDurationEnv(fmt.Sprintf("REST%s_TIMEOUT", _connName), time.Second)
	basicAuth := core.Getenv(fmt.Sprintf("REST%s_BASIC_AUTH", _connName), "")
	bearerAuth := core.Getenv(fmt.Sprintf("REST%s_BEARER_AUTH", _connName), "")
	retryCount := core.GetIntEnv(fmt.Sprintf("REST%s_RETRY_COUNT", _connName), 3)

	// Generate the default name as a key for the DB map
	if connName == "" {
		connName = "default"
	}

	err := OpenConnection(Config{
		ConnectionName: connName,
		BaseURL:        baseUrl,
		EnableDebug:    enableDebug,
		Timeout:        timeout,
		BasicAuth:      strings.Split(basicAuth, ":"),
		BearerAuth:     bearerAuth,
		RetryCount:     retryCount,
	})

	return err
}

func OpenConnection(config ...Config) error {
	for _, cfg := range config {
		client := resty.New()
		client.SetHeader("Content-Type", "application/json").
			SetBaseURL(cfg.BaseURL).
			SetDebug(cfg.EnableDebug).
			SetTimeout(cfg.Timeout).
			SetRetryCount(cfg.RetryCount)

		if cfg.BasicAuth != nil && len(cfg.BasicAuth) == 2 {
			client.SetBasicAuth(cfg.BasicAuth[0], cfg.BasicAuth[1])
		}

		if cfg.BearerAuth != "" {
			client.SetAuthToken(cfg.BearerAuth)
		}

		clients[cfg.ConnectionName] = client
	}

	return nil
}
