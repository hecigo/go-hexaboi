package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/hecigo/goutils"
)

type Compress struct {
	CompressLevel int
}

func (_comp *Compress) Enable(app *fiber.App) error {
	_comp.CompressLevel = goutils.Env("COMPRESS_LEVEL", 0)

	app.Use(compress.New(compress.Config{
		Level: compress.Level(_comp.CompressLevel),
	}))

	_comp.Print()

	return nil
}

func (_comp Compress) Print() {
	var compressLevelText string
	switch _comp.CompressLevel {
	case -1:
		compressLevelText = "LevelDisabled"
	case 1:
		compressLevelText = "LevelBestSpeed"
	case 2:
		compressLevelText = "LevelBestCompression"
	default:
		compressLevelText = "LevelDefault"
	}

	fmt.Println("\r\n┌─────── Middleware/Compress ────────")
	fmt.Printf("| COMPRESS_LEVEL: %d (%s)\r\n", _comp.CompressLevel, compressLevelText)
	fmt.Println("└────────────────────────────────────")
}
