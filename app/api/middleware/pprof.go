package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

type Pprof struct {
}

func (p *Pprof) Enable(app *fiber.App) error {
	app.Use(pprof.New())
	p.Print()
	return nil
}

func (p Pprof) Print() {
	fmt.Println("\r\n┌─────── Middleware/Pprof ─────────")
	fmt.Println("| ENABLE: true")
	fmt.Println("└──────────────────────────────────")
}
