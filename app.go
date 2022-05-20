package main

import (
	"flag"
	"log"

	"github.com/joho/godotenv"

	"hoangphuc.tech/hercules/app/api"
	"hoangphuc.tech/hercules/infra/postgres"
)

var (
	server  = flag.String("", "localhost:3000", "Host & port to listen on")
	profile = flag.String("profile", "", "Environment profile")
	migrate = flag.Bool("migrate", false, "Auto-migrate database")
)

// @title Hercules API Documentation
// @version 1.0
// @description Hercules API Documentation.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	// Parse command-line flags
	flag.Parse()

	// Environment profile
	env := *profile
	if env == "" {
		env = "development"
	}

	// Load .env (view more: https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use)
	godotenv.Load(".env." + env + ".local")
	if env != "test" {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env

	// Initialize API
	_api := api.Init(env)

	// Open PostgreSQL connections
	postgres.OpenDefaultConnection()

	// Detect migrations
	if *migrate {
		postgres.AutoMigrate()
	}

	// Listen on port 3000 as default
	log.Fatal(_api.App.Listen(*server)) // go run app.go -server=localhost:3000

}
