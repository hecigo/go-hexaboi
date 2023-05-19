package api

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/cobra"
	"go.elastic.co/apm/module/apmfiber"

	"github.com/hecigo/goredis"
	"github.com/hecigo/goutils"
	"hecigo.com/go-hexaboi/app/api/handler"
	"hecigo.com/go-hexaboi/app/api/middleware"
	"hecigo.com/go-hexaboi/app/api/router"
	"hecigo.com/go-hexaboi/infra/elasticsearch"
	"hecigo.com/go-hexaboi/infra/orientdb"
	"hecigo.com/go-hexaboi/infra/postgres"
	"hecigo.com/go-hexaboi/infra/redis"

	_ "hecigo.com/go-hexaboi/docs"
)

type API struct {
	App          *fiber.App
	Profile      string
	IsProduction bool
}

var listen string

func Register(rootApp string, env string, rootCmd *cobra.Command) {
	var selfCmd = &cobra.Command{
		Use:     "serve",
		Short:   "Start " + rootApp + " RESTful API",
		Long:    rootApp + ` RESTful API provides inventory data and requests for other services`,
		Example: "gohexaboi serve -l localhost:3000",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize new API
			api := New(env)

			// Open PostgreSQL connections
			postgres.OpenDefaultConnection()
			defer postgres.CloseAll()

			// Initialize Elasticsearch clients
			elasticsearch.OpenDefaultConnection()

			// Initialize Redis clients
			goredis.Open()
			redis.EnableSession() // Enable session store for Redis
			defer goredis.Close()

			// Open OrientDB connections
			orientdb.OpenDefaultConnection()
			defer orientdb.CloseAll()

			// Listen serves HTTP requests from the given addr
			return api.App.Listen(listen)
		},
	}
	selfCmd.Flags().StringVarP(&listen, "listen", "l", "localhost:3000", "Listen serves HTTP requests from the given addr")

	rootCmd.AddCommand(selfCmd)
}

func New(env string) *API {

	// go run app.go -prod
	isProduction := (env == "prod" || env == "production")

	// Create a new app
	app := fiber.New(fiber.Config{
		AppName:       goutils.AppName() + " " + goutils.AppVersion(),
		StrictRouting: false,
		CaseSensitive: false,
		Prefork:       isProduction,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		ErrorHandler:  handler.DefaultError,
	})

	// Initialize middlewares
	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
	}))

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	if goutils.Env("COMPRESS_ENABLE", false) {
		(&middleware.Compress{}).Enable(app)
	}

	if goutils.Env("CORS_ENABLE", true) {
		(&middleware.CORS{}).Enable(app)
	}

	// App health check. It must be by pass some middlewares.
	app.Get("/", handler.HealthCheck)

	if goutils.Env("CACHE_ENABLE", false) {
		(&middleware.Cache{}).Enable(app)
	}

	if goutils.Env("PPROF_ENABLE", false) {
		(&middleware.Pprof{}).Enable(app)
	}

	if goutils.Env("HTTP_LOG_ENABLE", false) {
		(&middleware.HttpLogger{}).Enable(app)
	}

	if goutils.Env("AUTH_ENABLE", true) {
		(&middleware.Auth{}).Enable(app)
	}

	// APM
	if goutils.Env("ELASTIC_APM_ENABLE", true) {
		app.Use(apmfiber.Middleware())
	}

	// Create a /v1 endpoint. Just replaces if the frontend is already.
	root := app.Group(goutils.APIRootPath())
	router.RegisterDefaultRouter(root)

	// Always response NotFound at the end of routes
	app.Use(handler.NotFound)

	return &API{App: app, Profile: env, IsProduction: isProduction}
}
