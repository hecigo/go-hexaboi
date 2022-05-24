package api

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"hoangphuc.tech/hercules/app/api/handler"
	"hoangphuc.tech/hercules/app/api/router"
	"hoangphuc.tech/hercules/app/middleware"
	"hoangphuc.tech/hercules/infra/core"

	_ "hoangphuc.tech/hercules/docs"
)

type Api struct {
	App          *fiber.App
	Profile      string
	IsProduction bool
}

func Init(env string) *Api {

	appVersion := core.Getenv("APP_VERSION", "v0.0.0")

	// go run app.go -prod
	isProduction := (env == "prod" || env == "production")

	// Create a new app
	app := fiber.New(fiber.Config{
		AppName:       core.Getenv("APP_NAME", "hercules") + " " + appVersion,
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

	if core.GetBoolEnv("CACHE_ENABLE", false) {
		(&middleware.Cache{}).Enable(app)
	}

	if core.GetBoolEnv("COMPRESS_ENABLE", false) {
		(&middleware.Compress{}).Enable(app)
	}

	if core.GetBoolEnv("CORS_ENABLE", true) {
		(&middleware.CORS{}).Enable(app)
	}

	if core.GetBoolEnv("PPROF_ENABLE", false) {
		(&middleware.Pprof{}).Enable(app)
	}

	// App health check
	app.Get("/", handler.HealthCheck)

	// Create a /v1 endpoint. Just replaces if the frontend is already.
	root := app.Group(core.Getenv("APP_ROOT_PATH", appVersion[0:2]))
	router.RegisterDefaultRouter(root)

	// Always response NotFound at the end of routes
	app.Use(handler.NotFound)

	if core.GetBoolEnv("HTTP_LOG_ENABLE", false) {
		(&middleware.HttpLogger{}).Enable(app)
	}

	return &Api{App: app, Profile: env, IsProduction: isProduction}
}
