package bigquery

import (
	"fmt"
	"log"

	"gorm.io/driver/bigquery"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"hoangphuc.tech/hercules/infra/core"
)

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

	db, err := gorm.Open(bigquery.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Always use singular table name
		},
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Append to DB map
	if databases == nil {
		databases = make(map[string]*gorm.DB)
	}

	// Generate the default name as a key for the DB map
	if connName == "" {
		connName = "default"
	}

	// Clear existed connection to renew
	if databases[connName] != nil {
		Close(databases[connName])
		databases[connName] = nil
	}

	databases[connName] = db
	return nil
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
