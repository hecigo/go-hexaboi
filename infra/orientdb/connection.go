package orientdb

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-resty/resty/v2"
	"hoangphuc.tech/go-hexaboi/infra/core"
)

type Config struct {
	ConnectionName      string
	BaseURL             string
	Database            string
	Authen              string
	Timeout             time.Duration
	MaxRetries          int
	RetryWaitTimeout    time.Duration
	RetryMaxWaitTimeout time.Duration
	IsDebug             bool
}

var clients map[string]*resty.Client = make(map[string]*resty.Client)

// Get the default OrientDB client
func Client() *resty.Client {
	if len(clients) == 0 {
		panic("No client found")
	}
	return clients["default"]
}

// Get the default OrientDB client by name
func ClientByName(name string) *resty.Client {
	if len(clients) == 0 {
		panic("No client found")
	}
	return clients[name]
}

// Open the default connection to OrientDB
func OpenDefaultConnection() error {
	err := OpenConnectionByName("")
	if err == nil {
		return nil
	}

	log.Error(err)
	log.Fatal("Force stop application, cause of the default OrientDB connection has an error.")
	return err
}

func OpenConnectionByName(connName string) error {
	_connName := ""     // Emtpy is default connection
	if connName != "" { // Add _ as a prefix to the connection name
		_connName = "_" + connName
	}

	baseURL := core.Getenv(fmt.Sprintf("ORIENTDB%s_URL", _connName), "")
	database := core.Getenv(fmt.Sprintf("ORIENTDB%s_DB", _connName), "")
	auth := core.Getenv(fmt.Sprintf("ORIENTDB%s_AUTH", _connName), "")
	timeout := core.GetDurationEnv(fmt.Sprintf("ORIENTDB%s_TIMEOUT", _connName), 10*time.Second)
	maxRetries := core.GetIntEnv(fmt.Sprintf("ORIENTDB%s_MAX_RETRIES", _connName), 3)
	retryWaitTimeout := core.GetDurationEnv(fmt.Sprintf("ORIENTDB%s_RETRY_WAIT_TIMEOUT", _connName), 250*time.Millisecond)
	retryMaxWaitTimeout := core.GetDurationEnv(fmt.Sprintf("ORIENTDB%s_RETRY_MAX_WAIT_TIMEOUT", _connName), 3*time.Second)
	isDebug := core.GetBoolEnv(fmt.Sprintf("ORIENTDB%s_DEBUG", _connName), false)

	// Generate the default name as a key for the DB map
	if connName == "" {
		connName = "default"
	}

	err := OpenConnection(Config{
		ConnectionName:      connName,
		BaseURL:             baseURL,
		Database:            database,
		Authen:              auth,
		Timeout:             timeout,
		MaxRetries:          maxRetries,
		RetryWaitTimeout:    retryWaitTimeout,
		RetryMaxWaitTimeout: retryMaxWaitTimeout,
		IsDebug:             isDebug,
	})

	return err
}

func OpenConnection(config ...Config) error {
	for _, cfg := range config {
		client := resty.New()
		client.SetBaseURL(cfg.BaseURL).
			SetTimeout(cfg.Timeout).
			SetRetryCount(cfg.MaxRetries).
			SetRetryWaitTime(cfg.RetryWaitTimeout).
			SetRetryMaxWaitTime(cfg.RetryMaxWaitTimeout).
			SetContentLength(true).
			SetHeaders(map[string]string{
				"Content-Type":    "application/json",
				"Accept-Encoding": "gzip,deflate",
			})

		if strings.Contains(cfg.Authen, ":") {
			basicAuth := strings.Split(cfg.Authen, ":")
			client.SetBasicAuth(basicAuth[0], basicAuth[1])
		} else if cfg.Authen != "" {
			client.SetAuthToken(cfg.Authen)
		}

		if cfg.IsDebug {
			client.EnableTrace().SetDebug(true)
		}

		// Connect
		resp, err := client.R().Get("/connect/" + cfg.Database)
		if err != nil {
			log.Fatal(err)
			continue
		} else {
			// Receives OSSESSIONID of OrientDB
			if resp.StatusCode() == http.StatusNoContent {
				client.SetCookies(resp.Cookies())
			} else {
				log.Fatalf("can not connect to OrientDB of %s", cfg.ConnectionName)
				if cfg.IsDebug {
					log.Debug(resp)
				}
				continue
			}
		}

		if clients[cfg.ConnectionName] != nil {
			Close(client)
			clients[cfg.ConnectionName] = nil
		}
		clients[cfg.ConnectionName] = client
		Print(cfg, *client)
	}

	return nil
}

func Close(client *resty.Client) error {
	_, err := client.R().SetHeader("Connection", "close").Get("/disconnect")
	if err != nil {
		return err
	}
	return nil
}

func CloseAll() error {
	for _, client := range clients {
		// By pass error to continue the next connection
		Close(client)
	}
	log.Info("Closed OrientDB connections...")
	return nil
}

func Print(cfg Config, client resty.Client) {
	_connName := cfg.ConnectionName
	if cfg.ConnectionName == "default" {
		_connName = ""
	}

	fmt.Printf("\r\n┌─────── OrientDB/%s: Connected ─────────\r\n", cfg.ConnectionName)
	fmt.Printf("| ORIENTDB%s_URL: %s\r\n", _connName, cfg.BaseURL)
	if cfg.Authen != "" {
		fmt.Printf("| ORIENTDB%s_AUTH: %s\r\n", _connName, cfg.Authen)
	}
	fmt.Printf("| ORIENTDB%s_COOKIES: %v\r\n", _connName, client.Cookies)
	fmt.Printf("| ORIENTDB%s_TIMEOUT: %d\r\n", _connName, cfg.Timeout)
	fmt.Printf("| ORIENTDB%s_MAX_RETRIES: %d\r\n", _connName, cfg.MaxRetries)
	fmt.Printf("| ORIENTDB%s_RETRY_WAIT_TIMEOUT: %d\r\n", _connName, cfg.RetryWaitTimeout)
	fmt.Printf("| ORIENTDB%s_RETRY_MAX_WAIT_TIMEOUT: %d\r\n", _connName, cfg.RetryMaxWaitTimeout)
	fmt.Printf("| ORIENTDB%s_DEBUG: %v\r\n", _connName, cfg.IsDebug)
	fmt.Println("└──────────────────────────────────────────────")
}
