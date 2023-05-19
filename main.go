package main

import (
	log "github.com/sirupsen/logrus"

	"strings"

	"github.com/hecigo/goutils"
	"github.com/spf13/cobra"

	"hecigo.com/go-hexaboi/app/api"
	"hecigo.com/go-hexaboi/app/cli"
)

func main() {
	env := goutils.QuickLoad()

	// Load client secret keys
	goutils.EnableAPISecretKeys()

	// rootCmd represents the base command when called without any subcommands
	appName := goutils.AppName()
	var rootCmd = &cobra.Command{
		Use:     strings.ToLower(appName),
		Short:   appName + " is HPI - Golang Hexagonal Boilerplate",
		Long:    appName + ` is the Golang Hexagonal Boilerplate Service built by HPI.Tech`,
		Version: goutils.AppVersion(),
	}
	rootCmd.PersistentFlags().StringVarP(&env, "env", "e", "", "environment profile name")

	// Register commands
	api.Register(appName, env, rootCmd)
	(&cli.Migrate{}).Register(appName, env, rootCmd)
	(&cli.Pull{}).Register(appName, env, rootCmd)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
