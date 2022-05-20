package postgres

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"hoangphuc.tech/hercules/infra/core"
)

type Config struct {
	ConnectionName        string
	DSN                   []string
	DSNRelicas            []string
	MaxIdleConns          int
	MaxOpenConns          int
	ConnectionMaxLifetime time.Duration
}

var Databases map[string]*gorm.DB = make(map[string]*gorm.DB)

// Get the default database
func DB() *gorm.DB {
	if len(Databases) == 0 {
		panic("No database found")
	}
	return Databases["default"]
}

// Get a database by name
func DBByName(name string) *gorm.DB {
	if len(Databases) == 0 {
		panic("No database found")
	}
	return Databases[name]
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

// Open a connection with name specified from ENV
func OpenConnectionByName(connName string) error {
	_connName := ""     // Emtpy is default connection
	if connName != "" { // Add _ as a prefix to the connection name
		_connName = "_" + connName
	}

	dsn := core.Getenv(fmt.Sprintf("DB_POSTGRES%s_DSN", _connName), "")
	dsnReplicas := core.Getenv(fmt.Sprintf("DB_POSTGRES%s_DSN_REPLICAS", _connName), "")

	// Connection pool config
	maxIdleConns := core.GetIntEnv(fmt.Sprintf("DB_POSTGRES%s_MAX_IDLE_CONNS", _connName), 5)
	maxOpenConns := core.GetIntEnv(fmt.Sprintf("DB_POSTGRES%s_MAX_OPEN_CONNS", _connName), 20)
	connMaxLifetime := core.GetDurationEnv(fmt.Sprintf("DB_POSTGRES%s_CONN_MAX_LIFETIME", _connName), 30*time.Minute)

	// Generate the default name as a key for the DB map
	if connName == "" {
		connName = "default"
	}

	err := OpenConnection(Config{
		ConnectionName:        connName,
		DSN:                   strings.Split(dsn, ";"),
		DSNRelicas:            strings.Split(dsnReplicas, ";"),
		MaxIdleConns:          maxIdleConns,
		MaxOpenConns:          maxOpenConns,
		ConnectionMaxLifetime: connMaxLifetime,
	})

	return err
}

func OpenConnection(config ...Config) error {
	for _, cfg := range config {
		if len(cfg.DSN) <= 0 {
			continue
		}

		db, err := gorm.Open(postgres.Open(cfg.DSN[0]), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // Always use singular table name
			},
		})
		if err != nil {
			return err
		}

		// Supports multi sources and replicas
		var dsnDialectors, dsnRelicasDialectors []gorm.Dialector
		var dbr *dbresolver.DBResolver
		if len(cfg.DSN) > 1 {
			for i := 1; i < len(cfg.DSN); i++ {
				if cfg.DSN[i] != "" {
					dsnDialectors = append(dsnDialectors, postgres.Open(cfg.DSN[i]))
				}
			}

			dbr = dbresolver.Register(dbresolver.Config{
				Sources: dsnDialectors,
			})
		}
		if len(cfg.DSNRelicas) > 0 {
			for _, dsn := range cfg.DSNRelicas {
				if dsn != "" {
					dsnRelicasDialectors = append(dsnRelicasDialectors, postgres.Open(dsn))
				}
			}

			dbr = dbresolver.Register(dbresolver.Config{
				Replicas: dsnRelicasDialectors,
			})
		}

		if len(dsnDialectors) > 0 || len(dsnRelicasDialectors) > 0 {
			dbresolverErr := db.Use(dbr)
			if dbresolverErr != nil {
				return dbresolverErr
			}
		}

		// Set connection pool
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}

		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlDB.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)

		// Append to DB map
		if Databases == nil {
			Databases = make(map[string]*gorm.DB)
		}

		// Clear existed connection to renew
		if Databases[cfg.ConnectionName] != nil {
			Close(Databases[cfg.ConnectionName])
			Databases[cfg.ConnectionName] = nil
		}

		Databases[cfg.ConnectionName] = db
		Print(cfg)
	}

	return nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	sqlDB.Close()
	return err
}

func CloseAll() {
	log.Println("Closing database...")
	for key, db := range Databases {
		sqlDB, _ := db.DB()
		sqlDB.Close()
		log.Printf("DB['%s']\r\n", key)
	}
}

func Print(cfg Config) {
	_connName := cfg.ConnectionName
	if cfg.ConnectionName == "default" {
		_connName = ""
	}

	fmt.Printf("\r\n┌─────── PostgreSQL/%s: Connected ─────────\r\n", cfg.ConnectionName)
	fmt.Printf("| %s: %s\r\n", fmt.Sprintf("DB_POSTGRES%s_DSN", _connName), cfg.DSN)
	fmt.Printf("| %s: %d\r\n", fmt.Sprintf("DB_POSTGRES%s_MAX_IDLE_CONNS", _connName), cfg.MaxIdleConns)
	fmt.Printf("| %s: %d\r\n", fmt.Sprintf("DB_POSTGRES%s_MAX_OPEN_CONNS", _connName), cfg.MaxOpenConns)
	fmt.Printf("| %s: %v\r\n", fmt.Sprintf("DB_POSTGRES%s_CONN_MAX_LIFETIME", _connName), cfg.ConnectionMaxLifetime)
	fmt.Println("└──────────────────────────────────")

}
