package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"hoangphuc.tech/go-hexaboi/app/api"
	"hoangphuc.tech/go-hexaboi/app/cli"
	"hoangphuc.tech/go-hexaboi/infra/core"
)

var (
	env = flag.String("env", "", "Environment profile")
)

func main() {
	// Load .env (view more: https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use)
	if env == nil || *env == "" {
		*env = "development"
	}
	godotenv.Load(".env." + *env + ".local")
	if *env != "test" {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + *env)
	godotenv.Load() // Load the default environment

	// Initialize logger
	core.InitLogger()

	appName := core.Getenv("APP_NAME", "Gohexaboi")
	appVersion := core.Getenv("APP_VERSION", "v0.0.0")

	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:     strings.ToLower(appName),
		Short:   appName + " is HPI - Golang Dora Service",
		Long:    appName + ` is the Golang Dora Service built by HPI.Tech`,
		Version: appVersion,
	}
	rootCmd.PersistentFlags().StringVarP(env, "env", "e", "", "environment profile name")

	// Register commands
	api.Register(appName, *env, rootCmd)
	(&cli.Migrate{}).Register(appName, *env, rootCmd)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
