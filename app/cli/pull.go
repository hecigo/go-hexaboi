package cli

import (
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/spf13/cobra"
	"hoangphuc.tech/hercules/infra/adapter"
	"hoangphuc.tech/hercules/infra/bigquery"
	"hoangphuc.tech/hercules/infra/postgres"
)

type Pull struct {
	db        *string
	model     *string
	pageIndex *uint
	pageSize  *uint
	repoBrand adapter.BrandRepository
	repoCate  adapter.CategoryRepository
	repoItem  adapter.ItemRepository
}

func (p *Pull) Register(rootApp string, env string, rootCmd *cobra.Command) {

	var selfCmd = &cobra.Command{
		Use:     "pull",
		Short:   "Pull master data",
		Long:    `Pull master data from a repository`,
		Example: "hercules pull -d bigquery -m item",
		Run: func(cmd *cobra.Command, args []string) {

			switch *p.db {
			case "bigquery":
				// Open connections
				bigquery.OpenDefaultConnection()
				defer bigquery.CloseAll()

				postgres.OpenDefaultConnection()
				defer postgres.CloseAll()

				log.Println("Pulling...")
				switch *p.model {
				case "brand":
					p.repoBrand = adapter.BrandRepository{}
					brands, err := p.repoBrand.BQFindAll()
					if err != nil {
						log.Fatal(err)
					}

					if len(brands) == 0 {
						log.Printf("No thing to pull")
					}
					log.Printf("Pulled %d records\n", len(brands))

					count, err := p.repoBrand.BatchCreate(brands)
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("Inserted %d records\n", count)
				case "category":
					p.repoCate = adapter.CategoryRepository{}
					categories, err := p.repoCate.BQFindAll()
					if err != nil {
						log.Fatal(err)
					}

					if len(categories) == 0 {
						log.Printf("No thing to pull")
					}
					log.Printf("Pulled %d records\n", len(categories))

					count, err := p.repoCate.BatchCreate(categories)
					if err != nil {
						log.Fatal(err)
					}
					log.Printf("Inserted %d records\n", count)
				case "item":
					p.repoItem = adapter.ItemRepository{}

					// Fix a number big enough to make it easier :D
					// 9999 x 100 = 999.900 records
					// If there are more than the above number, we need to find another way
					pageIndex := *p.pageIndex
					pageSize := *p.pageSize
					count, err := p.repoItem.BQCount()
					if err != nil {
						log.Fatal(err)
					}
					maxPages := uint(count) / pageSize
					if maxPages*pageSize < uint(count) {
						maxPages++
					}
					if pageIndex < 1 {
						pageIndex = 1
					} else if pageIndex > maxPages {
						pageIndex = maxPages
					}

					for i := pageIndex; i <= maxPages; i++ {
						fmt.Printf("\n---------- Page %d ----------\n", i)
						items, err := p.repoItem.BQFindAll(i, pageSize)
						if err != nil {
							log.Fatal(err)
						}

						l := len(items)
						if l == 0 {
							log.Printf("No thing to pull")
							break
						}
						log.Printf("Pulled %d records.\nInserting...", l)

						count, err := p.repoItem.BatchCreate(items)
						if err != nil {
							log.Println(err.(*pgconn.PgError).Detail)
							log.Fatal(err)
						}
						log.Printf("Inserted %d records\n", count)

						// Sleep to keep the connection safety
						time.Sleep(100 * time.Millisecond)
					}
				default:
					log.Fatal("Data model must be in: category, brand, item")
				}
			default:
				log.Fatal("Database type must be 'bigquery'")
			}

		},
	}

	p.db = new(string)
	p.model = new(string)
	p.pageIndex = new(uint)
	p.pageSize = new(uint)

	selfCmd.Flags().StringVarP(p.db, "database", "d", "", "Database type: bigquery (required)")
	selfCmd.Flags().StringVarP(p.model, "model", "m", "", "Data model, similar BigQuery table (required)")
	selfCmd.Flags().UintVarP(p.pageIndex, "pageIndex", "p", 1, "Page index")
	selfCmd.Flags().UintVarP(p.pageSize, "pageSize", "s", 10, "Page size")
	selfCmd.MarkFlagRequired("db")
	selfCmd.MarkFlagRequired("model")

	rootCmd.AddCommand(selfCmd)
}
