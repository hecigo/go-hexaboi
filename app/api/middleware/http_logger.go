package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"hoangphuc.tech/go-hexaboi/infra/core"
)

type HttpLogger struct {
	Format string
}

func (_logger *HttpLogger) Enable(app *fiber.App) error {
	_logger.Format = core.Getenv("HTTP_LOG_FORMAT",
		`[${time}] | PID: ${pid} | ${latency} | `+
			`${reqHeader:X-Hpi-App-Version} ${status} ${method} ${path} ?${queryParams} ${body}`) + "\r\n"

	app.Use(logger.New(logger.Config{

		// Why does use the following value?
		// View more: https://programming.guide/go/format-parse-string-time-date-example.html
		TimeFormat: "2006-01-02 15:04:05",
		Format:     _logger.Format,
	}))

	_logger.Print()

	return nil
}

func (_logger HttpLogger) Print() {
	fmt.Println("\r\n┌─────── Middleware/HTTP Log ────────")
	fmt.Printf("| HTTP_LOG_FORMAT: %s", _logger.Format)
	fmt.Println("└────────────────────────────────")
}
