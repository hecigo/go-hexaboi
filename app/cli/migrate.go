package cli

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"hoangphuc.tech/go-hexaboi/infra/postgres"
)

type Migrate struct {
	db    *string
	model *string
}

func (m *Migrate) Register(rootApp string, env string, rootCmd *cobra.Command) {

	var selfCmd = &cobra.Command{
		Use:     "migrate",
		Short:   "Migrate database",
		Long:    `A RESTful API mode of IMS built with gofiber`,
		Example: "gohexaboi migrate -d postgres",
		Run: func(cmd *cobra.Command, args []string) {
			switch *m.db {
			case "postgres":
				// Open PostgreSQL connections
				postgres.OpenDefaultConnection()
				defer postgres.CloseAll()

				postgres.AutoMigrate(*m.model)
			default:
				log.Fatal("Database type must be 'postgres'")
			}

		},
	}

	m.db = new(string)
	m.model = new(string)

	selfCmd.Flags().StringVarP(m.db, "database", "d", "postgres", "Database type: postgres (required)")
	selfCmd.Flags().StringVarP(m.model, "model", "m", "", "Data model name want to migrate")
	selfCmd.MarkFlagRequired("database")
	selfCmd.MarkFlagRequired("model")

	rootCmd.AddCommand(selfCmd)
}
