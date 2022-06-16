package bigquery

import (
	"fmt"
	"log"

	"gorm.io/driver/bigquery"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"hoangphuc.tech/go-hexaboi/infra/core"
)

type Config struct {
	ConnectionName string
	DSN            string
}

var databases map[string]*gorm.DB = make(map[string]*gorm.DB)

// Get default database
func DB() *gorm.DB {
	if len(databases) == 0 {
		panic("No database found")
	}
	return databases["default"]
}

// Get a database by name
func DBByName(name string) *gorm.DB {
	if len(databases) == 0 {
		panic("No database found")
	}
	return databases[name]
}

func OpenConnectionByName(connName string) error {
	_connName := ""     // Emtpy is default connection
	if connName != "" { // Add _ as a prefix to the connection name
		_connName = "_" + connName
	}

	dsn := core.Getenv(fmt.Sprintf("DB_BIGQUERY%s_DSN", _connName), "")

	// Generate the default name as a key for the DB map
	if connName == "" {
		connName = "default"
	}

	err := OpenConnection(Config{
		ConnectionName: connName,
		DSN:            dsn,
	})

	return err
}

// Open the default connection to main database.
func OpenDefaultConnection() error {
	err := OpenConnectionByName("")
	if err == nil {
		return nil
	}

	log.Println(fmt.Errorf("%w", err))
	log.Fatal("Force stop application, cause of the DB connection has an error.")
	return err
}

func OpenConnection(config ...Config) error {
	for _, cfg := range config {
		db, err := gorm.Open(bigquery.Open(cfg.DSN), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // Always use singular table name
			},
			PrepareStmt: true,
		})

		if err != nil {
			return err
		}

		// Append to DB map
		if databases == nil {
			databases = make(map[string]*gorm.DB)
		}

		// Clear existed connection to renew
		if databases[cfg.ConnectionName] != nil {
			Close(databases[cfg.ConnectionName])
			databases[cfg.ConnectionName] = nil
		}

		databases[cfg.ConnectionName] = db
		Print(cfg)
	}

	return nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func CloseAll() error {
	log.Println("Closing BigQuery connections...")
	for _, db := range databases {
		// By pass error to continue the next connection
		Close(db)
	}
	return nil
}

func Print(cfg Config) {
	_connName := cfg.ConnectionName
	if cfg.ConnectionName == "default" {
		_connName = ""
	}

	fmt.Printf("\r\n┌─────── BigQuery/%s: Connected ─────────\r\n", _connName)
	fmt.Printf("| BIGQUERY%s_DSN: %s\r\n", _connName, cfg.DSN)
	fmt.Println("└──────────────────────────────────")

}
