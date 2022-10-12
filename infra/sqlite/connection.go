package sqlite

import (
	"context"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	apmsqlite "go.elastic.co/apm/module/apmgormv2/v2/driver/sqlite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"hoangphuc.tech/go-hexaboi/infra/core"
)

type Config struct {
	ConnectionName        string
	DSN                   []string
	DSNRelicas            []string
	MaxIdleConns          int
	MaxOpenConns          int
	ConnectionMaxLifetime time.Duration
}

var databases map[string]*gorm.DB = make(map[string]*gorm.DB)

// Open a connection to the default database
func DB() *gorm.DB {
	if len(databases) == 0 {
		panic("No database found")
	}
	return databases["default"]
}

// Open a connection to the database with context. It is commonly used to monitor APM.
func DBWithContext(ctx context.Context) *gorm.DB {
	return DB().WithContext(ctx)
}

// Open a connection to the database with specified name.
func DBByName(name string) *gorm.DB {
	if len(databases) == 0 {
		panic("No database found")
	}
	return databases[name]
}

// Open a connection to the database with specified name and context. It is commonly used to monitor APM.
func DBByNameWithContext(ctx context.Context, name string) *gorm.DB {
	return DBByName(name).WithContext(ctx)
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

// Open a connection with name specified from ENV
func OpenConnectionByName(connName string) error {
	_connName := ""     // Emtpy is default connection
	if connName != "" { // Add _ as a prefix to the connection name
		_connName = "_" + connName
	}

	dsn := core.Getenv(fmt.Sprintf("DB_SQLITE%s_DSN", _connName), "")
	dsnReplicas := core.Getenv(fmt.Sprintf("DB_SQLITE%s_DSN_REPLICAS", _connName), "")

	// Connection pool config
	maxIdleConns := core.GetIntEnv(fmt.Sprintf("DB_SQLITE%s_MAX_IDLE_CONNS", _connName), 5)
	maxOpenConns := core.GetIntEnv(fmt.Sprintf("DB_SQLITE%s_MAX_OPEN_CONNS", _connName), 20)
	connMaxLifetime := core.GetDurationEnv(fmt.Sprintf("DB_SQLITE%s_CONN_MAX_LIFETIME", _connName), 30*time.Minute)

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

		db, err := gorm.Open(openDialector(cfg.DSN[0]), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // Always use singular table name
			},
			PrepareStmt: true,
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
					dsnDialectors = append(dsnDialectors, openDialector(cfg.DSN[i]))
				}
			}

			dbr = dbresolver.Register(dbresolver.Config{
				Sources: dsnDialectors,
			})
		}
		if len(cfg.DSNRelicas) > 0 {
			for _, dsn := range cfg.DSNRelicas {
				if dsn != "" {
					dsnRelicasDialectors = append(dsnRelicasDialectors, openDialector(dsn))
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
	for _, db := range databases {
		// By pass error to continue the next connection
		Close(db)
	}
	log.Info("Closed total SQL Server connections...")
	return nil
}

func Print(cfg Config) {
	_connName := cfg.ConnectionName
	if cfg.ConnectionName == "default" {
		_connName = ""
	}

	fmt.Printf("\r\n┌─────── SQLite/%s: Connected ─────────\r\n", cfg.ConnectionName)
	fmt.Printf("| DB_SQLITE%s_DSN: %s\r\n", _connName, cfg.DSN)
	fmt.Printf("| DB_SQLITE%s_MAX_IDLE_CONNS: %d\r\n", _connName, cfg.MaxIdleConns)
	fmt.Printf("| DB_SQLITE%s_MAX_OPEN_CONNS: %d\r\n", _connName, cfg.MaxOpenConns)
	fmt.Printf("| DB_SQLITE%s_CONN_MAX_LIFETIME: %v\r\n", _connName, cfg.ConnectionMaxLifetime)
	fmt.Println("└───────────────────────────────────────────────")

}

func openDialector(dsn string) gorm.Dialector {
	if core.GetBoolEnv("ELASTIC_APM_ENABLE", true) {
		return apmsqlite.Open(dsn)
	}
	return sqlite.Open(dsn)
}
