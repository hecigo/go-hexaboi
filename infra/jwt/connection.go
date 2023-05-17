package jwt

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/MicahParks/keyfunc"
	"hecigo.com/go-hexaboi/infra/core"
)

var jwksStore map[string]*keyfunc.JWKS = make(map[string]*keyfunc.JWKS)

type Config struct {
	ConnectionName string
	JwksUrl        string
}

// Open the default connection
func OpenDefaultConnection() error {
	err := OpenConnectionByName("")
	if err == nil {
		return nil
	}

	log.Println(fmt.Errorf("%w", err))
	log.Fatal("Force stop application, cause of the JWT connection has an error.")
	return err
}

func OpenConnectionByName(connName string) error {
	_connName := ""     // Emtpy is default connection
	if connName != "" { // Add _ as a prefix to the connection name
		_connName = "_" + connName
	}

	jwksUrl := core.Getenv(fmt.Sprintf("JWT%s_JWKS_URL", _connName), "")

	if connName == "" {
		connName = "default"
	}

	err := OpenConnection(Config{
		ConnectionName: connName,
		JwksUrl:        jwksUrl,
	})

	return err
}

func OpenConnection(config ...Config) error {
	for _, cfg := range config {
		// Create the keyfunc options. Use an error handler that logs. Refresh the JWKS when a JWT signed by an unknown KID
		// is found or at the specified interval. Rate limit these refreshes. Timeout the initial JWKS refresh request after
		// 10 seconds. This timeout is also used to create the initial context.Context for keyfunc.Get.
		options := keyfunc.Options{
			Ctx: context.Background(),
			RefreshErrorHandler: func(err error) {
				log.Printf("There was an error with the jwt.Keyfunc\nError: %s", err.Error())
			},
			RefreshInterval:   time.Hour,
			RefreshRateLimit:  time.Minute * 5,
			RefreshTimeout:    time.Second * 10,
			RefreshUnknownKID: true,
		}

		// Create the JWKS from the resource at the given URL.
		jwks, err := keyfunc.Get(cfg.JwksUrl, options)
		if err != nil {
			log.Fatalf("Failed to get the JWKS from the given URL.\nError: %s", err)
		}

		jwksStore[cfg.ConnectionName] = jwks
		Print(cfg)
	}
	return nil
}

// Get the default JWKS
func JWKS() *keyfunc.JWKS {
	if len(jwksStore) == 0 {
		panic("No JWKS found")
	}
	return jwksStore["default"]
}

// Get a JWKS by name
func JWKSByName(name string) *keyfunc.JWKS {
	if len(jwksStore) == 0 {
		panic("No JWKS found")
	}
	return jwksStore[name]
}

func Close(jwks *keyfunc.JWKS) {
	jwks.EndBackground()
}

func CloseAll() {
	for _, jwks := range jwksStore {
		// By pass error to continue the next connection
		Close(jwks)
	}
	log.Println("Closed JWKS connections...")
}

func Print(cfg Config) {
	_connName := cfg.ConnectionName
	if cfg.ConnectionName == "default" {
		_connName = ""
	}

	fmt.Printf("\r\n┌─────── JWT/%s: Connected ─────────\r\n", cfg.ConnectionName)
	fmt.Printf("│ %s: %s\r\n", fmt.Sprintf("JWT%s_JWKS_URL", _connName), cfg.JwksUrl)
	fmt.Println("└───────────────────────────────────────────────")
}
